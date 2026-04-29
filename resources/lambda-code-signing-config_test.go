//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

type TestLambdaCodeSigningConfigSuite struct {
	suite.Suite
	svc *lambda.Client
}

func (s *TestLambdaCodeSigningConfigSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = lambda.NewFromConfig(cfg)
}

func (s *TestLambdaCodeSigningConfigSuite) TestList() {
	assertions := assert.New(s.T())

	lister := LambdaCodeSigningConfigLister{}
	resources, err := lister.List(context.TODO(), testLambdaListerOpts)

	assertions.Nil(err)
	assertions.NotNil(resources)
}

func TestLambdaCodeSigningConfigIntegration(t *testing.T) {
	suite.Run(t, new(TestLambdaCodeSigningConfigSuite))
}
