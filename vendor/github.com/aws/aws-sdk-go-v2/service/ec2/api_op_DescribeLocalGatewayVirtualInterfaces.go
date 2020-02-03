// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package ec2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awsutil"
)

type DescribeLocalGatewayVirtualInterfacesInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `type:"boolean"`

	// One or more filters.
	Filters []Filter `locationName:"Filter" locationNameList:"Filter" type:"list"`

	// The IDs of the virtual interfaces.
	LocalGatewayVirtualInterfaceIds []string `locationName:"LocalGatewayVirtualInterfaceId" locationNameList:"item" type:"list"`

	// The maximum number of results to return with a single call. To retrieve the
	// remaining results, make another call with the returned nextToken value.
	MaxResults *int64 `min:"5" type:"integer"`

	// The token for the next page of results.
	NextToken *string `type:"string"`
}

// String returns the string representation
func (s DescribeLocalGatewayVirtualInterfacesInput) String() string {
	return awsutil.Prettify(s)
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *DescribeLocalGatewayVirtualInterfacesInput) Validate() error {
	invalidParams := aws.ErrInvalidParams{Context: "DescribeLocalGatewayVirtualInterfacesInput"}
	if s.MaxResults != nil && *s.MaxResults < 5 {
		invalidParams.Add(aws.NewErrParamMinValue("MaxResults", 5))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

type DescribeLocalGatewayVirtualInterfacesOutput struct {
	_ struct{} `type:"structure"`

	// Information about the virtual interfaces.
	LocalGatewayVirtualInterfaces []LocalGatewayVirtualInterface `locationName:"localGatewayVirtualInterfaceSet" locationNameList:"item" type:"list"`

	// The token to use to retrieve the next page of results. This value is null
	// when there are no more results to return.
	NextToken *string `locationName:"nextToken" type:"string"`
}

// String returns the string representation
func (s DescribeLocalGatewayVirtualInterfacesOutput) String() string {
	return awsutil.Prettify(s)
}

const opDescribeLocalGatewayVirtualInterfaces = "DescribeLocalGatewayVirtualInterfaces"

// DescribeLocalGatewayVirtualInterfacesRequest returns a request value for making API operation for
// Amazon Elastic Compute Cloud.
//
// Describes the specified local gateway virtual interfaces.
//
//    // Example sending a request using DescribeLocalGatewayVirtualInterfacesRequest.
//    req := client.DescribeLocalGatewayVirtualInterfacesRequest(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
//
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DescribeLocalGatewayVirtualInterfaces
func (c *Client) DescribeLocalGatewayVirtualInterfacesRequest(input *DescribeLocalGatewayVirtualInterfacesInput) DescribeLocalGatewayVirtualInterfacesRequest {
	op := &aws.Operation{
		Name:       opDescribeLocalGatewayVirtualInterfaces,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &DescribeLocalGatewayVirtualInterfacesInput{}
	}

	req := c.newRequest(op, input, &DescribeLocalGatewayVirtualInterfacesOutput{})
	return DescribeLocalGatewayVirtualInterfacesRequest{Request: req, Input: input, Copy: c.DescribeLocalGatewayVirtualInterfacesRequest}
}

// DescribeLocalGatewayVirtualInterfacesRequest is the request type for the
// DescribeLocalGatewayVirtualInterfaces API operation.
type DescribeLocalGatewayVirtualInterfacesRequest struct {
	*aws.Request
	Input *DescribeLocalGatewayVirtualInterfacesInput
	Copy  func(*DescribeLocalGatewayVirtualInterfacesInput) DescribeLocalGatewayVirtualInterfacesRequest
}

// Send marshals and sends the DescribeLocalGatewayVirtualInterfaces API request.
func (r DescribeLocalGatewayVirtualInterfacesRequest) Send(ctx context.Context) (*DescribeLocalGatewayVirtualInterfacesResponse, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &DescribeLocalGatewayVirtualInterfacesResponse{
		DescribeLocalGatewayVirtualInterfacesOutput: r.Request.Data.(*DescribeLocalGatewayVirtualInterfacesOutput),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// DescribeLocalGatewayVirtualInterfacesResponse is the response type for the
// DescribeLocalGatewayVirtualInterfaces API operation.
type DescribeLocalGatewayVirtualInterfacesResponse struct {
	*DescribeLocalGatewayVirtualInterfacesOutput

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// DescribeLocalGatewayVirtualInterfaces request.
func (r *DescribeLocalGatewayVirtualInterfacesResponse) SDKResponseMetdata() *aws.Response {
	return r.response
}
