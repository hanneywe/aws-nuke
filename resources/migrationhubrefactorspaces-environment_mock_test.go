package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/migrationhubrefactorspaces"
	migrationhubrefactorspacestypes "github.com/aws/aws-sdk-go-v2/service/migrationhubrefactorspaces/types"
)

func Test_Mock_MigrationHubRefactorSpacesEnvironment_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMigrationHubRefactorSpacesClient)
	mockClient.On("ListEnvironments", mock.Anything, mock.Anything).
		Return(&migrationhubrefactorspaces.ListEnvironmentsOutput{
			EnvironmentSummaryList: []migrationhubrefactorspacestypes.EnvironmentSummary{
				{EnvironmentId: ptr.String("env-12345"), Name: ptr.String("my-env")},
			},
		}, nil)
	lister := &MigrationHubRefactorSpacesEnvironmentLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMigrationHubRefactorSpacesListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	env := resources[0].(*MigrationHubRefactorSpacesEnvironment)
	a.Equal("env-12345", *env.EnvironmentID)
	a.Equal("my-env", *env.Name)
	mockClient.AssertExpectations(t)
}

func Test_Mock_MigrationHubRefactorSpacesEnvironment_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMigrationHubRefactorSpacesClient)
	mockClient.On("ListEnvironments", mock.Anything, mock.Anything).
		Return(&migrationhubrefactorspaces.ListEnvironmentsOutput{
			EnvironmentSummaryList: []migrationhubrefactorspacestypes.EnvironmentSummary{},
		}, nil)
	lister := &MigrationHubRefactorSpacesEnvironmentLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMigrationHubRefactorSpacesListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_MigrationHubRefactorSpacesEnvironment_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMigrationHubRefactorSpacesClient)
	env := &MigrationHubRefactorSpacesEnvironment{
		svc:           mockClient,
		EnvironmentID: ptr.String("env-12345"),
		Name:          ptr.String("my-env"),
	}
	mockClient.On("DeleteEnvironment", mock.Anything,
		&migrationhubrefactorspaces.DeleteEnvironmentInput{EnvironmentIdentifier: env.EnvironmentID}).
		Return(&migrationhubrefactorspaces.DeleteEnvironmentOutput{}, nil)
	a.NoError(env.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_MigrationHubRefactorSpacesEnvironment_Properties(t *testing.T) {
	a := assert.New(t)
	env := MigrationHubRefactorSpacesEnvironment{EnvironmentID: ptr.String("env-12345"), Name: ptr.String("my-env")}
	a.Equal("env-12345", env.Properties().Get("EnvironmentId"))
	a.Equal("my-env", env.Properties().Get("Name"))
}

func Test_Mock_MigrationHubRefactorSpacesEnvironment_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("my-env", (&MigrationHubRefactorSpacesEnvironment{Name: ptr.String("my-env")}).String())
}
