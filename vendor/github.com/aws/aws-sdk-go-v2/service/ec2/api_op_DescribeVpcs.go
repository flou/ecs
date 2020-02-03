// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package ec2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awsutil"
)

type DescribeVpcsInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// One or more filters.
	//
	//    * cidr - The primary IPv4 CIDR block of the VPC. The CIDR block you specify
	//    must exactly match the VPC's CIDR block for information to be returned
	//    for the VPC. Must contain the slash followed by one or two digits (for
	//    example, /28).
	//
	//    * cidr-block-association.cidr-block - An IPv4 CIDR block associated with
	//    the VPC.
	//
	//    * cidr-block-association.association-id - The association ID for an IPv4
	//    CIDR block associated with the VPC.
	//
	//    * cidr-block-association.state - The state of an IPv4 CIDR block associated
	//    with the VPC.
	//
	//    * dhcp-options-id - The ID of a set of DHCP options.
	//
	//    * ipv6-cidr-block-association.ipv6-cidr-block - An IPv6 CIDR block associated
	//    with the VPC.
	//
	//    * ipv6-cidr-block-association.ipv6-pool - The ID of the IPv6 address pool
	//    from which the IPv6 CIDR block is allocated.
	//
	//    * ipv6-cidr-block-association.association-id - The association ID for
	//    an IPv6 CIDR block associated with the VPC.
	//
	//    * ipv6-cidr-block-association.state - The state of an IPv6 CIDR block
	//    associated with the VPC.
	//
	//    * isDefault - Indicates whether the VPC is the default VPC.
	//
	//    * owner-id - The ID of the AWS account that owns the VPC.
	//
	//    * state - The state of the VPC (pending | available).
	//
	//    * tag:<key> - The key/value combination of a tag assigned to the resource.
	//    Use the tag key in the filter name and the tag value as the filter value.
	//    For example, to find all resources that have a tag with the key Owner
	//    and the value TeamA, specify tag:Owner for the filter name and TeamA for
	//    the filter value.
	//
	//    * tag-key - The key of a tag assigned to the resource. Use this filter
	//    to find all resources assigned a tag with a specific key, regardless of
	//    the tag value.
	//
	//    * vpc-id - The ID of the VPC.
	Filters []Filter `locationName:"Filter" locationNameList:"Filter" type:"list"`

	// The maximum number of results to return with a single call. To retrieve the
	// remaining results, make another call with the returned nextToken value.
	MaxResults *int64 `min:"5" type:"integer"`

	// The token for the next page of results.
	NextToken *string `type:"string"`

	// One or more VPC IDs.
	//
	// Default: Describes all your VPCs.
	VpcIds []string `locationName:"VpcId" locationNameList:"VpcId" type:"list"`
}

// String returns the string representation
func (s DescribeVpcsInput) String() string {
	return awsutil.Prettify(s)
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *DescribeVpcsInput) Validate() error {
	invalidParams := aws.ErrInvalidParams{Context: "DescribeVpcsInput"}
	if s.MaxResults != nil && *s.MaxResults < 5 {
		invalidParams.Add(aws.NewErrParamMinValue("MaxResults", 5))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

type DescribeVpcsOutput struct {
	_ struct{} `type:"structure"`

	// The token to use to retrieve the next page of results. This value is null
	// when there are no more results to return.
	NextToken *string `locationName:"nextToken" type:"string"`

	// Information about one or more VPCs.
	Vpcs []Vpc `locationName:"vpcSet" locationNameList:"item" type:"list"`
}

// String returns the string representation
func (s DescribeVpcsOutput) String() string {
	return awsutil.Prettify(s)
}

const opDescribeVpcs = "DescribeVpcs"

// DescribeVpcsRequest returns a request value for making API operation for
// Amazon Elastic Compute Cloud.
//
// Describes one or more of your VPCs.
//
//    // Example sending a request using DescribeVpcsRequest.
//    req := client.DescribeVpcsRequest(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
//
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DescribeVpcs
func (c *Client) DescribeVpcsRequest(input *DescribeVpcsInput) DescribeVpcsRequest {
	op := &aws.Operation{
		Name:       opDescribeVpcs,
		HTTPMethod: "POST",
		HTTPPath:   "/",
		Paginator: &aws.Paginator{
			InputTokens:     []string{"NextToken"},
			OutputTokens:    []string{"NextToken"},
			LimitToken:      "MaxResults",
			TruncationToken: "",
		},
	}

	if input == nil {
		input = &DescribeVpcsInput{}
	}

	req := c.newRequest(op, input, &DescribeVpcsOutput{})
	return DescribeVpcsRequest{Request: req, Input: input, Copy: c.DescribeVpcsRequest}
}

// DescribeVpcsRequest is the request type for the
// DescribeVpcs API operation.
type DescribeVpcsRequest struct {
	*aws.Request
	Input *DescribeVpcsInput
	Copy  func(*DescribeVpcsInput) DescribeVpcsRequest
}

// Send marshals and sends the DescribeVpcs API request.
func (r DescribeVpcsRequest) Send(ctx context.Context) (*DescribeVpcsResponse, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &DescribeVpcsResponse{
		DescribeVpcsOutput: r.Request.Data.(*DescribeVpcsOutput),
		response:           &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// NewDescribeVpcsRequestPaginator returns a paginator for DescribeVpcs.
// Use Next method to get the next page, and CurrentPage to get the current
// response page from the paginator. Next will return false, if there are
// no more pages, or an error was encountered.
//
// Note: This operation can generate multiple requests to a service.
//
//   // Example iterating over pages.
//   req := client.DescribeVpcsRequest(input)
//   p := ec2.NewDescribeVpcsRequestPaginator(req)
//
//   for p.Next(context.TODO()) {
//       page := p.CurrentPage()
//   }
//
//   if err := p.Err(); err != nil {
//       return err
//   }
//
func NewDescribeVpcsPaginator(req DescribeVpcsRequest) DescribeVpcsPaginator {
	return DescribeVpcsPaginator{
		Pager: aws.Pager{
			NewRequest: func(ctx context.Context) (*aws.Request, error) {
				var inCpy *DescribeVpcsInput
				if req.Input != nil {
					tmp := *req.Input
					inCpy = &tmp
				}

				newReq := req.Copy(inCpy)
				newReq.SetContext(ctx)
				return newReq.Request, nil
			},
		},
	}
}

// DescribeVpcsPaginator is used to paginate the request. This can be done by
// calling Next and CurrentPage.
type DescribeVpcsPaginator struct {
	aws.Pager
}

func (p *DescribeVpcsPaginator) CurrentPage() *DescribeVpcsOutput {
	return p.Pager.CurrentPage().(*DescribeVpcsOutput)
}

// DescribeVpcsResponse is the response type for the
// DescribeVpcs API operation.
type DescribeVpcsResponse struct {
	*DescribeVpcsOutput

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// DescribeVpcs request.
func (r *DescribeVpcsResponse) SDKResponseMetdata() *aws.Response {
	return r.response
}