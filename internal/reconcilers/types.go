package reconcilers

import (
	"context"

	"github.io/opdev/docling-operator/api/v1alpha1"
)

type Reconciler interface {
	Reconcile(ctx context.Context, doclingServ *v1alpha1.DoclingServ) (bool, error)
}
