package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Server   ServerConfig   `json:"server"`
	Security SecurityConfig `json:"security"`
	Transfer TransferConfig `json:"transfer"`
	UI       UIConfig       `json:"ui"`
}

type ServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	TLS  bool   `json:"tls"`
}

type SecurityConfig struct {
	SessionTimeout int      `json:"session_timeout_minutes"`
	MaxSessions    int      `json:"max_sessions"`
	PinLength      int      `json:"pin_length"`
	AllowedOrigins []string `json:"allowed_origins"`
}

type TransferConfig struct {
	MaxFileSizeMB   int64    `json:"max_file_size_mb"`
	AllowedTypes    []string `json:"allowed_types"`
	ChunkSizeKB     int      `json:"chunk_size_kb"`
	UploadDir       string   `json:"upload_dir"`
	AutoCleanupMins int      `json:"auto_cleanup_minutes"`
}

type UIConfig struct {
	Theme   string `json:"theme"`
	AppName string `json:"app_name"`
}

func Default() *Config {
	return &Config{
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: 3001,
			TLS:  false,
		},
		Security: SecurityConfig{
			SessionTimeout: 10,
			MaxSessions:    10,
			PinLength:      6,
			AllowedOrigins: []string{"*"},
		},
		Transfer: TransferConfig{
			MaxFileSizeMB:   500,
			AllowedTypes:    []string{},
			ChunkSizeKB:     256,
			UploadDir:       os.TempDir(),
			AutoCleanupMins: 30,
		},
		UI: UIConfig{
			Theme:   "dark",
			AppName: "LocalBeam",
		},
	}
}

func Load(path string) (*Config, error) {
	cfg := Default()
	if path == "" {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, ".localbeam", "config.json")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func Save(cfg *Config, path string) error {
	if path == "" {
		home, _ := os.UserHomeDir()
		dir := filepath.Join(home, ".localbeam")
		os.MkdirAll(dir, 0700)
		path = filepath.Join(dir, "config.json")
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}
