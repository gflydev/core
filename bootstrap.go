package core

import (
	"fmt"
	"github.com/gflydev/core/log"
	"github.com/gflydev/core/utils"
	"io"
	"os"
	"path/filepath"
)

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
	log.Info("-------------------------------------------")
	log.Info(fmt.Sprintf("	   ---- _=| gFly %s |=_ ----	   ", Version))
	log.Info("	   Laravel inspired web framework	  ")
	log.Info("-------------------------------------------")
	log.Infof("   * Server: %s", url)
	log.Infof("   * App Name: %s", name)
	log.Infof("   * Environment: %s", env)
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
		logFile := fmt.Sprintf("storage/logs/%s", utils.Getenv("LOG_FILE", "gfly.log"))

		// Set the output destination to the console and file.
		file, _ := os.OpenFile(filepath.Clean(logFile), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
		iw := io.MultiWriter(os.Stdout, file)
		log.SetOutput(iw)
	}

	// Set log level
	switch utils.Getenv("LOG_LEVEL", "Trace") {
	case "Trace":
		log.SetLevel(log.LevelTrace)
	case "Debug":
		log.SetLevel(log.LevelDebug)
	case "Info":
		log.SetLevel(log.LevelInfo)
	case "Warn":
		log.SetLevel(log.LevelWarn)
	case "Error":
		log.SetLevel(log.LevelError)
	case "Fatal":
		log.SetLevel(log.LevelFatal)
	case "Panic":
		log.SetLevel(log.LevelPanic)
	}

	log.Trace("Setup Logs")
}

// ====================================================================
//                           Serve Static File
// ====================================================================

// serveFiles configures static file serving for the application.
//
// Parameters:
//   - fly (*GFly): The instance of GFly which contains the router used to serve the static files.
func serveFiles(fly *GFly) {
	// Default static file path
	rootPath := utils.Getenv("STATIC_PATH", "public")

	fly.router.ServeFiles("/{filepath:*}", rootPath)
}
