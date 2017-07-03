package main

import(
	"fmt"
	_"github.com/go-sql-driver/mysql"
	"database/sql"
	"time"
	"log"
	"io/ioutil"
	"strconv"
)

//will hold entries in the post_db, or a row, or postData, or however tf
//i want to think about it. Point is its all the data from a post, which is
//an entry in the db. idk
type postEntry struct {
	userID int
	ups int
	doubleups int
	date string
	time string
	title string
	pagename string
	screenname string
}

//c gloabl will be created in main.go. Tossed around the main.go. holds the content. the 
//db is only used in here. 
//postHtml string is a string that holds the html of a single post. %v's are already typed
//in there, so its easy to use fmt.Sprintf to add values to it. 
//exhibits is a slice of pages, and pages is a slice of []byte that stores the content.
//useage: content.exhibits[num1][num2]
//num1 corresponds to global slice AllExhibits, and num2 will be the page, so:
//content.exhibits[1][0] is the literature exhibit, and the first page(indexing starts 0)
type content struct {
	db *sql.DB
	postHtml string
	//some mutex
	exhibits [EXHIBIT_AMT] pages
}

type pages [STORED_PAGES][]byte

//THESE HAVE TO BE THE SAME AS THE TABLE NAME IN MYSQL
//since its passed to getTopPosts as paramater. probably will use it some handle functions
//from main
var AllExhibits []string = []string{
	"mix",
	"literature",
	"tatoos",
	"skill_toys",
	"wall_art",//if you add/remove, you gotta change the constant EXHIBIT_AMT
	//this is used as is in getIndex() from main.go
}

const(
	EXHIBIT_AMT int = 5//just gonna have to count and do this by hand, makes it easier
	STORED_PAGES int = 3
	POSTS_PER_PAGE int = 15
)

//"root:test@/test_db1"
//initializes the database and the postHtml string, then calls update on the content
func (c *content) initialize(){

	var err error
	sqlopenString := osa.dbUser+":"+osa.dbPass+"@"+osa.dbHostAndPort+"/"+osa.db
	c.db, err = sql.Open("mysql", sqlopenString)	
	if err != nil {
		log.Fatal("(handleContent.initialize)shit isnt gonna work:", err)	
	}
	
	c.postHtml, err = loadPostHtml(osa.loaderFilesDir+"/"+osa.postHtmlFN)
	if err != nil {
		log.Fatal("(handleContent.intialize)loadPostHtml returned an error:", err)	
	}
	
	c.update(false)

}
//TODO need to somehow lock shit with mutex's when changing data, so that shit doesnt go
//wrong.
//loops through the AllExhibits slice and updates the contents in c. uses other functions
//to do that. spaceOut is pass as false only when c.initialize() calls this, so that it
//initializes everything at once. Other wise, c.updataForever is called as a go routine
//in main.go, which passes true to this function, so that it waits five seconds inbetween
//updating each exhibits. The idea is that waiting 5 seconds will reduce some strain on
//the server
func (c *content) update(spacedOut bool) {
	for i, exhibit := range AllExhibits {
		entries, err := c.getTopPosts(exhibit)
		if err != nil { //idk what to do with this error right now, maybe fucking email me?
			log.Println("(content.update) GETTOPPOSTS ERROR:", err, "-------------------")
		}
		c.createContent(i, entries)
		if spacedOut{
			time.Sleep(time.Second*5)
		}
	}
}

//called as 'go c.updateForever()' in main.go. updates the contents in an infinite loop,
//and passes true for spacedOut, which waits 5 seconds in between each exhibit update
func (c *content) updateForever(){
	for {
		c.update(true)
	}
}

//gets the top posts from the table passed. uses global variables STORED_PAGES and 
//POSTS_PER_PAGE to determine how many. reads in the rows from mysql database and returns
//an array of entries of the rows.
func(c *content)getTopPosts(table string)([STORED_PAGES * POSTS_PER_PAGE]postEntry, error){	
	
	var size string = strconv.Itoa(STORED_PAGES*POSTS_PER_PAGE)
	q := "select * from "+table+" order by upvotes desc limit "+size

	var entries [STORED_PAGES*POSTS_PER_PAGE]postEntry
	
	rows, err := (c.db).Query(q)
	if err != nil { 
		return entries, err 
	}
	defer rows.Close()

	e := new(postEntry)
	for i := 0; rows.Next(); i++{
		err = rows.Scan(&e.userID, &e.ups, &e.doubleups, &e.date, &e.time, 
						&e.title, &e.pagename, &e.screenname)
		if err != nil{ 
			return entries, err	
		}
		entries[i] = *e
	}
	if err = rows.Err(); err != nil{ 
		return entries, err 
	}
	return entries, nil

}

//gets passed index of the exhibit that is gonna get updated, and the array of entries
//does some modulo shit in the for loop that could be done better but isnt. 
//also i havent fully checked this function, it seems to work alright but im not sure
//if the site itself does the page suring properly. ie: exhibit/mix/0 works but idk about
//exhibit/mix/1, etc
//uses the post html string in content
func(c *content)createContent (exhInd int, entries [STORED_PAGES*POSTS_PER_PAGE]postEntry){
	page := ""
	pageNum := 0
	firstPass := true
	for i, e := range(entries){
		if i % POSTS_PER_PAGE == 0 && !firstPass{
			c.exhibits[exhInd][pageNum] = []byte(page)
			pageNum++
			page = ""
		}
		firstPass = false
		page = page + fmt.Sprintf(c.postHtml, e.ups, e.doubleups, e.screenname, 
								  e.pagename, e.title)
	}
}

//reads the file passed, which is the html of the post. assumes that the %v's are already
//typed in the correct place, because breatContent use fmt.Sprintf to add the values to it
func loadPostHtml(file string) (string, error) {
	data, err := ioutil.ReadFile(file)	
	if err != nil {
		return "", err
	}
	return string(data), nil
}

