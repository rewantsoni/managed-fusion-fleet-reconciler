package reconciler

import (
	"context"

	"github.com/red-hat-storage/managed-fusion-fleet-reconciler/pkg/db"
	"github.com/red-hat-storage/managed-fusion-fleet-reconciler/pkg/forman"

	"go.uber.org/zap"
)

func Reconcile(logger *zap.Logger, client *db.Database, req forman.Request) forman.Result {
	ctx := context.Background()
	logger.Info("Processing request", zap.String("name", req.Name))
	provider, err := client.GetProviderCluster(ctx, req.Name)
	if err != nil {
		logger.Error("Failed to get provider cluster", zap.Error(err))
		return forman.Result{}
	}
	logger.Info("Provider", zap.String("clusterId", provider.ClusterID))

	return forman.Result{}
}
