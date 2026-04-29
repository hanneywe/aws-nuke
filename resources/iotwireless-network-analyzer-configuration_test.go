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

type TestIoTWirelessNetworkAnalyzerConfigurationSuite struct {
	suite.Suite
	svc  *iotwireless.Client
	name *string
}

func (s *TestIoTWirelessNetworkAnalyzerConfigurationSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = iotwireless.NewFromConfig(cfg)

	name := fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano())
	_, err = s.svc.CreateNetworkAnalyzerConfiguration(ctx, &iotwireless.CreateNetworkAnalyzerConfigurationInput{
		Name: ptr.String(name),
	})
	if err != nil {
		s.T().Fatalf("failed to create network analyzer configuration: %v", err)
	}
	s.name = ptr.String(name)
}

func (s *TestIoTWirelessNetworkAnalyzerConfigurationSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = s.svc.DeleteNetworkAnalyzerConfiguration(ctx, &iotwireless.DeleteNetworkAnalyzerConfigurationInput{
		ConfigurationName: s.name,
	})
}

func (s *TestIoTWirelessNetworkAnalyzerConfigurationSuite) TestList() {
	a := assert.New(s.T())
	lister := &IoTWirelessNetworkAnalyzerConfigurationLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testIoTWirelessListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (s *TestIoTWirelessNetworkAnalyzerConfigurationSuite) TestRemove() {
	a := assert.New(s.T())
	cfg := &IoTWirelessNetworkAnalyzerConfiguration{svc: s.svc, Name: s.name}
	a.NoError(cfg.Remove(context.TODO()))
}

func TestIoTWirelessNetworkAnalyzerConfigurationIntegration(t *testing.T) {
	suite.Run(t, new(TestIoTWirelessNetworkAnalyzerConfigurationSuite))
}
