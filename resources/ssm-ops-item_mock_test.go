package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

var testSSMV2ListerOpts = &nuke.ListerOpts{
	Logger: logrus.WithField("test", true),
}

func Test_Mock_SSMOpsItem_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockSSMV2Client)

	mockClient.On("DescribeOpsItems", mock.Anything, mock.Anything).
		Return(&ssm.DescribeOpsItemsOutput{
			OpsItemSummaries: []ssmtypes.OpsItemSummary{
				{
					OpsItemId: ptr.String("oi-1234567890"),
					Title:     ptr.String("Test OpsItem"),
				},
				{
					OpsItemId: ptr.String("oi-0987654321"),
					Title:     ptr.String("Another OpsItem"),
				},
			},
		}, nil)

	lister := &SSMOpsItemLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSSMV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 2)

	opsItem := resources[0].(*SSMOpsItem)
	assertions.Equal("oi-1234567890", *opsItem.OpsItemID)
	assertions.Equal("Test OpsItem", *opsItem.Title)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SSMOpsItem_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockSSMV2Client)

	mockClient.On("DescribeOpsItems", mock.Anything, mock.Anything).
		Return(&ssm.DescribeOpsItemsOutput{
			OpsItemSummaries: []ssmtypes.OpsItemSummary{},
		}, nil)

	lister := &SSMOpsItemLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSSMV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SSMOpsItem_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockSSMV2Client)

	opsItem := &SSMOpsItem{
		svc:       mockClient,
		OpsItemID: ptr.String("oi-1234567890"),
		Title:     ptr.String("Test OpsItem"),
	}

	mockClient.On("UpdateOpsItem", mock.Anything, &ssm.UpdateOpsItemInput{
		OpsItemId: opsItem.OpsItemID,
		Status:    ssmtypes.OpsItemStatusResolved,
	}).Return(&ssm.UpdateOpsItemOutput{}, nil)

	err := opsItem.Remove(context.TODO())
	assertions.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SSMOpsItem_Properties(t *testing.T) {
	assertions := assert.New(t)

	opsItem := SSMOpsItem{
		OpsItemID: ptr.String("oi-1234567890"),
		Title:     ptr.String("Test OpsItem"),
	}

	properties := opsItem.Properties()
	assertions.Equal("oi-1234567890", properties.Get("OpsItemId"))
	assertions.Equal("Test OpsItem", properties.Get("Title"))
}

func Test_Mock_SSMOpsItem_String(t *testing.T) {
	assertions := assert.New(t)
	opsItem := SSMOpsItem{OpsItemID: ptr.String("oi-1234567890")}
	assertions.Equal("oi-1234567890", opsItem.String())
}
