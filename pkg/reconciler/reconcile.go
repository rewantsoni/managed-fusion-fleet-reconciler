package reconciler

import (
	"context"
	"fmt"
	cloudaws "github.com/red-hat-storage/managed-fusion-fleet-reconciler/pkg/cloud/aws"
	"github.com/red-hat-storage/managed-fusion-fleet-reconciler/pkg/db"
	"github.com/red-hat-storage/managed-fusion-fleet-reconciler/pkg/forman"
	"github.com/red-hat-storage/managed-fusion-fleet-reconciler/pkg/types"

	"go.uber.org/zap"
)

type AWS struct {
	TransitGatewayID                       string
	ControlPlaneTransitGatewayRouteTableID string
	ProviderTransitGatewayRouteTableID     string
	ControlPlanePublicRouteTableID         string
	ControlPlanePrivateRouteTableID0       string
	ControlPlanePrivateRouteTableID1       string
	ControlPlanePrivateRouteTableID2       string
	IgnitionScriptS3URL                    string
}

type FleetReconciler struct {
	ctx         context.Context
	logger      *zap.Logger
	dbClient    *db.Database
	awsProvider *cloudaws.AWSCloudProvider
}

func (r *FleetReconciler) Reconcile(ctx context.Context, logger *zap.Logger, dbClient *db.Database, provider *types.ProviderCluster) (forman.Result, error) {
	r.initReconciler(ctx, logger, dbClient)
	result, err := r.reconcilePhases(provider)
	if err != nil {
		return forman.Result{}, err
	}
	return result, nil
}

func (r *FleetReconciler) initReconciler(ctx context.Context, logger *zap.Logger, dbClient *db.Database) {
	r.ctx = ctx
	r.logger = logger
	r.dbClient = dbClient
}

func (r *FleetReconciler) reconcilePhases(provider *types.ProviderCluster) (forman.Result, error) {
	if result, err := r.reconcileCloudProvider(provider); err != nil {
		return result, err
	}
}

func (r *FleetReconciler) reconcileCloudProvider(provider *types.ProviderCluster) (forman.Result, error) {
	switch provider.Spec.Kind {
	case types.AWSCloudProvider:
		return r.awsProvider.Reconcile(r.ctx, provider)
	default:
		r.logger.Error("Failed to reconcile Cloud Provider Creation", zap.Error(fmt.Errorf("unsupported Cloud Provider")))
		return forman.Result{Requeue: false}, fmt.Errorf("unsupported Cloud Provider")
	}
}

func Reconcile(logger *zap.Logger, client *db.Database, req forman.Request) forman.Result {
	ctx := context.Background()
	logger.Info("Processing request", zap.String("name", req.Name))
	provider, err := client.GetProviderCluster(ctx, req.Name)
	if err != nil {
		logger.Error("Failed to get provider cluster", zap.Error(err))
		return forman.Result{}
	}
	logger.Info("Provider", zap.String("clusterId", provider.ClusterID))

	r := &FleetReconciler{}
	result, err := r.Reconcile(ctx, logger, client, provider)
	if err != nil {
		logger.Error("Failed to reconcile provider cluster", zap.Error(err))
		// decide a resonable requeue time based on error
		return forman.Result{}
	}

	return result
}
