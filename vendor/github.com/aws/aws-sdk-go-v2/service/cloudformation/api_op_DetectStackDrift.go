// Code generated by smithy-go-codegen DO NOT EDIT.

package cloudformation

import (
	"context"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Detects whether a stack's actual configuration differs, or has drifted, from
// its expected configuration, as defined in the stack template and any values
// specified as template parameters. For each resource in the stack that supports
// drift detection, CloudFormation compares the actual configuration of the
// resource with its expected template configuration. Only resource properties
// explicitly defined in the stack template are checked for drift. A stack is
// considered to have drifted if one or more of its resources differ from their
// expected template configurations. For more information, see Detecting
// Unregulated Configuration Changes to Stacks and Resources (https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-stack-drift.html)
// . Use DetectStackDrift to detect drift on all supported resources for a given
// stack, or DetectStackResourceDrift to detect drift on individual resources. For
// a list of stack resources that currently support drift detection, see Resources
// that Support Drift Detection (https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-stack-drift-resource-list.html)
// . DetectStackDrift can take up to several minutes, depending on the number of
// resources contained within the stack. Use DescribeStackDriftDetectionStatus to
// monitor the progress of a detect stack drift operation. Once the drift detection
// operation has completed, use DescribeStackResourceDrifts to return drift
// information about the stack and its resources. When detecting drift on a stack,
// CloudFormation doesn't detect drift on any nested stacks belonging to that
// stack. Perform DetectStackDrift directly on the nested stack itself.
func (c *Client) DetectStackDrift(ctx context.Context, params *DetectStackDriftInput, optFns ...func(*Options)) (*DetectStackDriftOutput, error) {
	if params == nil {
		params = &DetectStackDriftInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "DetectStackDrift", params, optFns, c.addOperationDetectStackDriftMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*DetectStackDriftOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type DetectStackDriftInput struct {

	// The name of the stack for which you want to detect drift.
	//
	// This member is required.
	StackName *string

	// The logical names of any resources you want to use as filters.
	LogicalResourceIds []string

	noSmithyDocumentSerde
}

type DetectStackDriftOutput struct {

	// The ID of the drift detection results of this operation. CloudFormation
	// generates new results, with a new drift detection ID, each time this operation
	// is run. However, the number of drift results CloudFormation retains for any
	// given stack, and for how long, may vary.
	//
	// This member is required.
	StackDriftDetectionId *string

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationDetectStackDriftMiddlewares(stack *middleware.Stack, options Options) (err error) {
	err = stack.Serialize.Add(&awsAwsquery_serializeOpDetectStackDrift{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsquery_deserializeOpDetectStackDrift{}, middleware.After)
	if err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddClientRequestIDMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddComputeContentLengthMiddleware(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = v4.AddComputePayloadSHA256Middleware(stack); err != nil {
		return err
	}
	if err = addRetryMiddlewares(stack, options); err != nil {
		return err
	}
	if err = addHTTPSignerV4Middleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = awsmiddleware.AddRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addClientUserAgent(stack, options); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addOpDetectStackDriftValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opDetectStackDrift(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = awsmiddleware.AddRecursionDetection(stack); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	return nil
}

func newServiceMetadataMiddleware_opDetectStackDrift(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		SigningName:   "cloudformation",
		OperationName: "DetectStackDrift",
	}
}
