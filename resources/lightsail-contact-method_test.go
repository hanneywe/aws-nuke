//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	lstypes "github.com/aws/aws-sdk-go-v2/service/lightsail/types"
)

type TestLightsailContactMethodSuite struct {
	suite.Suite
	svc *lightsail.Client
}

func (s *TestLightsailContactMethodSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = lightsail.NewFromConfig(cfg)

	_, err = s.svc.CreateContactMethod(ctx, &lightsail.CreateContactMethodInput{
		Protocol:        lstypes.ContactProtocolEmail,
		ContactEndpoint: ptr.String("aws-nuke-test@example.com"),
	})
	if err != nil {
		s.T().Fatalf("failed to create contact method: %v", err)
	}
}

func (s *TestLightsailContactMethodSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = s.svc.DeleteContactMethod(ctx, &lightsail.DeleteContactMethodInput{
		Protocol: lstypes.ContactProtocolEmail,
	})
}

func (s *TestLightsailContactMethodSuite) TestList() {
	a := assert.New(s.T())
	lister := &LightsailContactMethodLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testLightsailListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (s *TestLightsailContactMethodSuite) TestRemove() {
	a := assert.New(s.T())
	cm := &LightsailContactMethod{
		svc:             s.svc,
		Protocol:        ptr.String("Email"),
		ContactEndpoint: ptr.String("aws-nuke-test@example.com"),
	}
	a.NoError(cm.Remove(context.TODO()))
}

func TestLightsailContactMethodIntegration(t *testing.T) {
	suite.Run(t, new(TestLightsailContactMethodSuite))
}
