package config

import (
	"os"
	"path/filepath"
)

// ControllerConfig holds controller configuration.
type ControllerConfig struct {
	HTTPAddr      string `json:"http_addr"`
	WebSocketPath string `json:"ws_path"`
	DataDir       string `json:"data_dir"`
	SecretsFile   string `json:"secrets_file"`
	// WebDir is optional: path to built web dashboard (e.g. web/dist). If set, controller serves the dashboard so you run one process.
	WebDir    string `json:"web_dir"`
	AuthToken string `json:"-"` // loaded from secrets
}

// DefaultControllerConfig returns defaults. DataDir is set from env or current dir.
func DefaultControllerConfig() ControllerConfig {
	dataDir := os.Getenv("PYZANODE_DATA")
	if dataDir == "" {
		home, _ := os.UserHomeDir()
		dataDir = filepath.Join(home, ".pyzanode")
	}
	webDir := os.Getenv("PYZANODE_WEB")
	return ControllerConfig{
		HTTPAddr:      "0.0.0.0:9451",
		WebSocketPath: "/ws",
		DataDir:       dataDir,
		SecretsFile:   "secrets.json",
		WebDir:        webDir,
	}
}

// AgentConfig holds agent configuration.
type AgentConfig struct {
	ControllerURL string `json:"controller_url"`
	Token         string `json:"token"`
	HeartbeatSec  int    `json:"heartbeat_sec"`
}
