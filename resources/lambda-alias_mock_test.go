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

func Test_Mock_LambdaAlias_List(t *testing.T) {
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

	mockClient.On("ListAliases", mock.Anything, &lambda.ListAliasesInput{
		FunctionName: ptr.String("my-function"),
	}).Return(&lambda.ListAliasesOutput{
		Aliases: []lambdatypes.AliasConfiguration{
			{
				Name:     ptr.String("prod"),
				AliasArn: ptr.String("arn:aws:lambda:us-east-1:123456789012:function:my-function:prod"),
			},
			{
				Name:     ptr.String("staging"),
				AliasArn: ptr.String("arn:aws:lambda:us-east-1:123456789012:function:my-function:staging"),
			},
		},
	}, nil)

	lister := &LambdaAliasLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testLambdaListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 2)

	alias := resources[0].(*LambdaAlias)
	assertions.Equal("my-function", *alias.FunctionName)
	assertions.Equal("prod", *alias.Name)
	assertions.Equal("arn:aws:lambda:us-east-1:123456789012:function:my-function:prod", *alias.AliasArn)

	mockClient.AssertExpectations(t)
}

func Test_Mock_LambdaAlias_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := &mockLambdaClient{}

	mockClient.On("ListFunctions", mock.Anything, &lambda.ListFunctionsInput{}).
		Return(&lambda.ListFunctionsOutput{}, nil)

	lister := &LambdaAliasLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testLambdaListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_LambdaAlias_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := &mockLambdaClient{}

	alias := &LambdaAlias{
		svc:          mockClient,
		FunctionName: ptr.String("my-function"),
		Name:         ptr.String("prod"),
		AliasArn:     ptr.String("arn:aws:lambda:us-east-1:123456789012:function:my-function:prod"),
	}

	mockClient.On("DeleteAlias", mock.Anything, &lambda.DeleteAliasInput{
		FunctionName: alias.FunctionName,
		Name:         alias.Name,
	}).Return(&lambda.DeleteAliasOutput{}, nil)

	err := alias.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_LambdaAlias_Properties(t *testing.T) {
	assertions := assert.New(t)

	alias := &LambdaAlias{
		FunctionName: ptr.String("my-function"),
		Name:         ptr.String("prod"),
		AliasArn:     ptr.String("arn:aws:lambda:us-east-1:123456789012:function:my-function:prod"),
	}

	properties := alias.Properties()
	assertions.Equal("my-function", properties.Get("FunctionName"))
	assertions.Equal("prod", properties.Get("Name"))
	assertions.Equal("arn:aws:lambda:us-east-1:123456789012:function:my-function:prod", properties.Get("AliasArn"))
}

func Test_Mock_LambdaAlias_String(t *testing.T) {
	assertions := assert.New(t)

	alias := &LambdaAlias{
		Name: ptr.String("prod"),
	}

	assertions.Equal("prod", alias.String())
}
