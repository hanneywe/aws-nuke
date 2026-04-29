//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/configservice"
)

type TestConfigServiceStoredQuerySuite struct {
	suite.Suite
	svc *configservice.Client
}

func (s *TestConfigServiceStoredQuerySuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}
	s.svc = configservice.NewFromConfig(cfg)
}

func (s *TestConfigServiceStoredQuerySuite) TestList() {
	a := assert.New(s.T())
	lister := ConfigServiceStoredQueryLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testConfigServiceListerOpts)
	a.NoError(err)
	_ = resources
}

func TestConfigServiceStoredQueryIntegration(t *testing.T) {
	suite.Run(t, new(TestConfigServiceStoredQuerySuite))
}
