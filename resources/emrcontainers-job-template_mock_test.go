package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/emrcontainers"
	emrcontainerstypes "github.com/aws/aws-sdk-go-v2/service/emrcontainers/types"
)

func Test_Mock_EMRContainersJobTemplate_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockEMRContainersClient)
	mockClient.On("ListJobTemplates", mock.Anything, mock.Anything).
		Return(&emrcontainers.ListJobTemplatesOutput{
			Templates: []emrcontainerstypes.JobTemplate{
				{
					Id:   ptr.String("jt-123"),
					Name: ptr.String("my-template"),
					Arn:  ptr.String("arn:aws:emr-containers:us-east-1:123456789012:/jobtemplates/jt-123"),
					Tags: map[string]string{"env": "test"},
				},
			},
		}, nil)
	lister := &EMRContainersJobTemplateLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testEMRContainersListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	r := resources[0].(*EMRContainersJobTemplate)
	a.Equal("jt-123", *r.ID)
	a.Equal("my-template", *r.Name)
	a.Equal("arn:aws:emr-containers:us-east-1:123456789012:/jobtemplates/jt-123", *r.ARN)
	a.Equal("test", r.Tags["env"])
	mockClient.AssertExpectations(t)
}

func Test_Mock_EMRContainersJobTemplate_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockEMRContainersClient)
	mockClient.On("ListJobTemplates", mock.Anything, mock.Anything).
		Return(&emrcontainers.ListJobTemplatesOutput{
			Templates: []emrcontainerstypes.JobTemplate{},
		}, nil)
	lister := &EMRContainersJobTemplateLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testEMRContainersListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_EMRContainersJobTemplate_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockEMRContainersClient)
	r := &EMRContainersJobTemplate{svc: mockClient, ID: ptr.String("jt-123"), Name: ptr.String("my-template")}
	mockClient.On("DeleteJobTemplate", mock.Anything, &emrcontainers.DeleteJobTemplateInput{Id: r.ID}).
		Return(&emrcontainers.DeleteJobTemplateOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_EMRContainersJobTemplate_Properties(t *testing.T) {
	a := assert.New(t)
	r := EMRContainersJobTemplate{
		ID:   ptr.String("jt-123"),
		Name: ptr.String("my-template"),
		ARN:  ptr.String("arn:aws:emr-containers:us-east-1:123456789012:/jobtemplates/jt-123"),
		Tags: map[string]string{"env": "test"},
	}
	a.Equal("jt-123", r.Properties().Get("ID"))
	a.Equal("my-template", r.Properties().Get("Name"))
	a.Equal("arn:aws:emr-containers:us-east-1:123456789012:/jobtemplates/jt-123", r.Properties().Get("ARN"))
	a.Equal("test", r.Properties().Get("tag:env"))
}

func Test_Mock_EMRContainersJobTemplate_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("my-template", (&EMRContainersJobTemplate{Name: ptr.String("my-template")}).String())
}
