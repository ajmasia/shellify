package domain

// Config holds application configuration.
type Config struct {
	DefaultEditor    string `json:"defaultEditor"`
	DefaultTerminal  string `json:"defaultTerminal"`
	CLIIntegration   bool   `json:"cliIntegration"`
	MaximizeTerminal bool   `json:"maximizeTerminal"`
	ServerPort       int    `json:"serverPort"`
}

// DefaultConfig returns the default application configuration.
func DefaultConfig() Config {
	return Config{
		DefaultEditor:    "$EDITOR",
		DefaultTerminal:  "system",
		CLIIntegration:   false,
		MaximizeTerminal: true,
		ServerPort:       3777,
	}
}
