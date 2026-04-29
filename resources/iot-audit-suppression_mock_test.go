package resources

import (
	"context"
	"testing"
	"time"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/iot"
	iottypes "github.com/aws/aws-sdk-go-v2/service/iot/types"
)

func Test_Mock_IoTAuditSuppression_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIoTClient)

	now := time.Now()

	mockClient.On("ListAuditSuppressions", mock.Anything, mock.Anything).
		Return(&iot.ListAuditSuppressionsOutput{
			Suppressions: []iottypes.AuditSuppression{
				{
					CheckName: ptr.String("LOGGING_DISABLED_CHECK"),
					ResourceIdentifier: &iottypes.ResourceIdentifier{
						Account: ptr.String("123456789012"),
					},
					Description:          ptr.String("test suppression"),
					ExpirationDate:       &now,
					SuppressIndefinitely: ptr.Bool(true),
				},
			},
			NextToken: nil,
		}, nil)

	lister := &IoTAuditSuppressionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIoTListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	suppression := resources[0].(*IoTAuditSuppression)
	a.Equal("LOGGING_DISABLED_CHECK", *suppression.CheckName)
	a.Equal("123456789012", *suppression.ResourceIdentifierAccount)
	a.Equal("test suppression", *suppression.Description)
	a.Equal(now, *suppression.ExpirationDate)
	a.True(*suppression.SuppressIndefinitely)

	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTAuditSuppression_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIoTClient)

	mockClient.On("ListAuditSuppressions", mock.Anything, mock.Anything).
		Return(&iot.ListAuditSuppressionsOutput{
			Suppressions: []iottypes.AuditSuppression{},
		}, nil)

	lister := &IoTAuditSuppressionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIoTListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTAuditSuppression_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIoTClient)

	ri := &iottypes.ResourceIdentifier{
		Account: ptr.String("123456789012"),
	}

	suppression := &IoTAuditSuppression{
		svc:                mockClient,
		CheckName:          ptr.String("LOGGING_DISABLED_CHECK"),
		ResourceIdentifier: ri,
	}

	mockClient.On("DeleteAuditSuppression", mock.Anything, &iot.DeleteAuditSuppressionInput{
		CheckName:          suppression.CheckName,
		ResourceIdentifier: ri,
	}).Return(&iot.DeleteAuditSuppressionOutput{}, nil)

	a.NoError(suppression.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTAuditSuppression_Properties(t *testing.T) {
	a := assert.New(t)

	now := time.Now()

	suppression := IoTAuditSuppression{
		CheckName: ptr.String("LOGGING_DISABLED_CHECK"),
		ResourceIdentifier: &iottypes.ResourceIdentifier{
			Account: ptr.String("123456789012"),
		},
		Description:               ptr.String("test suppression"),
		ExpirationDate:            &now,
		SuppressIndefinitely:      ptr.Bool(true),
		ResourceIdentifierAccount: ptr.String("123456789012"),
	}

	props := suppression.Properties()
	a.Equal("LOGGING_DISABLED_CHECK", props.Get("CheckName"))
	a.Equal("test suppression", props.Get("Description"))
	a.Equal("true", props.Get("SuppressIndefinitely"))
	a.Equal("123456789012", props.Get("ResourceIdentifierAccount"))
	a.Equal("", props.Get("ResourceIdentifier"))
}

func Test_Mock_IoTAuditSuppression_String(t *testing.T) {
	a := assert.New(t)
	suppression := IoTAuditSuppression{CheckName: ptr.String("LOGGING_DISABLED_CHECK")}
	a.Equal("LOGGING_DISABLED_CHECK", suppression.String())
}
