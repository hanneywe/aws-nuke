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
	"github.com/aws/aws-sdk-go-v2/service/omics"
)

type TestOmicsReadSetSuite struct {
	suite.Suite
	svc             *omics.Client
	sequenceStoreID *string
}

func (s *TestOmicsReadSetSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = omics.NewFromConfig(cfg)

	storeName := ptr.String(fmt.Sprintf("aws-nuke-test-readset-%d", time.Now().UnixNano()))
	createOutput, err := s.svc.CreateSequenceStore(ctx, &omics.CreateSequenceStoreInput{
		Name: storeName,
	})
	if err != nil {
		s.T().Fatalf("failed to create sequence store: %v", err)
	}
	s.sequenceStoreID = createOutput.Id
}

func (s *TestOmicsReadSetSuite) TearDownSuite() {
	ctx := context.TODO()
	if s.sequenceStoreID != nil {
		_, _ = s.svc.DeleteSequenceStore(ctx, &omics.DeleteSequenceStoreInput{
			Id: s.sequenceStoreID,
		})
	}
}

func (s *TestOmicsReadSetSuite) TestList() {
	assertions := assert.New(s.T())

	lister := OmicsReadSetLister{}
	resources, err := lister.List(context.TODO(), testOmicsListerOpts)

	assertions.Nil(err)
	assertions.NotNil(resources)
}

func TestOmicsReadSetIntegration(t *testing.T) {
	suite.Run(t, new(TestOmicsReadSetSuite))
}
