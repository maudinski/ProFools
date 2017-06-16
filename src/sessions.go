//sessions manager
package main

//all speculative


type SessionManager struct {		
	activeSessions []session
}

type session struct {
	id int
}


//external functions
	//possibly internal
func (sm *SessionManager)VerifySesion(c *http.Cookie)(bool, error){
	_, ok := activeSessions[c.]
}
	//when signing in or signing up
func (sm *SessionManager) CreateSession(			) (cookie, error){
		
}
	//when hitting sign out OR when sm.CleanSessions() or whatever
func (sm *SessionManager) DestroySession(c *http.Cookie)(error){
		
}

/****************/


//internals
func (sm *SessionManager) newID(){
		
}

func (sm *SessionManager) cleanSessions(){
}






