package core

import (
	"fmt"
	"github.com/gflydev/core/log"
	"github.com/gflydev/core/utils"
	"io"
	"os"
	"path/filepath"
)

// ===========================================================================================================
//                                                Bootstrap
// ===========================================================================================================

// startupMessage Startup message.
func startupMessage(url, name, env string) {
	log.Info("-------------------------------------------")
	log.Info(fmt.Sprintf("       ---- _=| gFly %s |=_ ----       ", Version))
	log.Info("       Laravel inspired web framework      ")
	log.Info("-------------------------------------------")
	log.Infof("   * Server: %s", url)
	log.Infof("   * App Name: %s", name)
	log.Infof("   * Environment: %s", env)
}

// ===========================================================================================================
//                                                   Logs
// ===========================================================================================================

// serveFiles Serve static files from the given file system root is `./public`
// You can change parameter name STATIC_PATH.
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

// ===========================================================================================================
//                                             Serve Static File
// ===========================================================================================================

// serveFiles Serve static files from the given file system root is `./public`
// You can change parameter name STATIC_PATH.
func serveFiles(fly *GFly) {
	// Default static file path
	rootPath := utils.Getenv("STATIC_PATH", "public")

	fly.router.ServeFiles("/{filepath:*}", rootPath)
}
