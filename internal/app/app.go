package app

import (
	"strconv"

	"sigs.k8s.io/controller-runtime/pkg/manager"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

type (

	// Application configuration
	AppConfig struct {
		Port               string
		CertFolder         string
		MaxAllowedProjects int
	}
)

var logger = logf.Log.WithName("webhook_controller_app")

func NewWebHandler(config AppConfig, mgr manager.Manager) error {

	port, _ := strconv.Atoi(config.Port)

	logger.Info("Admission controller configuration", "info", config)

	// Create a webhook server.
	hookServer := &webhook.Server{
		Port:    port,
		CertDir: config.CertFolder,
	}
	hookServer.Register("/validate-project", &webhook.Admission{Handler: &projectValidator{AppConfig: &config}})

	logger.Info("Adding the hookServer to the manager")
	if err := mgr.Add(hookServer); err != nil {
		return err
	}

	return nil
}
