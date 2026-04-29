package resources

import (
	"context"
	"fmt"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/vpclattice"
	lattice "github.com/aws/aws-sdk-go-v2/service/vpclattice/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

func Test_Mock_VPCLatticeResourceConfiguration_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	mockClient.
		On("ListResourceConfigurations", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListResourceConfigurationsOutput{
				Items: []lattice.ResourceConfigurationSummary{
					{
						Id:   ptr.String("rcfg-1"),
						Arn:  ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:resourceconfiguration/rcfg-1"),
						Name: ptr.String("config-1"),
						Type: lattice.ResourceConfigurationTypeSingle,
					},
				},
			}, nil,
		)

	lister := &VPCLatticeResourceConfigurationLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), &nuke.ListerOpts{})
	assertions.NoError(err)
	assertions.Len(resources, 1)

	resourceConfig := resources[0].(*VPCLatticeResourceConfiguration)
	assertions.Equal("config-1", *resourceConfig.Name)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeResourceConfiguration_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	mockClient.
		On("ListResourceConfigurations", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListResourceConfigurationsOutput{
				Items: []lattice.ResourceConfigurationSummary{},
			}, nil,
		)

	lister := &VPCLatticeResourceConfigurationLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), &nuke.ListerOpts{})
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeResourceConfiguration_List_MultiPage(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	firstPageItems := make([]lattice.ResourceConfigurationSummary, 100)
	for i := range firstPageItems {
		firstPageItems[i] = lattice.ResourceConfigurationSummary{
			Id:   ptr.String(fmt.Sprintf("rcfg-%d", i)),
			Arn:  ptr.String(fmt.Sprintf("arn:aws:vpc-lattice:us-east-1:123456789012:resourceconfiguration/rcfg-%d", i)),
			Name: ptr.String(fmt.Sprintf("config-%d", i)),
			Type: lattice.ResourceConfigurationTypeSingle,
		}
	}

	mockClient.
		On(
			"ListResourceConfigurations",
			mock.Anything,
			mock.MatchedBy(func(input *vpclattice.ListResourceConfigurationsInput) bool {
				return input.NextToken == nil
			}),
		).
		Return(
			&vpclattice.ListResourceConfigurationsOutput{
				Items:     firstPageItems,
				NextToken: ptr.String("page2"),
			}, nil,
		).
		Once()

	mockClient.
		On(
			"ListResourceConfigurations",
			mock.Anything,
			mock.MatchedBy(func(input *vpclattice.ListResourceConfigurationsInput) bool {
				return input.NextToken != nil
			}),
		).
		Return(
			&vpclattice.ListResourceConfigurationsOutput{
				Items: []lattice.ResourceConfigurationSummary{
					{
						Id:   ptr.String("rcfg-100"),
						Arn:  ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:resourceconfiguration/rcfg-100"),
						Name: ptr.String("config-100"),
						Type: lattice.ResourceConfigurationTypeGroup,
					},
				},
			}, nil,
		).
		Once()

	lister := &VPCLatticeResourceConfigurationLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), &nuke.ListerOpts{})
	assertions.NoError(err)
	assertions.Len(resources, 101)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeResourceConfiguration_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	resourceConfig := &VPCLatticeResourceConfiguration{
		svc: mockClient,
		ARN: ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:resourceconfiguration/rcfg-1"),
	}

	mockClient.
		On(
			"DeleteResourceConfiguration",
			mock.Anything,
			&vpclattice.DeleteResourceConfigurationInput{
				ResourceConfigurationIdentifier: resourceConfig.ARN,
			},
		).
		Return(&vpclattice.DeleteResourceConfigurationOutput{}, nil)

	err := resourceConfig.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeResourceConfiguration_Properties(t *testing.T) {
	assertions := assert.New(t)

	resourceConfig := VPCLatticeResourceConfiguration{
		ID:   ptr.String("rcfg-12345"),
		ARN:  ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:resourceconfiguration/rcfg-12345"),
		Name: ptr.String("my-resource-config"),
		Type: ptr.String("SINGLE"),
	}

	properties := resourceConfig.Properties()

	assertions.Equal("rcfg-12345", properties.Get("ID"))
	assertions.Equal("my-resource-config", properties.Get("Name"))
	assertions.Equal("SINGLE", properties.Get("Type"))
}

func Test_Mock_VPCLatticeResourceConfiguration_String(t *testing.T) {
	assertions := assert.New(t)

	resourceConfig := VPCLatticeResourceConfiguration{
		Name: ptr.String("my-resource-config"),
	}

	assertions.Equal("my-resource-config", resourceConfig.String())
}
