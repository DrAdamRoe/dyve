package reconciler

import (
	"errors"
	"github.com/joscha-alisch/dyve/internal/core/database"
	"github.com/joscha-alisch/dyve/internal/core/provider"
	recon "github.com/joscha-alisch/dyve/internal/reconciliation"
	"time"
)

func NewReconciler(db database.Database, m provider.Manager, olderThan time.Duration) recon.Reconciler {
	if olderThan == 0 {
		olderThan = time.Minute
	}

	r := &reconciler{
		Reconciler: recon.NewReconciler(db, olderThan),
		db:         db,
		m:          m,
	}
	r.Handler(database.ReconcileAppProvider, r.reconcileAppProvider)
	r.Handler(database.ReconcilePipelineProvider, r.reconcilePipelineProvider)

	return r
}

type reconciler struct {
	recon.Reconciler
	db database.Database
	m  provider.Manager
}

func (r *reconciler) reconcileAppProvider(j recon.Job) error {
	p, err := r.m.GetAppProvider(j.Guid)
	if errors.Is(err, provider.ErrNotFound) {
		r.db.DeleteAppProvider(j.Guid)
		return nil
	}
	if err != nil {
		return err
	}

	apps, err := p.ListApps()
	if err != nil {
		return err
	}

	r.db.UpdateApps(j.Guid, apps)

	return nil
}

func (r *reconciler) reconcilePipelineProvider(j recon.Job) error {
	p, err := r.m.GetPipelineProvider(j.Guid)
	if errors.Is(err, provider.ErrNotFound) {
		r.db.DeletePipelineProvider(j.Guid)
		return nil
	}
	if err != nil {
		return err
	}

	pipelines, err := p.ListPipelines()
	if err != nil {
		return err
	}

	updates, err := p.ListUpdates(j.LastUpdated)

	r.db.UpdatePipelines(j.Guid, pipelines)
	r.db.AddPipelineVersions(j.Guid, updates.Versions)
	r.db.AddPipelineRuns(j.Guid, updates.Runs)
	return nil
}
