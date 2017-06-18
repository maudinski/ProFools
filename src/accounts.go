package main

import (
	_"github.com/go-sql-driver/mysql"
	"databases/sql"
)



type AccountManager struct {
	
	db *sql.DB
	getByEmail string
	getByHandle string
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
	am.db, err := sql.Open("mysql", sqlopenString)
	if err != nil {
		log.Fatal("(accountmanage.initialize) db didnt open:", err)	
	}
								//incase i forget, db.Prepare uses ? as placeholder
	am.getByEmail = "select * from "+osa.users_table+" where email=?"
	am.getByHandle = "select * from "+osa.users_table+" where handle=?"

}



//either handle or email will be null, so they can log in with either
//sometimes duplicate code runs more efficiently than non duplicate code. also lazy
func (am *AccountManager) VerifyLogin(handle string, email string, pass string) bool {

	if handle == ""{
		stmt, err := db.Prepare(am.getByEmail)
		if err != nil {
			log.Println("(accountmanager.VERIFYLOGIN DB.PREPARE ERR:", err) // for debug
			return false
		}
		var retrievedPass string
		row, err2 := stmt.Query(email).Scan(_, retrievedPass, _, _, _)
		if err2 != nil {
			log.Println("(accountmanager.VERIFYLOGIN STMT.QUERY ERR:", err) // for debug
			return false
		}
		return retrievedPass == pass

	} else {
		stmt, err := db.Prepare(am.getByHandle)
		if err != nil {
			log.Println("(accountmanager.VERIFYLOGIN DB.PREPARE ERR:", err) // for debug
			return false
		}
		var retrievedPass string
		row, err2 := stmt.Query(pass).Scan(_, retrievedPass, _, _, _)
		if err2 != nil {
			log.Println("(accountmanager.VERIFYLOGIN STMT.QUERY ERR:", err) // for debug
			return false
		}
		return retrievedPass == pass
	}

}















