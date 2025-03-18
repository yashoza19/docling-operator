package reconcilers

import (
	"context"

	"github.com/go-logr/logr"
	routev1 "github.com/openshift/api/route/v1"
	"github.io/opdev/docling-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
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
	foundRoute := &routev1.Route{}
	routeName := doclingServ.Name + "-route"

	err := r.Get(ctx, types.NamespacedName{Name: routeName, Namespace: doclingServ.Namespace}, foundRoute)
	if err != nil && errors.IsNotFound(err) {
		r.Log.Info("Creating a new Route", "Route", routeName)
		route := r.routeForDoclingServ(doclingServ)

		err = r.Create(ctx, route)
		if err != nil {
			r.Log.Error(err, "Failed to create new Route", "Route.Namespace", route.Namespace, "Route.Name", route.Name)
			return true, err
		}
	}

	return false, nil
}

func (r *RouteReconciler) routeForDoclingServ(doclingServ *v1alpha1.DoclingServ) *routev1.Route {
	labels := labelsForDocling(doclingServ.Name)
	route := &routev1.Route{
		ObjectMeta: metav1.ObjectMeta{
			Name:      doclingServ.Name + "-route",
			Namespace: doclingServ.Namespace,
			Labels:    labels,
		},
		Spec: routev1.RouteSpec{
			Path: "/", // Pass all traffic
			To: routev1.RouteTargetReference{
				Kind: "Service",
				Name: doclingServ.Name + "-service",
			},
			Port: &routev1.RoutePort{
				TargetPort: intstr.FromString("http"),
			},
			TLS: &routev1.TLSConfig{ // Optional: Enable TLS
				Termination: routev1.TLSTerminationEdge,
			},
		},
	}

	_ = ctrl.SetControllerReference(doclingServ, route, r.Scheme)

	return route
}
