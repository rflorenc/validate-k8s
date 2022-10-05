package webhook

import (
	"errors"
	"net/http"
	"net/url"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// Builder builds a Webhook.
type Builder struct {
	mgr     manager.Manager
	apiType runtime.Object
}

// NewGenericWebhookManagedBy returns a new webhook builder that will be started by the provided Manager.
func NewGenericWebhookManagedBy(mgr manager.Manager) *Builder {
	return &Builder{mgr: mgr}
}

// For takes a runtime.Object which should be a CR.
func (blder *Builder) For(apiType runtime.Object) *Builder {
	blder.apiType = apiType
	return blder
}

// Complete builds the webhook.
// If the given object implements the Mutator interface, a MutatingWebhook will be created.
// If the given object implements the Validator interface, a ValidatingWebhook will be created.
func (blder *Builder) Complete(i interface{}) error {
	w := &admission.Webhook{
		Handler:         &handler{Handler: i},
		WithContextFunc: nil,
	}

	if err := w.InjectScheme(blder.mgr.GetScheme()); err != nil {
		return err
	}

	if err := w.InjectFunc(func(i interface{}) error {
		if injector, ok := i.(inject.Client); ok {
			return injector.InjectClient(blder.mgr.GetClient())
		}

		return nil
	}); err != nil {
		return err
	}

	if _, ok := i.(Validator); ok {
		return blder.registerValidatingWebhook(w)
	}
	if _, ok := i.(Mutator); ok {
		return blder.registerMutatingWebhook(w)
	}

	return errors.New("")
}

func (blder *Builder) registerValidatingWebhook(w *admission.Webhook) error {
	// see helm validation.yaml
	path := generatePath("/validate-netpol")
	if !isAlreadyHandled(blder.mgr, path) {
		blder.mgr.GetWebhookServer().Register(path, w)
	}

	return nil
}

func (blder *Builder) registerMutatingWebhook(w *admission.Webhook) error {

	path := generatePath("/mutate-netpol")
	if !isAlreadyHandled(blder.mgr, path) {
		blder.mgr.GetWebhookServer().Register(path, w)
	}

	return nil
}

func isAlreadyHandled(mgr ctrl.Manager, path string) bool {
	if mgr.GetWebhookServer().WebhookMux == nil {
		return false
	}
	h, p := mgr.GetWebhookServer().WebhookMux.Handler(&http.Request{URL: &url.URL{Path: path}})
	if p == path && h != nil {
		return true
	}
	return false
}

func generatePath(path string) string {
	return path
}
