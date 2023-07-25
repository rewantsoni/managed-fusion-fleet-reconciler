package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cloudformationType "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"

	"github.com/red-hat-storage/managed-fusion-fleet-reconciler/pkg/forman"
	"github.com/red-hat-storage/managed-fusion-fleet-reconciler/pkg/types"
)

const (
	TransitGatewayID                       string = "TransitGatewayID"
	ControlPlaneTransitGatewayRouteTableID string = "ControlPlaneTransitGatewayRouteTableID"
	ProviderTransitGatewayRouteTableID     string = "ProviderTransitGatewayRouteTableID"
	ControlPlanePublicRouteTableID         string = "ControlPlanePublicRouteTableID"
	ControlPlanePrivateRouteTableID0       string = "ControlPlanePrivateRouteTableID0"
	ControlPlanePrivateRouteTableID1       string = "ControlPlanePrivateRouteTableID1"
	ControlPlanePrivateRouteTableID2       string = "ControlPlanePrivateRouteTableID2"
	IgnitionScriptS3URL                    string = "IgnitionScriptS3URL"
	PrimaryCIDR                            string = "PrimaryCIDR"
	SecondaryCIDR                          string = "SecondaryCIDR"
	ImageID                                string = "ImageID"
	InstanceType                           string = "InstanceType"
	InstanceCount                          string = "InstanceCount"
)

type awsconfig struct {
	config              aws.Config
	providerTemplateURL string
}

type AWSCloudProvider struct {
	config               aws.Config
	cloudFormationClient *cloudformation.Client
	providerTemplateURL  string
}

func NewAWSCloudProvider(config interface{}) (*AWSCloudProvider, error) {
	conf := config.(awsconfig)
	return &AWSCloudProvider{
		config:               conf.config,
		providerTemplateURL:  conf.providerTemplateURL,
		cloudFormationClient: cloudformation.NewFromConfig(conf.config),
	}, nil
}

func (cloudProvider *AWSCloudProvider) Reconcile(ctx context.Context, provider *types.ProviderCluster) (forman.Result, error) {
	if result, err := reconcilePhases(cloudProvider, ctx, provider); err != nil {
		return result, err
	}

	return forman.Result{Requeue: false}, nil
}

func reconcilePhases(cloudProvider *AWSCloudProvider, ctx context.Context, provider *types.ProviderCluster) (forman.Result, error) {

	if result, err := reconcileAWSCloudFormation(cloudProvider, ctx, provider); err != nil {
		return result, err
	}

	return forman.Result{}, nil
}

func reconcileAWSCloudFormation(cloudProvider *AWSCloudProvider, ctx context.Context, provider *types.ProviderCluster) (forman.Result, error) {
	// 1. If stackID is empty, create the stack and requeue
	// 	2. If not empty, check if the stack is ready, if no requeue
	//		3 if ready verify the last known hash with current hash for any changes
	//			3.1 if no change return
	//			3.2 if change then update the stack

	if provider.MetaData.StackID == "" {
		if result, err := createOrUpdateStack(cloudProvider, ctx, provider, cloudformationType.ChangeSetTypeCreate); err != nil {
			return result, err
		}
	}

	stackStatus, err := stackStatus(cloudProvider, ctx, provider.MetaData.StackID)
	if err != nil {
		return forman.Result{}, err
	}
	if (stackStatus != cloudformationType.StackStatusCreateComplete) && (stackStatus != cloudformationType.StackStatusUpdateComplete) {
		return forman.Result{Requeue: true, After: time.Minute * 10}, nil
	}

	// TODO: 3, how do we know if parameter has changed?
	return forman.Result{}, nil
}

func createOrUpdateStack(cloudProvider *AWSCloudProvider, ctx context.Context, provider *types.ProviderCluster, changeSetType cloudformationType.ChangeSetType) (forman.Result, error) {

	if provider.MetaData.ChangeSetID == "" {

		controlPlaneStackOutput, err := getStackOutput(cloudProvider, provider.MetaData.ControlPlaneStackID)
		if err != nil {
			return forman.Result{}, err
		}
		createChangeSetInput := &cloudformation.CreateChangeSetInput{
			StackName:   aws.String(provider.ClusterID),
			TemplateURL: aws.String(cloudProvider.providerTemplateURL),
			Parameters: []cloudformationType.Parameter{
				{
					ParameterKey:   aws.String(TransitGatewayID),
					ParameterValue: aws.String(controlPlaneStackOutput[TransitGatewayID]),
				},
				{
					ParameterKey:   aws.String(ControlPlanePublicRouteTableID),
					ParameterValue: aws.String(controlPlaneStackOutput[ControlPlanePublicRouteTableID]),
				},
				{
					ParameterKey:   aws.String(ControlPlanePrivateRouteTableID0),
					ParameterValue: aws.String(controlPlaneStackOutput[ControlPlanePrivateRouteTableID0]),
				},
				{
					ParameterKey:   aws.String(ControlPlanePrivateRouteTableID1),
					ParameterValue: aws.String(controlPlaneStackOutput[ControlPlanePrivateRouteTableID1]),
				},
				{
					ParameterKey:   aws.String(ControlPlanePrivateRouteTableID2),
					ParameterValue: aws.String(controlPlaneStackOutput[ControlPlanePrivateRouteTableID2]),
				},
				{
					ParameterKey:   aws.String(ProviderTransitGatewayRouteTableID),
					ParameterValue: aws.String(controlPlaneStackOutput[ProviderTransitGatewayRouteTableID]),
				},
				{
					ParameterKey:   aws.String(ControlPlaneTransitGatewayRouteTableID),
					ParameterValue: aws.String(controlPlaneStackOutput[ControlPlaneTransitGatewayRouteTableID]),
				},
			},
			Capabilities: []cloudformationType.Capability{
				cloudformationType.CapabilityCapabilityIam,
			},
			ChangeSetName: aws.String(fmt.Sprintf("%s-%d", provider.ClusterID, time.Now().Unix())),
			ChangeSetType: changeSetType,
		}

		createChangeSetOutput, err := cloudProvider.cloudFormationClient.CreateChangeSet(ctx, createChangeSetInput)
		if err != nil {
			fmt.Println("Error creating change set:", err)
			return forman.Result{}, err
		}
		provider.MetaData.ChangeSetID = *createChangeSetOutput.Id
		provider.MetaData.StackID = *createChangeSetOutput.StackId
	}

	describeChangeSetInput := &cloudformation.DescribeChangeSetInput{
		ChangeSetName: aws.String(provider.MetaData.ChangeSetID),
	}

	describeChangeSetOutput, err := cloudProvider.cloudFormationClient.DescribeChangeSet(ctx, describeChangeSetInput)
	if err != nil {
		return forman.Result{}, err
	}

	if describeChangeSetOutput.ExecutionStatus != cloudformationType.ExecutionStatusAvailable {
		return forman.Result{Requeue: true, After: time.Minute * 10}, fmt.Errorf(string(describeChangeSetOutput.ExecutionStatus))
	}

	_, err = cloudProvider.cloudFormationClient.ExecuteChangeSet(ctx,
		&cloudformation.ExecuteChangeSetInput{ChangeSetName: aws.String(provider.MetaData.ChangeSetID)})
	if err != nil {
		return forman.Result{}, err
	}
	return forman.Result{Requeue: false}, nil
}

func getStackOutput(cloudProvider *AWSCloudProvider, stackName string) (map[string]string, error) {

	// Set up input parameters for describing the stack
	describeParams := &cloudformation.DescribeStacksInput{
		StackName: aws.String(stackName),
	}

	// Describe the stack to get its outputs
	describeResp, err := cloudProvider.cloudFormationClient.DescribeStacks(context.TODO(), describeParams)
	if err != nil {
		return nil, fmt.Errorf("error describing stack: %w", err)
	}

	// Check if the stack exists
	if len(describeResp.Stacks) == 0 {
		return nil, fmt.Errorf("stack not found: %s", stackName)
	}

	// Extract the outputs from the describe response
	outputs := make(map[string]string)
	for _, output := range describeResp.Stacks[0].Outputs {
		outputs[*output.OutputKey] = *output.OutputValue
	}

	return outputs, nil
}

func stackStatus(cloudProvider *AWSCloudProvider, ctx context.Context, stackID string) (cloudformationType.StackStatus, error) {
	describeStackInput := &cloudformation.DescribeStacksInput{
		StackName: aws.String(stackID),
	}
	describeStackOutput, err := cloudProvider.cloudFormationClient.DescribeStacks(ctx, describeStackInput)
	if err != nil {
		fmt.Println("Error describing CloudFormation stack:", err)
		return "", err
	}
	return describeStackOutput.Stacks[0].StackStatus, nil
}
