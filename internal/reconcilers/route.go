package reconcilers

import (
	"context"

	"github.com/go-logr/logr"
	"github.io/opdev/docling-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type RouteReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func NewRouteReconciler(client client.Client, log logr.Logger, scheme *runtime.Scheme) *RouteReconciler {
	return &RouteReconciler{
		Client: client,
		Log:    log,
		Scheme: scheme,
	}
}

func (r *RouteReconciler) Reconcile(ctx context.Context, doclingServ *v1alpha1.DoclingServ) (bool, error) {

	return false, nil
}
