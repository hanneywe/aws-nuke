package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
	lambdatypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

func Test_Mock_LambdaVersion_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := &mockLambdaClient{}

	mockClient.On("ListFunctions", mock.Anything, &lambda.ListFunctionsInput{}).
		Return(&lambda.ListFunctionsOutput{
			Functions: []lambdatypes.FunctionConfiguration{
				{
					FunctionName: ptr.String("my-function"),
				},
			},
		}, nil)

	mockClient.On("ListVersionsByFunction", mock.Anything, &lambda.ListVersionsByFunctionInput{
		FunctionName: ptr.String("my-function"),
	}).Return(&lambda.ListVersionsByFunctionOutput{
		Versions: []lambdatypes.FunctionConfiguration{
			{
				FunctionName: ptr.String("my-function"),
				Version:      ptr.String("$LATEST"),
				FunctionArn:  ptr.String("arn:aws:lambda:us-east-1:123456789012:function:my-function:$LATEST"),
			},
			{
				FunctionName: ptr.String("my-function"),
				Version:      ptr.String("1"),
				FunctionArn:  ptr.String("arn:aws:lambda:us-east-1:123456789012:function:my-function:1"),
			},
		},
	}, nil)

	lister := &LambdaVersionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testLambdaListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 2)

	version := resources[1].(*LambdaVersion)
	assertions.Equal("my-function", *version.FunctionName)
	assertions.Equal("1", *version.Version)
	assertions.Equal("arn:aws:lambda:us-east-1:123456789012:function:my-function:1", *version.FunctionArn)

	mockClient.AssertExpectations(t)
}

func Test_Mock_LambdaVersion_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := &mockLambdaClient{}

	mockClient.On("ListFunctions", mock.Anything, &lambda.ListFunctionsInput{}).
		Return(&lambda.ListFunctionsOutput{}, nil)

	lister := &LambdaVersionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testLambdaListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_LambdaVersion_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := &mockLambdaClient{}

	version := &LambdaVersion{
		svc:          mockClient,
		FunctionName: ptr.String("my-function"),
		Version:      ptr.String("1"),
		FunctionArn:  ptr.String("arn:aws:lambda:us-east-1:123456789012:function:my-function:1"),
	}

	mockClient.On("DeleteFunction", mock.Anything, &lambda.DeleteFunctionInput{
		FunctionName: version.FunctionName,
		Qualifier:    version.Version,
	}).Return(&lambda.DeleteFunctionOutput{}, nil)

	err := version.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_LambdaVersion_Properties(t *testing.T) {
	assertions := assert.New(t)

	version := &LambdaVersion{
		FunctionName: ptr.String("my-function"),
		Version:      ptr.String("1"),
		FunctionArn:  ptr.String("arn:aws:lambda:us-east-1:123456789012:function:my-function:1"),
	}

	properties := version.Properties()
	assertions.Equal("my-function", properties.Get("FunctionName"))
	assertions.Equal("1", properties.Get("Version"))
	assertions.Equal("arn:aws:lambda:us-east-1:123456789012:function:my-function:1", properties.Get("FunctionArn"))
}

func Test_Mock_LambdaVersion_String(t *testing.T) {
	assertions := assert.New(t)

	version := &LambdaVersion{
		Version: ptr.String("1"),
	}

	assertions.Equal("1", version.String())
}

func Test_Mock_LambdaVersion_Filter_Latest(t *testing.T) {
	assertions := assert.New(t)

	version := &LambdaVersion{
		FunctionName: ptr.String("my-function"),
		Version:      ptr.String("$LATEST"),
		FunctionArn:  ptr.String("arn:aws:lambda:us-east-1:123456789012:function:my-function:$LATEST"),
	}

	err := version.Filter()
	assertions.Error(err)
	assertions.Contains(err.Error(), "$LATEST")
}

func Test_Mock_LambdaVersion_Filter_NonLatest(t *testing.T) {
	assertions := assert.New(t)

	version := &LambdaVersion{
		FunctionName: ptr.String("my-function"),
		Version:      ptr.String("1"),
		FunctionArn:  ptr.String("arn:aws:lambda:us-east-1:123456789012:function:my-function:1"),
	}

	err := version.Filter()
	assertions.NoError(err)
}
