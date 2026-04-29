//go:build integration

package resources

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iot"
)

type TestIoTBillingGroupSuite struct {
	suite.Suite
	svc  *iot.Client
	name *string
}

func (s *TestIoTBillingGroupSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = iot.NewFromConfig(cfg)

	s.name = ptr.String(fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano()))
	_, err = s.svc.CreateBillingGroup(ctx, &iot.CreateBillingGroupInput{
		BillingGroupName: s.name,
	})
	if err != nil {
		s.T().Fatalf("failed to create billing group: %v", err)
	}
}

func (s *TestIoTBillingGroupSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = s.svc.DeleteBillingGroup(ctx, &iot.DeleteBillingGroupInput{
		BillingGroupName: s.name,
	})
}

func (s *TestIoTBillingGroupSuite) TestList() {
	a := assert.New(s.T())
	lister := &IoTBillingGroupLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testIoTListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (s *TestIoTBillingGroupSuite) TestRemove() {
	a := assert.New(s.T())
	bg := &IoTBillingGroup{svc: s.svc, BillingGroupName: s.name}
	a.NoError(bg.Remove(context.TODO()))
}

func TestIoTBillingGroupIntegration(t *testing.T) {
	suite.Run(t, new(TestIoTBillingGroupSuite))
}
