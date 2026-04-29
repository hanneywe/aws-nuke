package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/glue"
	gluetypes "github.com/aws/aws-sdk-go-v2/service/glue/types"
)

func Test_Mock_GlueDataCatalogEncryptionSettings_List_WithConfig(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGlueV2Client)

	mockClient.On("GetDataCatalogEncryptionSettings", mock.Anything, mock.Anything).
		Return(&glue.GetDataCatalogEncryptionSettingsOutput{
			DataCatalogEncryptionSettings: &gluetypes.DataCatalogEncryptionSettings{
				EncryptionAtRest: &gluetypes.EncryptionAtRest{
					CatalogEncryptionMode: gluetypes.CatalogEncryptionModeSsekms,
					SseAwsKmsKeyId:        ptr.String("arn:aws:kms:us-east-1:123456789012:key/test-key"),
				},
				ConnectionPasswordEncryption: &gluetypes.ConnectionPasswordEncryption{
					ReturnConnectionPasswordEncrypted: true,
					AwsKmsKeyId:                       ptr.String("arn:aws:kms:us-east-1:123456789012:key/test-key"),
				},
			},
		}, nil)

	lister := &GlueDataCatalogEncryptionSettingsLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGlueV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*GlueDataCatalogEncryptionSettings)
	a.True(*r.EncryptionAtRest)
	a.True(*r.ConnectionPasswordEncryption)
	mockClient.AssertExpectations(t)
}

func Test_Mock_GlueDataCatalogEncryptionSettings_List_DefaultConfig(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGlueV2Client)

	mockClient.On("GetDataCatalogEncryptionSettings", mock.Anything, mock.Anything).
		Return(&glue.GetDataCatalogEncryptionSettingsOutput{
			DataCatalogEncryptionSettings: &gluetypes.DataCatalogEncryptionSettings{
				EncryptionAtRest: &gluetypes.EncryptionAtRest{
					CatalogEncryptionMode: gluetypes.CatalogEncryptionModeDisabled,
				},
				ConnectionPasswordEncryption: &gluetypes.ConnectionPasswordEncryption{
					ReturnConnectionPasswordEncrypted: false,
				},
			},
		}, nil)

	lister := &GlueDataCatalogEncryptionSettingsLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGlueV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_GlueDataCatalogEncryptionSettings_List_NilSettings(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGlueV2Client)

	mockClient.On("GetDataCatalogEncryptionSettings", mock.Anything, mock.Anything).
		Return(&glue.GetDataCatalogEncryptionSettingsOutput{
			DataCatalogEncryptionSettings: nil,
		}, nil)

	lister := &GlueDataCatalogEncryptionSettingsLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGlueV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_GlueDataCatalogEncryptionSettings_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGlueV2Client)

	r := &GlueDataCatalogEncryptionSettings{
		svc:                          mockClient,
		EncryptionAtRest:             ptr.Bool(true),
		ConnectionPasswordEncryption: ptr.Bool(true),
	}

	mockClient.On("PutDataCatalogEncryptionSettings", mock.Anything,
		&glue.PutDataCatalogEncryptionSettingsInput{
			DataCatalogEncryptionSettings: &gluetypes.DataCatalogEncryptionSettings{
				ConnectionPasswordEncryption: &gluetypes.ConnectionPasswordEncryption{
					ReturnConnectionPasswordEncrypted: false,
				},
				EncryptionAtRest: &gluetypes.EncryptionAtRest{
					CatalogEncryptionMode: gluetypes.CatalogEncryptionModeDisabled,
				},
			},
		}).Return(&glue.PutDataCatalogEncryptionSettingsOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_GlueDataCatalogEncryptionSettings_Properties(t *testing.T) {
	a := assert.New(t)
	r := &GlueDataCatalogEncryptionSettings{
		EncryptionAtRest:             ptr.Bool(true),
		ConnectionPasswordEncryption: ptr.Bool(true),
	}
	props := r.Properties()
	a.Equal("true", props.Get("EncryptionAtRest"))
	a.Equal("true", props.Get("ConnectionPasswordEncryption"))
}

func Test_Mock_GlueDataCatalogEncryptionSettings_String(t *testing.T) {
	a := assert.New(t)
	r := &GlueDataCatalogEncryptionSettings{
		EncryptionAtRest:             ptr.Bool(true),
		ConnectionPasswordEncryption: ptr.Bool(true),
	}
	a.Equal("GlueDataCatalogEncryptionSettings", r.String())
}
