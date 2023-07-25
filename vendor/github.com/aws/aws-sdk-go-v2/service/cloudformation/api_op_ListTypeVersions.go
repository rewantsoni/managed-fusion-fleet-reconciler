// Code generated by smithy-go-codegen DO NOT EDIT.

package cloudformation

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Returns summary information about the versions of an extension.
func (c *Client) ListTypeVersions(ctx context.Context, params *ListTypeVersionsInput, optFns ...func(*Options)) (*ListTypeVersionsOutput, error) {
	if params == nil {
		params = &ListTypeVersionsInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "ListTypeVersions", params, optFns, c.addOperationListTypeVersionsMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*ListTypeVersionsOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type ListTypeVersionsInput struct {

	// The Amazon Resource Name (ARN) of the extension for which you want version
	// summary information. Conditional: You must specify either TypeName and Type , or
	// Arn .
	Arn *string

	// The deprecation status of the extension versions that you want to get summary
	// information about. Valid values include:
	//   - LIVE : The extension version is registered and can be used in CloudFormation
	//   operations, dependent on its provisioning behavior and visibility scope.
	//   - DEPRECATED : The extension version has been deregistered and can no longer
	//   be used in CloudFormation operations.
	// The default is LIVE .
	DeprecatedStatus types.DeprecatedStatus

	// The maximum number of results to be returned with a single call. If the number
	// of available results exceeds this maximum, the response includes a NextToken
	// value that you can assign to the NextToken request parameter to get the next
	// set of results.
	MaxResults *int32

	// If the previous paginated request didn't return all of the remaining results,
	// the response object's NextToken parameter value is set to a token. To retrieve
	// the next set of results, call this action again and assign that token to the
	// request object's NextToken parameter. If there are no remaining results, the
	// previous response object's NextToken parameter is set to null .
	NextToken *string

	// The publisher ID of the extension publisher. Extensions published by Amazon
	// aren't assigned a publisher ID.
	PublisherId *string

	// The kind of the extension. Conditional: You must specify either TypeName and
	// Type , or Arn .
	Type types.RegistryType

	// The name of the extension for which you want version summary information.
	// Conditional: You must specify either TypeName and Type , or Arn .
	TypeName *string

	noSmithyDocumentSerde
}

type ListTypeVersionsOutput struct {

	// If the request doesn't return all of the remaining results, NextToken is set to
	// a token. To retrieve the next set of results, call this action again and assign
	// that token to the request object's NextToken parameter. If the request returns
	// all results, NextToken is set to null .
	NextToken *string

	// A list of TypeVersionSummary structures that contain information about the
	// specified extension's versions.
	TypeVersionSummaries []types.TypeVersionSummary

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationListTypeVersionsMiddlewares(stack *middleware.Stack, options Options) (err error) {
	err = stack.Serialize.Add(&awsAwsquery_serializeOpListTypeVersions{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsquery_deserializeOpListTypeVersions{}, middleware.After)
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
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opListTypeVersions(options.Region), middleware.Before); err != nil {
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

// ListTypeVersionsAPIClient is a client that implements the ListTypeVersions
// operation.
type ListTypeVersionsAPIClient interface {
	ListTypeVersions(context.Context, *ListTypeVersionsInput, ...func(*Options)) (*ListTypeVersionsOutput, error)
}

var _ ListTypeVersionsAPIClient = (*Client)(nil)

// ListTypeVersionsPaginatorOptions is the paginator options for ListTypeVersions
type ListTypeVersionsPaginatorOptions struct {
	// The maximum number of results to be returned with a single call. If the number
	// of available results exceeds this maximum, the response includes a NextToken
	// value that you can assign to the NextToken request parameter to get the next
	// set of results.
	Limit int32

	// Set to true if pagination should stop if the service returns a pagination token
	// that matches the most recent token provided to the service.
	StopOnDuplicateToken bool
}

// ListTypeVersionsPaginator is a paginator for ListTypeVersions
type ListTypeVersionsPaginator struct {
	options   ListTypeVersionsPaginatorOptions
	client    ListTypeVersionsAPIClient
	params    *ListTypeVersionsInput
	nextToken *string
	firstPage bool
}

// NewListTypeVersionsPaginator returns a new ListTypeVersionsPaginator
func NewListTypeVersionsPaginator(client ListTypeVersionsAPIClient, params *ListTypeVersionsInput, optFns ...func(*ListTypeVersionsPaginatorOptions)) *ListTypeVersionsPaginator {
	if params == nil {
		params = &ListTypeVersionsInput{}
	}

	options := ListTypeVersionsPaginatorOptions{}
	if params.MaxResults != nil {
		options.Limit = *params.MaxResults
	}

	for _, fn := range optFns {
		fn(&options)
	}

	return &ListTypeVersionsPaginator{
		options:   options,
		client:    client,
		params:    params,
		firstPage: true,
		nextToken: params.NextToken,
	}
}

// HasMorePages returns a boolean indicating whether more pages are available
func (p *ListTypeVersionsPaginator) HasMorePages() bool {
	return p.firstPage || (p.nextToken != nil && len(*p.nextToken) != 0)
}

// NextPage retrieves the next ListTypeVersions page.
func (p *ListTypeVersionsPaginator) NextPage(ctx context.Context, optFns ...func(*Options)) (*ListTypeVersionsOutput, error) {
	if !p.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	params := *p.params
	params.NextToken = p.nextToken

	var limit *int32
	if p.options.Limit > 0 {
		limit = &p.options.Limit
	}
	params.MaxResults = limit

	result, err := p.client.ListTypeVersions(ctx, &params, optFns...)
	if err != nil {
		return nil, err
	}
	p.firstPage = false

	prevToken := p.nextToken
	p.nextToken = result.NextToken

	if p.options.StopOnDuplicateToken &&
		prevToken != nil &&
		p.nextToken != nil &&
		*prevToken == *p.nextToken {
		p.nextToken = nil
	}

	return result, nil
}

func newServiceMetadataMiddleware_opListTypeVersions(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		SigningName:   "cloudformation",
		OperationName: "ListTypeVersions",
	}
}
