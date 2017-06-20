// TODO need to get rid of majin bu symbols from html files created with python. Use clarks shit maybe or idk
// some logic in EndSessioion/signout needs to be re-thought
// no sure how but need a favicon
// if specific content type isnt working, it may be the
// the log.Println()'s are obviously a bottle neck, get rid of them (not the error ones)
// in pageCreater, make it scrape out comments and white unnessecarry white space
// might be good to change file permisions for all these files to add an extra layer of security, not sure though
// vegeta load testing utility
// check out the web slot in useful shit.txt
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
)
//all the globals used around the program. All are only used in this file for now,
//except osa
var files Files
var pageContent content
var osa outsideStructureAbstractor
var sm SessionManager
var am AccountManager

//meant so if i change a password, or a directory structure, or my table/database
//names in sql, or anything that is outside the scope of go code, i wont have to 
//sift through a fuckton of code to change it. just gotta change what osa.initialize() 
//initializes shit it. used all over the place
type outsideStructureAbstractor struct{
	loaderFilesDir string
	picturesDir string	
	postHtmlFN string
	
	dbUser string
	dbPass string // probably not the safest thing to do
	dbHostAndPort string
	
	db string  //change this later so less confusion cause 2 databases. this is for posts

	users_db string //for users and passwords
	users_table string

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
}

//cutsom handler. Nothing in it, just necessary for http.Handle
type handler struct {
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
//	representitive of directories, theyre just used for internal shit, so we know what
//	to send back. something like "/exhbiit/mix/7" will be parsed into a slice like this:
//	["exhbibit", "mix", "7"] (thats pathParts). switch statement in here evaluates
// 	pathParts[0]
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	log.Println(r.URL.Path, r.Method)
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
		default:
			h.sendExhibit(w, r, 0, 0)
			log.Println("(servehttp) hit the deafult:", r.URL.Path)
	}
	
}

//handles shit with method=="POST". evaluates pathParts, passed to it, in the switch
//statemnet
func (h *handler) handlePOST(w http.ResponseWriter, r *http.Request, pathParts []string){
	
	switch  pathParts[0]{
		case "verify-signin":
			h.handleSignin(w, r)
		case "logout":
			h.handleSignout(w, r)
		case "postpost":
		case "signup":
		default:
			log.Println("Hit the default for POST")
			h.sendExhibit(w, r, 0, 0)
	}
}

//this works for now. Just ends the session, then sends back the home page.
//TODO delete cookie, or set the value to "". then handle sign in will first check for
//a current cookie, and modify that
//TODO make a sm.handleSignoutCookie function, that verifies the cookie by calling 
//verifysessioncookie (maybe)
//TODO needs to redirect to homepage, so that the url changes and they dont accientally
//try and refresh the signout
func (h *handler) handleSignout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session")
	if err == nil {
		id, _, err := sm.ParseCookie(c)
		if err == nil {
			sm.EndSession(id)
		}
	}	
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
	handle := r.Form.Get("handle")
	password := r.Form.Get("password")
	valid := am.VerifySignin(handle, password)
	if valid {
		id := sm.StartSession(handle)
		val := strconv.Itoa(id)+"|"+handle
		c := &http.Cookie{Name: "session", Value: val}
		http.SetCookie(w, c)
		//cookie is stored in request, so sending a page now will still send
		//a exhibitTopSignedOut.html, since they check r.Cookie. so its dirty
		//but just manually send the pages
		w.Write(files.fd["exhibitTopSignedIn.html"].data)
		w.Write(pageContent.exhibits[0][0])
		w.Write(files.fd["exhibitBottom.html"].data)
	} else {
		w.Write(files.fd["failedSignIn.html"].data)	
	}
}
//these requests are not manually made by the browser, but are linked in the html files
//error checking anyways cause technically they could still manually request for files 
//that arent there. send back nothing if they ask for random shit
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

//all if statments are just error checking. Basically just taking the pathParts passed
//to it, making sure they are valid, and passing that to h.sendExhibit
//browser requests are made like this: /exhibit/literature/2
//TODO do a redirect, so as to send back a full url for the navbar javascript
func(h *handler)handleExhibit(w http.ResponseWriter, r *http.Request, pathParts []string) {
	length := len(pathParts)
	log.Println("length is :", length, "|path is:",pathParts)
	if length == 1{
		h.sendExhibit(w, r, 0, 0)//send /mix/0
		return
	}
	exhibit := pathParts[1]
	exhibitNum, er := getIndex(exhibit)
	if er != nil { //if they fucked up the exhibit, send 404
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
//TODO 
//for now im doing these if else but they might be unecessary. See if i can use javascript
//to tell if a cookie is on browser or not then dynamically change the top right to either
//say "sign up sign in" or "sabio667 sign out" or something. maybe if signout sets the 
//value "" then it can just make sure the cookie isnt blank
//and YEAH javascript can handle cookie shit https://www.w3schools.com/js/js_cookies.asp
//-------------------------------------------------------------------------------------
//(below sounds expensive on request times. if it ever gets a lot of traffic then yeah)
//...unless i impliment that cool side panel of "things that might interest you"
//which i fucking should
//unless unless i request that info with javascript too

//sends 404 page
func (h *handler) send404(w http.ResponseWriter, r *http.Request){
	cookie, err := r.Cookie("session")
	
	if err != nil || !sm.VerifySessionCookie(cookie){
		w.Write(files.fd["exhibitTopSignedOut.html"].data)
	} else {
		w.Write(files.fd["exhibitTopSignedIn.html"].data)
	}
	
	w.Write(files.fd["404.html"].data)
}

//when creating cookies you give them names, so that you can give multiple ones.
//by saying cookie, err := r.Cookie("session"), im asking for the cookie named
//session. Returns an error if its not there
//if theres an error or the session is invalid, send back TopSignedOut, else
//send signedIntTop. then send the rest
func(h *handler)sendExhibit(w http.ResponseWriter, r *http.Request, exh int, page int){
	cookie, err := r.Cookie("session")

	if err != nil || !sm.VerifySessionCookie(cookie) {
		w.Write(files.fd["exhibitTopSignedOut.html"].data)
	} else {
		w.Write(files.fd["exhibitTopSignedOut.html"].data)
	}

	w.Write(pageContent.exhibits[exh][page])
	w.Write(files.fd["exhibitBottom.html"].data)
}

func main() {
	
	//initialize the globals
	//TODO make these return errors and handle the log.Fatals out here
	osa.initialize()
	files.initialize()
	pageContent.initialize()
	am.initialize()
	sm.initialize()
	
	//run the background processes
	//go sm.SessionSweep()
	go pageContent.updateForever()
	
	//start the server
	http.Handle("/", new(handler))
	log.Println("listening....")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

//gets index from AllExhibits global (i think its in content.go) and returns it. used
//in h.handleExhibit
func getIndex(exhibit string) (int, error){

	for i, val := range AllExhibits{
		if exhibit == val {
			return i, nil	
		}	
	}
	
	return 0, errors.New("fuck")
}
