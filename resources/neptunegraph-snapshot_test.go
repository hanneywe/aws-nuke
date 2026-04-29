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
	"github.com/aws/aws-sdk-go-v2/service/neptunegraph"
)

type TestNeptuneGraphSnapshotSuite struct {
	suite.Suite
	svc        *neptunegraph.Client
	snapshotID *string
	graphID    *string
}

func (s *TestNeptuneGraphSnapshotSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = neptunegraph.NewFromConfig(cfg)

	graphName := fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano())
	graphResp, err := s.svc.CreateGraph(ctx, &neptunegraph.CreateGraphInput{
		GraphName:          ptr.String(graphName),
		ProvisionedMemory:  ptr.Int32(128),
		DeletionProtection: ptr.Bool(false),
	})
	if err != nil {
		s.T().Fatalf("failed to create graph: %v", err)
	}
	s.graphID = graphResp.Id

	snapshotName := fmt.Sprintf("aws-nuke-snap-%d", time.Now().UnixNano())
	snapResp, err := s.svc.CreateGraphSnapshot(ctx, &neptunegraph.CreateGraphSnapshotInput{
		GraphIdentifier: s.graphID,
		SnapshotName:    ptr.String(snapshotName),
	})
	if err != nil {
		s.T().Fatalf("failed to create snapshot: %v", err)
	}
	s.snapshotID = snapResp.Id
}

func (s *TestNeptuneGraphSnapshotSuite) TearDownSuite() {
	ctx := context.TODO()
	if s.snapshotID != nil {
		_, _ = s.svc.DeleteGraphSnapshot(ctx, &neptunegraph.DeleteGraphSnapshotInput{
			SnapshotIdentifier: s.snapshotID,
		})
	}
	if s.graphID != nil {
		_, _ = s.svc.DeleteGraph(ctx, &neptunegraph.DeleteGraphInput{
			GraphIdentifier: s.graphID,
			SkipSnapshot:    ptr.Bool(true),
		})
	}
}

func (s *TestNeptuneGraphSnapshotSuite) TestList() {
	a := assert.New(s.T())
	lister := &NeptuneGraphSnapshotLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testNeptuneGraphListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (s *TestNeptuneGraphSnapshotSuite) TestRemove() {
	a := assert.New(s.T())
	snap := &NeptuneGraphSnapshot{svc: s.svc, SnapshotID: s.snapshotID}
	a.NoError(snap.Remove(context.TODO()))
}

func TestNeptuneGraphSnapshotIntegration(t *testing.T) {
	suite.Run(t, new(TestNeptuneGraphSnapshotSuite))
}
