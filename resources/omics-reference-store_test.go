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

type TestOmicsReferenceStoreSuite struct {
	suite.Suite
	svc              *omics.Client
	referenceStoreID *string
}

func (s *TestOmicsReferenceStoreSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = omics.NewFromConfig(cfg)

	storeName := ptr.String(fmt.Sprintf("aws-nuke-test-refstore-%d", time.Now().UnixNano()))
	createOutput, err := s.svc.CreateReferenceStore(ctx, &omics.CreateReferenceStoreInput{
		Name: storeName,
	})
	if err != nil {
		s.T().Fatalf("failed to create reference store: %v", err)
	}
	s.referenceStoreID = createOutput.Id
}

func (s *TestOmicsReferenceStoreSuite) TearDownSuite() {
	ctx := context.TODO()
	if s.referenceStoreID != nil {
		_, _ = s.svc.DeleteReferenceStore(ctx, &omics.DeleteReferenceStoreInput{
			Id: s.referenceStoreID,
		})
	}
}

func (s *TestOmicsReferenceStoreSuite) TestList() {
	assertions := assert.New(s.T())

	lister := OmicsReferenceStoreLister{}
	resources, err := lister.List(context.TODO(), testOmicsListerOpts)

	assertions.Nil(err)
	assertions.Greater(len(resources), 0)
}

func (s *TestOmicsReferenceStoreSuite) TestRemove() {
	assertions := assert.New(s.T())

	referenceStore := OmicsReferenceStore{
		svc: s.svc,
		ID:  s.referenceStoreID,
	}

	err := referenceStore.Remove(context.TODO())
	assertions.Nil(err)
}

func TestOmicsReferenceStoreIntegration(t *testing.T) {
	suite.Run(t, new(TestOmicsReferenceStoreSuite))
}
