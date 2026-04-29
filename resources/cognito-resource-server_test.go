//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

type TestCognitoResourceServerSuite struct {
	suite.Suite
	svc *cognitoidentityprovider.Client
}

func (suite *TestCognitoResourceServerSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = cognitoidentityprovider.NewFromConfig(cfg)
}

func (suite *TestCognitoResourceServerSuite) TestList() {
	a := assert.New(suite.T())
	lister := CognitoResourceServerLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testCognitoListerOpts)
	a.NoError(err)
	_ = resources
}

func TestCognitoResourceServerIntegration(t *testing.T) {
	suite.Run(t, new(TestCognitoResourceServerSuite))
}
