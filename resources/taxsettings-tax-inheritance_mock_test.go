package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/taxsettings"
	taxsettingstypes "github.com/aws/aws-sdk-go-v2/service/taxsettings/types"
)

func Test_Mock_TaxSettingsTaxInheritance_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockTaxSettingsClient)

	mockClient.On("GetTaxInheritance", mock.Anything, mock.Anything).
		Return(&taxsettings.GetTaxInheritanceOutput{
			HeritageStatus: taxsettingstypes.HeritageStatusOptIn,
		}, nil)

	lister := &TaxSettingsTaxInheritanceLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*TaxSettingsTaxInheritance)
	a.Equal(taxsettingstypes.HeritageStatusOptIn, r.Status)
	mockClient.AssertExpectations(t)
}

func Test_Mock_TaxSettingsTaxInheritance_List_OptOut(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockTaxSettingsClient)

	mockClient.On("GetTaxInheritance", mock.Anything, mock.Anything).
		Return(&taxsettings.GetTaxInheritanceOutput{
			HeritageStatus: taxsettingstypes.HeritageStatusOptOut,
		}, nil)

	lister := &TaxSettingsTaxInheritanceLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_TaxSettingsTaxInheritance_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockTaxSettingsClient)

	r := &TaxSettingsTaxInheritance{
		svc:    mockClient,
		Status: taxsettingstypes.HeritageStatusOptIn,
	}

	mockClient.On("PutTaxInheritance", mock.Anything,
		&taxsettings.PutTaxInheritanceInput{
			HeritageStatus: taxsettingstypes.HeritageStatusOptOut,
		}).Return(&taxsettings.PutTaxInheritanceOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_TaxSettingsTaxInheritance_Properties(t *testing.T) {
	a := assert.New(t)
	r := &TaxSettingsTaxInheritance{
		Status: taxsettingstypes.HeritageStatusOptIn,
	}
	props := r.Properties()
	a.Equal("OptIn", props.Get("Status"))
}

func Test_Mock_TaxSettingsTaxInheritance_String(t *testing.T) {
	a := assert.New(t)
	r := &TaxSettingsTaxInheritance{
		Status: taxsettingstypes.HeritageStatusOptIn,
	}
	a.Equal("OptIn", r.String())
}
