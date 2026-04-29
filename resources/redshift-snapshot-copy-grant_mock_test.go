package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/redshift"
	redshifttypes "github.com/aws/aws-sdk-go-v2/service/redshift/types"
)

func Test_Mock_RedshiftSnapshotCopyGrant_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockRedshiftClient)

	mockClient.On("DescribeSnapshotCopyGrants", mock.Anything, mock.Anything).
		Return(&redshift.DescribeSnapshotCopyGrantsOutput{
			SnapshotCopyGrants: []redshifttypes.SnapshotCopyGrant{
				{
					SnapshotCopyGrantName: ptr.String("my-grant"),
				},
			},
		}, nil)

	lister := &RedshiftSnapshotCopyGrantLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testRedshiftListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	grant := resources[0].(*RedshiftSnapshotCopyGrant)
	a.Equal("my-grant", *grant.SnapshotCopyGrantName)

	mockClient.AssertExpectations(t)
}

func Test_Mock_RedshiftSnapshotCopyGrant_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockRedshiftClient)

	mockClient.On("DescribeSnapshotCopyGrants", mock.Anything, mock.Anything).
		Return(&redshift.DescribeSnapshotCopyGrantsOutput{
			SnapshotCopyGrants: []redshifttypes.SnapshotCopyGrant{},
		}, nil)

	lister := &RedshiftSnapshotCopyGrantLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testRedshiftListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_RedshiftSnapshotCopyGrant_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockRedshiftClient)

	grant := &RedshiftSnapshotCopyGrant{
		svc:                   mockClient,
		SnapshotCopyGrantName: ptr.String("my-grant"),
	}

	mockClient.On("DeleteSnapshotCopyGrant", mock.Anything, &redshift.DeleteSnapshotCopyGrantInput{
		SnapshotCopyGrantName: grant.SnapshotCopyGrantName,
	}).Return(&redshift.DeleteSnapshotCopyGrantOutput{}, nil)

	a.NoError(grant.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_RedshiftSnapshotCopyGrant_Properties(t *testing.T) {
	a := assert.New(t)

	grant := RedshiftSnapshotCopyGrant{
		SnapshotCopyGrantName: ptr.String("my-grant"),
	}

	props := grant.Properties()
	a.Equal("my-grant", props.Get("SnapshotCopyGrantName"))
}

func Test_Mock_RedshiftSnapshotCopyGrant_String(t *testing.T) {
	a := assert.New(t)
	grant := RedshiftSnapshotCopyGrant{SnapshotCopyGrantName: ptr.String("my-grant")}
	a.Equal("my-grant", grant.String())
}
