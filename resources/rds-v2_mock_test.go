package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/rds"
)

type mockRDSV2Client struct {
	mock.Mock
}

func (m *mockRDSV2Client) DescribeDBSecurityGroups(ctx context.Context, params *rds.DescribeDBSecurityGroupsInput,
	_ ...func(*rds.Options)) (*rds.DescribeDBSecurityGroupsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*rds.DescribeDBSecurityGroupsOutput), args.Error(1)
}

func (m *mockRDSV2Client) DeleteDBSecurityGroup(ctx context.Context, params *rds.DeleteDBSecurityGroupInput,
	_ ...func(*rds.Options)) (*rds.DeleteDBSecurityGroupOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*rds.DeleteDBSecurityGroupOutput), args.Error(1)
}
