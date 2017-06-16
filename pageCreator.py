# this script just concatenates specified html pieces together, or
# css pieces together, shit like that


# hardcoded in all the sub-sources for the files, oh well
# the key is the file that will hold the concatenated pieces, value is an array of the pieces

# "key": [[list of "filename"], [list of "links"]]
# key = ending file name
# list of file names = names of the pieces for ending file
# links = any additional links(mainly for html files to link to js and css). blank list if not anything

pages = {
	"exhibitTop.html": [["htmlfiles/header.html", "htmlfiles/navbar.html", "htmlfiles/sortbar.html"], 
			["js/exhibit.js", "css/exhibit.css"]], # note 1
	"exhibitBottom.html":[["htmlfiles/footer.html"],[]],# nextButtons.html
	"post.html": [["htmlfiles/post.html"],[]],
	"exhibit.css": [["cssfiles/header.css", "cssfiles/navbar.css", "cssfiles/sortbar.css", 
						"cssfiles/footer.css", "cssfiles/post.css"],[]],# nextButtons.css
	"exhibit.js":[["jsfiles/navbar.js", "jsfiles/header.js"],[]]#keep header.js at end
}

# source of pieces, destination of pieced-together file
sourceDir= "htmlcssjs"
destinationDir = "loaderfiles"
cssSite = "//localhost:8080/"
jsSite = cssSite

# open the file that will hold the result. get the file extension. if its an html page, 
# then write the beginning of html pages. for each piece in passed in array,
# add that to the final file. if its html, write the closing html shit
def create(finalName, pieces, links):

	f = open(destinationDir+"/"+finalName, "w")
	
	fileExtension = finalName[finalName.index('.')+1:]
	if fileExtension == "html":#this block does the html linking shit
		f.write("<!DOCTYPE html><html><head><title>ProFools</title>")
		for link in links:
			fileExtension2 = link[link.index('.')+1:]
			if fileExtension2 == "css":
				f.write('<link rel="stylesheet" type="text/css" href="'+cssSite+link+'">')
			elif fileExtension2 == "js":
				f.write('<script src="'+jsSite+link+'"></script>')	
		f.write("</head><body>")

	for fileName in pieces:
		with open(sourceDir+"/"+fileName) as f2:
			f.write(f2.read())

	if fileExtension == "html":
		f.write("</body></html>")
	
	f.close()



# for each element in the dictionary, pass the key and value to create()
for k, v in pages.items():
	create(k, v[0], v[1])


# note 1: js/exhbiti and css/exhibit not because theyre not going to be stored in files call
# css or js, but because thats how the serveHTTP in main.go rationalizes what its looking for
# the browser will call for those (preloaded)files, and that method searches the cached files
# misleading cause urls are patterned like directories, but thats not whats gonna happen here
