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
	"github.com/aws/aws-sdk-go-v2/service/keyspaces"
)

type TestKeyspacesKeyspaceSuite struct {
	suite.Suite
	svc  *keyspaces.Client
	name *string
}

func (s *TestKeyspacesKeyspaceSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = keyspaces.NewFromConfig(cfg)

	name := fmt.Sprintf("aws_nuke_test_%d", time.Now().UnixNano())
	_, err = s.svc.CreateKeyspace(ctx, &keyspaces.CreateKeyspaceInput{
		KeyspaceName: ptr.String(name),
	})
	if err != nil {
		s.T().Fatalf("failed to create keyspace: %v", err)
	}
	s.name = ptr.String(name)
}

func (s *TestKeyspacesKeyspaceSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = s.svc.DeleteKeyspace(ctx, &keyspaces.DeleteKeyspaceInput{
		KeyspaceName: s.name,
	})
}

func (s *TestKeyspacesKeyspaceSuite) TestList() {
	a := assert.New(s.T())
	lister := &KeyspacesKeyspaceLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testKeyspacesListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (s *TestKeyspacesKeyspaceSuite) TestRemove() {
	a := assert.New(s.T())
	ks := &KeyspacesKeyspace{svc: s.svc, KeyspaceName: s.name}
	a.NoError(ks.Remove(context.TODO()))
}

func TestKeyspacesKeyspaceIntegration(t *testing.T) {
	suite.Run(t, new(TestKeyspacesKeyspaceSuite))
}
