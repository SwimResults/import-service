package service

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// ProgressEvent represents a progress update event
type ProgressEvent struct {
	Type     string  `json:"type"`     // "progress"
	Progress float64 `json:"progress"` // 0-100
	Message  string  `json:"message"`  // optional message
}

// LogEvent represents a log message event
type LogEvent struct {
	Type      string `json:"type"`      // "log"
	Message   string `json:"message"`   // log message
	Level     string `json:"level"`     // "info", "error", "warning", "success"
	Timestamp string `json:"timestamp"` // ISO 8601 timestamp
}

// StreamSession represents an active SSE connection
type StreamSession struct {
	ID      string
	Channel chan string
	Done    chan bool
}

var (
	sessions      = make(map[string]*StreamSession)
	sessionsMutex sync.RWMutex
)

// CreateSession creates a new SSE session
func CreateSession(sessionID string) *StreamSession {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()

	session := &StreamSession{
		ID:      sessionID,
		Channel: make(chan string, 100), // buffered channel to prevent blocking
		Done:    make(chan bool),
	}

	sessions[sessionID] = session
	return session
}

// GetSession retrieves an existing session
func GetSession(sessionID string) (*StreamSession, bool) {
	sessionsMutex.RLock()
	defer sessionsMutex.RUnlock()

	session, exists := sessions[sessionID]
	return session, exists
}

// CloseSession closes and removes a session
func CloseSession(sessionID string) {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()

	if session, exists := sessions[sessionID]; exists {
		close(session.Done)
		close(session.Channel)
		delete(sessions, sessionID)
	}
}

// SendProgress sends a progress update to the client
// progress should be between 0 and 100
// message is optional additional context
func SendProgress(sessionID string, progress float64, message string) {
	session, exists := GetSession(sessionID)
	if !exists {
		return // session not found or closed
	}

	event := ProgressEvent{
		Type:     "progress",
		Progress: progress,
		Message:  message,
	}

	data, err := json.Marshal(event)
	if err != nil {
		return
	}

	select {
	case session.Channel <- string(data):
		// sent successfully
	default:
		// channel full, skip this update
	}
}

// SendLog sends a log message to the client
// level can be: "info", "error", "warning", "success"
func SendLog(sessionID string, message string, level string) {
	session, exists := GetSession(sessionID)
	if !exists {
		return // session not found or closed
	}

	if level == "" {
		level = "info"
	}

	event := LogEvent{
		Type:      "log",
		Message:   message,
		Level:     level,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	data, err := json.Marshal(event)
	if err != nil {
		return
	}

	select {
	case session.Channel <- string(data):
		// sent successfully
	default:
		// channel full, skip this message
	}
}

// SendComplete sends a completion event and closes the session
func SendComplete(sessionID string) {
	SendProgress(sessionID, 100, "Import completed")
	SendLog(sessionID, "Import process finished successfully", "success")

	// Give a moment for messages to be sent before closing
	go func() {
		// Short delay to ensure last messages are sent
		CloseSession(sessionID)
	}()
}

// SendError sends an error event
func SendError(sessionID string, err error) {
	SendLog(sessionID, fmt.Sprintf("Error: %s", err.Error()), "error")
}
