package sessionMaker

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/gotd/td/session"
	"github.com/gotd/td/session/tdesktop"
	"github.com/yaroslavvasilenko/gotgproto/functions"
	"github.com/yaroslavvasilenko/gotgproto/storage"
)

// SessionName object consists of name and SessionType.
type SessionName struct {
	Name        string
	FileName    string
	Path        string
	SessionType SessionType
	Data        []byte
	Err         error
}

// SessionType is the type of session you want to log in through.
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
		Name:        sessionName,
		FileName:    sessionFileName,
		Path:        sessionPath,
		SessionType: sessionType,
	}
	s.Data, s.Err = s.OptimizeSessionData()
	return &s
}

// GetName is used for retrieving the name of the session.
func (s *SessionName) GetName() string {
	return s.Name
}

// GetData is used for retrieving session data.
func (s *SessionName) GetData() ([]byte, error) {
	return s.Data, s.Err
}

// OptimizeSessionData optimizes session data based on the session type.
func (s *SessionName) OptimizeSessionData() ([]byte, error) {
	var err error

	switch s.SessionType {
	case PyrogramSession, TelethonSession, TDataSession:
		err = storage.Load(filepath.Join(s.Path, s.FileName+".session"), false)
		if err != nil {
			return nil, err
		}

		var sd *session.Data

		switch s.SessionType {
		case PyrogramSession:
			sd, err = DecodePyrogramSession(s.Name)
		case TelethonSession:
			sd, err = session.TelethonSession(s.Name)
		case TDataSession:
			accounts, err := tdesktop.Read(s.Name, nil)
			if err == nil && len(accounts) > 0 {
				auth := accounts[0]
				sd, err = session.TDesktopSession(auth)
			}
		}

		if err != nil {
			return nil, err
		}

		data, err := json.Marshal(jsonData{
			Version: storage.LatestVersion,
			Data:    *sd,
		})

		return data, err

	case StringSession:
		err = storage.Load(filepath.Join(s.Path, s.FileName+".session"), false)
		if err != nil {
			return nil, err
		}
		sd, err := functions.DecodeStringToSession(s.Name)
		if err != nil {
			return nil, err
		}
		return sd.Data, err

	default:
		if s.Name == "" {
			s.Name = "new_session"
		}
		err = storage.Load(s.Path+s.FileName+".session", false)
		if err != nil {
			return nil, err
		}
		sFD := storage.GetSession()
		return sFD.Data, nil
	}
}
