//sessions manager
package main

import(
	"time"
	"log"
	"strconv"
	"strings"
	"net/http"
)

//some power of 2, just cause
//TODO for testing functionality of growing slice, start it at 2
var initialSize int = 1024

//should eventuall reimpliment the sessions slice like this
//sessions [][]session
//
//  [ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ]
//   |  |  |  |  |  |  |  |  |  |  |
//  [ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ]
//  [ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ]
//  [ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ]
//  [ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ]
//  [ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ]
//  [ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ]
//  [ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ]
//
//like clark was saying that one time about array lists/file systems
//that way you just resize the top array, and copy over a pointer to
//that data. in fact get it working this current way then do this right
//after
//id will have to have 2 parts to it, the first with which list its in,
//second with position in that list


//session ID is just the index in the sessions slice + 1million

//this is not implimented with a hash table (well it TECCHNICALLY is). This will dish 
//out numbers in order starting at 1,000,000. So the first person to ever log in gets 
//1,000,000 as their session id, the second person gets 1,000,001, and so on. These 
//functions will then subtract 1 million from each of those and that will be the index
//position in the SessionManager.sessions slice. 
//
//the sessions slice will start at a size powe of 2, and continually grow by powers of 2 as
//needed, but will never shrink. sm.SessionSweep will go through and check all sessions
//every 48 hours (i guess at around 5am when traffic would be dead) and check the 
//lastActive time for each session in the slice. If it has been longer than 24 hours (or
//whatever), then it will be cleared/ended/deleted, and the lastDeleted sessions 
//nextBlankSpot will be set to that spot, and the current session spots nextBlankSpot
//will be set to -1. it then keeps traversing, and finds the next
//session that is overdue, and does the same. So now, whenever StartSession uses 
//nextRawId, if the next sessions nextBlankSpot is -1, then it will assume the next 
//bkank spot is farthestPlacementYet, and set it accordingly, and incriment it and
//lowestPlacement. 
//
//Whole process is still a little ambiguous but it'll make more sense
//once i start typing it out
//
//farthestPlacementYet will never go down in value, it can only go up
//
//heavily fucking commented for now, will delete some later, this is fucking annoying

type SessionManager struct {		

	//this will be the master slice of sessions
	sessions []session

	//holds current size of sessions. power of 2
	currentSize int
	
	//where newID or nextSpot or whatever will first look for a free sessions id/spot
	lowestPlacement int
	
	//farthest placement a session has been, never get decrimented(word?)
	farthestPlacementYet int 
	
	//needs to be reset to -1 everytime SweepSessions is done, so the next time
	//it starts it wont cause some huge ass bugs
	lastDeleted int //for use in session sweep to chain together the sweeped spots
	
	//holds the amount of time that a session is allowed to be inactive for
	//end the session if its longer
	inactiveAllowance time.Time

	//hould probably toss a mutex somewhere in here for SessionSweep
	
	//so that shit doesnt go sour if being resized by 2 routines
	beingResized bool
}

type session struct {
	//for use with sweeper and chaining empty spots before farthestPlacementYet
	active bool
	nextBlankSpot int

	handle string
	lastActive time.Time
}

func (sm *SessionManager) initialize(){
	sm.sessions = make([]session, initialSize)
	sm.currentSize = initialSize
	sm.lowestPlacement = -1
	sm.farthestPlacementYet = 0
	sm.lastDeleted = -1
	sm.inactiveAllowance = time.Now() //BUG BUG BUG setting this for compiler
}

//recieves the id and the handle that were stored in the users cookie and returns
//if they are currently in a session
//
//check for not only if the sessions is active, but also if the handle is the same one.
//session spot may be active because some other user took that session spot
func (sm *SessionManager)VerifySession(id int, handle string) bool {
	s := sm.sessions[id-1000000]
	if s.active && s.handle == handle{
		return true	
	}
	log.Println("---------------verify sess", id, handle, false)
	return false
}

//takes the cookie as a parameter, call ParseCookie on it, then calls VErifySEssion on it
//made a seperate function so that handler in main call either this or VerifySession, 
//depending on wether it already called parseCookie
func (sm *SessionManager) VerifySessionCookie(c *http.Cookie) bool {
	id, handle, err := sm.ParseCookie(c)
	if err != nil{
		log.Println("sm,verifySession) smParseCookie returned err", err)
		return false	
	}
	return sm.VerifySession(id, handle)	
}


//retuns the session id. the cookie is made with this id in h.handleLogin
//calls sm.nextRawId() to get the index, sets all its values in the sessions
//slice, then returns the un-raw or cooked (ha) id, which is just the index
//+ 1000000
func (sm *SessionManager) StartSession(handle string) int{
	rid := sm.nextRawId()
	sm.sessions[rid].active = true
	sm.sessions[rid].nextBlankSpot = -1 //for easy sweeping, just do it now
	sm.sessions[rid].handle = handle
	sm.sessions[rid].lastActive = time.Now()
	return rid+1000000
}

//sets the sessions at id-1000000(which is the index/rawid) false
//also sets the nextBlankSpot to -1, which is pre-emptive for sm.SessionSweep
func (sm *SessionManager) EndSession(id int){
	index := id-1000000
	sm.sessions[index].active = false
	sm.sessions[index].nextBlankSpot = -1
}

//this function will return the ID and also update the chains 
//the calling function (StartSession) will set the rest of the values in the session
//RawID means staight index. Unraw would then mean with the 1000000 added
//
//lowestPlacement is only -1 if there are no free spots under the farthestPlacemnetYet,
//so, first thing this does is check for it to be -1, then sets it accordingly to either
//frathestPlacemnetYet or lowestPlaceMent. updates whichever is set
//
//if setting to farthestPlacementYet, check if the sessions slice needs to be resized
func (sm *SessionManager) nextRawId() int{
	index := 0
	if sm.lowestPlacement == -1 {
		index = sm.farthestPlacementYet
		sm.farthestPlacementYet++
		sm.resizeCheck()
	} else {											
		index = sm.lowestPlacement  
		sm.lowestPlacement = sm.sessions[index].nextBlankSpot// this assumes the chaining
	}									//is done correctly, and the last open spot
	return index						//underneath farthestPlacemnetYet is set to -1
}										//(the nextBlankSpot is -1)



//first checks if the sm is already being resized (meaning this was called and passed
// and called sm.resize()). then if its not, check if 75 percent of the slice is full, 
//and resize if it is.
//
//TODO should find a more efficient way than converting both to floats (which is needed by
//go to do the arithmetic), although maybe its not that inefficient
func (sm *SessionManager) resizeCheck(){	
	if sm.beingResized{
		return	
	}
	if .75*float64(sm.currentSize)<= float64(sm.farthestPlacementYet){
		sm.beingResized = true
		go sm.resize()	
	}
}

///if the resize is needed, first copy the sessions slice into a new slice, THEN
//lock the mutex, change the pointer to point to secind one, then unlock
//check if sm.farthestPlacementYet is 75% through slice, then resize if it is
//has to copy all data, will probably be really fucking taxing later on
func (sm *SessionManager) resize(){	
	log.Println("-----(sm.resize) havent written this, nothing being done-------------")
}

//call this as a goroutine
//if planning to keep people logged in indefinitly then this should only
//check for the session.active to be false, other wise, if logging people
//out after a certain time, this also has to compare the time.Time shit
//it should do the time.Time thing anyways around everyweek or so, since
//people sign in from random computers some times, then ill just have
//active sessions that will never get used
//but the session.active sweep (for actual sign outs) should happen often
func (sm *SessionManager) SessionSweep(){
	
}

//TODO not perfect error checking, if someone fucks with the cookie and sends
//some messed up stuff back this wont work
//
//parses the cookie, which is set like so: "193736594|jockster", id|handle
func (sm *SessionManager) ParseCookie(c *http.Cookie) (int, string, error){
	parts := strings.Split(c.Value, "|")
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", err	
	}
	log.Println("-------------------------parsecookie, ", id, parts[1])		
	return id, parts[1], nil
}












