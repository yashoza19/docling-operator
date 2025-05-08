package reconcilers

import (
	"context"

	"github.io/opdev/docling-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type ServiceAccountReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func NewServiceAccountReconciler(client client.Client, scheme *runtime.Scheme) *ServiceAccountReconciler {
	return &ServiceAccountReconciler{
		Client: client,
		Scheme: scheme,
	}
}

func (r *ServiceAccountReconciler) Reconcile(ctx context.Context, doclingServe *v1alpha1.DoclingServe) (bool, error) {
	log := logf.FromContext(ctx)
	serviceaccount := &corev1.ServiceAccount{ObjectMeta: metav1.ObjectMeta{Name: doclingServe.Name + "-serviceaccount", Namespace: doclingServe.Namespace}}
	_, err := controllerutil.CreateOrUpdate(ctx, r.Client, serviceaccount, func() error {
		labels := labelsForDocling(doclingServe.Name)
		serviceaccount.Labels = labels
		_ = ctrl.SetControllerReference(doclingServe, serviceaccount, r.Scheme)
		return nil
	})
	if err != nil {
		log.Error(err, "Error creating Service Account", "ServiceAccount.Namespace", serviceaccount.Namespace, "ServiceAccount.Name", serviceaccount.Name)
		return true, err
	}

	log.Info("Successfully created ServiceAccount", "ServiceAccount.Namespace", serviceaccount.Namespace, "ServiceAccount.Name", serviceaccount.Name)
	return false, nil
}
