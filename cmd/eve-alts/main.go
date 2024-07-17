package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/kava-forge/eve-alts/lib/deferutil"
	"github.com/kava-forge/eve-alts/lib/errors"
	"github.com/kava-forge/eve-alts/lib/logging/level"

	"github.com/kava-forge/eve-alts/pkg/app"
)

// build time variables
var (
	AppName      string
	BuildDate    string
	BuildVersion string
	BuildSHA     string
)

func main() {
	home := app.GetConfigDir()
	crashlog := path.Join(home, "eve-alts-crash.log")
	defer func() {
		if err := recover(); err != nil {
			_ = os.WriteFile(crashlog, []byte(fmt.Sprintf("%v\n", err)), 0o600)
		}
	}()

	maybeEnableFyneLog(home)

	ctx := context.Background()
	err := run(ctx, os.Args[1:])
	if err != nil {
		panic(err)
	}
}

func run(ctx context.Context, args []string) error {
	var configFile string

	fs := flag.NewFlagSet(AppName, flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage of %s:\n\tVersion=%s (%s)\n\tBuildDate=%s\n", AppName, BuildVersion, BuildSHA, BuildDate)
		fs.PrintDefaults()
	}
	fs.StringVar(&configFile, "config", "", "The config file to use")
	if err := fs.Parse(args); err != nil {
		return err
	}

	conf, err := app.SetupConfig(configFile)
	if err != nil {
		return errors.Wrap(err, "could not set up config")
	}

	deps, err := app.CreateDependencies(ctx, AppName, BuildVersion, conf)
	if err != nil {
		return errors.Wrap(err, "could not create dependencies")
	}
	defer deferutil.CheckDeferLog(deps.Logger(), func() error { return deps.Close(ctx, 2*time.Second) })

	logger := deps.Logger()

	level.Debug(logger).Message("config dump", "config", fmt.Sprintf("%+v", conf))

	coreApp, err := app.New(AppName, deps, conf)
	if err != nil {
		return errors.Wrap(err, "could not construct core app")
	}

	if err := coreApp.Start(); err != nil {
		level.Error(logger).Err("program exited", err)
	}
	return nil
}
