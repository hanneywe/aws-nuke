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
	"github.com/aws/aws-sdk-go-v2/service/redshift"
)

type TestRedshiftSnapshotCopyGrantSuite struct {
	suite.Suite
	svc  *redshift.Client
	name *string
}

func (s *TestRedshiftSnapshotCopyGrantSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = redshift.NewFromConfig(cfg)

	name := fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano())
	_, err = s.svc.CreateSnapshotCopyGrant(ctx, &redshift.CreateSnapshotCopyGrantInput{
		SnapshotCopyGrantName: ptr.String(name),
	})
	if err != nil {
		s.T().Fatalf("failed to create snapshot copy grant: %v", err)
	}
	s.name = ptr.String(name)
}

func (s *TestRedshiftSnapshotCopyGrantSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = s.svc.DeleteSnapshotCopyGrant(ctx, &redshift.DeleteSnapshotCopyGrantInput{
		SnapshotCopyGrantName: s.name,
	})
}

func (s *TestRedshiftSnapshotCopyGrantSuite) TestList() {
	a := assert.New(s.T())
	lister := &RedshiftSnapshotCopyGrantLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testRedshiftListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (s *TestRedshiftSnapshotCopyGrantSuite) TestRemove() {
	a := assert.New(s.T())
	grant := &RedshiftSnapshotCopyGrant{svc: s.svc, SnapshotCopyGrantName: s.name}
	a.NoError(grant.Remove(context.TODO()))
}

func TestRedshiftSnapshotCopyGrantIntegration(t *testing.T) {
	suite.Run(t, new(TestRedshiftSnapshotCopyGrantSuite))
}
