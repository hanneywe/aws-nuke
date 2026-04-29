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
	"github.com/aws/aws-sdk-go-v2/service/kafkaconnect"
)

type TestKafkaConnectWorkerConfigurationSuite struct {
	suite.Suite
	svc *kafkaconnect.Client
	arn *string
}

func (s *TestKafkaConnectWorkerConfigurationSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = kafkaconnect.NewFromConfig(cfg)

	name := fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano())
	resp, err := s.svc.CreateWorkerConfiguration(ctx, &kafkaconnect.CreateWorkerConfigurationInput{
		Name:                  ptr.String(name),
		PropertiesFileContent: ptr.String("key.converter=org.apache.kafka.connect.storage.StringConverter\nvalue.converter=org.apache.kafka.connect.storage.StringConverter"),
	})
	if err != nil {
		s.T().Fatalf("failed to create worker configuration: %v", err)
	}
	s.arn = resp.WorkerConfigurationArn
}

func (s *TestKafkaConnectWorkerConfigurationSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = s.svc.DeleteWorkerConfiguration(ctx, &kafkaconnect.DeleteWorkerConfigurationInput{
		WorkerConfigurationArn: s.arn,
	})
}

func (s *TestKafkaConnectWorkerConfigurationSuite) TestList() {
	a := assert.New(s.T())
	lister := &KafkaConnectWorkerConfigurationLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testKafkaConnectListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (s *TestKafkaConnectWorkerConfigurationSuite) TestRemove() {
	a := assert.New(s.T())
	wc := &KafkaConnectWorkerConfiguration{svc: s.svc, WorkerConfigurationArn: s.arn}
	a.NoError(wc.Remove(context.TODO()))
}

func TestKafkaConnectWorkerConfigurationIntegration(t *testing.T) {
	suite.Run(t, new(TestKafkaConnectWorkerConfigurationSuite))
}
