//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ivschat"
)

type TestIVSChatLoggingConfigurationSuite struct {
	suite.Suite
	svc *ivschat.Client
}

func (s *TestIVSChatLoggingConfigurationSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}
	s.svc = ivschat.NewFromConfig(cfg)
}

func (s *TestIVSChatLoggingConfigurationSuite) TestList() {
	assertions := assert.New(s.T())
	lister := IVSChatLoggingConfigurationLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testIVSChatListerOpts)
	assertions.NoError(err)
	_ = resources
}

func TestIVSChatLoggingConfigurationIntegration(t *testing.T) {
	suite.Run(t, new(TestIVSChatLoggingConfigurationSuite))
}
