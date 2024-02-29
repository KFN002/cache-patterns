package main

import (
	"math/rand"
	"sync"
	"time"
)

type Session struct {
	ID        string
	UserID    string
	ExpiresAt time.Time
}

type SessionManager struct {
	sessions map[string]Session
	mutex    sync.RWMutex
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]Session),
	}
}

func (sm *SessionManager) StartSession(userID string) string {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	sessionID := generateSessionID()
	expiresAt := time.Now().Add(2 * time.Minute)
	session := Session{
		ID:        sessionID,
		UserID:    userID,
		ExpiresAt: expiresAt,
	}
	sm.sessions[sessionID] = session
	return sessionID
}

func (sm *SessionManager) GetSession(sessionID string) (Session, bool) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	session, ok := sm.sessions[sessionID]
	if !ok {
		return Session{}, false
	}
	if time.Now().After(session.ExpiresAt) {
		delete(sm.sessions, sessionID)
		return Session{}, false
	}
	return session, true
}

func generateSessionID() string {
	return "session_" + randomString(10)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
