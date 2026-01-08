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
