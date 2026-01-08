package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"

	"github.com/ajmasia/shellify/internal/domain"
)

// Theme for all prompts - Catppuccin
var theme = huh.ThemeCatppuccin()

// ProjectPrompt collects project information interactively.
func ProjectPrompt() (name, description string, err error) {
	err = huh.NewInput().
		Title("Project name").
		Value(&name).
		Validate(func(s string) error {
			if strings.TrimSpace(s) == "" {
				return fmt.Errorf("name is required")
			}
			return nil
		}).
		WithTheme(theme).
		Run()
	if err != nil {
		return "", "", err
	}

	err = huh.NewInput().
		Title("Description (optional)").
		Value(&description).
		WithTheme(theme).
		Run()
	if err != nil {
		return "", "", err
	}

	return name, description, nil
}

// UpdateProjectPrompt collects updated project information interactively.
func UpdateProjectPrompt(current domain.Project) (name, description string, err error) {
	name = current.Name
	err = huh.NewInput().
		Title("Project name").
		Value(&name).
		Placeholder(current.Name).
		Validate(func(s string) error {
			if strings.TrimSpace(s) == "" {
				return fmt.Errorf("name is required")
			}
			return nil
		}).
		WithTheme(theme).
		Run()
	if err != nil {
		return "", "", err
	}

	description = current.Description
	err = huh.NewInput().
		Title("Description").
		Value(&description).
		Placeholder(current.Description).
		WithTheme(theme).
		Run()
	if err != nil {
		return "", "", err
	}

	return name, description, nil
}

// Confirm shows a simple y/N confirmation prompt.
func Confirm(message string) (bool, error) {
	var result bool
	err := huh.NewConfirm().
		Title(message).
		Value(&result).
		WithTheme(theme).
		Run()
	return result, err
}

// SelectProject shows an interactive project selector.
func SelectProject(projects []domain.Project) (*domain.Project, error) {
	if len(projects) == 0 {
		return nil, fmt.Errorf("no projects available")
	}

	var selected int

	options := make([]huh.Option[int], len(projects))
	for i, p := range projects {
		label := p.Name
		if p.Description != "" {
			label = fmt.Sprintf("%s - %s", p.Name, p.Description)
		}
		options[i] = huh.NewOption(label, i)
	}

	err := huh.NewSelect[int]().
		Title("Select project").
		Options(options...).
		Value(&selected).
		WithTheme(theme).
		Run()
	if err != nil {
		return nil, err
	}

	return &projects[selected], nil
}

// SelectUpdateFields prompts user to select which fields to update.
func SelectUpdateFields() (updateName, updateDescription bool, err error) {
	var selected []string

	options := []huh.Option[string]{
		huh.NewOption("Name", "name"),
		huh.NewOption("Description", "description"),
	}

	err = huh.NewMultiSelect[string]().
		Title("What do you want to update?").
		Options(options...).
		Value(&selected).
		WithTheme(theme).
		Run()
	if err != nil {
		return false, false, err
	}

	for _, s := range selected {
		if s == "name" {
			updateName = true
		}
		if s == "description" {
			updateDescription = true
		}
	}

	return updateName, updateDescription, nil
}

// PromptText shows a text input prompt with a default value.
func PromptText(label string, defaultValue string) (string, error) {
	result := defaultValue
	err := huh.NewInput().
		Title(label).
		Value(&result).
		Placeholder(defaultValue).
		WithTheme(theme).
		Run()
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(result) == "" {
		return defaultValue, nil
	}
	return result, nil
}

// Session prompts

// SessionPromptResult contains the result of SessionPrompt.
type SessionPromptResult struct {
	Name             string
	Description      string
	WorkingDirectory string
	Multiplexer      domain.MultiplexerType
	Environment      map[string]string
	PreCommands      []string
	PostCommands     []string
	Windows          []WindowPromptResult
	DefaultWindowID  string
}

// WindowPromptResult contains the result of WindowPrompt.
type WindowPromptResult struct {
	Name          string
	RootDirection domain.Direction
	Panes         []PanePromptResult
}

// PanePromptResult contains the result of PanePrompt.
type PanePromptResult struct {
	Name             string
	Command          string
	WorkingDirectory string
}

// SessionPrompt collects session information interactively with feedback.
func SessionPrompt() (SessionPromptResult, error) {
	var result SessionPromptResult

	// Session name
	err := huh.NewInput().
		Title("Session name").
		Value(&result.Name).
		Validate(func(s string) error {
			if strings.TrimSpace(s) == "" {
				return fmt.Errorf("name is required")
			}
			return nil
		}).
		WithTheme(theme).
		Run()
	if err != nil {
		return result, err
	}
	fmt.Printf("Name: %s\n", result.Name)

	// Multiplexer
	result.Multiplexer, err = SelectMultiplexer()
	if err != nil {
		return result, err
	}
	fmt.Printf("Multiplexer: %s\n", result.Multiplexer)

	// Working directory
	result.WorkingDirectory = "~"
	err = huh.NewInput().
		Title("Working directory").
		Value(&result.WorkingDirectory).
		Placeholder("~").
		WithTheme(theme).
		Run()
	if err != nil {
		return result, err
	}
	if result.WorkingDirectory == "" {
		result.WorkingDirectory = "~"
	}
	fmt.Printf("Directory: %s\n", result.WorkingDirectory)

	// Description
	err = huh.NewInput().
		Title("Description (optional)").
		Value(&result.Description).
		WithTheme(theme).
		Run()
	if err != nil {
		return result, err
	}
	if result.Description != "" {
		fmt.Printf("Description: %s\n", result.Description)
	}

	// Environment variables
	result.Environment, err = promptEnvironmentVariables()
	if err != nil {
		return result, err
	}

	// Pre-session commands
	result.PreCommands, err = promptCommandList("pre-session")
	if err != nil {
		return result, err
	}

	// Post-session commands
	result.PostCommands, err = promptCommandList("post-session")
	if err != nil {
		return result, err
	}

	// Windows loop
	fmt.Println()
	windowNum := 1
	for {
		window, err := WindowPrompt(windowNum)
		if err != nil {
			return result, err
		}
		result.Windows = append(result.Windows, window)
		fmt.Printf("Window %d: %s (%d panes)\n", windowNum, window.Name, len(window.Panes))
		windowNum++

		var addMore bool
		err = huh.NewConfirm().
			Title("Add another window?").
			Value(&addMore).
			WithTheme(theme).
			Run()
		if err != nil {
			return result, err
		}

		if !addMore {
			break
		}
	}

	// Default window selection (if multiple windows)
	if len(result.Windows) > 1 {
		result.DefaultWindowID, err = promptDefaultWindow(result.Windows)
		if err != nil {
			return result, err
		}
	}

	return result, nil
}

// promptEnvironmentVariables collects environment variables interactively.
func promptEnvironmentVariables() (map[string]string, error) {
	var addEnv bool
	err := huh.NewConfirm().
		Title("Add environment variables?").
		Value(&addEnv).
		WithTheme(theme).
		Run()
	if err != nil {
		return nil, err
	}

	if !addEnv {
		return nil, nil
	}

	envVars := make(map[string]string)

	for {
		var key, value string

		err := huh.NewInput().
			Title("Variable name").
			Value(&key).
			Validate(func(s string) error {
				if strings.TrimSpace(s) == "" {
					return fmt.Errorf("variable name is required")
				}
				return nil
			}).
			WithTheme(theme).
			Run()
		if err != nil {
			return nil, err
		}

		err = huh.NewInput().
			Title("Value (optional)").
			Value(&value).
			WithTheme(theme).
			Run()
		if err != nil {
			return nil, err
		}

		envVars[key] = value
		fmt.Printf("Env: %s=%s\n", key, value)

		var addMore bool
		err = huh.NewConfirm().
			Title("Add another environment variable?").
			Value(&addMore).
			WithTheme(theme).
			Run()
		if err != nil {
			return nil, err
		}

		if !addMore {
			break
		}
	}

	return envVars, nil
}

// promptCommandList collects a list of commands (pre or post session).
func promptCommandList(cmdType string) ([]string, error) {
	var addCmd bool
	err := huh.NewConfirm().
		Title(fmt.Sprintf("Add %s commands?", cmdType)).
		Value(&addCmd).
		WithTheme(theme).
		Run()
	if err != nil {
		return nil, err
	}

	if !addCmd {
		return nil, nil
	}

	var commands []string

	for {
		var cmd string

		err := huh.NewInput().
			Title("Command").
			Value(&cmd).
			Validate(func(s string) error {
				if strings.TrimSpace(s) == "" {
					return fmt.Errorf("command is required")
				}
				return nil
			}).
			WithTheme(theme).
			Run()
		if err != nil {
			return nil, err
		}

		commands = append(commands, cmd)
		prefix := "Pre"
		if cmdType == "post-session" {
			prefix = "Post"
		}
		fmt.Printf("%s: %s\n", prefix, cmd)

		var addMore bool
		err = huh.NewConfirm().
			Title("Add another command?").
			Value(&addMore).
			WithTheme(theme).
			Run()
		if err != nil {
			return nil, err
		}

		if !addMore {
			break
		}
	}

	return commands, nil
}

// promptDefaultWindow prompts user to select the default window.
func promptDefaultWindow(windows []WindowPromptResult) (string, error) {
	options := make([]huh.Option[string], len(windows)+1)
	options[0] = huh.NewOption("First window (default)", "")

	for i, w := range windows {
		options[i+1] = huh.NewOption(fmt.Sprintf("%s (window %d)", w.Name, i+1), fmt.Sprintf("w%d", i+1))
	}

	var selected string
	err := huh.NewSelect[string]().
		Title("Select default window").
		Options(options...).
		Value(&selected).
		WithTheme(theme).
		Run()
	if err != nil {
		return "", err
	}

	if selected != "" {
		// Find window name for feedback
		for i, w := range windows {
			if fmt.Sprintf("w%d", i+1) == selected {
				fmt.Printf("Default: %s (window %d)\n", w.Name, i+1)
				break
			}
		}
	}

	return selected, nil
}

// WindowPrompt collects window information with multiple panes support.
func WindowPrompt(windowNum int) (WindowPromptResult, error) {
	var result WindowPromptResult
	result.Name = fmt.Sprintf("window-%d", windowNum)

	// Window name
	err := huh.NewInput().
		Title("Window name").
		Value(&result.Name).
		Placeholder(fmt.Sprintf("window-%d", windowNum)).
		WithTheme(theme).
		Run()
	if err != nil {
		return result, err
	}
	if result.Name == "" {
		result.Name = fmt.Sprintf("window-%d", windowNum)
	}

	// Split direction
	result.RootDirection, err = SelectDirection()
	if err != nil {
		return result, err
	}

	// Panes loop
	paneNum := 1
	for {
		pane, err := PanePrompt(paneNum)
		if err != nil {
			return result, err
		}
		result.Panes = append(result.Panes, pane)
		paneNum++

		var addMore bool
		err = huh.NewConfirm().
			Title("Add another pane?").
			Value(&addMore).
			WithTheme(theme).
			Run()
		if err != nil {
			return result, err
		}

		if !addMore {
			break
		}
	}

	return result, nil
}

// PanePrompt collects pane information interactively.
func PanePrompt(paneNum int) (PanePromptResult, error) {
	var result PanePromptResult

	// Pane name (optional, defaults to pane-N)
	err := huh.NewInput().
		Title(fmt.Sprintf("Pane %d name (optional)", paneNum)).
		Value(&result.Name).
		Placeholder(fmt.Sprintf("pane-%d", paneNum)).
		WithTheme(theme).
		Run()
	if err != nil {
		return result, err
	}
	if result.Name == "" {
		result.Name = fmt.Sprintf("pane-%d", paneNum)
	}

	// Command (optional)
	err = huh.NewInput().
		Title("Command (optional)").
		Value(&result.Command).
		Placeholder("e.g., nvim .").
		WithTheme(theme).
		Run()
	if err != nil {
		return result, err
	}

	// Working directory (optional)
	err = huh.NewInput().
		Title("Working directory (optional)").
		Value(&result.WorkingDirectory).
		Placeholder("Leave empty for session default").
		WithTheme(theme).
		Run()
	if err != nil {
		return result, err
	}

	return result, nil
}

// SelectDirection shows a split direction selector.
func SelectDirection() (domain.Direction, error) {
	var selected string

	err := huh.NewSelect[string]().
		Title("Split direction").
		Options(
			huh.NewOption("Horizontal (side by side)", "horizontal"),
			huh.NewOption("Vertical (stacked)", "vertical"),
		).
		Value(&selected).
		WithTheme(theme).
		Run()
	if err != nil {
		return "", err
	}

	return domain.Direction(selected), nil
}

// SelectMultiplexer shows a multiplexer type selector.
func SelectMultiplexer() (domain.MultiplexerType, error) {
	var selected string

	err := huh.NewSelect[string]().
		Title("Target multiplexer").
		Options(
			huh.NewOption("tmux", "tmux"),
			huh.NewOption("zellij", "zellij"),
		).
		Value(&selected).
		WithTheme(theme).
		Run()
	if err != nil {
		return "", err
	}

	return domain.MultiplexerType(selected), nil
}

// SelectSession shows an interactive session selector.
func SelectSession(sessions []domain.SessionMetadata) (*domain.SessionMetadata, error) {
	if len(sessions) == 0 {
		return nil, fmt.Errorf("no sessions available")
	}

	var selected int

	options := make([]huh.Option[int], len(sessions))
	for i, s := range sessions {
		label := fmt.Sprintf("%s (%s)", s.Name, s.TargetMultiplexer)
		if s.Description != "" {
			label = fmt.Sprintf("%s - %s (%s)", s.Name, s.Description, s.TargetMultiplexer)
		}
		options[i] = huh.NewOption(label, i)
	}

	err := huh.NewSelect[int]().
		Title("Select session").
		Options(options...).
		Value(&selected).
		WithTheme(theme).
		Run()
	if err != nil {
		return nil, err
	}

	return &sessions[selected], nil
}

// SessionOption represents a session with its project info for selection.
type SessionOption struct {
	SessionID   string
	ProjectID   string
	ProjectName string
	Session     domain.SessionMetadata
}

// SelectSessionFromAll shows a session selector from all projects.
func SelectSessionFromAll(options []SessionOption) (sessionID, projectID string, err error) {
	if len(options) == 0 {
		return "", "", fmt.Errorf("no sessions available")
	}

	var selected int

	huhOptions := make([]huh.Option[int], len(options))
	for i, opt := range options {
		huhOptions[i] = huh.NewOption(
			fmt.Sprintf("[%s] %s (%s)", opt.ProjectName, opt.Session.Name, opt.Session.TargetMultiplexer),
			i,
		)
	}

	err = huh.NewSelect[int]().
		Title("Select session").
		Options(huhOptions...).
		Value(&selected).
		WithTheme(theme).
		Run()
	if err != nil {
		return "", "", err
	}

	return options[selected].SessionID, options[selected].ProjectID, nil
}

// SelectSessionUpdateFields prompts user to select which session fields to update.
func SelectSessionUpdateFields() (name, desc, workDir, multiplexer bool, err error) {
	var selected []string

	options := []huh.Option[string]{
		huh.NewOption("Name", "name"),
		huh.NewOption("Description", "description"),
		huh.NewOption("Working directory", "workdir"),
		huh.NewOption("Multiplexer", "multiplexer"),
	}

	err = huh.NewMultiSelect[string]().
		Title("What do you want to update?").
		Options(options...).
		Value(&selected).
		WithTheme(theme).
		Run()
	if err != nil {
		return false, false, false, false, err
	}

	for _, s := range selected {
		switch s {
		case "name":
			name = true
		case "description":
			desc = true
		case "workdir":
			workDir = true
		case "multiplexer":
			multiplexer = true
		}
	}

	return name, desc, workDir, multiplexer, nil
}

// PrintSessionSummary prints a complete summary of the session configuration.
func PrintSessionSummary(result SessionPromptResult, sessionName string) {
	fmt.Println()
	fmt.Println("Session Summary")
	fmt.Println("---------------")
	fmt.Printf("Name:        %s\n", result.Name)
	fmt.Printf("Session:     %s\n", sessionName)
	if result.Description != "" {
		fmt.Printf("Description: %s\n", result.Description)
	}
	fmt.Printf("Multiplexer: %s\n", result.Multiplexer)
	fmt.Printf("Directory:   %s\n", result.WorkingDirectory)
	fmt.Printf("Windows:     %d\n", len(result.Windows))

	// Count total panes
	totalPanes := 0
	for _, w := range result.Windows {
		totalPanes += len(w.Panes)
	}
	fmt.Printf("Panes:       %d\n", totalPanes)

	// Default window
	if result.DefaultWindowID != "" {
		for i, w := range result.Windows {
			if fmt.Sprintf("w%d", i+1) == result.DefaultWindowID {
				fmt.Printf("Default:     %s (window %d)\n", w.Name, i+1)
				break
			}
		}
	}

	// Environment variables
	if len(result.Environment) > 0 {
		fmt.Println()
		fmt.Println("Environment")
		for k, v := range result.Environment {
			fmt.Printf("  %s=%s\n", k, v)
		}
	}

	// Pre-commands
	if len(result.PreCommands) > 0 {
		fmt.Println()
		fmt.Println("Pre-commands")
		for _, cmd := range result.PreCommands {
			fmt.Printf("  %s\n", cmd)
		}
	}

	// Post-commands
	if len(result.PostCommands) > 0 {
		fmt.Println()
		fmt.Println("Post-commands")
		for _, cmd := range result.PostCommands {
			fmt.Printf("  %s\n", cmd)
		}
	}

	// Layout
	fmt.Println()
	fmt.Println("Layout")
	for i, w := range result.Windows {
		dirLabel := "horizontal"
		if w.RootDirection == domain.DirectionVertical {
			dirLabel = "vertical"
		}
		fmt.Printf("  Window %d: %s (%s split)\n", i+1, w.Name, dirLabel)
		for _, p := range w.Panes {
			name := p.Name
			if name == "" {
				name = "(unnamed)"
			}
			cmd := p.Command
			if cmd == "" {
				cmd = "(no command)"
			}
			fmt.Printf("    - %s: %s\n", name, cmd)
		}
	}
}

// ConfirmSessionCreation shows a confirmation prompt after the summary.
func ConfirmSessionCreation() (bool, error) {
	fmt.Println()
	return Confirm("Create this session?")
}

// SelectMultipleSessions shows an interactive multi-select for sessions.
func SelectMultipleSessions(options []SessionOption) ([]SessionOption, error) {
	if len(options) == 0 {
		return nil, fmt.Errorf("no sessions available")
	}

	var selected []int

	huhOptions := make([]huh.Option[int], len(options))
	for i, opt := range options {
		huhOptions[i] = huh.NewOption(
			fmt.Sprintf("[%s] %s (%s)", opt.ProjectName, opt.Session.Name, opt.Session.TargetMultiplexer),
			i,
		)
	}

	err := huh.NewMultiSelect[int]().
		Title("Select sessions (space to select, enter to confirm)").
		Options(huhOptions...).
		Value(&selected).
		WithTheme(theme).
		Run()
	if err != nil {
		return nil, err
	}

	if len(selected) == 0 {
		return nil, fmt.Errorf("no sessions selected")
	}

	result := make([]SessionOption, len(selected))
	for i, idx := range selected {
		result[i] = options[idx]
	}

	return result, nil
}
