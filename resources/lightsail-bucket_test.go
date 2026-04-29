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
)

type TestLightsailBucketSuite struct {
	suite.Suite
	svc  *lightsail.Client
	name *string
}

func (s *TestLightsailBucketSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = lightsail.NewFromConfig(cfg)

	s.name = ptr.String(fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano()))
	_, err = s.svc.CreateBucket(ctx, &lightsail.CreateBucketInput{
		BucketName: s.name,
		BundleId:   ptr.String("small_1_0"),
	})
	if err != nil {
		s.T().Fatalf("failed to create bucket: %v", err)
	}
}

func (s *TestLightsailBucketSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = s.svc.DeleteBucket(ctx, &lightsail.DeleteBucketInput{
		BucketName: s.name,
	})
}

func (s *TestLightsailBucketSuite) TestList() {
	a := assert.New(s.T())
	lister := &LightsailBucketLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testLightsailListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (s *TestLightsailBucketSuite) TestRemove() {
	a := assert.New(s.T())
	b := &LightsailBucket{svc: s.svc, BucketName: s.name}
	a.NoError(b.Remove(context.TODO()))
}

func TestLightsailBucketIntegration(t *testing.T) {
	suite.Run(t, new(TestLightsailBucketSuite))
}
