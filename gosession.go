package gosession

import (
	"github.com/pkg/errors"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var sessionIDHashLength int = 64

// Name of the cookie checked by gosession, and sent to the browser
var cookieName string = "GosessionCookie"

type Session struct {
	// Fields for internal use
	sessionID string
	ge GosessionEngine
	ges *EngineState
}

type EngineState struct {
	initialised bool
	closed bool
}

type GosessionEngine interface {

	// Tells the engine to prepare everything it needs to manage sessions and
	// read/write their keys, and return an error if something goes wrong
	Init() error

	// Returns a pointer to EngineState so gosession can view and modify it on
	// the engine's behalf
	GetEngineStatePointer() *EngineState

	// If the engine is using resources on background, this function should save
	// and close everything, it'll be called once the engine is no longer of use
	Close()

	// Returns true if this session does exist in the database, mainly for
	// checking SessionID duplicates before creating new one
	SessionExists(sessionID string) bool

	// Creates new session with all keys empty
	CreateSession(sessionID string)

	// Destroys the session and keys
	DestroySession(sessionID string)

	// Destroys ALL sessions and keys
	DestroyAllSessions()

	// Reads a key's value. If the key doesn't exist, or it doesn't have a value,
	// the returned value is nil
	ReadKey(sessionID string, key string) interface{}

	// Writes a value to a key. If the key doesnt exist, it's created
	// automatically and the value is written to it. If it does exist, the
	// original value is replaced by the new value.
	WriteKey(sessionID string, key string, value interface{})

	// Deletes a key from the database
	DeleteKey(sessionID string, key string)

}

/*
func LoadSessionFromID(sessionID string, ge GosessionEngine) *Session {
	if ge.SessionExists(sessionID) {
		return &Session{
			sessionID: sessionID,
			ge:        ge,
		}
	}
	ge.CreateSession(sessionID)
	return &Session{
		sessionID: sessionID,
		ge:        ge,
	}
}expires, _ := time.Parse("2.1.expires, _ := time.Parse("2.1.2006", "1.1.2037")2006", "1.1.2037")
  */


func addSessionCookie(w http.ResponseWriter, r *http.Request, value string){
	expires := time.Now().AddDate(10,0,0)
	cookie := http.Cookie{
		Name:       cookieName,
		Value:      value,
		Path:       "/",
		Expires:    expires,
	}
	http.SetCookie(w, &cookie)
}

// Tells an engine to initialise itself, and returns an error if any problems
// are encountered
func Init(ge GosessionEngine) error {
	var es *EngineState = ge.GetEngineStatePointer()
	if es.closed {
		return errors.New("GosessionEngine has been closed by CloseEngine().")
	}
	if es.initialised {
		return nil
	}
	es.initialised = true
	return ge.Init()
}

func LoadSessionFromRequest(w http.ResponseWriter, r *http.Request, ge GosessionEngine) (*Session, error) {
	browserHasTheCookie := false
	cookieValue := ""
	for _, cookie := range r.Cookies() {
		if cookie.Name == cookieName && len(cookie.Value) == sessionIDHashLength {
			browserHasTheCookie = true
			cookieValue = cookie.Value
		}
	}
	if browserHasTheCookie && len(cookieValue) > 0 {
		return LoadSession(cookieValue, ge)
	}

	if err := Init(ge); err != nil {
		return nil, err
	}
	session, sessionID := NewSession(ge)
	addSessionCookie(w, r, sessionID)
	return session, nil
}

// Returns name of the created session
func NewSession(ge GosessionEngine) (*Session, string) {
	sessionID := generateHashForSession(ge)
	ge.CreateSession(sessionID)
	session := &Session{
		sessionID: sessionID,
		ge:        ge,
	}
	return session, sessionID
}

// Tries to load *Session specified by parameter, and if it fails
// (session doesn't exist) it makes automatically a new one
func LoadSession(sessionID string, ge GosessionEngine) (*Session, error) {
	if err := Init(ge); err != nil {
		return nil, err
	}

	if ge.SessionExists(sessionID) {
		return &Session{
			sessionID: sessionID,
			ge:        ge,
		}, nil
	} else {
		session, _ := NewSession(ge)
		return session, nil
	}
}

func DestroyAllSessions(ge GosessionEngine) {
		ge.DestroyAllSessions()
}

// Tells the engine to save everything and exit, also frees up some resources
func CloseEngine(ge GosessionEngine) {
	ge.Close()
	ge.GetEngineStatePointer().closed = true
}

// Destroys session
func (s *Session) DestroySession() {
	s.ge.DestroySession(s.sessionID)
}

// Destroys session and makes a new one using the same session ID, in other
// words, clears all the keys and values
func (s *Session) ClearSession() {
	s.DestroySession()
	s.ge.CreateSession(s.sessionID)
}

// Retrieves value by key
func (s *Session) ReadKey(key string) interface{} {
	return s.ge.ReadKey(s.sessionID, key)
}

// Retrieves value by key and makes sure the value is not nil by returning value
// of a specified type rather than nil
func (s *Session) ReadKeyAs(key string, valueReturnedIfNil interface{}) interface{} {
	raw := s.ge.ReadKey(s.sessionID, key)
	if raw == nil {
		return valueReturnedIfNil
	}
	return raw
}

func (s *Session) ReadKeyAsString(key string) string {
	return s.ReadKeyAs(key, "").(string)
}

// Deletes key and its value
func (s *Session) DeleteKey(key string) {
	s.ge.DeleteKey(s.sessionID, key)
}

// Writes key:value
func (s *Session) WriteKey(key string, value interface{}) {
	s.ge.WriteKey(s.sessionID, key, value)
}

// Generates a hash for SessionID, and makes sure it's unique to
// prevent disasters
func generateHashForSession(ge GosessionEngine) string {
	var hash string
	for {
		hash = ""

		for i := 0; i < sessionIDHashLength; i++ {
			r := rand.Intn(59)
			if r < 25 {
				hash+=string(rune(r+97))
			} else if r < 50 {
				hash+=string(rune(r+65-25))
			} else {
				hash+=strconv.Itoa(r-50)
			}
		}

		if ge.SessionExists(hash) {
			continue
		} else {
			break
		}
	}
	return hash
}