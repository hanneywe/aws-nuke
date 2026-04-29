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

type TestIoTWirelessServiceProfileSuite struct {
	suite.Suite
	svc *iotwireless.Client
	id  *string
}

func (s *TestIoTWirelessServiceProfileSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = iotwireless.NewFromConfig(cfg)

	name := fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano())
	resp, err := s.svc.CreateServiceProfile(ctx, &iotwireless.CreateServiceProfileInput{
		Name: ptr.String(name),
	})
	if err != nil {
		s.T().Fatalf("failed to create service profile: %v", err)
	}
	s.id = resp.Id
}

func (s *TestIoTWirelessServiceProfileSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = s.svc.DeleteServiceProfile(ctx, &iotwireless.DeleteServiceProfileInput{
		Id: s.id,
	})
}

func (s *TestIoTWirelessServiceProfileSuite) TestList() {
	a := assert.New(s.T())
	lister := &IoTWirelessServiceProfileLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testIoTWirelessListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (s *TestIoTWirelessServiceProfileSuite) TestRemove() {
	a := assert.New(s.T())
	sp := &IoTWirelessServiceProfile{svc: s.svc, Id: s.id}
	a.NoError(sp.Remove(context.TODO()))
}

func TestIoTWirelessServiceProfileIntegration(t *testing.T) {
	suite.Run(t, new(TestIoTWirelessServiceProfileSuite))
}
