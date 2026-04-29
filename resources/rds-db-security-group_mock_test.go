package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/rds"
	rdstypes "github.com/aws/aws-sdk-go-v2/service/rds/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

var testRDSV2ListerOpts = &nuke.ListerOpts{
	Logger: logrus.WithField("test", true),
}

func Test_Mock_RDSDBSecurityGroup_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockRDSV2Client)

	mockClient.On("DescribeDBSecurityGroups", mock.Anything, mock.Anything).
		Return(&rds.DescribeDBSecurityGroupsOutput{
			DBSecurityGroups: []rdstypes.DBSecurityGroup{
				{
					DBSecurityGroupName: ptr.String("my-sg"),
					DBSecurityGroupArn:  ptr.String("arn:aws:rds:us-east-1:123456789012:secgrp:my-sg"),
				},
			},
		}, nil)

	lister := &RDSDBSecurityGroupLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testRDSV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	securityGroup := resources[0].(*RDSDBSecurityGroup)
	assertions.Equal("my-sg", *securityGroup.DBSecurityGroupName)
	assertions.Equal("arn:aws:rds:us-east-1:123456789012:secgrp:my-sg", *securityGroup.DBSecurityGroupArn)
	mockClient.AssertExpectations(t)
}

func Test_Mock_RDSDBSecurityGroup_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockRDSV2Client)

	mockClient.On("DescribeDBSecurityGroups", mock.Anything, mock.Anything).
		Return(&rds.DescribeDBSecurityGroupsOutput{
			DBSecurityGroups: []rdstypes.DBSecurityGroup{},
		}, nil)

	lister := &RDSDBSecurityGroupLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testRDSV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_RDSDBSecurityGroup_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockRDSV2Client)

	securityGroup := &RDSDBSecurityGroup{
		svc:                 mockClient,
		DBSecurityGroupName: ptr.String("my-sg"),
	}

	mockClient.On("DeleteDBSecurityGroup", mock.Anything, &rds.DeleteDBSecurityGroupInput{
		DBSecurityGroupName: securityGroup.DBSecurityGroupName,
	}).Return(&rds.DeleteDBSecurityGroupOutput{}, nil)

	err := securityGroup.Remove(context.TODO())
	assertions.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_RDSDBSecurityGroup_Properties(t *testing.T) {
	assertions := assert.New(t)

	securityGroup := RDSDBSecurityGroup{
		DBSecurityGroupName: ptr.String("my-sg"),
		DBSecurityGroupArn:  ptr.String("arn:aws:rds:us-east-1:123456789012:secgrp:my-sg"),
	}

	properties := securityGroup.Properties()
	assertions.Equal("my-sg", properties.Get("DBSecurityGroupName"))
	assertions.Equal("arn:aws:rds:us-east-1:123456789012:secgrp:my-sg", properties.Get("DBSecurityGroupArn"))
}

func Test_Mock_RDSDBSecurityGroup_String(t *testing.T) {
	assertions := assert.New(t)
	securityGroup := RDSDBSecurityGroup{DBSecurityGroupName: ptr.String("my-sg")}
	assertions.Equal("my-sg", securityGroup.String())
}

func Test_Mock_RDSDBSecurityGroup_Filter(t *testing.T) {
	assertions := assert.New(t)

	defaultSG := RDSDBSecurityGroup{DBSecurityGroupName: ptr.String("default")}
	assertions.Error(defaultSG.Filter())

	customSG := RDSDBSecurityGroup{DBSecurityGroupName: ptr.String("my-sg")}
	assertions.NoError(customSG.Filter())
}
