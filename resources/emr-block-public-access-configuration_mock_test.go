package resources

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/emr"
	emrtypes "github.com/aws/aws-sdk-go-v2/service/emr/types"
	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Mock_EMRBlockPublicAccessConfiguration_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockEMRV2Client)

	mockClient.On("GetBlockPublicAccessConfiguration", mock.Anything, mock.Anything).
		Return(&emr.GetBlockPublicAccessConfigurationOutput{
			BlockPublicAccessConfiguration: &emrtypes.BlockPublicAccessConfiguration{
				BlockPublicSecurityGroupRules: ptr.Bool(true),
			},
		}, nil)

	lister := &EMRBlockPublicAccessConfigurationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testEMRV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*EMRBlockPublicAccessConfiguration)
	a.True(*r.BlockPublicSecurityGroupRules)
	mockClient.AssertExpectations(t)
}

func Test_Mock_EMRBlockPublicAccessConfiguration_List_NoConfig(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockEMRV2Client)

	mockClient.On("GetBlockPublicAccessConfiguration", mock.Anything, mock.Anything).
		Return(&emr.GetBlockPublicAccessConfigurationOutput{
			BlockPublicAccessConfiguration: nil,
		}, nil)

	lister := &EMRBlockPublicAccessConfigurationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testEMRV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_EMRBlockPublicAccessConfiguration_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockEMRV2Client)

	r := &EMRBlockPublicAccessConfiguration{
		svc:                           mockClient,
		BlockPublicSecurityGroupRules: ptr.Bool(true),
	}

	mockClient.On("PutBlockPublicAccessConfiguration", mock.Anything,
		&emr.PutBlockPublicAccessConfigurationInput{
			BlockPublicAccessConfiguration: &emrtypes.BlockPublicAccessConfiguration{
				BlockPublicSecurityGroupRules: aws.Bool(false),
			},
		}).Return(&emr.PutBlockPublicAccessConfigurationOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_EMRBlockPublicAccessConfiguration_Properties(t *testing.T) {
	a := assert.New(t)
	r := &EMRBlockPublicAccessConfiguration{
		BlockPublicSecurityGroupRules: ptr.Bool(true),
	}
	props := r.Properties()
	a.Equal("true", props.Get("BlockPublicSecurityGroupRules"))
}

func Test_Mock_EMRBlockPublicAccessConfiguration_String(t *testing.T) {
	a := assert.New(t)
	r := &EMRBlockPublicAccessConfiguration{
		BlockPublicSecurityGroupRules: ptr.Bool(true),
	}
	a.Equal("EMRBlockPublicAccessConfiguration", r.String())
}
