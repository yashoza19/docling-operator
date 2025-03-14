package reconcilers

import (
	"context"

	"github.com/go-logr/logr"
	"github.io/opdev/docling-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ServiceReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func NewServiceReconciler(client client.Client, log logr.Logger, scheme *runtime.Scheme) *ServiceReconciler {
	return &ServiceReconciler{
		Client: client,
		Log:    log,
		Scheme: scheme,
	}
}

func (r *ServiceReconciler) Reconcile(ctx context.Context, doclingServ *v1alpha1.DoclingServ) (bool, error) {

	return false, nil
}
