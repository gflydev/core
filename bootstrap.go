package core

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/gflydev/core/log"
	"github.com/gflydev/core/utils"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ansiStripper is a writer that strips ANSI color codes before writing to the underlying writer.
type ansiStripper struct {
	writer io.Writer
}

// Write implements the io.Writer interface.
// It strips ANSI color codes from the input before writing to the underlying writer.
func (s *ansiStripper) Write(p []byte) (n int, err error) {
	// Regular expression to match ANSI color codes
	re := regexp.MustCompile("\033\\[[0-9;]*m")

	// Strip ANSI color codes
	clean := re.ReplaceAll(p, []byte(""))

	// Write the cleaned content to the underlying writer
	_, err = s.writer.Write(clean)

	// Return the original length to satisfy the Writer interface
	return len(p), err
}

// ====================================================================
//                              Bootstrap
// ====================================================================

// startupMessage prints the application startup message to the logs.
//
// Parameters:
//   - url (string): The server URL the application is running on.
//   - name (string): The name of the application.
//   - env (string): The application's runtime environment (e.g., development, production).
func startupMessage(url, name, env string) {
	// Create a pretty box for the startup message
	headerText := fmt.Sprintf("gFly Framework %s", Version)

	// Calculate box width based on the version string length
	minWidth := 70 // Minimum width to accommodate the framework name line
	versionLineLen := len(fmt.Sprintf("---- _=| %s |=_ ----", headerText))
	boxWidth := minWidth
	if versionLineLen+10 > minWidth {
		boxWidth = versionLineLen + 10
	} // Add padding

	// Create horizontal borders
	topBorder := "╔" + strings.Repeat("═", boxWidth-2) + "╗"
	midBorder := "╠" + strings.Repeat("═", boxWidth-2) + "╣"
	bottomBorder := "╚" + strings.Repeat("═", boxWidth-2) + "╝"

	// Center text in the box
	centerText := func(text string) string {
		padding := boxWidth - 2 - len(text)
		leftPad := padding / 2
		rightPad := padding - leftPad
		return "║" + strings.Repeat(" ", leftPad) + text + strings.Repeat(" ", rightPad) + "║"
	}

	// Create the box with borders
	boxLines := []string{
		topBorder,
		centerText(headerText),
		centerText("Laravel inspired web framework"),
		midBorder,
		centerText(fmt.Sprintf("App Name: %s", name)),
		centerText(fmt.Sprintf("Server: %s | Environment: %s", url, env)),
		bottomBorder,
	}

	// Log the box and additional information
	for _, line := range boxLines {
		color.Blue(line)
	}
}

// ====================================================================
//                                 Logs
// ====================================================================

// setupLog configures the log output destination and log level
// based on environment variables.
func setupLog() {
	logChannel := utils.Getenv("LOG_CHANNEL", "file")

	// Log channel file
	if logChannel == "file" {
		logDir := utils.Getenv("LOG_DIR", "storage/logs")
		logFileName := utils.Getenv("LOG_FILE", "gfly.log")
		logFile := filepath.Join(logDir, logFileName)

		// Ensure log directory exists
		if err := os.MkdirAll(logDir, 0o755); err != nil {
			fmt.Printf("Error creating log directory: %v\n", err)
			// Continue with stdout only if directory creation fails
		} else {
			// Set the output destination to the console and file
			file, err := os.OpenFile(filepath.Clean(logFile), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
			if err != nil {
				fmt.Printf("Error opening log file: %v\n", err)
				// Continue with stdout only if file opening fails
			} else {
				// Create a writer that strips ANSI color codes before writing to the file
				noColorFile := &ansiStripper{writer: file}
				iw := io.MultiWriter(os.Stdout, noColorFile)
				log.SetOutput(iw)
			}
		}
	}

	// Set log level (case-insensitive)
	logLevel := strings.ToLower(utils.Getenv("LOG_LEVEL", "trace"))

	switch logLevel {
	case "trace":
		log.SetLevel(log.LevelTrace)
	case "debug":
		log.SetLevel(log.LevelDebug)
	case "info":
		log.SetLevel(log.LevelInfo)
	case "warn", "warning":
		log.SetLevel(log.LevelWarn)
	case "error":
		log.SetLevel(log.LevelError)
	case "fatal":
		log.SetLevel(log.LevelFatal)
	case "panic":
		log.SetLevel(log.LevelPanic)
	default:
		// Default to trace level if unrecognized
		log.SetLevel(log.LevelTrace)
		fmt.Printf("Unrecognized log level: %s, defaulting to Trace\n", logLevel)
	}

	log.Info("Setup Logs")
}

// ====================================================================
//                           Serve Static File
// ====================================================================

// serveFiles configures static files serving for the application.
//
// Parameters:
//   - fly (*GFly): The instance of GFly which contains the router used to serve the static files.
func serveFiles(fly *GFly) {
	// Default static file path
	rootPath := utils.Getenv("STATIC_PATH", "public")

	fly.router.ServeFiles("/{filepath:*}", rootPath)
}
