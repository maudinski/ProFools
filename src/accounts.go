package main

import (
	_"github.com/go-sql-driver/mysql"
	"database/sql"
	"log"
)



type AccountManager struct {
	
	db *sql.DB
	verifyQuery string
}

type userEntry struct {

	handle string 
	screenname string
	password string //is this unsafe?
	email string
	portpicDir string //idk why i think this should be around
	
}


func (am *AccountManager) initialize(){
	var err error
	sqlopenString := osa.dbUser+":"+osa.dbPass+"@"+osa.dbHostAndPort+"/"+osa.users_db
	am.db, err = sql.Open("mysql", sqlopenString)
	if err != nil {
		log.Fatal("(accountmanage.initialize) db didnt open:", err)	
	}
								//incase i forget, db.Prepare uses ? as placeholder
	am.verifyQuery = "select * from "+osa.users_table+" where handle = "

}



//either handle or email will be null, so they can log in with either
//sometimes duplicate code runs more efficiently than non duplicate code. also lazy
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















