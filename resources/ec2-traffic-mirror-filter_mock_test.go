package resources

import (
	"context"
	"fmt"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func Test_Mock_EC2TrafficMirrorFilter_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeTrafficMirrorFilters", mock.Anything, mock.Anything).
		Return(
			&ec2.DescribeTrafficMirrorFiltersOutput{
				TrafficMirrorFilters: []ec2types.TrafficMirrorFilter{
					{
						TrafficMirrorFilterId: ptr.String("tmf-11111111111111111"),
						Description:           ptr.String("test filter"),
						Tags: []ec2types.Tag{
							{Key: ptr.String("Name"), Value: ptr.String("test-filter")},
						},
					},
				},
			}, nil,
		)

	lister := &EC2TrafficMirrorFilterLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	filter := resources[0].(*EC2TrafficMirrorFilter)
	assertions.Equal("tmf-11111111111111111", *filter.TrafficMirrorFilterID)
	assertions.Equal("test filter", *filter.Description)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2TrafficMirrorFilter_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeTrafficMirrorFilters", mock.Anything, mock.Anything).
		Return(
			&ec2.DescribeTrafficMirrorFiltersOutput{
				TrafficMirrorFilters: []ec2types.TrafficMirrorFilter{},
			}, nil,
		)

	lister := &EC2TrafficMirrorFilterLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2TrafficMirrorFilter_List_MultiPage(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	firstPageItems := make([]ec2types.TrafficMirrorFilter, 100)
	for i := range firstPageItems {
		firstPageItems[i] = ec2types.TrafficMirrorFilter{
			TrafficMirrorFilterId: ptr.String(fmt.Sprintf("tmf-%d", i)),
			Description:           ptr.String(fmt.Sprintf("filter-%d", i)),
		}
	}

	mockClient.
		On(
			"DescribeTrafficMirrorFilters",
			mock.Anything,
			mock.MatchedBy(func(input *ec2.DescribeTrafficMirrorFiltersInput) bool {
				return input.NextToken == nil
			}),
		).
		Return(
			&ec2.DescribeTrafficMirrorFiltersOutput{
				TrafficMirrorFilters: firstPageItems,
				NextToken:            ptr.String("page2"),
			}, nil,
		).
		Once()

	mockClient.
		On(
			"DescribeTrafficMirrorFilters",
			mock.Anything,
			mock.MatchedBy(func(input *ec2.DescribeTrafficMirrorFiltersInput) bool {
				return input.NextToken != nil
			}),
		).
		Return(
			&ec2.DescribeTrafficMirrorFiltersOutput{
				TrafficMirrorFilters: []ec2types.TrafficMirrorFilter{
					{
						TrafficMirrorFilterId: ptr.String("tmf-100"),
						Description:           ptr.String("filter-100"),
					},
				},
			}, nil,
		).
		Once()

	lister := &EC2TrafficMirrorFilterLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 101)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2TrafficMirrorFilter_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	filter := &EC2TrafficMirrorFilter{
		svc:                   mockClient,
		TrafficMirrorFilterID: ptr.String("tmf-11111111111111111"),
	}

	mockClient.
		On(
			"DeleteTrafficMirrorFilter",
			mock.Anything,
			&ec2.DeleteTrafficMirrorFilterInput{
				TrafficMirrorFilterId: filter.TrafficMirrorFilterID,
			},
		).
		Return(&ec2.DeleteTrafficMirrorFilterOutput{}, nil)

	err := filter.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2TrafficMirrorFilter_Properties(t *testing.T) {
	assertions := assert.New(t)

	filter := EC2TrafficMirrorFilter{
		TrafficMirrorFilterID: ptr.String("tmf-11111111111111111"),
		Description:           ptr.String("my traffic mirror filter"),
		Tags: []ec2types.Tag{
			{Key: ptr.String("Environment"), Value: ptr.String("production")},
		},
	}

	properties := filter.Properties()

	assertions.Equal("tmf-11111111111111111", properties.Get("TrafficMirrorFilterId"))
	assertions.Equal("my traffic mirror filter", properties.Get("Description"))
	assertions.Equal("production", properties.Get("tag:Environment"))
}

func Test_Mock_EC2TrafficMirrorFilter_Properties_EmptyTags(t *testing.T) {
	assertions := assert.New(t)

	filter := EC2TrafficMirrorFilter{
		TrafficMirrorFilterID: ptr.String("tmf-99999999999999999"),
		Description:           ptr.String("no tags filter"),
		Tags:                  []ec2types.Tag{},
	}

	properties := filter.Properties()

	assertions.Equal("tmf-99999999999999999", properties.Get("TrafficMirrorFilterId"))
	assertions.Equal("no tags filter", properties.Get("Description"))
}

func Test_Mock_EC2TrafficMirrorFilter_String(t *testing.T) {
	assertions := assert.New(t)

	filter := EC2TrafficMirrorFilter{
		TrafficMirrorFilterID: ptr.String("tmf-11111111111111111"),
	}

	assertions.Equal("tmf-11111111111111111", filter.String())
}
