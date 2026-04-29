package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/bedrockdataautomation"
	bedrocktypes "github.com/aws/aws-sdk-go-v2/service/bedrockdataautomation/types"
)

func Test_Mock_BedrockDataAutomationProject_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockBedrockDataAutomationClient)

	mockClient.
		On("ListDataAutomationProjects", mock.Anything, mock.Anything).
		Return(&bedrockdataautomation.ListDataAutomationProjectsOutput{
			Projects: []bedrocktypes.DataAutomationProjectSummary{
				{
					ProjectArn:  ptr.String("arn:aws:bedrock:us-east-1:123456789012:data-automation-project/proj-1"),
					ProjectName: ptr.String("my-project"),
				},
			},
		}, nil)

	lister := &BedrockDataAutomationProjectLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testBedrockDataAutomationListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	project := resources[0].(*BedrockDataAutomationProject)
	a.Equal("my-project", *project.ProjectName)

	mockClient.AssertExpectations(t)
}

func Test_Mock_BedrockDataAutomationProject_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockBedrockDataAutomationClient)

	mockClient.
		On("ListDataAutomationProjects", mock.Anything, mock.Anything).
		Return(&bedrockdataautomation.ListDataAutomationProjectsOutput{
			Projects: []bedrocktypes.DataAutomationProjectSummary{},
		}, nil)

	lister := &BedrockDataAutomationProjectLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testBedrockDataAutomationListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_BedrockDataAutomationProject_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockBedrockDataAutomationClient)

	project := &BedrockDataAutomationProject{
		svc:        mockClient,
		ProjectArn: ptr.String("arn:aws:bedrock:us-east-1:123456789012:data-automation-project/proj-1"),
	}

	mockClient.
		On("DeleteDataAutomationProject", mock.Anything, &bedrockdataautomation.DeleteDataAutomationProjectInput{
			ProjectArn: project.ProjectArn,
		}).
		Return(&bedrockdataautomation.DeleteDataAutomationProjectOutput{}, nil)

	err := project.Remove(context.TODO())
	a.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_BedrockDataAutomationProject_Properties(t *testing.T) {
	a := assert.New(t)

	project := BedrockDataAutomationProject{
		ProjectArn:  ptr.String("arn:aws:bedrock:us-east-1:123456789012:data-automation-project/proj-1"),
		ProjectName: ptr.String("my-project"),
	}

	props := project.Properties()
	a.Equal("arn:aws:bedrock:us-east-1:123456789012:data-automation-project/proj-1", props.Get("ProjectArn"))
	a.Equal("my-project", props.Get("ProjectName"))
}

func Test_Mock_BedrockDataAutomationProject_String(t *testing.T) {
	a := assert.New(t)

	project := BedrockDataAutomationProject{
		ProjectArn: ptr.String("arn:aws:bedrock:us-east-1:123456789012:data-automation-project/proj-1"),
	}

	a.Equal("arn:aws:bedrock:us-east-1:123456789012:data-automation-project/proj-1", project.String())
}
