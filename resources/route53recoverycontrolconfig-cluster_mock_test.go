package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/route53recoverycontrolconfig"
	r53rcctypes "github.com/aws/aws-sdk-go-v2/service/route53recoverycontrolconfig/types"
)

func Test_Mock_Route53RecoveryControlConfigCluster_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockRoute53RecoveryControlConfigClient)
	mockClient.On("ListClusters", mock.Anything, mock.Anything).
		Return(&route53recoverycontrolconfig.ListClustersOutput{
			Clusters: []r53rcctypes.Cluster{
				{ClusterArn: ptr.String("arn:aws:route53-recovery-control::123456789012:cluster/my-cluster"), Name: ptr.String("my-cluster")},
			},
		}, nil)
	lister := &Route53RecoveryControlConfigClusterLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testRoute53RecoveryControlConfigListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	c := resources[0].(*Route53RecoveryControlConfigCluster)
	a.Equal("my-cluster", *c.Name)
	mockClient.AssertExpectations(t)
}

func Test_Mock_Route53RecoveryControlConfigCluster_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockRoute53RecoveryControlConfigClient)
	mockClient.On("ListClusters", mock.Anything, mock.Anything).
		Return(&route53recoverycontrolconfig.ListClustersOutput{Clusters: []r53rcctypes.Cluster{}}, nil)
	lister := &Route53RecoveryControlConfigClusterLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testRoute53RecoveryControlConfigListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_Route53RecoveryControlConfigCluster_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockRoute53RecoveryControlConfigClient)
	clusterArn := "arn:aws:route53-recovery-control::123456789012:cluster/my-cluster"
	c := &Route53RecoveryControlConfigCluster{
		svc:        mockClient,
		ClusterArn: ptr.String(clusterArn),
	}
	mockClient.On("DeleteCluster", mock.Anything,
		&route53recoverycontrolconfig.DeleteClusterInput{ClusterArn: c.ClusterArn}).
		Return(&route53recoverycontrolconfig.DeleteClusterOutput{}, nil)
	a.NoError(c.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_Route53RecoveryControlConfigCluster_Properties(t *testing.T) {
	a := assert.New(t)
	c := Route53RecoveryControlConfigCluster{
		ClusterArn: ptr.String("arn:aws:route53-recovery-control::123456789012:cluster/my-cluster"),
		Name:       ptr.String("my-cluster"),
	}
	a.Equal("my-cluster", c.Properties().Get("Name"))
}

func Test_Mock_Route53RecoveryControlConfigCluster_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("my-cluster", (&Route53RecoveryControlConfigCluster{Name: ptr.String("my-cluster")}).String())
}
