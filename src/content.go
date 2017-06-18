// TODO a fuckton of hardcoded 10's all around this program. Use slices 
// 'ts about it
//bouta change this bad boy the FUCK up, after writing some go scripts to
// fill a fuckton of data in mysql. Snag some user name csv files from the internet,
// use a random number generator, prolly find some way to make a bunch of titles
// that are only 20 chars long, then loop loop loop my nigga. shits about to get real

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

/******************************
 will hold entries in the post_db, or a row, or postData, or however tf
i want to think about it. Point is its all the data from a post, which is
an entry in the db. idk
*******************************/
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

func (e postEntry) String() string {
	return fmt.Sprintf("userID: %v -- ups: %v -- doubleups: %v -- date: %v -- time: %v -- title: %v -- pagename: %v -- screenname: %v",
						e.userID, e.ups, e.doubleups, e.date, e.time, e.title, e.pagename, e.screenname)
}

/**************************************************************
this type will eventually store top 100 (or 200 or whatever) posts for ALL
subs/stages/whatever the fuck im gonna call it. Right now, only storing for the mix.
The idea is that the main.go will create a content object Called PageContent, that will
be tossed around the program. These functions are method for use by this object. main 
function
will call PageContent.initialize() or something then call go PageContent.runForever()
***************************************************************/
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
}

const(
	EXHIBIT_AMT int = 5//just gonna have to count and do this by hand, makes it easier
	STORED_PAGES int = 3
	POSTS_PER_PAGE int = 15
)

/*************************************************
"root:test@/test_db1"
*************************************************/
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
/*************************************************
*************************************************/
func (c *content) update(spacedOut bool) {

	
	for i, exhibit := range AllExhibits {
		entries, err := c.getTopPosts(exhibit)
		//idk what to do with this error right now, maybe fucking email me?
		if err != nil {
			log.Println("(content.update) GETTOPPOSTS ERROR:", err)
		}
		c.createContent(i, entries)
		if spacedOut{
			time.Sleep(time.Second*5)
		}
	}
}
/*
	entries, err := c.getTop10Posts()
	if err != nil {
		return err
	}
	c.createContent(entries)
	return nil
*/


/*****************************************************************
purpose: inifinite loop of updating the content pages

details: run like this: go c.updateForever(). just calls c.update over and over
		and sleeps for a minute inbetween. eventually needs to call all subs seperately,
		sleeping like 5 seconds inbetween or something. 
*******************************************************************/
func (c *content) updateForever(){
	for {
		c.update(true)
	}
}
//beautiful. I'm loving life. But I hate that I think everything through so much.
//just fucking BE Mauricio, just fucking be
/**********************************************************************
purpose: to get the top 10 posts based on upvotes from the db passed to it,
		and return them

notes: theres a lot of hardcoded 10's in this function, change them all if changing. 
		entries being returned even if error is returned for compiler reasons

details: sets up a query and an array to hold the entries that will be returned, 
		then calls db.Query of that query. loops through the rows object that 
		db.Query returns. Scans the elements in that row into the address of each of
		entry's fields, then stores it in the array. returns the array
***********************************************************************/
// update to slices and/or a global variable for size isntead of all these 10's
func(c *content)getTopPosts(table string)([STORED_PAGES * POSTS_PER_PAGE]postEntry, error){	
	
	var size string = strconv.Itoa(STORED_PAGES*POSTS_PER_PAGE)
	q := "select * from "+table+" order by upvotes desc limit "+size

	var entries [STORED_PAGES*POSTS_PER_PAGE]entry
	
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

/**********************************************************
purpose: to accept an array of entries and return an html string that will be
		assigned to c.mix(eventually this while program needs to create all pages
		not just mix)

notes: only updates mix as of now. DONT just strictly have it update everything. Some how 
		work out some logic that it will update each seperate sub in like 5 sec intervals.
		maybe not in this function, but in another, and pass the mix as parameter, or make 
		this a method for a sub object, idk

details: creates and empty string and starts looping through the array of entries.
		concatenates(spelling?) the html generated for each entry by Sprintf
		each entry is a seperate post in the db

TODO: currently makes page full of every thing in the entry array, so whenever im 
	grabbing more data, i cant just pass in a huge ass array of it, unless i handle 
	that in here. also, needs to eventually deal with port pics and shit (idk how yet)
**********************************************************/
// hard coded 10 here, update to slice or global variable for size
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

/*********************************************************************
TODO but i also need to keep a string version of postHtml for continuous updates
purpose: to load the html data for a post  and return a string version of that data

details: Reads file in from passed in. returns a string typecast
		of the []byte that ioutil.ReadFile returns. The post html should already have the 
		%v's in it (typed into the html file)
*********************************************************************/
func loadPostHtml(file string) (string, error) {
	data, err := ioutil.ReadFile(file)	
	if err != nil {
		return "", err
	}
	return string(data), nil
}
