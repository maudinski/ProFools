//handles storing all the sent html pieces, css files, and js files, before the server
//starts running. Doing this to make the process faster. Now theres not constant
//(actually now theres not ANY) file reading during server running.
package main

import (
	"io/ioutil"
	"log"
	"strings"
	"fmt"
)
//keeping the bottom 2 structs as structs just incase they eventually need more relevant
//fields, so that all the code doesnt need to be refactored

//declaring a global of this in main.go. holds all the html pieces, css files, & js files
//might just not make this a struct to make it a little more clear, since its only one 
//field in it
//general useage: w.Write(files.fd["exhibitBottom.html"].data)
//the files are stored as they are named, for easy reading. because the object stored 
//is a FileData struct, and im sending the data field, which is an array of bytes.
//w.Write (which writes back to the requesting browser) takes []byte as its parameter
type Files struct {
	fd map[string] FileData	
}

//this holds the data of a read in file and its details
type FileData struct {
	contentType string//not really used
	name string//not really used
	data []byte	
}

//So that FileData impliments the Stringer interface (naked p___ to fmt.Println)
func (fd FileData) String() string{
	return fmt.Sprintf("File name: %v -- Content Type: %v", fd.name, fd.contentType)
}

//makes the map of FileData, then initiates the files
func (filesObject *Files) initialize() {
	filesObject.fd = make(map[string] FileData)
	filesObject.initializeFiles()	
}

//goes through the directory of loaderFiles(files to be loaded/cached) and saves
//theyre []byte data in  filesObject.fd map, with they key being the file name.
//uses iouitl.ReadDir and ioutil.ReadFile, which makes it pretty easy. For now,
//___umes that there wont be any directories in loader files, just all the files
//in there. Doesnt do anything if there is a directory, but shouldnt be hard to
//change that if needed
func (filesObject *Files) initializeFiles() {
	files, err := ioutil.ReadDir(osa.loaderFilesDir) //returns []os.FileInfo
	if err != nil {log.Fatal("(init.initiateFIles) Directory not read right:", err)}
	for _, f := range files{
		name := f.Name()
		if f.IsDir() {
			log.Println("NOT DOING ANYTHING WITH THIS DIRECTORY:", name)
			continue	
		}
		data, err2 := ioutil.ReadFile(osa.loaderFilesDir+"/"+name)
		if err2 != nil {log.Fatal("(init.initiateFiles) File not read right:", err2)}
		filesObject.fd[name] = FileData{getContentType(name), name, data}
	}
}

//gets the content type based on the file extension. used in itializeFiles. 
//content type not really used but its there just in case
func getContentType(fileName string) string{
	contentTypes := map[string]string {
		"html": "text/html",
		"css": "text/css",
		"js": "application/javascript",
	}
	index := strings.Index(fileName, ".")
	fileExt := fileName[index+1:]
	ct, ok := contentTypes[fileExt]
	if ok {
		return ct
	}
	log.Println("(getContentType)Content type set to plain text for file", fileName)
	return "text/plain"
}


