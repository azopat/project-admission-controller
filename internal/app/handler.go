package app

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// projectValidator validates namespace creation
type projectValidator struct {
	client    client.Client
	decoder   *admission.Decoder
	AppConfig *AppConfig
}

// This allows project creation for users who haven't yet reached the maximum allowed project count
func (v *projectValidator) Handle(ctx context.Context, req admission.Request) admission.Response {

	actualProjectCount := 0

	logger.Info("New admission request", "info", map[string]interface{}{
		"requester": ":" + req.UserInfo.Username,
		"UID":       req.UID,
	})

	// Getting the amount of projects owned by the requester
	nsList := &corev1.NamespaceList{}
	_ = v.client.List(context.Background(), nsList)
	for _, ns := range nsList.Items {

		if req.UserInfo.Username == ns.ObjectMeta.Annotations["openshift.io/requester"] {
			actualProjectCount = actualProjectCount + 1
		}

	}

	// Rejecting if user has reached the limit
	if actualProjectCount >= v.AppConfig.MaxAllowedProjects {
		logger.Info("Admission request rejected, User actual projects count have reached or exceeds the allowed maximum", "details", map[string]interface{}{
			"requester": ":" + req.UserInfo.Username,
			"userCount": actualProjectCount,
			"maximum":   v.AppConfig.MaxAllowedProjects,
		})
		return admission.Denied(fmt.Sprintf("User actual projects count have reached or exceeds the allowed maximum of %d", v.AppConfig.MaxAllowedProjects))
	}

	return admission.Allowed("")

}

// InjectClient injects the client.
func (v *projectValidator) InjectClient(c client.Client) error {
	v.client = c
	return nil
}

// InjectDecoder injects the decoder.
func (v *projectValidator) InjectDecoder(d *admission.Decoder) error {
	v.decoder = d
	return nil
}
