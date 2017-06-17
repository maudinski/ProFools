// TODO need to get rid of majin bu symbols from html files created with python. Use clarks shit maybe or idk
//IDEA TODO these are notes and todos
// no sure how but need a favicon
// if specific content type isnt working, it may be the
// the log.Println()'s are obviously a bottle neck, get rid of them (not the error ones)
// in pageCreater, make it scrape out comments and white unnessecarry white space
// might be good to change file permisions for all these files to add an extra layer of security, not sure though
// vegeta load testing utility
// check out the web slot in useful shit.txt
// unless i can combat every case, learn to recover from panic...
// idk if this method of doing to urls is efficient. gotta load test
//TODO IDEA
package main

import (
	"log"
	"net/http"
	"strings"
	"io/ioutil"
	"strconv"
	"errors"
)

/*for now the key will be the file name, with extension*/
var files Files
var pageContent content
var osa outsideStructureAbstractor


/************************************************
meant so if i change a password, or a directory structure, or my table/database
names in sql, or anything that is outside the scope of go code, i wont have to 
sift through a fuckton of code to change it. just gotta change what newOSA() 
initializes shit it

to be used in content.go for db, table, and sql password. init.go for file loading
************************************************/
type outsideStructureAbstractor struct{
	loaderFilesDir string
	picturesDir string	
	postHtmlFN string
	
	dbUser string
	dbPass string // probably not the safest thing to do
	dbHostAndPort string
	db string
	mixPostTable string

}

func (osa *outsideStructureAbstractor) initialize(){
	osa.loaderFilesDir = "loaderfiles"
	osa.picturesDir = "picturesamples"
	osa.postHtmlFN = "post.html"
	osa.dbUser = "root"
	osa.dbPass = "test"
	osa.dbHostAndPort = ""
	osa.db = "test_db2"
	osa.mixPostTable = "post_tb"
}

type customHandler struct {
}
/*******************************************
paths will be like this: 
	/exhibit/mix/1
	/exhibit/literature/0
	/js/exhibit.js
	/css/exhbit.css
*********************************************/
//still not real error checking done here
func (h *customHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	if r.URL.Path == "/"{								
		h.sendExhibit(w, 0, 0)
		return;
	}	
	pathParts := strings.Split(r.URL.Path, "/")[1:]//Split returns a blank at spot [0] if
	switch pathParts[0] {
		case "exhibit":			
			h.handleExhibit(w, pathParts)
		case "port":
		case "post":
		case "css":	
			h.handleJsAndCss(w, pathParts[1], "text/css")
		case "picture":	
			h.handlePic(w, osa.picturesDir+"/"+pathParts[1])
		case "js":
			h.handleJsAndCss(w, pathParts[1], "application/js")
		default:
			log.Println("(servehttp) hit the deafult:", r.URL.Path)
	}
	
}
/************************************************
TODO i think the blank still technically gets all the data copied over, 
so larg-ish amounts of data are getting stored then trashed for this.
maybe chane files.fd to hold *fileData instead of actual file data. Not 
sure how that would work, but if it did, then this would only copy over 
an 8 byte pointer. Better than the entire []byte
***********************************************/

//some sites use seperate servers for js and css
func (h *customHandler) handleJsAndCss(w http.ResponseWriter, f string, ctype string){	
	w.Header().Add("Content-Type", ctype)
	if _, ok := files.fd[f]; ok { //check above for worry on this
		w.Write(files.fd[f].data)
		return;
	} else {
		log.Println("js or css file error:", f)	
	}
}

//reddit has a seperate server just for pics
func (h *customHandler) handlePic(w http.ResponseWriter, path string){
	picData, err := ioutil.ReadFile(path)
	if err == nil {
		w.Write(picData)
	} else {
		log.Println("picData not read right:", err)
	}
}

//all if blocks are error checks
//seems alright
//this error checking works fine but it fucks with the java script
//TODO do a redirect, so as to send back a full url for the navbar javascript
func (h *customHandler) handleExhibit(w http.ResponseWriter, pathParts []string) {
	length := len(pathParts)
	log.Println("length is :", length, "|path is:",pathParts)
	if length == 1{
		h.sendExhibit(w, 0, 0)//send /mix/0
		return
	}
	exhibit := pathParts[1]
	exhibitNum, er := getIndex(exhibit)
	if er != nil { //if they fucked up the exhibit, send 404
		h.send404(w)			//everything else just gets the 'index' (mix/0)
		return	
	}
	if length == 2 {
		h.sendExhibit(w, exhibitNum, 0)//send /whateveritis/0
		return
	}
	pageNum, err := strconv.Atoi(pathParts[2])
	if err != nil {
		h.sendExhibit(w, exhibitNum, 0)	
		return
	}
	w.Write(files.fd["exhibitTop.html"].data)
	w.Write(pageContent.exhibits[exhibitNum][pageNum])
	w.Write(files.fd["exhibitBottom.hmtl"].data)
}	

//no error check needed
func (h *customHandler) send404(w http.ResponseWriter){
	w.Write(files.fd["404.html"].data)
}

//nees to take r *http.Request eventually for the cookies
func (h *customHandler) sendExhibit(w http.ResponseWriter, exhibitNum int, pageNum int){
	w.Write(files.fd["exhibitTop.html"].data)
	w.Write(pageContent.exhibits[exhibitNum][pageNum])
	w.Write(files.fd["exhibitTop.html"].data)
}

func main() {
	
	osa.initialize()
	files.initialize()
	pageContent.initialize()//need to go through and use osa in this
	
	go pageContent.updateForever()//same ^^

	http.Handle("/", new(customHandler))
	log.Println("Starting listening....")
	err := http.ListenAndServe(":8080", nil)

	if err != nil{
		log.Fatal("(main)ListenAndServe error:", err)
	}
}

func getIndex(exhibit string) (int, error){

	for i, val := range AllExhibits{
		if exhibit == val {
			return i, nil	
		}	
	}
	
	return 0, errors.New("fuck")
}
