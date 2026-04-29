package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/iot"
	iottypes "github.com/aws/aws-sdk-go-v2/service/iot/types"
)

func Test_Mock_IoTJobTemplate_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIoTClient)

	mockClient.On("ListJobTemplates", mock.Anything, mock.Anything).
		Return(&iot.ListJobTemplatesOutput{
			JobTemplates: []iottypes.JobTemplateSummary{
				{
					JobTemplateId:  ptr.String("jt-12345"),
					JobTemplateArn: ptr.String("arn:aws:iot:us-east-1:123456789012:jobtemplate/jt-12345"),
				},
			},
		}, nil)

	lister := &IoTJobTemplateLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIoTListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	jt := resources[0].(*IoTJobTemplate)
	a.Equal("jt-12345", *jt.JobTemplateID)
	a.Equal("arn:aws:iot:us-east-1:123456789012:jobtemplate/jt-12345", *jt.JobTemplateArn)

	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTJobTemplate_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIoTClient)

	mockClient.On("ListJobTemplates", mock.Anything, mock.Anything).
		Return(&iot.ListJobTemplatesOutput{
			JobTemplates: []iottypes.JobTemplateSummary{},
		}, nil)

	lister := &IoTJobTemplateLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIoTListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTJobTemplate_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIoTClient)

	jt := &IoTJobTemplate{
		svc:           mockClient,
		JobTemplateID: ptr.String("jt-12345"),
	}

	mockClient.On("DeleteJobTemplate", mock.Anything, &iot.DeleteJobTemplateInput{
		JobTemplateId: jt.JobTemplateID,
	}).Return(&iot.DeleteJobTemplateOutput{}, nil)

	a.NoError(jt.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTJobTemplate_Properties(t *testing.T) {
	a := assert.New(t)

	jt := IoTJobTemplate{
		JobTemplateID:  ptr.String("jt-12345"),
		JobTemplateArn: ptr.String("arn:aws:iot:us-east-1:123456789012:jobtemplate/jt-12345"),
	}

	props := jt.Properties()
	a.Equal("jt-12345", props.Get("JobTemplateId"))
	a.Equal("arn:aws:iot:us-east-1:123456789012:jobtemplate/jt-12345", props.Get("JobTemplateArn"))
}

func Test_Mock_IoTJobTemplate_String(t *testing.T) {
	a := assert.New(t)
	jt := IoTJobTemplate{JobTemplateID: ptr.String("jt-12345")}
	a.Equal("jt-12345", jt.String())
}
