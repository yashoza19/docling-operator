package reconcilers

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/go-logr/logr"
	"github.io/opdev/docling-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
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
	foundService := &corev1.Service{}
	serviceName := doclingServ.Name + "-service"

	err := r.Get(ctx, types.NamespacedName{Name: serviceName, Namespace: doclingServ.Namespace}, foundService)
	if err != nil && errors.IsNotFound(err) {
		r.Log.Info("Creating a new Service", "Service", serviceName)
		service := r.serviceForDoclingServ(doclingServ)

		err = r.Create(ctx, service)
		if err != nil {
			r.Log.Error(err, "Failed to create new Service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
			return true, err
		}
	}

	return false, nil
}

func (r *ServiceReconciler) serviceForDoclingServ(doclingServ *v1alpha1.DoclingServ) *corev1.Service {
	labels := labelsForDocling(doclingServ.Name)
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      doclingServ.Name + "-service",
			Namespace: doclingServ.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Port:       5001,
					TargetPort: intstr.FromInt(5001),
				},
			},
		},
	}

	_ = ctrl.SetControllerReference(doclingServ, service, r.Scheme)

	return service
}
