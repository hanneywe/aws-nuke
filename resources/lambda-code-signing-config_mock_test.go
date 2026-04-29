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

const TestLambdaCodeSigningConfigArn = "arn:aws:lambda:us-east-1:123456789012:code-signing-config:csc-1234567890abcdef0"

func Test_Mock_LambdaCodeSigningConfig_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockLambdaClient)

	mockClient.
		On("ListCodeSigningConfigs", mock.Anything, mock.Anything).
		Return(&lambda.ListCodeSigningConfigsOutput{
			CodeSigningConfigs: []lambdatypes.CodeSigningConfig{
				{
					CodeSigningConfigId:  ptr.String("csc-1234567890abcdef0"),
					CodeSigningConfigArn: ptr.String(TestLambdaCodeSigningConfigArn),
				},
			},
		}, nil)

	lister := &LambdaCodeSigningConfigLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testLambdaListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	codeSigningConfig := resources[0].(*LambdaCodeSigningConfig)
	assertions.Equal("csc-1234567890abcdef0", *codeSigningConfig.CodeSigningConfigID)
	assertions.Equal(TestLambdaCodeSigningConfigArn, *codeSigningConfig.CodeSigningConfigArn)

	mockClient.AssertExpectations(t)
}

func Test_Mock_LambdaCodeSigningConfig_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockLambdaClient)

	mockClient.
		On("ListCodeSigningConfigs", mock.Anything, mock.Anything).
		Return(&lambda.ListCodeSigningConfigsOutput{
			CodeSigningConfigs: []lambdatypes.CodeSigningConfig{},
		}, nil)

	lister := &LambdaCodeSigningConfigLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testLambdaListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_LambdaCodeSigningConfig_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockLambdaClient)

	codeSigningConfig := &LambdaCodeSigningConfig{
		svc:                  mockClient,
		CodeSigningConfigID:  ptr.String("csc-1234567890abcdef0"),
		CodeSigningConfigArn: ptr.String(TestLambdaCodeSigningConfigArn),
	}

	mockClient.
		On("DeleteCodeSigningConfig", mock.Anything, &lambda.DeleteCodeSigningConfigInput{
			CodeSigningConfigArn: codeSigningConfig.CodeSigningConfigArn,
		}).
		Return(&lambda.DeleteCodeSigningConfigOutput{}, nil)

	err := codeSigningConfig.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_LambdaCodeSigningConfig_Properties(t *testing.T) {
	assertions := assert.New(t)

	codeSigningConfig := LambdaCodeSigningConfig{
		CodeSigningConfigID:  ptr.String("csc-1234567890abcdef0"),
		CodeSigningConfigArn: ptr.String(TestLambdaCodeSigningConfigArn),
	}

	properties := codeSigningConfig.Properties()

	assertions.Equal("csc-1234567890abcdef0", properties.Get("CodeSigningConfigID"))
	assertions.Equal(TestLambdaCodeSigningConfigArn, properties.Get("CodeSigningConfigArn"))
}

func Test_Mock_LambdaCodeSigningConfig_String(t *testing.T) {
	assertions := assert.New(t)

	codeSigningConfig := LambdaCodeSigningConfig{
		CodeSigningConfigID: ptr.String("csc-1234567890abcdef0"),
	}

	assertions.Equal("csc-1234567890abcdef0", codeSigningConfig.String())
}
