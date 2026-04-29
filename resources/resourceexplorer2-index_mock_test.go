package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/resourceexplorer2"
	"github.com/aws/aws-sdk-go-v2/service/resourceexplorer2/types"
)

func Test_Mock_ResourceExplorer2Index_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockResourceExplorer2Client)

	mockClient.On("ListIndexes", mock.Anything, mock.Anything).
		Return(&resourceexplorer2.ListIndexesOutput{
			Indexes: []types.Index{
				{
					Arn:    ptr.String("arn:aws:resource-explorer-2:us-east-1:123456789012:index/12345"),
					Region: ptr.String("us-east-1"),
					Type:   types.IndexTypeLocal,
				},
			},
		}, nil)

	mockClient.On("ListTagsForResource", mock.Anything, mock.Anything).
		Return(&resourceexplorer2.ListTagsForResourceOutput{
			Tags: map[string]string{"env": "test"},
		}, nil)

	lister := &ResourceExplorer2IndexLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testResourceExplorer2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*ResourceExplorer2Index)
	a.Equal("arn:aws:resource-explorer-2:us-east-1:123456789012:index/12345", *r.ARN)
	a.Equal("us-east-1", *r.Region)
	a.Equal(types.IndexTypeLocal, r.Type)
	a.Equal("test", r.Tags["env"])
	mockClient.AssertExpectations(t)
}

func Test_Mock_ResourceExplorer2Index_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockResourceExplorer2Client)

	mockClient.On("ListIndexes", mock.Anything, mock.Anything).
		Return(&resourceexplorer2.ListIndexesOutput{
			Indexes: []types.Index{},
		}, nil)

	lister := &ResourceExplorer2IndexLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testResourceExplorer2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_ResourceExplorer2Index_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockResourceExplorer2Client)

	r := &ResourceExplorer2Index{
		svc: mockClient,
		ARN: ptr.String("arn:aws:resource-explorer-2:us-east-1:123456789012:index/12345"),
	}

	mockClient.On("DeleteIndex", mock.Anything,
		&resourceexplorer2.DeleteIndexInput{
			Arn: r.ARN,
		}).Return(&resourceexplorer2.DeleteIndexOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_ResourceExplorer2Index_Properties(t *testing.T) {
	a := assert.New(t)
	r := &ResourceExplorer2Index{
		ARN:    ptr.String("arn:aws:resource-explorer-2:us-east-1:123456789012:index/12345"),
		Region: ptr.String("us-east-1"),
		Type:   types.IndexTypeLocal,
		Tags:   map[string]string{"env": "test"},
	}
	props := r.Properties()
	a.Equal("arn:aws:resource-explorer-2:us-east-1:123456789012:index/12345", props.Get("ARN"))
	a.Equal("us-east-1", props.Get("Region"))
	a.Equal("LOCAL", props.Get("Type"))
	a.Equal("test", props.Get("tag:env"))
}

func Test_Mock_ResourceExplorer2Index_String(t *testing.T) {
	a := assert.New(t)
	r := &ResourceExplorer2Index{
		ARN: ptr.String("arn:aws:resource-explorer-2:us-east-1:123456789012:index/12345"),
	}
	a.Equal("arn:aws:resource-explorer-2:us-east-1:123456789012:index/12345", r.String())
}
