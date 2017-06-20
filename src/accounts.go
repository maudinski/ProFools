package main

import (
	_"github.com/go-sql-driver/mysql"
	"database/sql"
	"log"
)


//a global will be used in main.go
//used to manage account related things, like VErifying sign ins and sign ups
type AccountManager struct {
	//the pointer to the database
	db *sql.DB
	//will hold the query for looking things up by handle
	//initialized to be a partial query string in initialize(). just append the handle
	verifyQuery string
}

//struct to hold what the queries return
type userEntry struct {
	handle string 
	screenname string
	password string //is this unsafe?
	email string
	portpicDir string //idk why i think this should be around
}

//initializes the variables in the struct. uses osa a lot to abstract some things
func (am *AccountManager) initialize(){
	var err error
	sqlopenString := osa.dbUser+":"+osa.dbPass+"@"+osa.dbHostAndPort+"/"+osa.users_db
	am.db, err = sql.Open("mysql", sqlopenString)
	if err != nil {
		log.Fatal("(accountmanage.initialize) db didnt open:", err)	
	}
	am.verifyQuery = "select * from "+osa.users_table+" where handle = "

}

//you can only log in if you know your handle name, so that must be passed
//looks up in the data base that handle and verifies that the passed in passwor
//is the same as the stored password. should really encrypt this shit
func (am *AccountManager) VerifySignin(handle string, pass string) bool {

	var e userEntry
	fullQuery := am.verifyQuery + "'" + handle + "'" +";"
	err := am.db.QueryRow(fullQuery).Scan(&e.handle, &e.password, &e.screenname, &e.email, &e.portpicDir)
	if err == sql.ErrNoRows{
		log.Println("am.VerifySignin", err)
		return false	
	}
	if err != nil {
		log.Println("am.VerifySignin", err)
		return false	
	}
	return e.password == pass

}














