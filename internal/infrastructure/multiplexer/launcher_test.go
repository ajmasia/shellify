package multiplexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStripAnsiCodes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no ansi codes",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "with color code",
			input:    "\x1b[32mgreen\x1b[0m",
			expected: "green",
		},
		{
			name:     "with bold",
			input:    "\x1b[1mbold\x1b[0m",
			expected: "bold",
		},
		{
			name:     "mixed content",
			input:    "prefix\x1b[31mred\x1b[0msuffix",
			expected: "prefixredsuffix",
		},
		{
			name:     "multiple codes",
			input:    "\x1b[1;32mbold green\x1b[0m normal",
			expected: "bold green normal",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stripAnsiCodes(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestContainsSession(t *testing.T) {
	tests := []struct {
		name        string
		output      string
		sessionName string
		expected    bool
	}{
		{
			name:        "exact match",
			output:      "my-session\nother-session\n",
			sessionName: "my-session",
			expected:    true,
		},
		{
			name:        "match with status",
			output:      "my-session (attached)\nother-session\n",
			sessionName: "my-session",
			expected:    true,
		},
		{
			name:        "no match",
			output:      "other-session\nanother-session\n",
			sessionName: "my-session",
			expected:    false,
		},
		{
			name:        "exited session should not match",
			output:      "my-session (EXITED - attach to resurrect)\n",
			sessionName: "my-session",
			expected:    false,
		},
		{
			name:        "partial name should not match",
			output:      "my-session-extended\n",
			sessionName: "my-session",
			expected:    false,
		},
		{
			name:        "empty output",
			output:      "",
			sessionName: "my-session",
			expected:    false,
		},
		{
			name:        "single line exact",
			output:      "my-session",
			sessionName: "my-session",
			expected:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsSession(tt.output, tt.sessionName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewLauncher(t *testing.T) {
	launcher, err := NewLauncher()
	assert.NoError(t, err)
	assert.NotNil(t, launcher)
	assert.NotEmpty(t, launcher.tempDir)
}
