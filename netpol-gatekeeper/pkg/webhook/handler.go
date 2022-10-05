package webhook

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// handler is wrapper type for Validator and Mutator.
type handler struct {
	Handler interface{}
}

// Handle implements the admission.Handler interface.
func (h *handler) Handle(ctx context.Context, req admission.Request) admission.Response {
	if validator, ok := h.Handler.(Validator); ok {
		return validator.Validate(ctx, req)
	}

	if mutator, ok := h.Handler.(Mutator); ok {
		return mutator.Mutate(ctx, req)
	}

	return admission.Denied("")
}

// InjectDecoder implements the admission.DecoderInjector interface.
func (h *handler) InjectDecoder(decoder *admission.Decoder) error {
	if injector, ok := h.Handler.(admission.DecoderInjector); ok {
		return injector.InjectDecoder(decoder)
	}

	return nil
}

// InjectClient implements the inject.Client interface.
func (h *handler) InjectClient(client client.Client) error {
	if injector, ok := h.Handler.(inject.Client); ok {
		return injector.InjectClient(client)
	}

	return nil
}
