// TODO need to get rid of majin bu symbols from html files created with python. Use clarks shit maybe or idk
//IDEA TODO these are notes and todos
// no sure how but need a favicon
// if specific content type isnt working, it may be the
// the log.Println()'s are obviously a bottle neck, get rid of them (not the error ones)
// put all css files into one file or use the templating for go. point is, less dependent files makes less requests from browser
// in pageCreater, make it scrape out comments and white unnessecarry white space
// might be good to change file permisions for all these files to add an extra layer of security, not sure though
// vegeta load testing utility
// check out the web slot in useful shit.txt
//TODO IDEA
package main

import (
	"log"
	"net/http"
	"strings"
	"io/ioutil"
	"strconv"
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
	osa.picturesDir = "pictureSamples"
	osa.postHtmlFN = "post.html"
	osa.dbUser = "root"
	osa.dbPass = "test"
	osa.dbHostAndPort = ""
	osa.db = "test_db2"
	osa.mixPostTable = "post_tb"
}

type customHandler struct {
}
// paths will be like this: 
/*
	/exhibit/mix/1
	/exhibit/literature/0
	/js/exhibit.js
	/css/exhbit.css
*/
//still not real error checking done here
func (h *customHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	log.Println(r.URL.Path)
	if r.URL.Path == "/"{								
		w.Write(files.fd["exhibitTop.html"].data)	
		w.Write(pageContent.exhibits[0][0])		
		w.Write(files.fd["exhibitBottom.hmtl"].data)
		return;
	}	
	pathParts := strings.Split(r.URL.Path, "/")[1:]//Split returns a blank at spot [0] if
	switch pathParts[0] {
		case "exhibit":			
			h.handleExhibit(w, pathParts)
		case "port":
		case "post":
		case "css":	
			h.handleCss(w, pathParts[1])
		case "picture":	
			h.handlePic(w, osa.picturesDir+"/"+pathParts[1])
		case "js":
			h.handleJs(w, pathParts[1])
		default:
			log.Println("(servehttp) hit the deafult:", r.URL.Path)
	}
	
}

func (h *customHandler) handleJs(w http.ResponseWriter, jsfile string){	
	w.Header().Add("Content-Type", "application/js")
	w.Write(files.fd[jsfile].data)
}

func (h *customHandler) handlePic(w http.ResponseWriter, path string){
	picData, err := ioutil.ReadFile(path)
	if err == nil {
		w.Write(picData)
	} else {
		log.Println("picData not read right:", err)
	}
}

func (h *customHandler) handleExhibit(w http.ResponseWriter, pathParts []string) {
	exhibit := pathParts[1]
	pageNum, _ := strconv.Atoi(pathParts[2])
	exhibitNum := getIndex(exhibit)
	w.Write(files.fd["exhibitTop.html"].data)
	w.Write(pageContent.exhibits[exhibitNum][pageNum])
	w.Write(files.fd["exhibitBottom.hmtl"].data)
}	

func (h *customHandler) handleCss(w http.ResponseWriter, cssFile string){
	w.Header().Add("Content-Type", "text/css")
	w.Write(files.fd[cssFile].data)
}

func (h *customHandler) write404(w http.ResponseWriter){
	w.Write(files.fd["404.html"].data)
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
//ommenting out so that compiler reminds me i have to write this1
func getIndex(exhibit string) int {

	for i, val := range AllExhibits{
		if exhibit == val {
			return i	
		}	
	}
	
	return 0
}
