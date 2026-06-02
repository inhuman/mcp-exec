package config

import "github.com/inhuman/config"

// Config holds the server configuration loaded from environment variables.
// All values are documented in CLAUDE.md.
type Config struct {
	Transport       string `env:"MCP_EXEC_TRANSPORT" env-default:"stdio"`
	Addr            string `env:"MCP_EXEC_ADDR" env-default:":8080"`
	DefaultTimeoutS int    `env:"MCP_EXEC_DEFAULT_TIMEOUT_S" env-default:"30"`
	MaxTimeoutS     int    `env:"MCP_EXEC_MAX_TIMEOUT_S" env-default:"300"`
	MaxOutputBytes  int    `env:"MCP_EXEC_MAX_OUTPUT_BYTES" env-default:"1048576"`
	MaxStdinBytes   int    `env:"MCP_EXEC_MAX_STDIN_BYTES" env-default:"1048576"`
	Python          string `env:"MCP_EXEC_PYTHON" env-default:"python3"`
	// AuthToken, when non-empty, requires HTTP/SSE requests to carry a matching
	// X-MCP-AUTH header. Empty disables auth. Not applicable to stdio.
	AuthToken string `env:"MCP_EXEC_AUTH_TOKEN" env-default:""`
}

func Load() (Config, error) {
	var c Config
	if err := config.Load(&c); err != nil {
		return Config{}, err
	}
	return c, nil
}
