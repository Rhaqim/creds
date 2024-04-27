package authentication

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
)

type PasskeyUser interface {
	webauthn.User
	AddCredential(*webauthn.Credential)
	UpdateCredential(*webauthn.Credential)
}

type PasskeyStore interface {
	GetUser(userName string) PasskeyUser
	SaveUser(PasskeyUser)
	GetSession(token string) webauthn.SessionData
	SaveSession(token string, data webauthn.SessionData)
	DeleteSession(token string)
}

var (
	webAuthn *webauthn.WebAuthn
	err      error

	datastore PasskeyStore
	l         = log.New(os.Stdout, "webauthn", log.LstdFlags)
)

func Stuff() {
	// Your initialization function
	host := "localhost"
	origin := "http://localhost:8080"

	wconfig := &webauthn.Config{
		RPDisplayName: "Go Webauthn",    // Display Name for your site
		RPID:          host,             // Generally the FQDN for your site
		RPOrigins:     []string{origin}, // The origin URLs allowed for WebAuthn
	}

	if webAuthn, err = webauthn.New(wconfig); err != nil {
		fmt.Printf("[FATA] %s", err.Error())
		os.Exit(1)
	}

	datastore = NewInMem(l)
}

type InMem struct {
	l    *log.Logger
	data map[string]interface{}
}

func (i *InMem) GetUser(userName string) PasskeyUser {
	return nil
}

func (i *InMem) SaveUser(user PasskeyUser) {
}

func (i *InMem) GetSession(token string) webauthn.SessionData {
	return webauthn.SessionData{}
}

func (i *InMem) SaveSession(token string, data webauthn.SessionData) {
}

func (i *InMem) DeleteSession(token string) {
}

func NewInMem(l *log.Logger) *InMem {
	return &InMem{
		l:    l,
		data: make(map[string]interface{}),
	}
}

// JSONResponse is a helper function to send json responsefunc JSONResponse(w http.ResponseWriter, sessionKey string, data interface{}, status int) {
func JSONResponse(w http.ResponseWriter, sessionKey string, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Session-Key", sessionKey)
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func BeginRegistration(w http.ResponseWriter, r *http.Request) {
	l.Printf("[INFO] begin registration ----------------------\\")

	username, err := getUsername(r)
	if err != nil {
		l.Printf("[ERRO] can't get user name: %s", err.Error())
		panic(err)
	}

	user := datastore.GetUser(username) // Find or create the new user

	options, session, err := webAuthn.BeginRegistration(user)
	if err != nil {
		msg := fmt.Sprintf("can't begin registration: %s", err.Error())
		l.Printf("[ERRO] %s", msg)
		JSONResponse(w, "", msg, http.StatusBadRequest)

		return
	}

	// Make a session key and store the sessionData values
	t := uuid.New().String()
	datastore.SaveSession(t, *session)

	JSONResponse(w, t, options, http.StatusOK) // return the options generated with the session key
	// options.publicKey contain our registration options
}

func FinishRegistration(w http.ResponseWriter, r *http.Request) {
	user := datastore.GetUser("username") // Get the user

	// Get the session key from the header
	t := r.Header.Get("Session-Key")
	// Get the session data stored from the function above
	session := datastore.GetSession(t) // FIXME: cover invalid session

	// In out example username == userID, but in real world it should be different   user := datastore.GetUser(string(session.UserID)) // Get the user

	credential, err := webAuthn.FinishRegistration(user, session, r)
	if err != nil {
		msg := fmt.Sprintf("can't finish registration: %s", err.Error())
		l.Printf("[ERRO] %s", msg)
		JSONResponse(w, "", msg, http.StatusBadRequest)

		return
	}

	// If creation was successful, store the credential object
	user.AddCredential(credential)
	datastore.SaveUser(user)
	// Delete the session data
	datastore.DeleteSession(t)

	l.Printf("[INFO] finish registration ----------------------/")
	JSONResponse(w, "", "Registration Success", http.StatusOK) // Handle next steps
}

func BeginLogin(w http.ResponseWriter, r *http.Request) {
	l.Printf("[INFO] begin login ----------------------\\")

	username, err := getUsername(r)
	if err != nil {
		l.Printf("[ERRO]can't get user name: %s", err.Error())
		panic(err)
	}

	user := datastore.GetUser(username) // Find the user

	options, session, err := webAuthn.BeginLogin(user)
	if err != nil {
		msg := fmt.Sprintf("can't begin login: %s", err.Error())
		l.Printf("[ERRO] %s", msg)
		JSONResponse(w, "", msg, http.StatusBadRequest)

		return
	}

	// Make a session key and store the sessionData values
	t := uuid.New().String()
	datastore.SaveSession(t, *session)

	JSONResponse(w, t, options, http.StatusOK) // return the options generated with the session key
	// options.publicKey contain our registration options
}

func FinishLogin(w http.ResponseWriter, r *http.Request) {
	user := datastore.GetUser("username") // Get the user

	// Get the session key from the header
	t := r.Header.Get("Session-Key")
	// Get the session data stored from the function above
	session := datastore.GetSession(t) // FIXME: cover invalid session

	// In out example username == userID, but in real world it should be different   user := datastore.GetUser(string(session.UserID)) // Get the user

	credential, err := webAuthn.FinishLogin(user, session, r)
	if err != nil {
		l.Printf("[ERRO] can't finish login %s", err.Error())
		panic(err)
	}

	// Handle credential.Authenticator.CloneWarning
	if credential.Authenticator.CloneWarning {
		l.Printf("[WARN] can't finish login: %s", "CloneWarning")
	}

	// If login was successful, update the credential object
	user.UpdateCredential(credential)
	datastore.SaveUser(user)
	// Delete the session data
	datastore.DeleteSession(t)

	l.Printf("[INFO] finish login ----------------------/")
	JSONResponse(w, "", "Login Success", http.StatusOK)
}

// getUsername is a helper function to extract the username from json requestfunc getUsername(r *http.Request) (string, error) {
func getUsername(r *http.Request) (string, error) {
	type Username struct {
		Username string `json:"username"`
	}
	var u Username
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		return "", err
	}

	return u.Username, nil
}

// var (
// 	webAuthn *webauthn.WebAuthn
// 	err      error
// )

// // Your initialization function
// func main() {
// 	wconfig := &webauthn.Config{
// 		RPDisplayName: "Go Webauthn",                               // Display Name for your site
// 		RPID:          "go-webauthn.local",                         // Generally the FQDN for your site
// 		RPOrigins:     []string{"https://login.go-webauthn.local"}, // The origin URLs allowed for WebAuthn requests
// 	}

// 	if webAuthn, err = webauthn.New(wconfig); err != nil {
// 		fmt.Println(err)
// 	}
// }

// // Register
// func BeginRegistration(w http.ResponseWriter, r *http.Request) {
// 	user := datastore.GetUser() // Find or create the new user
// 	options, session, err := webAuthn.BeginRegistration(user)
// 	// handle errors if present
// 	// store the sessionData values
// 	JSONResponse(w, options, http.StatusOK) // return the options generated
// 	// options.publicKey contain our registration options
// }

// func FinishRegistration(w http.ResponseWriter, r *http.Request) {
// 	user := datastore.GetUser() // Get the user

// 	// Get the session data stored from the function above
// 	session := datastore.GetSession()

// 	credential, err := webAuthn.FinishRegistration(user, session, r)
// 	if err != nil {
// 		// Handle Error and return.

// 		return
// 	}

// 	// If creation was successful, store the credential object
// 	// Pseudocode to add the user credential.
// 	user.AddCredential(credential)
// 	datastore.SaveUser(user)

// 	JSONResponse(w, "Registration Success", http.StatusOK) // Handle next steps
// }

// // Login
// func BeginLogin(w http.ResponseWriter, r *http.Request) {
// 	user := datastore.GetUser() // Find the user

// 	options, session, err := webAuthn.BeginLogin(user)
// 	if err != nil {
// 		// Handle Error and return.

// 		return
// 	}

// 	// store the session values
// 	datastore.SaveSession(session)

// 	JSONResponse(w, options, http.StatusOK) // return the options generated
// 	// options.publicKey contain our registration options
// }

// func FinishLogin(w http.ResponseWriter, r *http.Request) {
// 	user := datastore.GetUser() // Get the user

// 	// Get the session data stored from the function above
// 	session := datastore.GetSession()

// 	credential, err := webAuthn.FinishLogin(user, session, r)
// 	if err != nil {
// 		// Handle Error and return.

// 		return
// 	}

// 	// Handle credential.Authenticator.CloneWarning

// 	// If login was successful, update the credential object
// 	// Pseudocode to update the user credential.
// 	user.UpdateCredential(credential)
// 	datastore.SaveUser(user)

// 	JSONResponse(w, "Login Success", http.StatusOK)
// }
