package reconcilers

import (
	"context"

	"github.io/opdev/docling-operator/api/v1alpha1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type RoleBindingReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func NewRoleBindingReconciler(client client.Client, scheme *runtime.Scheme) *RoleBindingReconciler {
	return &RoleBindingReconciler{
		Client: client,
		Scheme: scheme,
	}
}

func (r *RoleBindingReconciler) Reconcile(ctx context.Context, doclingServe *v1alpha1.DoclingServe) (bool, error) {
	log := logf.FromContext(ctx)
	rolebinding := &rbacv1.RoleBinding{ObjectMeta: metav1.ObjectMeta{Name: doclingServe.Name + "-roleBinding", Namespace: doclingServe.Namespace}}
	_, err := controllerutil.CreateOrUpdate(ctx, r.Client, rolebinding, func() error {
		labels := labelsForDocling(doclingServe.Name)
		rolebinding.Labels = labels
		rolebinding.RoleRef = rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     doclingServe.Name + "-role",
		}
		rolebinding.Subjects = []rbacv1.Subject{{
			Kind:      "ServiceAccount",
			Name:      doclingServe.Name + "-serviceaccount",
			Namespace: doclingServe.Namespace,
		}}

		_ = ctrl.SetControllerReference(doclingServe, rolebinding, r.Scheme)
		return nil
	})
	if err != nil {
		log.Error(err, "Error creating RoleBinding", "RoleBinding.Namespace", rolebinding.Namespace, "RoleBinding.Name", rolebinding.Name)
		return true, err
	}

	log.Info("Successfully created RoleBinding", "RoleBinding.Namespace", rolebinding.Namespace, "RoleBinding.Name", rolebinding.Name)
	return false, nil
}
