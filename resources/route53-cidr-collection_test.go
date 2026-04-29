//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
)

type TestRoute53CidrCollectionSuite struct {
	suite.Suite
	svc *route53.Client
}

func (s *TestRoute53CidrCollectionSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = route53.NewFromConfig(cfg)
}

func (s *TestRoute53CidrCollectionSuite) TestList() {
	assertions := assert.New(s.T())

	lister := Route53CidrCollectionLister{}
	resources, err := lister.List(context.TODO(), testRoute53ListerOpts)

	assertions.Nil(err)
	assertions.NotNil(resources)
}

func TestRoute53CidrCollectionIntegration(t *testing.T) {
	suite.Run(t, new(TestRoute53CidrCollectionSuite))
}
