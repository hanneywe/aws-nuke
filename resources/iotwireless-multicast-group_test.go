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
	"github.com/aws/aws-sdk-go-v2/service/iotwireless"
)

type TestIoTWirelessMulticastGroupSuite struct {
	suite.Suite
	svc *iotwireless.Client
	id  *string
}

func (s *TestIoTWirelessMulticastGroupSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = iotwireless.NewFromConfig(cfg)

	name := fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano())
	resp, err := s.svc.CreateMulticastGroup(ctx, &iotwireless.CreateMulticastGroupInput{
		Name: ptr.String(name),
	})
	if err != nil {
		s.T().Fatalf("failed to create multicast group: %v", err)
	}
	s.id = resp.Id
}

func (s *TestIoTWirelessMulticastGroupSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = s.svc.DeleteMulticastGroup(ctx, &iotwireless.DeleteMulticastGroupInput{
		Id: s.id,
	})
}

func (s *TestIoTWirelessMulticastGroupSuite) TestList() {
	a := assert.New(s.T())
	lister := &IoTWirelessMulticastGroupLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testIoTWirelessListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (s *TestIoTWirelessMulticastGroupSuite) TestRemove() {
	a := assert.New(s.T())
	mg := &IoTWirelessMulticastGroup{svc: s.svc, Id: s.id}
	a.NoError(mg.Remove(context.TODO()))
}

func TestIoTWirelessMulticastGroupIntegration(t *testing.T) {
	suite.Run(t, new(TestIoTWirelessMulticastGroupSuite))
}
