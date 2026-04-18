package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	// Version is set at build time via ldflags
	Version = "dev"
	// BuildDate is set at build time via ldflags
	BuildDate = "unknown"
)

func main() {
	// Configure pretty logging for development
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	rootCmd := &cobra.Command{
		Use:   "neko",
		Short: "neko - self-hosted virtual browser",
		Long: `neko is a self-hosted virtual browser that runs in Docker
and uses WebRTC to stream the browser to the client.`,
		Version: fmt.Sprintf("%s (built %s)", Version, BuildDate),
		RunE:    run,
	}

	// Global flags
	rootCmd.PersistentFlags().String("log-level", "debug", "log level (trace, debug, info, warn, error)")
	rootCmd.PersistentFlags().String("bind", "0.0.0.0:8080", "address to bind the HTTP server")
	rootCmd.PersistentFlags().String("static", "/var/www", "path to static files to serve")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("failed to execute command")
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	// Parse log level
	levelStr, _ := cmd.Flags().GetString("log-level")
	level, err := zerolog.ParseLevel(levelStr)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	bind, _ := cmd.Flags().GetString("bind")
	staticPath, _ := cmd.Flags().GetString("static")

	log.Info().
		Str("version", Version).
		Str("build_date", BuildDate).
		Msg("starting neko server")

	log.Info().
		Str("bind", bind).
		Str("static", staticPath).
		Msg("configuration loaded")

	// TODO: initialize and start server components
	server := NewServer(bind, staticPath)
	return server.Start()
}
