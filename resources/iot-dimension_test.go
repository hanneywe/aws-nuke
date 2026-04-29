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
	"github.com/aws/aws-sdk-go-v2/service/iot"
	iottypes "github.com/aws/aws-sdk-go-v2/service/iot/types"
)

type TestIoTDimensionSuite struct {
	suite.Suite
	svc  *iot.Client
	name *string
}

func (s *TestIoTDimensionSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = iot.NewFromConfig(cfg)

	s.name = ptr.String(fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano()))
	_, err = s.svc.CreateDimension(ctx, &iot.CreateDimensionInput{
		Name:               s.name,
		Type:               iottypes.DimensionType("TOPIC_FILTER"),
		StringValues:       []string{"test/topic"},
		ClientRequestToken: ptr.String(fmt.Sprintf("token-%d", time.Now().UnixNano())),
	})
	if err != nil {
		s.T().Fatalf("failed to create dimension: %v", err)
	}
}

func (s *TestIoTDimensionSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = s.svc.DeleteDimension(ctx, &iot.DeleteDimensionInput{
		Name: s.name,
	})
}

func (s *TestIoTDimensionSuite) TestList() {
	assertions := assert.New(s.T())
	lister := &IoTDimensionLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testIoTListerOpts)
	assertions.NoError(err)
	assertions.Greater(len(resources), 0)
}

func (s *TestIoTDimensionSuite) TestRemove() {
	assertions := assert.New(s.T())
	dimension := &IoTDimension{svc: s.svc, Name: s.name}
	assertions.NoError(dimension.Remove(context.TODO()))
}

func TestIoTDimensionIntegration(t *testing.T) {
	suite.Run(t, new(TestIoTDimensionSuite))
}
