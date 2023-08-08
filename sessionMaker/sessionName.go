package sessionMaker

import (
	"encoding/json"
	"fmt"

	"github.com/gotd/td/session"
	"github.com/gotd/td/session/tdesktop"
	"github.com/jaskaur18/gotgproto/functions"
	"github.com/jaskaur18/gotgproto/storage"
)

// SessionName object consists of name and SessionType.
type SessionName struct {
	name        string
	fileName    string
	path        string
	sessionType SessionType
	data        []byte
	err         error
}

// SessionType is the type of session you want to log in through.
// It consists of three types: Session, StringSession, TelethonSession.
type SessionType int

const (
	// Session should be used for authorizing into session with default settings.
	Session SessionType = iota
	// StringSession is used as SessionType when you want to log in through the string session made by gotgproto.
	StringSession
	// TelethonSession is used as SessionType when you want to log in through the string session made by telethon - a Python MTProto library.
	TelethonSession
	// PyrogramSession is used as SessionType when you want to log in through the string session made by pyrogram - a Python MTProto library.
	PyrogramSession
	// TDataSession is used as SessionType when you want to log in through the string session made by Telegram Client - a C++ MTProto library.
	TDataSession
)

//git

// NewSessionOpts is the options for creating a new session.
type NewSessionOpts struct {
	SessionName string
	SessionPath string
}

// NewSession creates a new session with provided name string and SessionType.
func NewSession(sessionName string, sessionType SessionType, newSessionOpts ...NewSessionOpts) *SessionName {
	var sessionFileName string
	var sessionPath string
	if len(newSessionOpts) > 0 && newSessionOpts[0].SessionName != "" {
		if newSessionOpts[0].SessionPath == "" {
			sessionPath = "./sessions"
		} else {
			sessionPath = newSessionOpts[0].SessionPath
		}
		sessionFileName = newSessionOpts[0].SessionName
	} else {
		sessionFileName = fmt.Sprintf("%s_%s", sessionName, "telegram")
		sessionPath = "./sessions"
	}
	s := SessionName{
		name:        sessionName,
		fileName:    sessionFileName,
		path:        sessionPath,
		sessionType: sessionType,
	}
	s.data, s.err = s.load()
	return &s
}

func (s *SessionName) load() ([]byte, error) {

	fileName := fmt.Sprintf("%s/%s.session", s.path, s.fileName)

	switch s.sessionType {
	case PyrogramSession:
		storage.Load(fileName, false)
		sd, err := DecodePyrogramSession(s.name)
		if err != nil {
			return nil, err
		}
		data, err := json.Marshal(jsonData{
			Version: storage.LatestVersion,
			Data:    *sd,
		})
		return data, err
	case TelethonSession:
		storage.Load(fileName, false)
		sd, err := session.TelethonSession(s.name)
		if err != nil {
			return nil, err
		}
		data, err := json.Marshal(jsonData{
			Version: storage.LatestVersion,
			Data:    *sd,
		})
		return data, err
	case TDataSession:
		storage.Load(fileName, false)
		accounts, err := tdesktop.Read(s.name, nil)
		if err != nil {
			return nil, err
		}
		if len(accounts) == 0 {
			return nil, fmt.Errorf("no accounts found")
		}
		auth := accounts[0]
		sd, err := session.TDesktopSession(auth)
		if err != nil {
			return nil, err
		}

		data, err := json.Marshal(jsonData{
			Version: storage.LatestVersion,
			Data:    *sd,
		})

		return data, err

	case StringSession:
		storage.Load(fileName, false)
		sd, err := functions.DecodeStringToSession(s.name)
		if err != nil {
			return nil, err
		}

		// data, err := json.Marshal(jsonData{
		// 	Version: latestVersion,
		// 	Data:    *sd,
		// })
		return sd.Data, err
	default:
		if s.name == "" {
			s.name = "new"
		}
		storage.Load(fileName, false)
		sFD := storage.GetSession()
		return sFD.Data, nil
	}
}

// GetName is used for retrieving name of the session.
func (s *SessionName) GetName() string {
	return s.name
}

// GetData is used for retrieving session data through provided SessionName type.
func (s *SessionName) GetData() ([]byte, error) {
	return s.data, s.err
}
