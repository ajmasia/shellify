package cli

import (
	"bytes"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/ajmasia/shellify/internal/domain"
)

// yamlSession is the user-facing YAML representation of a session.
// It contains no internal IDs or storage wrapper fields.
type yamlSession struct {
	Name             string            `yaml:"name"`
	Description      string            `yaml:"description,omitempty"`
	Multiplexer      string            `yaml:"multiplexer,omitempty"`
	WorkingDirectory string            `yaml:"workingDirectory,omitempty"`
	Env              map[string]string `yaml:"env,omitempty"`
	PreCommands      []string          `yaml:"preCommands,omitempty"`
	PostCommands     []string          `yaml:"postCommands,omitempty"`
	DefaultWindow    string            `yaml:"defaultWindow,omitempty"`
	Windows          []yamlWindow      `yaml:"windows"`
}

type yamlWindow struct {
	Name             string     `yaml:"name"`
	WorkingDirectory string     `yaml:"workingDirectory,omitempty"`
	Direction        string     `yaml:"direction,omitempty"`
	Panes            []yamlPane `yaml:"panes"`
}

// yamlPane represents either a leaf pane (has command/workingDirectory) or a
// container pane (has panes). The presence of panes makes it a container.
type yamlPane struct {
	Name             string     `yaml:"name,omitempty"`
	Command          string     `yaml:"command,omitempty"`
	WorkingDirectory string     `yaml:"workingDirectory,omitempty"`
	Size             float64    `yaml:"size,omitempty"`
	Direction        string     `yaml:"direction,omitempty"`
	Panes            []yamlPane `yaml:"panes,omitempty"`
}

// SessionToYAML converts a domain.Session to a clean, user-facing YAML
// representation with no internal IDs or storage wrapper fields.
func SessionToYAML(s domain.Session) ([]byte, error) {
	ys := yamlSession{
		Name:             s.Name,
		Description:      s.Description,
		Multiplexer:      string(s.TargetMultiplexer),
		WorkingDirectory: s.WorkingDirectory,
		Env:              s.Environment,
		PreCommands:      s.PreCommands,
		PostCommands:     s.PostCommands,
		Windows:          make([]yamlWindow, len(s.Windows)),
	}

	// Resolve default window name
	for _, w := range s.Windows {
		if w.ID == s.DefaultWindowID {
			ys.DefaultWindow = w.Name
			break
		}
	}

	for i, w := range s.Windows {
		ys.Windows[i] = domainWindowToYAML(w)
	}

	return yaml.Marshal(ys)
}

func domainWindowToYAML(w domain.Window) yamlWindow {
	yw := yamlWindow{
		Name:             w.Name,
		WorkingDirectory: w.WorkingDirectory,
	}

	if len(w.Panes) == 0 {
		return yw
	}

	root := w.Panes[0]

	if root.IsContainer() && len(w.Panes) == 1 {
		// Single root container: unwrap it so the YAML shows a flat panes list
		// with the direction at the window level.
		yw.Direction = string(root.Direction)
		yw.Panes = domainPanesToYAML(root.Children)
	} else {
		// Single leaf pane, or unusual multi-root case.
		yw.Direction = string(w.RootDirection)
		yw.Panes = domainPanesToYAML(w.Panes)
	}

	return yw
}

func domainPanesToYAML(panes []domain.Pane) []yamlPane {
	result := make([]yamlPane, len(panes))
	for i, p := range panes {
		result[i] = domainPaneToYAML(p)
	}
	return result
}

func domainPaneToYAML(p domain.Pane) yamlPane {
	yp := yamlPane{
		Name:             p.Name,
		Command:          p.Command,
		WorkingDirectory: p.WorkingDirectory,
		Size:             p.Size,
	}

	if p.IsContainer() {
		yp.Direction = string(p.Direction)
		yp.Panes = domainPanesToYAML(p.Children)
		// Clear leaf fields on containers
		yp.Name = ""
		yp.Command = ""
		yp.WorkingDirectory = ""
	}

	return yp
}

// YAMLToSession parses user-authored YAML and returns a domain.Session.
// Stable IDs are generated from window and pane positions so that
// round-tripping through export → edit → import does not create duplicate
// sessions or break running-session detection.
func YAMLToSession(data []byte) (domain.Session, error) {
	var ys yamlSession

	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	if err := dec.Decode(&ys); err != nil {
		return domain.Session{}, fmt.Errorf("invalid YAML: %w", err)
	}

	if strings.TrimSpace(ys.Name) == "" {
		return domain.Session{}, fmt.Errorf("session name is required")
	}

	if len(ys.Windows) == 0 {
		return domain.Session{}, fmt.Errorf("session must have at least one window")
	}

	multiplexer := domain.MultiplexerType(ys.Multiplexer)
	if multiplexer == "" {
		multiplexer = domain.MultiplexerTmux
	}
	if multiplexer != domain.MultiplexerTmux {
		return domain.Session{}, fmt.Errorf("unsupported multiplexer %q: only tmux is supported in this version", ys.Multiplexer)
	}

	s := domain.Session{
		Name:              strings.TrimSpace(ys.Name),
		Description:       strings.TrimSpace(ys.Description),
		TargetMultiplexer: multiplexer,
		WorkingDirectory:  strings.TrimSpace(ys.WorkingDirectory),
		Environment:       ys.Env,
		PreCommands:       ys.PreCommands,
		PostCommands:      ys.PostCommands,
		Windows:           make([]domain.Window, len(ys.Windows)),
	}

	for i, yw := range ys.Windows {
		w, err := yamlWindowToDomain(yw, i+1)
		if err != nil {
			return domain.Session{}, fmt.Errorf("window %q: %w", yw.Name, err)
		}
		s.Windows[i] = w
	}

	// Resolve default window by name; fall back to first window if not found.
	defaultWindowName := strings.TrimSpace(ys.DefaultWindow)
	if defaultWindowName != "" {
		for _, w := range s.Windows {
			if strings.EqualFold(w.Name, defaultWindowName) {
				s.DefaultWindowID = w.ID
				break
			}
		}
	}
	if s.DefaultWindowID == "" && len(s.Windows) > 0 {
		s.DefaultWindowID = s.Windows[0].ID
	}

	return s, nil
}

func yamlWindowToDomain(yw yamlWindow, windowIdx int) (domain.Window, error) {
	name := strings.TrimSpace(yw.Name)
	if name == "" {
		name = fmt.Sprintf("window-%d", windowIdx)
	}

	if len(yw.Panes) == 0 {
		return domain.Window{}, fmt.Errorf("must have at least one pane")
	}

	direction := domain.DirectionHorizontal
	if yw.Direction != "" {
		direction = domain.Direction(yw.Direction)
		if !direction.IsValid() {
			return domain.Window{}, fmt.Errorf("invalid direction %q: must be horizontal or vertical", yw.Direction)
		}
	}

	windowID := fmt.Sprintf("w%d", windowIdx)
	basePath := windowID

	w := domain.Window{
		ID:               windowID,
		Name:             name,
		WorkingDirectory: strings.TrimSpace(yw.WorkingDirectory),
		RootDirection:    direction,
	}

	if len(yw.Panes) == 1 && len(yw.Panes[0].Panes) == 0 {
		// Single leaf pane — store directly without a root container.
		p := yamlPaneToDomain(yw.Panes[0], basePath+"-p1", 100)
		w.Panes = []domain.Pane{p}
	} else {
		// Multiple panes or a single container — wrap in root container.
		children := yamlPanesToDomain(yw.Panes, basePath, direction)
		w.Panes = []domain.Pane{
			{
				ID:        basePath + "-root",
				Size:      100,
				Direction: direction,
				Children:  children,
			},
		}
	}

	return w, nil
}

func yamlPanesToDomain(panes []yamlPane, parentPath string, parentDir domain.Direction) []domain.Pane {
	n := len(panes)
	result := make([]domain.Pane, n)
	for i, yp := range panes {
		paneID := fmt.Sprintf("%s-p%d", parentPath, i+1)
		size := yp.Size
		if size == 0 {
			size = 100.0 / float64(n)
		}
		result[i] = yamlPaneToDomainWithDir(yp, paneID, size, parentDir)
	}
	return result
}

func yamlPaneToDomain(yp yamlPane, id string, size float64) domain.Pane {
	return yamlPaneToDomainWithDir(yp, id, size, domain.DirectionHorizontal)
}

func yamlPaneToDomainWithDir(yp yamlPane, id string, size float64, parentDir domain.Direction) domain.Pane {
	if len(yp.Panes) > 0 {
		// Container pane
		dir := parentDir
		if yp.Direction != "" {
			dir = domain.Direction(yp.Direction)
		}
		return domain.Pane{
			ID:        id,
			Size:      size,
			Direction: dir,
			Children:  yamlPanesToDomain(yp.Panes, id, dir),
		}
	}

	// Leaf pane
	return domain.Pane{
		ID:               id,
		Name:             strings.TrimSpace(yp.Name),
		Command:          strings.TrimSpace(yp.Command),
		WorkingDirectory: strings.TrimSpace(yp.WorkingDirectory),
		Size:             size,
	}
}
