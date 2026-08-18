package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/localbeam/localbeam/internal/config"
	"github.com/localbeam/localbeam/internal/qr"
	"github.com/localbeam/localbeam/internal/transfer"
)

type Server struct {
	cfg     *config.Config
	manager *transfer.Manager
	mux     *http.ServeMux
}

func New(cfg *config.Config) *Server {
	timeout := time.Duration(cfg.Security.SessionTimeout) * time.Minute
	manager := transfer.NewManager(cfg.Transfer.UploadDir, timeout)

	s := &Server{
		cfg:     cfg,
		manager: manager,
		mux:     http.NewServeMux(),
	}

	s.registerRoutes()
	return s
}

func (s *Server) registerRoutes() {
	// API routes
	s.mux.HandleFunc("/api/session/create", s.handleCreateSession)
	s.mux.HandleFunc("/api/session/", s.handleSession)
	s.mux.HandleFunc("/api/join", s.handleJoinByPIN)
	s.mux.HandleFunc("/api/upload/", s.handleUpload)
	s.mux.HandleFunc("/api/download/", s.handleDownload)
	s.mux.HandleFunc("/api/text/", s.handleText)
	s.mux.HandleFunc("/api/qr/", s.handleQR)
	s.mux.HandleFunc("/api/info", s.handleInfo)

	// Static files (embedded HTML/CSS/JS)
	s.mux.HandleFunc("/", s.handleStatic)
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Server.Host, s.cfg.Server.Port)
	log.Printf("🚀 LocalBeam server running at http://%s", addr)

	localIP := getLocalIP()
	if localIP != "" {
		log.Printf("📱 Local network: http://%s:%d", localIP, s.cfg.Server.Port)
		url := fmt.Sprintf("http://%s:%d", localIP, s.cfg.Server.Port)
		qrStr, err := qr.GenerateTerminal(url)
		if err == nil {
			fmt.Println("\n📷 Scan to open on your device:")
			fmt.Println(qrStr)
			fmt.Printf("   URL: %s\n\n", url)
		}
	}

	return http.ListenAndServe(addr, s.corsMiddleware(s.mux))
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// POST /api/session/create
func (s *Server) handleCreateSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body struct {
		Type      string `json:"type"`
		Direction string `json:"direction"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		body.Type = "file"
		body.Direction = "push"
	}

	sessionType := transfer.TypeFile
	if body.Type == "text" {
		sessionType = transfer.TypeText
	}
	direction := transfer.DirPush
	if body.Direction == "pull" {
		direction = transfer.DirPull
	}

	session, err := s.manager.CreateSession(sessionType, direction)
	if err != nil {
		httpError(w, "failed to create session", http.StatusInternalServerError)
		return
	}

	localIP := getLocalIP()
	baseURL := fmt.Sprintf("http://%s:%d", localIP, s.cfg.Server.Port)
	receiveURL := fmt.Sprintf("%s/receive/%s", baseURL, session.ID)

	jsonOK(w, map[string]interface{}{
		"session_id":  session.ID,
		"pin":         session.PIN,
		"receive_url": receiveURL,
		"expires_at":  session.ExpiresAt,
	})
}

// GET /api/session/{id}
func (s *Server) handleSession(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/session/")
	session, err := s.manager.GetSession(id)
	if err != nil {
		httpError(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(s.manager.SessionToJSON(session))
}

// POST /api/join
func (s *Server) handleJoinByPIN(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		PIN string `json:"pin"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpError(w, "invalid body", http.StatusBadRequest)
		return
	}
	session, err := s.manager.GetSessionByPIN(body.PIN)
	if err != nil {
		httpError(w, "invalid or expired PIN", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(s.manager.SessionToJSON(session))
}

// POST /api/upload/{sessionID}
func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionID := strings.TrimPrefix(r.URL.Path, "/api/upload/")
	maxSize := s.cfg.Transfer.MaxFileSizeMB * 1024 * 1024
	r.ParseMultipartForm(maxSize)

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		// Single file
		file, header, err := r.FormFile("file")
		if err != nil {
			httpError(w, "no file provided", http.StatusBadRequest)
			return
		}
		defer file.Close()
		info, err := s.manager.AddFile(sessionID, header, file)
		if err != nil {
			httpError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonOK(w, map[string]interface{}{"files": []interface{}{info}})
		return
	}

	var infos []interface{}
	for _, fh := range files {
		file, err := fh.Open()
		if err != nil {
			continue
		}
		info, err := s.manager.AddFile(sessionID, fh, file)
		file.Close()
		if err != nil {
			continue
		}
		infos = append(infos, info)
	}
	jsonOK(w, map[string]interface{}{"files": infos})
}

// GET /api/download/{sessionID}/{fileID}
func (s *Server) handleDownload(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/download/"), "/")
	if len(parts) < 2 {
		httpError(w, "invalid path", http.StatusBadRequest)
		return
	}
	sessionID := parts[0]
	fileID := parts[1]

	fileInfo, err := s.manager.GetFile(sessionID, fileID)
	if err != nil {
		httpError(w, err.Error(), http.StatusNotFound)
		return
	}

	f, err := os.Open(fileInfo.Path)
	if err != nil {
		httpError(w, "file not available", http.StatusNotFound)
		return
	}
	defer f.Close()

	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileInfo.Name))
	w.Header().Set("Content-Type", fileInfo.MimeType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size))
	w.Header().Set("X-File-Hash", fileInfo.Hash)
	io.Copy(w, f)
}

// POST /api/text/{sessionID}
func (s *Server) handleText(w http.ResponseWriter, r *http.Request) {
	sessionID := strings.TrimPrefix(r.URL.Path, "/api/text/")

	if r.Method == http.MethodPost {
		var body struct {
			Content string `json:"content"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			httpError(w, "invalid body", http.StatusBadRequest)
			return
		}
		if err := s.manager.SetTextContent(sessionID, body.Content); err != nil {
			httpError(w, err.Error(), http.StatusNotFound)
			return
		}
		jsonOK(w, map[string]string{"status": "ok"})
		return
	}

	if r.Method == http.MethodGet {
		session, err := s.manager.GetSession(sessionID)
		if err != nil {
			httpError(w, err.Error(), http.StatusNotFound)
			return
		}
		jsonOK(w, map[string]string{"content": session.Content})
		return
	}

	httpError(w, "method not allowed", http.StatusMethodNotAllowed)
}

// GET /api/qr/{sessionID}
func (s *Server) handleQR(w http.ResponseWriter, r *http.Request) {
	sessionID := strings.TrimPrefix(r.URL.Path, "/api/qr/")

	localIP := getLocalIP()
	receiveURL := fmt.Sprintf("http://%s:%d/receive/%s", localIP, s.cfg.Server.Port, sessionID)

	pngData, err := qr.GeneratePNG(receiveURL, 400)
	if err != nil {
		// Return base64 JSON fallback
		httpError(w, "qr generation failed", http.StatusInternalServerError)
		return
	}

	format := r.URL.Query().Get("format")
	if format == "base64" {
		jsonOK(w, map[string]string{
			"data": base64.StdEncoding.EncodeToString(pngData),
			"url":  receiveURL,
		})
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(pngData)
}

// GET /api/info
func (s *Server) handleInfo(w http.ResponseWriter, r *http.Request) {
	localIP := getLocalIP()
	jsonOK(w, map[string]interface{}{
		"version":    "1.0.0",
		"local_ip":   localIP,
		"port":       s.cfg.Server.Port,
		"app_name":   s.cfg.UI.AppName,
		"base_url":   fmt.Sprintf("http://%s:%d", localIP, s.cfg.Server.Port),
		"max_file_mb": s.cfg.Transfer.MaxFileSizeMB,
	})
}

// Static files handler
func (s *Server) handleStatic(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// Route /receive/{id} to SPA
	if strings.HasPrefix(path, "/receive/") || path == "/" || path == "/send" || path == "/join" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(indexHTML))
		return
	}

	// Try to serve embedded file
	switch filepath.Ext(path) {
	case ".css":
		w.Header().Set("Content-Type", "text/css")
	case ".js":
		w.Header().Set("Content-Type", "application/javascript")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	case ".ico":
		w.Header().Set("Content-Type", "image/x-icon")
	}

	http.NotFound(w, r)
}

func httpError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func jsonOK(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func getLocalIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "localhost"
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			if ip = ip.To4(); ip != nil {
				return ip.String()
			}
		}
	}
	return "localhost"
}
