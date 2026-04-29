package resources

import (
	"context"
	"fmt"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/vpclattice"
	lattice "github.com/aws/aws-sdk-go-v2/service/vpclattice/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

var testVPCLatticeListerOpts = &nuke.ListerOpts{
	Logger: logrus.WithField("test", true),
}

func Test_Mock_VPCLatticeTargetGroup_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	mockClient.
		On("ListTargetGroups", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListTargetGroupsOutput{
				Items: []lattice.TargetGroupSummary{
					{
						Id:   ptr.String("tg-1"),
						Arn:  ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:targetgroup/tg-1"),
						Name: ptr.String("group-1"),
						Type: lattice.TargetGroupTypeInstance,
					},
				},
			}, nil,
		)

	mockClient.
		On("ListTagsForResource", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListTagsForResourceOutput{
				Tags: map[string]string{"env": "test"},
			}, nil,
		)

	lister := &VPCLatticeTargetGroupLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testVPCLatticeListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	targetGroup := resources[0].(*VPCLatticeTargetGroup)
	assertions.Equal("tg-1", *targetGroup.ID)
	assertions.Equal("group-1", *targetGroup.Name)
	assertions.Equal("test", targetGroup.Tags["env"])

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeTargetGroup_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	mockClient.
		On("ListTargetGroups", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListTargetGroupsOutput{
				Items: []lattice.TargetGroupSummary{},
			}, nil,
		)

	lister := &VPCLatticeTargetGroupLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testVPCLatticeListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeTargetGroup_List_MultiPage(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	// Build 101 items across 2 pages to verify pagination handling
	firstPageItems := make([]lattice.TargetGroupSummary, 100)
	for i := range firstPageItems {
		firstPageItems[i] = lattice.TargetGroupSummary{
			Id:   ptr.String(fmt.Sprintf("tg-%d", i)),
			Arn:  ptr.String(fmt.Sprintf("arn:aws:vpc-lattice:us-east-1:123456789012:targetgroup/tg-%d", i)),
			Name: ptr.String(fmt.Sprintf("group-%d", i)),
			Type: lattice.TargetGroupTypeInstance,
		}
	}

	secondPageItems := []lattice.TargetGroupSummary{
		{
			Id:   ptr.String("tg-100"),
			Arn:  ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:targetgroup/tg-100"),
			Name: ptr.String("group-100"),
			Type: lattice.TargetGroupTypeIp,
		},
	}

	mockClient.
		On(
			"ListTargetGroups",
			mock.Anything,
			mock.MatchedBy(func(input *vpclattice.ListTargetGroupsInput) bool {
				return input.NextToken == nil
			}),
		).
		Return(
			&vpclattice.ListTargetGroupsOutput{
				Items:     firstPageItems,
				NextToken: ptr.String("page2"),
			}, nil,
		).
		Once()

	mockClient.
		On(
			"ListTargetGroups",
			mock.Anything,
			mock.MatchedBy(func(input *vpclattice.ListTargetGroupsInput) bool {
				return input.NextToken != nil && *input.NextToken == "page2"
			}),
		).
		Return(
			&vpclattice.ListTargetGroupsOutput{
				Items: secondPageItems,
			}, nil,
		).
		Once()

	mockClient.
		On("ListTagsForResource", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListTagsForResourceOutput{
				Tags: map[string]string{},
			}, nil,
		)

	lister := &VPCLatticeTargetGroupLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testVPCLatticeListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 101)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeTargetGroup_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	targetGroup := &VPCLatticeTargetGroup{
		svc: mockClient,
		ARN: ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:targetgroup/tg-1"),
	}

	mockClient.
		On("ListTargets", mock.Anything, mock.Anything).
		Return(&vpclattice.ListTargetsOutput{
			Items: []lattice.TargetSummary{},
		}, nil)

	mockClient.
		On(
			"DeleteTargetGroup",
			mock.Anything,
			&vpclattice.DeleteTargetGroupInput{
				TargetGroupIdentifier: targetGroup.ARN,
			},
		).
		Return(&vpclattice.DeleteTargetGroupOutput{}, nil)

	err := targetGroup.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeTargetGroup_Properties(t *testing.T) {
	assertions := assert.New(t)

	targetGroup := VPCLatticeTargetGroup{
		ID:   ptr.String("tg-12345"),
		ARN:  ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:targetgroup/tg-12345"),
		Name: ptr.String("my-target-group"),
		Type: ptr.String("INSTANCE"),
		Tags: map[string]string{
			"Environment": "production",
			"Team":        "platform",
		},
	}

	properties := targetGroup.Properties()

	assertions.Equal("tg-12345", properties.Get("ID"))
	assertions.Equal("arn:aws:vpc-lattice:us-east-1:123456789012:targetgroup/tg-12345", properties.Get("ARN"))
	assertions.Equal("my-target-group", properties.Get("Name"))
	assertions.Equal("INSTANCE", properties.Get("Type"))
	assertions.Equal("production", properties.Get("tag:Environment"))
	assertions.Equal("platform", properties.Get("tag:Team"))
}

func Test_Mock_VPCLatticeTargetGroup_Properties_EmptyTags(t *testing.T) {
	assertions := assert.New(t)

	targetGroup := VPCLatticeTargetGroup{
		ID:   ptr.String("tg-99999"),
		ARN:  ptr.String("arn:aws:vpc-lattice:us-west-2:111111111111:targetgroup/tg-99999"),
		Name: ptr.String("no-tags-group"),
		Type: ptr.String("IP"),
		Tags: map[string]string{},
	}

	properties := targetGroup.Properties()

	assertions.Equal("tg-99999", properties.Get("ID"))
	assertions.Equal("no-tags-group", properties.Get("Name"))
	assertions.Equal("IP", properties.Get("Type"))
}

func Test_Mock_VPCLatticeTargetGroup_String(t *testing.T) {
	assertions := assert.New(t)

	targetGroup := VPCLatticeTargetGroup{
		Name: ptr.String("my-target-group"),
	}

	assertions.Equal("my-target-group", targetGroup.String())
}
