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
	Windows          []WindowPromptResult
}

// WindowPromptResult contains the result of WindowPrompt.
type WindowPromptResult struct {
	Name    string
	Command string
}

// SessionPrompt collects session information interactively (simplified for v0.3.0).
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

	// Multiplexer
	result.Multiplexer, err = SelectMultiplexer()
	if err != nil {
		return result, err
	}

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

	// Description
	err = huh.NewInput().
		Title("Description (optional)").
		Value(&result.Description).
		WithTheme(theme).
		Run()
	if err != nil {
		return result, err
	}

	// Windows loop
	windowNum := 1
	for {
		window, err := WindowPrompt(windowNum)
		if err != nil {
			return result, err
		}
		result.Windows = append(result.Windows, window)
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

	return result, nil
}

// WindowPrompt collects window information (simplified: one pane per window).
func WindowPrompt(windowNum int) (WindowPromptResult, error) {
	var result WindowPromptResult
	result.Name = fmt.Sprintf("window-%d", windowNum)

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

	err = huh.NewInput().
		Title("Command (optional)").
		Value(&result.Command).
		Placeholder("e.g., nvim .").
		WithTheme(theme).
		Run()
	if err != nil {
		return result, err
	}

	return result, nil
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
