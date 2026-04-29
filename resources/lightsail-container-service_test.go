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
	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	lstypes "github.com/aws/aws-sdk-go-v2/service/lightsail/types"
)

type TestLightsailContainerServiceSuite struct {
	suite.Suite
	svc  *lightsail.Client
	name *string
}

func (s *TestLightsailContainerServiceSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = lightsail.NewFromConfig(cfg)

	s.name = ptr.String(fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano()))
	_, err = s.svc.CreateContainerService(ctx, &lightsail.CreateContainerServiceInput{
		ServiceName: s.name,
		Power:       lstypes.ContainerServicePowerNameNano,
		Scale:       ptr.Int32(1),
	})
	if err != nil {
		s.T().Fatalf("failed to create container service: %v", err)
	}
}

func (s *TestLightsailContainerServiceSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = s.svc.DeleteContainerService(ctx, &lightsail.DeleteContainerServiceInput{
		ServiceName: s.name,
	})
}

func (s *TestLightsailContainerServiceSuite) TestList() {
	a := assert.New(s.T())
	lister := &LightsailContainerServiceLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testLightsailListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (s *TestLightsailContainerServiceSuite) TestRemove() {
	a := assert.New(s.T())
	cs := &LightsailContainerService{svc: s.svc, ContainerServiceName: s.name}
	a.NoError(cs.Remove(context.TODO()))
}

func TestLightsailContainerServiceIntegration(t *testing.T) {
	suite.Run(t, new(TestLightsailContainerServiceSuite))
}
