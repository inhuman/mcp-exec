package exectool

import (
	"fmt"

	"github.com/inhuman/mcp-exec/internal/config"
)

// Description honestly lists the sandbox environment so the agent knows the
// available capabilities and limits (FR-016).
func Description(cfg config.Config) string {
	return fmt.Sprintf(
		"Execute Python code in a network-isolated sandbox and return its output. "+
			"Environment: Python 3 with stdlib plus PyYAML and Jinja2. "+
			"Limits: no network access; wall-clock timeout defaults to %ds (max %ds); "+
			"combined stdout+stderr is capped at %d bytes then truncated; stdin is capped at %d bytes. "+
			"A non-zero exit_code or timed_out=true is a normal result, not a tool error.",
		cfg.DefaultTimeoutS, cfg.MaxTimeoutS, cfg.MaxOutputBytes, cfg.MaxStdinBytes,
	)
}
