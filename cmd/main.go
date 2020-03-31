package main

import (
	"os"

	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"

	"github.com/azopat/project-admission-validation/internal/app"
	"github.com/urfave/cli/v2"
)

const (
	FLAG_PORT                = "port"
	FLAG_CERT_FOLDER         = "cert_folder"
	FLAG_MAX_ALLOWED_PROJECT = "max_allowed_projects"
)

var logger = logf.Log.WithName("ctrl_cmd")

func main() {

	logf.SetLogger(zap.Logger(false))
	entryLog := logger.WithName("main")

	cliApp := &cli.App{}
	cliApp.Action = runIt
	cliApp.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    FLAG_MAX_ALLOWED_PROJECT,
			Usage:   "Total amount of projects a user can create",
			EnvVars: []string{"MAX_ALLOWED_PROJECT"},
		},
		&cli.StringFlag{
			Name:    FLAG_PORT,
			Usage:   "service port, default is 443",
			Value:   "8443",
			EnvVars: []string{"PORT"},
		},
		&cli.StringFlag{
			Name:    FLAG_CERT_FOLDER,
			Usage:   "folder container ssl certificate for the admission controller",
			EnvVars: []string{"CERT_FOLDER"},
		},
	}

	err := cliApp.Run(os.Args)
	if err != nil {
		entryLog.Error(err, "Could not start execution")
	}

}

func runIt(c *cli.Context) error {

	// Logger initiation
	entryLog := logger.WithName("runIt")

	appConfig := app.AppConfig{}

	appConfig.Port = c.String(FLAG_PORT)
	appConfig.CertFolder = c.String(FLAG_CERT_FOLDER)
	appConfig.MaxAllowedProjects = c.Int(FLAG_MAX_ALLOWED_PROJECT)

	// Setup a Manager
	entryLog.Info("setting up manager")
	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{})
	if err != nil {
		entryLog.Error(err, "unable to set up overall controller manager")
		return err
	}

	entryLog.Info("Instantiating the admission controller webhook")
	if err := app.NewWebHandler(appConfig, mgr); err != nil {
		entryLog.Error(err, "failed to initiate webhook")
		return err
	}

	entryLog.Info("starting manager")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		entryLog.Error(err, "unable to run manager")
		return err
	}

	return nil

}
