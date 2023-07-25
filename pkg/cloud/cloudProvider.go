package cloud

import (
	"context"
	"fmt"
	"github.com/red-hat-storage/managed-fusion-fleet-reconciler/pkg/cloud/aws"
	"github.com/red-hat-storage/managed-fusion-fleet-reconciler/pkg/forman"
	"github.com/red-hat-storage/managed-fusion-fleet-reconciler/pkg/types"
	"go.uber.org/zap"
)

type CloudProvider interface {
	Reconcile(ctx context.Context, provider *types.ProviderCluster) (forman.Result, error)
}

func NewCloudProvider(provider types.CloudProvider, logger *zap.Logger, config interface{}) (CloudProvider, error) {
	switch provider {
	case types.AWSCloudProvider:
		return aws.NewAWSCloudProvider(config)
	default:
		logger.Error(fmt.Sprintf("Failed to create Cloud Provider for %s", provider), zap.Error(fmt.Errorf("unsupported Cloud Provider")))
		return nil, fmt.Errorf("unsupported Cloud Provider")
	}
}
