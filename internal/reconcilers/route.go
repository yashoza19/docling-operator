package reconcilers

import (
	"context"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/go-logr/logr"
	routev1 "github.com/openshift/api/route/v1"
	"github.io/opdev/docling-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	route := &routev1.Route{ObjectMeta: metav1.ObjectMeta{Name: doclingServ.Name + "-route", Namespace: doclingServ.Namespace}}
	_, err := controllerutil.CreateOrUpdate(ctx, r.Client, route, func() error {
		labels := labelsForDocling(doclingServ.Name)
		route.Labels = labels
		route.Spec = routev1.RouteSpec{
			Path: "/",
			To: routev1.RouteTargetReference{
				Kind: "Service",
				Name: doclingServ.Name + "-service",
			},
			Port: &routev1.RoutePort{
				TargetPort: intstr.FromString("http"),
			},
			TLS: &routev1.TLSConfig{
				Termination: routev1.TLSTerminationEdge,
			},
		}
		_ = ctrl.SetControllerReference(doclingServ, route, r.Scheme)
		return nil
	})

	if err != nil {
		r.Log.Error(err, "Error reconciling Route", "Route.Namespace", route.Namespace, "Route.Name", route.Name)
		return true, err
	}

	r.Log.Info("Successfully reconciled Route", "Route.Namespace", route.Namespace, "Route.Name", route.Name)
	return false, nil
}
