package transfer

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type SessionType string

const (
	TypeText SessionType = "text"
	TypeFile SessionType = "file"
)

type Direction string

const (
	DirPush Direction = "push" // sender -> receiver
	DirPull Direction = "pull" // receiver pulls from sender
)

type Session struct {
	ID        string      `json:"id"`
	PIN       string      `json:"pin"`
	Type      SessionType `json:"type"`
	Direction Direction   `json:"direction"`
	CreatedAt time.Time   `json:"created_at"`
	ExpiresAt time.Time   `json:"expires_at"`
	Content   string      `json:"content,omitempty"`
	Files     []FileInfo  `json:"files,omitempty"`
	Downloaded bool       `json:"downloaded"`
	mu        sync.Mutex
}

type FileInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	MimeType string `json:"mime_type"`
	Path     string `json:"-"`
	Hash     string `json:"hash"`
}

type Manager struct {
	sessions  map[string]*Session
	mu        sync.RWMutex
	uploadDir string
	timeout   time.Duration
}

func NewManager(uploadDir string, timeout time.Duration) *Manager {
	m := &Manager{
		sessions:  make(map[string]*Session),
		uploadDir: uploadDir,
		timeout:   timeout,
	}
	go m.cleanupLoop()
	return m
}

func (m *Manager) CreateSession(sessionType SessionType, dir Direction) (*Session, error) {
	id, err := generateID(16)
	if err != nil {
		return nil, err
	}
	pin, err := generatePIN(6)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	session := &Session{
		ID:        id,
		PIN:       pin,
		Type:      sessionType,
		Direction: dir,
		CreatedAt: now,
		ExpiresAt: now.Add(m.timeout),
	}

	m.mu.Lock()
	m.sessions[id] = session
	m.mu.Unlock()

	return session, nil
}

func (m *Manager) GetSession(id string) (*Session, error) {
	m.mu.RLock()
	session, ok := m.sessions[id]
	m.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("session not found")
	}
	if time.Now().After(session.ExpiresAt) {
		m.DeleteSession(id)
		return nil, fmt.Errorf("session expired")
	}
	return session, nil
}

func (m *Manager) GetSessionByPIN(pin string) (*Session, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, s := range m.sessions {
		if s.PIN == pin && time.Now().Before(s.ExpiresAt) {
			return s, nil
		}
	}
	return nil, fmt.Errorf("invalid or expired PIN")
}

func (m *Manager) SetTextContent(id, content string) error {
	session, err := m.GetSession(id)
	if err != nil {
		return err
	}
	session.mu.Lock()
	defer session.mu.Unlock()
	session.Type = TypeText
	session.Content = content
	return nil
}

func (m *Manager) AddFile(id string, header *multipart.FileHeader, file multipart.File) (*FileInfo, error) {
	session, err := m.GetSession(id)
	if err != nil {
		return nil, err
	}

	fileID, _ := generateID(8)
	dir := filepath.Join(m.uploadDir, "localbeam", id)
	os.MkdirAll(dir, 0700)

	destPath := filepath.Join(dir, fileID+"_"+header.Filename)
	dest, err := os.Create(destPath)
	if err != nil {
		return nil, err
	}
	defer dest.Close()

	hasher := sha256.New()
	writer := io.MultiWriter(dest, hasher)

	size, err := io.Copy(writer, file)
	if err != nil {
		os.Remove(destPath)
		return nil, err
	}

	info := &FileInfo{
		ID:       fileID,
		Name:     header.Filename,
		Size:     size,
		MimeType: header.Header.Get("Content-Type"),
		Path:     destPath,
		Hash:     hex.EncodeToString(hasher.Sum(nil)),
	}

	session.mu.Lock()
	session.Type = TypeFile
	session.Files = append(session.Files, *info)
	session.mu.Unlock()

	return info, nil
}

func (m *Manager) GetFile(sessionID, fileID string) (*FileInfo, error) {
	session, err := m.GetSession(sessionID)
	if err != nil {
		return nil, err
	}
	session.mu.Lock()
	defer session.mu.Unlock()

	for _, f := range session.Files {
		if f.ID == fileID {
			return &f, nil
		}
	}
	return nil, fmt.Errorf("file not found")
}

func (m *Manager) MarkDownloaded(id string) {
	if s, err := m.GetSession(id); err == nil {
		s.mu.Lock()
		s.Downloaded = true
		s.mu.Unlock()
	}
}

func (m *Manager) DeleteSession(id string) {
	m.mu.Lock()
	_, ok := m.sessions[id]
	if ok {
		delete(m.sessions, id)
	}
	m.mu.Unlock()

	if ok {
		// Cleanup files
		dir := filepath.Join(m.uploadDir, "localbeam", id)
		os.RemoveAll(dir)
	}
}

func (m *Manager) SessionToJSON(s *Session) []byte {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, _ := json.Marshal(s)
	return data
}

func (m *Manager) cleanupLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		m.mu.Lock()
		for id, s := range m.sessions {
			if time.Now().After(s.ExpiresAt) {
				delete(m.sessions, id)
				dir := filepath.Join(m.uploadDir, "localbeam", id)
				os.RemoveAll(dir)
			}
		}
		m.mu.Unlock()
	}
}

func generateID(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b)[:n], nil
}

func generatePIN(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	pin := ""
	for i := 0; i < n; i++ {
		pin += fmt.Sprintf("%d", int(b[i])%10)
	}
	return pin, nil
}
