package webhook

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// Mutator specifies the interface for a mutating webhook.
type Mutator interface {
	// Mutate yields a response to an mutating AdmissionRequest.
	Mutate(context.Context, admission.Request) admission.Response
}

// MutatingWebhook is a generic mutating admission webhook.
type MutatingWebhook struct {
	Client  client.Client
	Decoder *admission.Decoder
}

// Mutate implements the Mutator interface.
func (v *MutatingWebhook) Mutate(_ context.Context, _ admission.Request) admission.Response {
	return admission.Allowed("")
}

// InjectDecoder implements the admission.DecoderInjector interface.
func (v *MutatingWebhook) InjectDecoder(decoder *admission.Decoder) error {
	v.Decoder = decoder
	return nil
}

// InjectClient implements the inject.Client interface.
func (v *MutatingWebhook) InjectClient(client client.Client) error {
	v.Client = client
	return nil
}
