package webhook

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// Validator specifies the interface for a validating webhook.
type Validator interface {
	// Validate yields a response to an validating AdmissionRequest.
	Validate(context.Context, admission.Request) admission.Response
}

// ValidatingWebhook is a generic validating admission webhook.
type ValidatingWebhook struct {
	Client  client.Client
	Decoder *admission.Decoder
}

// Validate implements the Validator interface.
func (v *ValidatingWebhook) Validate(_ context.Context, _ admission.Request) admission.Response {
	return admission.Allowed("")
}

// InjectDecoder implements the admission.DecoderInjector interface.
func (v *ValidatingWebhook) InjectDecoder(decoder *admission.Decoder) error {
	v.Decoder = decoder
	return nil
}

// InjectClient implements the inject.Client interface.
func (v *ValidatingWebhook) InjectClient(client client.Client) error {
	v.Client = client
	return nil
}
