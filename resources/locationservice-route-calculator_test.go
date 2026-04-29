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
	"github.com/aws/aws-sdk-go-v2/service/location"
)

type TestLocationServiceRouteCalculatorSuite struct {
	suite.Suite
	svc  *location.Client
	name *string
}

func (s *TestLocationServiceRouteCalculatorSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = location.NewFromConfig(cfg)

	name := fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano())
	_, err = s.svc.CreateRouteCalculator(ctx, &location.CreateRouteCalculatorInput{
		CalculatorName: ptr.String(name),
		DataSource:     ptr.String("Esri"),
	})
	if err != nil {
		s.T().Fatalf("failed to create route calculator: %v", err)
	}
	s.name = ptr.String(name)
}

func (s *TestLocationServiceRouteCalculatorSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = s.svc.DeleteRouteCalculator(ctx, &location.DeleteRouteCalculatorInput{CalculatorName: s.name})
}

func (s *TestLocationServiceRouteCalculatorSuite) TestList() {
	a := assert.New(s.T())
	lister := &LocationServiceRouteCalculatorLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testLocationServiceListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (s *TestLocationServiceRouteCalculatorSuite) TestRemove() {
	a := assert.New(s.T())
	rc := &LocationServiceRouteCalculator{svc: s.svc, CalculatorName: s.name}
	a.NoError(rc.Remove(context.TODO()))
}

func TestLocationServiceRouteCalculatorIntegration(t *testing.T) {
	suite.Run(t, new(TestLocationServiceRouteCalculatorSuite))
}
