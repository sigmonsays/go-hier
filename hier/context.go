package main

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

type Context struct {
	DataDir string
}

func (me *Context) Close() {
}

func NewContext(cmd *cobra.Command, args []string) *Context {
	datadir, _ := cmd.Flags().GetString("data-dir")
	loglevel, _ := cmd.Flags().GetString("loglevel")
	slog.Debug("new context", "datadir", datadir)
	c := &Context{
		DataDir: datadir,
	}

	if loglevel == "" {
		loglevel = "info"
	}

	// setup logging
	slogOpts := &slog.HandlerOptions{
		// AddSource: true,
	}
	// slogOpts.Level = slog.LevelDebug
	SetLevel(slogOpts, loglevel)
	logger := slog.New(slog.NewJSONHandler(os.Stderr, slogOpts))
	// With("app", "rollout")
	slog.SetDefault(logger)
	slog.SetLogLoggerLevel(slogOpts.Level.Level())

	return c
}

func SetLevel(opts *slog.HandlerOptions, level_name string) error {

	// parse log level from string into proper value
	var level slog.Level
	err := level.UnmarshalText([]byte(level_name))
	if err != nil {
		return err
	}
	opts.Level = level

	return nil
}
