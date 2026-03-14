package debuglog

import (
	"log"
	"os"
	"strings"
)

// Enabled is true when PYZANODE_DEBUG is set to 1, true, or yes (case-insensitive).
// Set PYZANODE_DEBUG=1 (or unset) to turn debug logging on or off without recompiling.
var Enabled bool

func init() {
	v := strings.ToLower(strings.TrimSpace(os.Getenv("PYZANODE_DEBUG")))
	Enabled = v == "1" || v == "true" || v == "yes"
}

// Printf logs to the standard logger with a [debug] prefix when Enabled is true.
// No-op when PYZANODE_DEBUG is unset or not 1/true/yes.
func Printf(format string, args ...interface{}) {
	if Enabled {
		log.Printf("[debug] "+format, args...)
	}
}
