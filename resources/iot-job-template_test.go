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
	"github.com/aws/aws-sdk-go-v2/service/iot"
)

type TestIoTJobTemplateSuite struct {
	suite.Suite
	svc *iot.Client
	id  *string
}

func (s *TestIoTJobTemplateSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = iot.NewFromConfig(cfg)

	jtId := fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano())
	_, err = s.svc.CreateJobTemplate(ctx, &iot.CreateJobTemplateInput{
		JobTemplateId: ptr.String(jtId),
		Description:   ptr.String("aws-nuke integration test"),
	})
	if err != nil {
		s.T().Fatalf("failed to create job template: %v", err)
	}
	s.id = ptr.String(jtId)
}

func (s *TestIoTJobTemplateSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = s.svc.DeleteJobTemplate(ctx, &iot.DeleteJobTemplateInput{
		JobTemplateId: s.id,
	})
}

func (s *TestIoTJobTemplateSuite) TestList() {
	a := assert.New(s.T())
	lister := &IoTJobTemplateLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testIoTListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (s *TestIoTJobTemplateSuite) TestRemove() {
	a := assert.New(s.T())
	jt := &IoTJobTemplate{svc: s.svc, JobTemplateId: s.id}
	a.NoError(jt.Remove(context.TODO()))
}

func TestIoTJobTemplateIntegration(t *testing.T) {
	suite.Run(t, new(TestIoTJobTemplateSuite))
}
