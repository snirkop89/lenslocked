package models

import "database/sql"

type Session struct {
	ID     int
	UserID int
	// Token is only set when creating a new session. when looking up a session
	// this iwll be left empty, as we only store the hash of a session token
	// in our database and we cannot reverse it into a raw token.
	Token     string
	TokenHash string
}

type SessionService struct {
	DB *sql.DB
}

func (ss *SessionService) Create(userID int) (*Session, error) {
	// TODO: Create new session token
	// TODO: Implement sessionService.Create
	return nil, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	// TODO: Implement SessionService.User
	return nil, nil
}
