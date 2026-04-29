package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

func Test_Mock_SSMDefaultPatchBaseline_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSSMV2Client)

	mockClient.On("DescribePatchBaselines", mock.Anything, mock.Anything).
		Return(&ssm.DescribePatchBaselinesOutput{
			BaselineIdentities: []ssmtypes.PatchBaselineIdentity{
				{
					BaselineId:      ptr.String("pb-custom-123"),
					BaselineName:    ptr.String("MyCustomBaseline"),
					DefaultBaseline: true,
					OperatingSystem: ssmtypes.OperatingSystemWindows,
				},
				{
					BaselineId:      ptr.String("pb-custom-456"),
					BaselineName:    ptr.String("NonDefaultBaseline"),
					DefaultBaseline: false,
					OperatingSystem: ssmtypes.OperatingSystemAmazonLinux2,
				},
			},
		}, nil)

	lister := &SSMDefaultPatchBaselineLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*SSMDefaultPatchBaseline)
	a.Equal("pb-custom-123", *r.BaselineID)
	a.Equal("MyCustomBaseline", *r.BaselineName)
	a.Equal(ssmtypes.OperatingSystemWindows, r.OperatingSystem)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SSMDefaultPatchBaseline_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSSMV2Client)

	mockClient.On("DescribePatchBaselines", mock.Anything, mock.Anything).
		Return(&ssm.DescribePatchBaselinesOutput{
			BaselineIdentities: []ssmtypes.PatchBaselineIdentity{},
		}, nil)

	lister := &SSMDefaultPatchBaselineLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SSMDefaultPatchBaseline_List_NoDefaults(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSSMV2Client)

	mockClient.On("DescribePatchBaselines", mock.Anything, mock.Anything).
		Return(&ssm.DescribePatchBaselinesOutput{
			BaselineIdentities: []ssmtypes.PatchBaselineIdentity{
				{
					BaselineId:      ptr.String("pb-custom-789"),
					BaselineName:    ptr.String("NonDefault"),
					DefaultBaseline: false,
					OperatingSystem: ssmtypes.OperatingSystemUbuntu,
				},
			},
		}, nil)

	lister := &SSMDefaultPatchBaselineLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SSMDefaultPatchBaseline_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSSMV2Client)

	r := &SSMDefaultPatchBaseline{
		svc:             mockClient,
		BaselineID:      ptr.String("pb-custom-123"),
		BaselineName:    ptr.String("MyCustomBaseline"),
		OperatingSystem: ssmtypes.OperatingSystemWindows,
	}

	mockClient.On("DescribePatchBaselines", mock.Anything, mock.Anything).
		Return(&ssm.DescribePatchBaselinesOutput{
			BaselineIdentities: []ssmtypes.PatchBaselineIdentity{
				{
					BaselineId:   ptr.String("pb-aws-default"),
					BaselineName: ptr.String("AWS-DefaultPatchBaseline"),
				},
			},
		}, nil)

	mockClient.On("RegisterDefaultPatchBaseline", mock.Anything,
		&ssm.RegisterDefaultPatchBaselineInput{
			BaselineId: ptr.String("pb-aws-default"),
		}).Return(&ssm.RegisterDefaultPatchBaselineOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_SSMDefaultPatchBaseline_Properties(t *testing.T) {
	a := assert.New(t)
	r := &SSMDefaultPatchBaseline{
		BaselineID:      ptr.String("pb-custom-123"),
		BaselineName:    ptr.String("MyCustomBaseline"),
		OperatingSystem: ssmtypes.OperatingSystemWindows,
	}
	props := r.Properties()
	a.Equal("pb-custom-123", props.Get("BaselineID"))
	a.Equal("MyCustomBaseline", props.Get("BaselineName"))
}

func Test_Mock_SSMDefaultPatchBaseline_String(t *testing.T) {
	a := assert.New(t)
	r := &SSMDefaultPatchBaseline{
		BaselineID: ptr.String("pb-custom-123"),
	}
	a.Equal("pb-custom-123", r.String())
}
