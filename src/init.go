package main

import (
	"io/ioutil"
	"log"
	"strings"
	"fmt"
)

type Files struct {
	fd map[string] FileData	
}

//this holds the data of a read in file and its details
type FileData struct {
	contentType string
	name string
	data []byte	
}
//So that FileData impliments the Stringer interface (naked pass to fmt.Println)
func (fd FileData) String() string{
	return fmt.Sprintf("File name: %v -- Content Type: %v", fd.name, fd.contentType)
}

/**********************************
stores the data and contentType of files that will be sent (html, css)
so no constant directory searching and file reading
***********************************
/*looks right*/
func (filesObject *Files) initialize() {
	filesObject.fd = make(map[string] FileData)
	filesObject.initializeFiles()	
}

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

/*********************************************************************
purpose: to return to content type of the entered file.

notes: might just make it a method for a FileData object
		should still work
details: enter the file, it grabs the file extension and returns an appropriate
	content type to write in the header for the browser
**********************************************************************/
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


