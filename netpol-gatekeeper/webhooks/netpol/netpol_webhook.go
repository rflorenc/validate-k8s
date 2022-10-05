package netpol

import (
	"context"
	"fmt"

	"github.com/rflorenc/validate-k8s/netpol-gatekeeper/pkg/webhook"
	admissionv1 "k8s.io/api/admission/v1"
	authorizationv1 "k8s.io/api/authorization/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	netpolv1 "k8s.io/api/networking/v1"
)

// Webhook to validate netpols.
type Webhook struct {
	webhook.ValidatingWebhook
}

// SetupWithManager sets up the webhook with the Manager.
func (w *Webhook) SetupWithManager(mgr ctrl.Manager) error {
	return webhook.NewGenericWebhookManagedBy(mgr).
		For(&netpolv1.NetworkPolicy{}).
		Complete(w)
}

// Validate validates admission.Request for a  netpol object.
func (w *Webhook) Validate(ctx context.Context, req admission.Request) admission.Response {
	logger := log.FromContext(ctx).WithValues("networking.k8s.io", types.NamespacedName{Name: req.Name, Namespace: req.Namespace}.String())

	// get netpol object
	netpol := netpolv1.NetworkPolicy{}
	switch req.Operation {
	case admissionv1.Create, admissionv1.Update:
		err := w.Decoder.Decode(req, &netpol)
		if err != nil {
			logger.Error(err, "failed to decode netpol object")
			return admission.Denied("webhook failed to decode netpol object.")
		}
	case admissionv1.Delete:
		err := w.Decoder.DecodeRaw(req.OldObject, &netpol)
		if err != nil {
			logger.Error(err, "failed to decode netpol object")
			return admission.Denied("webhook failed to decode old netpol object.")
		}
	}

	// check user's access rights of namespaces read by the netpol
	sar := authorizationv1.SubjectAccessReview{
		Spec: authorizationv1.SubjectAccessReviewSpec{
			User:   req.UserInfo.Username,
			UID:    req.UserInfo.UID,
			Groups: req.UserInfo.Groups,
			ResourceAttributes: &authorizationv1.ResourceAttributes{
				Namespace: "",
				Verb:      "delete",
				Group:     "networking.k8s.io",
				Version:   "v1",
				Resource:  "networkpolicies",
			},
		},
	}

	if err := w.Client.Create(ctx, &sar); err != nil {
		logger.Error(err, "failed to create subject access review")
		return admission.Denied("webhook failed to create subject access review.")
	}

	if !sar.Status.Allowed || sar.Status.Denied {
		return admission.Denied(fmt.Sprintf("subject \"%s\" is not allowed to delete netpol, only subjects who have the permission to delete\"networkpolicies\" in that particular namespace are allowed to delete the network policy.", req.UserInfo.Username))
	}

	return admission.Allowed("")
}
