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

type TestRedshiftHsmClientCertificateSuite struct {
	suite.Suite
	svc  *redshift.Client
	name *string
}

func (s *TestRedshiftHsmClientCertificateSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = redshift.NewFromConfig(cfg)

	name := fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano())
	_, err = s.svc.CreateHsmClientCertificate(ctx, &redshift.CreateHsmClientCertificateInput{
		HsmClientCertificateIdentifier: ptr.String(name),
	})
	if err != nil {
		s.T().Fatalf("failed to create HSM client certificate: %v", err)
	}
	s.name = ptr.String(name)
}

func (s *TestRedshiftHsmClientCertificateSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = s.svc.DeleteHsmClientCertificate(ctx, &redshift.DeleteHsmClientCertificateInput{
		HsmClientCertificateIdentifier: s.name,
	})
}

func (s *TestRedshiftHsmClientCertificateSuite) TestList() {
	a := assert.New(s.T())
	lister := &RedshiftHsmClientCertificateLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testRedshiftListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (s *TestRedshiftHsmClientCertificateSuite) TestRemove() {
	a := assert.New(s.T())
	cert := &RedshiftHsmClientCertificate{svc: s.svc, HsmClientCertificateIdentifier: s.name}
	a.NoError(cert.Remove(context.TODO()))
}

func TestRedshiftHsmClientCertificateIntegration(t *testing.T) {
	suite.Run(t, new(TestRedshiftHsmClientCertificateSuite))
}
