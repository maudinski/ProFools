// TODO need to get rid of majin bu symbols from html files created with python. Use clarks ___ maybe or idk
// some logic in EndSessioion/signout needs to be re-thought
// no sure how but need a favicon
// if specific content type isnt working, it may be the
// the log.Println()'s are obviously a bottle neck, get rid of them (not the error ones)
// in pageCreater, make it s___e out comments and white unnessecarry white space
// might be good to change file permisions for all these files to add an extra layer of security, not sure though
// vegeta load testing utility
// check out the web slot in useful ___.txt
// unless i can combat every case, learn to recover from panic...
// idk if this method of doing to urls is efficient. gotta load test
//get rid of all log.fatals
//end TODO s

//main function is at the bottom
package main

import (
	"log"
	"net/http"
	"strings"
	"io/ioutil"
	"strconv"
	"errors"
	"github.com/maudinski/sesh"
	"github.com/maudinski/mysqlTM"
)
//all the globals used around the program. All are only used in this file for now,
//except osa
var files Files
var pc content
var osa outsideStructureAbstractor
var sm *sesh.SessionManager
var am *mysqlTM.TableManager //stands for account manager, since it will be used to manage
					//the users

//meant so if i change a p___word, or a directory structure, or my table/database
//names in sql, or anything that is outside the scope of go code, i wont have to 
//sift through a ___ton of code to change it. just gotta change what osa.initialize() 
//initializes ___ it. used all over the place
type outsideStructureAbstractor struct{
	loaderFilesDir string
	picturesDir string	
	postHtmlFN string
	
	dbUser string
	dbPass string // probably not the safest thing to do
	dbHostAndPort string
	
	db string  //change this later so less confusion cause 2 databases. this is for posts

	users_db string //for users and p___words
	users_table string
	
	defaultPortPic string
}

//initializes the osa
func (osa *outsideStructureAbstractor) initialize(){
	osa.loaderFilesDir = "loaderfiles"
	osa.picturesDir = "picturesamples"
	osa.postHtmlFN = "post.html"
	osa.dbUser = "root"
	osa.dbPass = "test"
	osa.dbHostAndPort = ""
	osa.db = "test_db2"
	osa.users_db = "test_users_db2"
	osa.users_table = "users"
	osa.defaultPortPic = "nothing"
}

//cutsom handler. Nothing in it, just necessary for http.Handle
type handler struct {
	//TODO should probably have the globals in here, except OSA
}

//paths will be like this: 
//	/exhibit/mix/1
//	/exhibit/literature/0
//	/js/exhibit.js
//	/css/exhbit.css

//this makes handler impliment the interface that http.Handle needs
//this is the main logic of all requests. All requests start here. 
//3 stages to this function so far:
//if theyre requesting "/", then they want the homepage. send it and return.
//if theyre using the "POST" method, then theyre trying to sign in or post a post or 
//	something like that. Handle that in another function cause thats not the most common
//	method of request
//if its not the other 2, then its a simple get. Note: url patterns are not
//	representitive of directories, theyre just used for internal ___, so we know what
//	to send back. something like "/exhbiit/mix/7" will be parsed into a slice like this:
//	["exhbibit", "mix", "7"] (thats pathParts). switch statement in here evaluates
// 	pathParts[0]
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	//log.Println(r.URL.Path, r.Method)
	if r.URL.Path == "/"{								
		h.sendExhibit(w, r, 0, 0)
		return;
	}
	pathParts := strings.Split(r.URL.Path, "/")[1:]//Split returns a blank at spot [0] if
	//factor this out since most request wont be POST
	if r.Method == "POST" {
		h.handlePOST(w, r, pathParts)
		return
	}
	//i wont factor this out since most requests will be GET
	switch pathParts[0] {
		case "exhibit":			
			h.handleExhibit(w, r, pathParts)
		case "port":
		case "post"://this is looking for a post
		case "css":	
			h.handleJsAndCss(w, pathParts[1], "text/css")
		case "picture":	
			h.handlePic(w, osa.picturesDir+"/"+pathParts[1])
		case "js":
			h.handleJsAndCss(w, pathParts[1], "application/js")
		case "signin":
			w.Write(files.fd["signin.html"].data)
		case "signout":
			h.handleSignout(w, r)
		case "signup":
			h.handleSignup(w, r)
		default:
			h.sendExhibit(w, r, 0, 0)
			log.Println("(servehttp) hit the deafult:", r.URL.Path)
	}
	
}

//handles ___ with method=="POST". evaluates pathParts, p___ed to it, in the switch
//statemnet.
//there is 2 different outcomes of localhost:8080/signup, depending on wether the method
//is POST or GET. this is the Post outcome
func (h *handler) handlePOST(w http.ResponseWriter, r *http.Request, pathParts []string){
	
	switch  pathParts[0]{
		case "verify-signin":
			h.handleSignin(w, r)
		case "logout":
			h.handleSignout(w, r)
		case "postpost":
		case "signup":
			log.Println("signup called")
			h.verifySignup(w, r)
		default:
			log.Println("Hit the default for POST")
			h.sendExhibit(w, r, 0, 0)
	}
}

//not to be confused with verifySingup, all this does is if they click the signup link at
//the top of the page, it sends back the signup form. verifySignup does all the logistical
//stuff, this function is its ___
func (h *handler) handleSignup(w http.ResponseWriter, r *http.Request) {
	w.Write(files.fd["signup.html"].data)	
}

//TODO change comments, not like david is gonna do anything tho
//handles signup. First checks if there is any stored cookies. If there is, end any session
//___ociated with it, and set isCookie to true (used a little later). Parse the form, and
//get the handle, p___word, screen name, and email, and p___ them to am.Signup. If that 
//returns an error, send signupFailed.html and return. Otherwise the signup was succesful.
//start a new session, get a new value for cookie (should really make that a method for
//sessiosManager or something). If there was a cookie already then reset its value, or else
//create a new cookie and set it. Send back exhibitSignedIn.html manually
func (h *handler) verifySignup(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	ha, p, e := r.Form.Get("handle"), r.Form.Get("password"), r.Form.Get("email")
	s := ha //TODO write a method for TableManager that takes in the request
	err := am.Insert(ha, p, s, e)
	if err != nil {	
		log.Println(err)
		w.Write(files.fd["signupFailed.html"].data)
		return
	}
	sm.StartSession(w, ha)
	w.Write(files.fd["exhibitTopSignedIn.html"].data)
	w.Write(pc.exhibits[0][0])//TODO write that ___ing javascript man
	w.Write(files.fd["exhibitBottom.html"].data)
}

//this works for now. Just ends the session, then sends back the home page.
//TODO needs to redirect to homepage, so that the url changes and they dont accientally
//try and refresh the signout
func (h *handler) handleSignout(w http.ResponseWriter, r *http.Request) {
	sm.EndSession(r)
	h.sendExhibit(w, r, 0, 0)
}

//dirty but works. verifies the signins using am, then if its valid, starts a new session,
//creates a new cookie, and manually sends back the home page. I want to try and find a 
//way to send back the page they were previously on. Maybe store that url with javascript
//right when they hit the "signin" link, idk
//TODO needs to redirect with home page AND signed in (so cant just call sendExhibit)
//so that the url changes and they dont try and refresh
func (h *handler) handleSignin(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	handle, password := r.Form.Get("handle"), r.Form.Get("password")
	valid, err := am.Verify(handle, password)
	if err != nil {
		w.Write([]byte("internal problem"))//TODO this is bad
	}
	if !valid {
		w.Write(files.fd["failedSignIn.html"].data)
		return
	}
	sm.StartSession(w, handle)
	w.Write(files.fd["exhibitTopSignedIn.html"].data)
	w.Write(pc.exhibits[0][0])//TODO entire html change with javascript so this wont happen
	w.Write(files.fd["exhibitBottom.html"].data)
}

//these requests are not manually made by the browser, but are linked in the html files
//error checking anyways cause technically they could still manually request for files 
//that arent there. send back nothing if they ask for random ___
//some sites use seperate servers for js and css
func (h *handler) handleJsAndCss(w http.ResponseWriter, f string, ctype string){	
	w.Header().Add("Content-Type", ctype)
	if _, ok := files.fd[f]; ok { //check above for worry on this
		w.Write(files.fd[f].data)
		return;
	} else {
		log.Println("js or css file error:", f)	
	}
}

//same basic comments as handleJsAndCss, except this one has to make a file read every
//time. Would be stupid to try and keep a million pictures cached. Maybe sometime, keep
//just the top posts pictures cached. who knows
func (h *handler) handlePic(w http.ResponseWriter, path string){
	picData, err := ioutil.ReadFile(path)
	if err == nil {
		w.Write(picData)
	} else {
		log.Println("picData not read right:", err)
	}
}

//all if statments are just error checking. Basically just taking the pathParts p___ed
//to it, making sure they are valid, and p___ing that to h.sendExhibit
//browser requests are made like this: /exhibit/literature/2
//TODO do a redirect, so as to send back a full url for the navbar javascript
//TODO this function could look a lot better
func(h *handler)handleExhibit(w http.ResponseWriter, r *http.Request, pathParts []string) {
	length := len(pathParts)
	log.Println("length is :", length, "|path is:",pathParts)
	if length == 1{
		h.sendExhibit(w, r, 0, 0)//send /mix/0
		return
	}
	exhibit := pathParts[1]
	exhibitNum, er := getIndex(exhibit)
	if er != nil { //if they ___ed up the exhibit, send 404
		h.send404(w, r)			//everything else just gets the 'index' (mix/0)
		return	
	}
	if length == 2 {
		h.sendExhibit(w, r, exhibitNum, 0)//send /whateveritis/0
		return
	}
	pageNum, err := strconv.Atoi(pathParts[2])
	if err != nil {
		h.sendExhibit(w, r, exhibitNum, 0)	
		return
	}
	h.sendExhibit(w, r, exhibitNum, pageNum)
}	

//this comment is random af
//TODO  change comments
//for now im doing these if else but they might be unecessary. See if i can use javascript
//to tell if a cookie is on browser or not then dynamically change the top right to either
//say "sign up sign in" or "sabio667 sign out" or something. maybe if signout sets the 
//value "" then it can just make sure the cookie isnt blank
//and YEAH javascript can handle cookie ___ https://www.w3schools.com/js/js_cookies.asp
//-------------------------------------------------------------------------------------
//(below sounds expensive on request times. if it ever gets a lot of traffic then yeah)
//...unless i impliment that cool side panel of "things that might interest you"
//which i ___ing should
//unless unless i request that info with javascript too

//sends 404 page
func (h *handler) send404(w http.ResponseWriter, r *http.Request){
	if sm.VerifySession(r) != nil{
		w.Write(files.fd["exhibitTopSignedOut.html"].data)
	} else { //TODO ___ing gross man, change the html/javascript already
		w.Write(files.fd["exhibitTopSignedIn.html"].data)
	}
	w.Write(files.fd["404.html"].data)
}

//TODO change comments
//when creating cookies you give them names, so that you can give multiple ones.
//by saying cookie, err := r.Cookie("session"), im asking for the cookie named
//session. Returns an error if its not there
//if theres an error or the session is invalid, send back TopSignedOut, else
//send signedIntTop. then send the rest
func(h *handler)sendExhibit(w http.ResponseWriter, r *http.Request, exh int, page int){
	if sm.VerifySession(r) != nil {
		w.Write(files.fd["exhibitTopSignedOut.html"].data)
	} else {//TODO ___ing disgusting, change the ___ing html/javascript format god dammit
		w.Write(files.fd["exhibitTopSignedIn.html"].data)
	}

	w.Write(pc.exhibits[exh][page])
	w.Write(files.fd["exhibitBottom.html"].data)
}

func main() {
	
	//TODO make these return errors and handle the log.Fatals out here
	osa.initialize()
	
	sm = sesh.NewSM()
	am, err := mysqlTM.NewTM(osa.dbUser, osa.dbPass, osa.dbHostAndPort, osa.users_db, 
				osa.users_table, "handle", "password", "screenname", "email")//BUG 
	if err != nil {
		log.Fatal(err)	
	}
	err = am.SetupVerify("handle", "password")
	if err != nil {
		log.Fatal(err)	
	}

	files.initialize()
	//TODO pc still kind of nasty, maybe revert it? idk how this could work un-nastily
	pc, err := NewC(osa.dbUser, osa.dbPass, osa.dbHostAndPort, osa.db, 
									osa.loaderFilesDir+"/"+osa.postHtmlFN)
	if err != nil {
		log.Fatal("someshit",err)	
	}


	//TODO these should take the variables from osa as parameters, not use it globally
	//that means change all those files
	//runs as a go routine
	pc.updateForever()
	
	//start the server
	http.Handle("/", new(handler))
	log.Println("listening....")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

//gets index from AllExhibits global (i think its in content.go) and returns it. used
//in h.handleExhibit
//TODO this is kinda ___ed up man. Figure out a better way. This should be a method in
//content.go. Fucking prick
func getIndex(exhibit string) (int, error){

	for i, val := range AllExhibits{
		if exhibit == val {
			return i, nil	
		}	
	}
	
	return 0, errors.New("fuck")
}
