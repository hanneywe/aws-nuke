//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/transfer"
)

type TestTransferProfileSuite struct {
	suite.Suite
	svc *transfer.Client
}

func (s *TestTransferProfileSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = transfer.NewFromConfig(cfg)
}

func (s *TestTransferProfileSuite) TestList() {
	assertions := assert.New(s.T())

	lister := TransferProfileLister{}
	resources, err := lister.List(context.TODO(), testTransferListerOpts)

	assertions.Nil(err)
	assertions.NotNil(resources)
}

func TestTransferProfileIntegration(t *testing.T) {
	suite.Run(t, new(TestTransferProfileSuite))
}
