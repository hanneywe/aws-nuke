package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/wellarchitected"
	wellarchitectedtypes "github.com/aws/aws-sdk-go-v2/service/wellarchitected/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

var testWellArchitectedListerOpts = &nuke.ListerOpts{
	Logger: logrus.WithField("test", true),
}

func Test_Mock_WellArchitectedWorkload_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockWellArchitectedClient)

	mockClient.
		On("ListWorkloads", mock.Anything, mock.Anything).
		Return(
			&wellarchitected.ListWorkloadsOutput{
				WorkloadSummaries: []wellarchitectedtypes.WorkloadSummary{
					{
						WorkloadId:   ptr.String("wl-12345"),
						WorkloadName: ptr.String("test-workload"),
					},
				},
			}, nil,
		)

	lister := &WellArchitectedWorkloadLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testWellArchitectedListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	workload := resources[0].(*WellArchitectedWorkload)
	assertions.Equal("wl-12345", *workload.WorkloadID)
	assertions.Equal("test-workload", *workload.WorkloadName)

	mockClient.AssertExpectations(t)
}

func Test_Mock_WellArchitectedWorkload_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockWellArchitectedClient)

	mockClient.
		On("ListWorkloads", mock.Anything, mock.Anything).
		Return(
			&wellarchitected.ListWorkloadsOutput{
				WorkloadSummaries: []wellarchitectedtypes.WorkloadSummary{},
			}, nil,
		)

	lister := &WellArchitectedWorkloadLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testWellArchitectedListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_WellArchitectedWorkload_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockWellArchitectedClient)

	workload := &WellArchitectedWorkload{
		svc:          mockClient,
		WorkloadID:   ptr.String("wl-12345"),
		WorkloadName: ptr.String("test-workload"),
	}

	mockClient.
		On(
			"DeleteWorkload",
			mock.Anything,
			mock.MatchedBy(func(input *wellarchitected.DeleteWorkloadInput) bool {
				return *input.WorkloadId == "wl-12345"
			}),
		).
		Return(&wellarchitected.DeleteWorkloadOutput{}, nil)

	err := workload.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_WellArchitectedWorkload_Properties(t *testing.T) {
	assertions := assert.New(t)

	workload := WellArchitectedWorkload{
		WorkloadID:   ptr.String("wl-12345"),
		WorkloadName: ptr.String("test-workload"),
	}

	properties := workload.Properties()

	assertions.Equal("wl-12345", properties.Get("WorkloadId"))
	assertions.Equal("test-workload", properties.Get("WorkloadName"))
}

func Test_Mock_WellArchitectedWorkload_String(t *testing.T) {
	assertions := assert.New(t)

	workload := WellArchitectedWorkload{
		WorkloadID:   ptr.String("wl-12345"),
		WorkloadName: ptr.String("test-workload"),
	}

	assertions.Equal("wl-12345", workload.String())
}
