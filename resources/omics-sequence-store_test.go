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

type TestOmicsSequenceStoreSuite struct {
	suite.Suite
	svc             *omics.Client
	sequenceStoreID *string
}

func (s *TestOmicsSequenceStoreSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = omics.NewFromConfig(cfg)

	storeName := ptr.String(fmt.Sprintf("aws-nuke-test-seqstore-%d", time.Now().UnixNano()))
	createOutput, err := s.svc.CreateSequenceStore(ctx, &omics.CreateSequenceStoreInput{
		Name: storeName,
	})
	if err != nil {
		s.T().Fatalf("failed to create sequence store: %v", err)
	}
	s.sequenceStoreID = createOutput.Id
}

func (s *TestOmicsSequenceStoreSuite) TearDownSuite() {
	ctx := context.TODO()
	if s.sequenceStoreID != nil {
		_, _ = s.svc.DeleteSequenceStore(ctx, &omics.DeleteSequenceStoreInput{
			Id: s.sequenceStoreID,
		})
	}
}

func (s *TestOmicsSequenceStoreSuite) TestList() {
	assertions := assert.New(s.T())

	lister := OmicsSequenceStoreLister{}
	resources, err := lister.List(context.TODO(), testOmicsListerOpts)

	assertions.Nil(err)
	assertions.Greater(len(resources), 0)
}

func (s *TestOmicsSequenceStoreSuite) TestRemove() {
	assertions := assert.New(s.T())

	sequenceStore := OmicsSequenceStore{
		svc: s.svc,
		ID:  s.sequenceStoreID,
	}

	err := sequenceStore.Remove(context.TODO())
	assertions.Nil(err)
}

func TestOmicsSequenceStoreIntegration(t *testing.T) {
	suite.Run(t, new(TestOmicsSequenceStoreSuite))
}
