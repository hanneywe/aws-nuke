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

type TestLightsailCertificateSuite struct {
	suite.Suite
	svc  *lightsail.Client
	name *string
}

func (s *TestLightsailCertificateSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = lightsail.NewFromConfig(cfg)

	s.name = ptr.String(fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano()))
	_, err = s.svc.CreateCertificate(ctx, &lightsail.CreateCertificateInput{
		CertificateName: s.name,
		DomainName:      ptr.String("example.com"),
	})
	if err != nil {
		s.T().Fatalf("failed to create certificate: %v", err)
	}
}

func (s *TestLightsailCertificateSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = s.svc.DeleteCertificate(ctx, &lightsail.DeleteCertificateInput{
		CertificateName: s.name,
	})
}

func (s *TestLightsailCertificateSuite) TestList() {
	a := assert.New(s.T())
	lister := &LightsailCertificateLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testLightsailListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (s *TestLightsailCertificateSuite) TestRemove() {
	a := assert.New(s.T())
	c := &LightsailCertificate{svc: s.svc, CertificateName: s.name}
	a.NoError(c.Remove(context.TODO()))
}

func TestLightsailCertificateIntegration(t *testing.T) {
	suite.Run(t, new(TestLightsailCertificateSuite))
}
