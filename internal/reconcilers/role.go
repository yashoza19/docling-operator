package reconcilers

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	rbacv1 "k8s.io/api/rbac/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type RoleReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func        NewRoleReconciler(client client.Client, scheme *runtime.Scheme) *RoleReconciler {
	return &RoleReconciler{
		Client: client,
		Scheme: scheme,
	}
}

func (r *RoleReconciler) Reconcile(ctx context.Context, doclingServe *v1alpha1.DoclingServe) (bool, error) {
	log := logf.FromContext(ctx)
	role := &rbacv1.Role{ObjectMeta: metav1.ObjectMeta{Name: doclingServe.Name + "-role", Namespace: doclingServe.Namespace,}}
	_, err := controllerutil.CreateOrUpdate(ctx, r.Client, role, func() error {
		labels := labelsForDocling(doclingServe.Name)
		role.Labels = labels
		role.Rules := []rbacv1.PolicyRule(
			APIGroups: []string{"security.openshift.io"}, 
			Resources: []string{"securitycontextcontraints"},
			ResourceName: []string{"restricted-v2"},
			Verbs: []string{"use"},
		)
		_ = ctrl.SetControllerReference(doclingServe, role, r.Scheme)
		return nil
	})
	if err != nil {
		log.Error(err, "Error creating Role", "Role.Namespace", role.Namespace, "Role.Name", role.Name)
		return true, err
	}

	log.Info("Successfully created Role", "Role.Namespace", role.Namespace, "Role.Name", Role.Name)
	return false, nil
}