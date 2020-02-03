// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package ec2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awsutil"
)

type CreateLaunchTemplateInput struct {
	_ struct{} `type:"structure"`

	// Unique, case-sensitive identifier you provide to ensure the idempotency of
	// the request. For more information, see Ensuring Idempotency (https://docs.aws.amazon.com/AWSEC2/latest/APIReference/Run_Instance_Idempotency.html).
	//
	// Constraint: Maximum 128 ASCII characters.
	ClientToken *string `type:"string"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `type:"boolean"`

	// The information for the launch template.
	//
	// LaunchTemplateData is a required field
	LaunchTemplateData *RequestLaunchTemplateData `type:"structure" required:"true"`

	// A name for the launch template.
	//
	// LaunchTemplateName is a required field
	LaunchTemplateName *string `min:"3" type:"string" required:"true"`

	// The tags to apply to the launch template during creation.
	TagSpecifications []TagSpecification `locationName:"TagSpecification" locationNameList:"item" type:"list"`

	// A description for the first version of the launch template.
	VersionDescription *string `type:"string"`
}

// String returns the string representation
func (s CreateLaunchTemplateInput) String() string {
	return awsutil.Prettify(s)
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *CreateLaunchTemplateInput) Validate() error {
	invalidParams := aws.ErrInvalidParams{Context: "CreateLaunchTemplateInput"}

	if s.LaunchTemplateData == nil {
		invalidParams.Add(aws.NewErrParamRequired("LaunchTemplateData"))
	}

	if s.LaunchTemplateName == nil {
		invalidParams.Add(aws.NewErrParamRequired("LaunchTemplateName"))
	}
	if s.LaunchTemplateName != nil && len(*s.LaunchTemplateName) < 3 {
		invalidParams.Add(aws.NewErrParamMinLen("LaunchTemplateName", 3))
	}
	if s.LaunchTemplateData != nil {
		if err := s.LaunchTemplateData.Validate(); err != nil {
			invalidParams.AddNested("LaunchTemplateData", err.(aws.ErrInvalidParams))
		}
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

type CreateLaunchTemplateOutput struct {
	_ struct{} `type:"structure"`

	// Information about the launch template.
	LaunchTemplate *LaunchTemplate `locationName:"launchTemplate" type:"structure"`
}

// String returns the string representation
func (s CreateLaunchTemplateOutput) String() string {
	return awsutil.Prettify(s)
}

const opCreateLaunchTemplate = "CreateLaunchTemplate"

// CreateLaunchTemplateRequest returns a request value for making API operation for
// Amazon Elastic Compute Cloud.
//
// Creates a launch template. A launch template contains the parameters to launch
// an instance. When you launch an instance using RunInstances, you can specify
// a launch template instead of providing the launch parameters in the request.
//
//    // Example sending a request using CreateLaunchTemplateRequest.
//    req := client.CreateLaunchTemplateRequest(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
//
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/CreateLaunchTemplate
func (c *Client) CreateLaunchTemplateRequest(input *CreateLaunchTemplateInput) CreateLaunchTemplateRequest {
	op := &aws.Operation{
		Name:       opCreateLaunchTemplate,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &CreateLaunchTemplateInput{}
	}

	req := c.newRequest(op, input, &CreateLaunchTemplateOutput{})
	return CreateLaunchTemplateRequest{Request: req, Input: input, Copy: c.CreateLaunchTemplateRequest}
}

// CreateLaunchTemplateRequest is the request type for the
// CreateLaunchTemplate API operation.
type CreateLaunchTemplateRequest struct {
	*aws.Request
	Input *CreateLaunchTemplateInput
	Copy  func(*CreateLaunchTemplateInput) CreateLaunchTemplateRequest
}

// Send marshals and sends the CreateLaunchTemplate API request.
func (r CreateLaunchTemplateRequest) Send(ctx context.Context) (*CreateLaunchTemplateResponse, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &CreateLaunchTemplateResponse{
		CreateLaunchTemplateOutput: r.Request.Data.(*CreateLaunchTemplateOutput),
		response:                   &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// CreateLaunchTemplateResponse is the response type for the
// CreateLaunchTemplate API operation.
type CreateLaunchTemplateResponse struct {
	*CreateLaunchTemplateOutput

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// CreateLaunchTemplate request.
func (r *CreateLaunchTemplateResponse) SDKResponseMetdata() *aws.Response {
	return r.response
}
