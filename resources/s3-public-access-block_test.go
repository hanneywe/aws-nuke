//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3control"
)

type TestS3PublicAccessBlockSuite struct {
	suite.Suite
	svc *s3control.Client
}

func (suite *TestS3PublicAccessBlockSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = s3control.NewFromConfig(cfg)
}

func (suite *TestS3PublicAccessBlockSuite) TestList() {
	a := assert.New(suite.T())
	lister := S3PublicAccessBlockLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testS3AccountListerOpts)
	a.NoError(err)
	_ = resources
}

func TestS3PublicAccessBlockIntegration(t *testing.T) {
	suite.Run(t, new(TestS3PublicAccessBlockSuite))
}
