package cli

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/ajmasia/shellify/internal/application"
	"github.com/ajmasia/shellify/internal/infrastructure/storage"
	httpserver "github.com/ajmasia/shellify/internal/interfaces/http"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the HTTP API server",
	Long: `Start the HTTP server for the REST API.

The server provides a REST API for managing projects and sessions.
Use --static flag to serve a GUI directory.

Examples:
  sfy server              # Start on default port 3777
  sfy server -p 8080      # Start on port 8080
  sfy server --open       # Start and open browser
  sfy server -d           # Start in background (daemon mode)
  sfy server stop         # Stop the background server
  sfy server status       # Check if server is running`,
	RunE: runServer,
}

var serverStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the background server",
	RunE:  runServerStop,
}

var serverStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check if the server is running",
	RunE:  runServerStatus,
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.AddCommand(serverStopCmd)
	serverCmd.AddCommand(serverStatusCmd)

	serverCmd.Flags().IntP("port", "p", 3777, "Port to listen on")
	serverCmd.Flags().StringP("host", "H", "localhost", "Host to bind to")
	serverCmd.Flags().BoolP("open", "o", false, "Open browser automatically")
	serverCmd.Flags().StringP("static", "s", "", "Path to static files directory")
	serverCmd.Flags().BoolP("daemon", "d", false, "Run in background (daemon mode)")
}

func runServer(cmd *cobra.Command, args []string) error {
	port, _ := cmd.Flags().GetInt("port")
	host, _ := cmd.Flags().GetString("host")
	openBrowser, _ := cmd.Flags().GetBool("open")
	staticDir, _ := cmd.Flags().GetString("static")
	daemon, _ := cmd.Flags().GetBool("daemon")

	// If daemon mode, fork the process
	if daemon {
		return startDaemon(port, host, staticDir, openBrowser)
	}

	// Initialize storage
	store, err := storage.NewStorage()
	if err != nil {
		return fmt.Errorf("initializing storage: %w", err)
	}

	// Initialize services
	projectSvc := application.NewProjectService(store)
	sessionSvc := application.NewSessionService(store, store)
	generatorSvc := application.NewGeneratorService(store, store)
	launcherSvc, err := application.NewLauncherService(store, store, generatorSvc)
	if err != nil {
		return fmt.Errorf("initializing launcher service: %w", err)
	}

	// Create server config
	config := httpserver.Config{
		Port:      port,
		Host:      host,
		StaticDir: staticDir,
	}

	if staticDir != "" {
		fmt.Printf("Serving static files from: %s\n", staticDir)
	}

	// Create server
	srv := httpserver.New(config, projectSvc, sessionSvc, launcherSvc, store)

	// Write PID file
	if err := writePIDFile(); err != nil {
		fmt.Printf("Warning: could not write PID file: %v\n", err)
	}
	defer removePIDFile()

	// Handle graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start server in goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- srv.Start()
	}()

	// Open browser if requested
	if openBrowser {
		url := fmt.Sprintf("http://%s:%d", host, port)
		openURL(url)
	}

	// Wait for shutdown signal or error
	select {
	case <-ctx.Done():
		fmt.Println("\nShutting down server...")
		return srv.Shutdown(context.Background())
	case err := <-errChan:
		return err
	}
}

func runServerStop(cmd *cobra.Command, args []string) error {
	pid, err := readPIDFile()
	if err != nil {
		return fmt.Errorf("server is not running (no PID file found)")
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		removePIDFile()
		return fmt.Errorf("server process not found")
	}

	// Send SIGTERM for graceful shutdown
	if err := process.Signal(syscall.SIGTERM); err != nil {
		removePIDFile()
		return fmt.Errorf("failed to stop server: %w", err)
	}

	removePIDFile()
	fmt.Println("Server stopped")
	return nil
}

func runServerStatus(cmd *cobra.Command, args []string) error {
	pid, err := readPIDFile()
	if err != nil {
		fmt.Println("Server is not running")
		return nil
	}

	// Check if process is still alive
	process, err := os.FindProcess(pid)
	if err != nil {
		fmt.Println("Server is not running")
		removePIDFile()
		return nil
	}

	// On Unix, FindProcess always succeeds, so we need to send signal 0 to check
	if err := process.Signal(syscall.Signal(0)); err != nil {
		fmt.Println("Server is not running")
		removePIDFile()
		return nil
	}

	fmt.Printf("Server is running (PID: %d)\n", pid)
	return nil
}

func startDaemon(port int, host, staticDir string, openBrowser bool) error {
	// Check if already running
	if pid, err := readPIDFile(); err == nil {
		process, _ := os.FindProcess(pid)
		if process != nil && process.Signal(syscall.Signal(0)) == nil {
			return fmt.Errorf("server is already running (PID: %d)", pid)
		}
	}

	// Get the executable path
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Build arguments for the background process (without -d flag)
	args := []string{"server", "-p", strconv.Itoa(port), "-H", host}
	if staticDir != "" {
		args = append(args, "-s", staticDir)
	}
	if openBrowser {
		args = append(args, "-o")
	}

	// Start the server in background
	cmd := exec.Command(exe, args...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil

	// Detach from parent process
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start daemon: %w", err)
	}

	fmt.Printf("Server started in background (PID: %d)\n", cmd.Process.Pid)
	fmt.Printf("URL: http://%s:%d\n", host, port)
	fmt.Println("Use 'sfy server stop' to stop the server")

	return nil
}

// PID file helpers
func getPIDFilePath() string {
	return filepath.Join(os.TempDir(), "shellify-server.pid")
}

func writePIDFile() error {
	return os.WriteFile(getPIDFilePath(), []byte(strconv.Itoa(os.Getpid())), 0644)
}

func readPIDFile() (int, error) {
	data, err := os.ReadFile(getPIDFilePath())
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(data))
}

func removePIDFile() {
	_ = os.Remove(getPIDFilePath())
}

// openURL opens a URL in the default browser
func openURL(url string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default: // linux, freebsd, etc.
		cmd = exec.Command("xdg-open", url)
	}

	_ = cmd.Start()
}
