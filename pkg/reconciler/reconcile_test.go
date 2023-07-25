package reconciler

import (
	"context"
	"fmt"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/red-hat-storage/managed-fusion-fleet-reconciler/pkg/db"
	"github.com/red-hat-storage/managed-fusion-fleet-reconciler/pkg/forman"
	"github.com/red-hat-storage/managed-fusion-fleet-reconciler/pkg/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestNewFleetReconciler(t *testing.T) {
	mockdb := db.MockProviderDatabase{}
	logger, err := zap.NewDevelopment()
	assert.NoError(t, err)
	awsconf, err := awsConfig.LoadDefaultConfig(context.Background())
	assert.NoError(t, err)
	templateURL := "mydummys3URL"
	actualReconciler := NewFleetReconciler(logger, &mockdb, &awsconf, templateURL)
	expectedReconciler := &FleetReconciler{logger: logger, dbClient: &mockdb}
	assert.Equal(t, expectedReconciler, actualReconciler)
}

func TestFleetReconciler_Reconcile(t *testing.T) {

	mockdb := db.MockProviderDatabase{}
	logger, err := zap.NewDevelopment()
	assert.NoError(t, err)
	awsconf, err := awsConfig.LoadDefaultConfig(context.Background())
	assert.NoError(t, err)
	templateURL := "mydummys3URL"
	reconciler := NewFleetReconciler(logger, &mockdb, &awsconf, templateURL)

	mockdb.On("GetProviderCluster", context.Background(), "cluster1").Return(&types.ProviderCluster{}, fmt.Errorf("no provider with this cluster ID"))
	mockdb.On("GetProviderCluster", context.Background(), "cluster2").Return(&types.ProviderCluster{
		ClusterID: "cluster2",
		Spec: types.ProviderClusterSpec{
			Kind: types.AWSCloudProvider,
		},
	}, nil)

	testCases := []struct {
		request forman.Request
		result  forman.Result
	}{
		{request: forman.Request{Name: "cluster1"}, result: forman.Result{}},
		{request: forman.Request{Name: "cluster2"}, result: forman.Result{Requeue: false}},
	}

	for _, testCase := range testCases {
		actualResult := reconciler.Reconcile(testCase.request)
		assert.Equal(t, testCase.result, actualResult)
	}

}
