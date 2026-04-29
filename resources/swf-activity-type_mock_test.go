package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/swf"
	swftypes "github.com/aws/aws-sdk-go-v2/service/swf/types"
)

func Test_Mock_SWFActivityType_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSWFClient)
	mockClient.On("ListDomains", mock.Anything, mock.Anything).
		Return(&swf.ListDomainsOutput{
			DomainInfos: []swftypes.DomainInfo{
				{Name: ptr.String("my-domain"), Status: swftypes.RegistrationStatusRegistered},
			},
		}, nil)
	mockClient.On("ListActivityTypes", mock.Anything, mock.Anything).
		Return(&swf.ListActivityTypesOutput{
			TypeInfos: []swftypes.ActivityTypeInfo{
				{
					ActivityType: &swftypes.ActivityType{
						Name:    ptr.String("my-activity"),
						Version: ptr.String("1.0"),
					},
					Status: swftypes.RegistrationStatusRegistered,
				},
			},
		}, nil)
	lister := &SWFActivityTypeLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSWFListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	r := resources[0].(*SWFActivityType)
	a.Equal("my-activity", *r.Name)
	a.Equal("1.0", *r.Version)
	a.Equal("my-domain", *r.Domain)
	a.Equal("REGISTERED", *r.Status)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SWFActivityType_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSWFClient)
	mockClient.On("ListDomains", mock.Anything, mock.Anything).
		Return(&swf.ListDomainsOutput{DomainInfos: []swftypes.DomainInfo{}}, nil)
	lister := &SWFActivityTypeLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSWFListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SWFActivityType_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSWFClient)
	r := &SWFActivityType{
		svc:     mockClient,
		Domain:  ptr.String("my-domain"),
		Name:    ptr.String("my-activity"),
		Version: ptr.String("1.0"),
	}
	mockClient.On("DeprecateActivityType", mock.Anything, &swf.DeprecateActivityTypeInput{
		Domain: r.Domain,
		ActivityType: &swftypes.ActivityType{
			Name:    r.Name,
			Version: r.Version,
		},
	}).Return(&swf.DeprecateActivityTypeOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_SWFActivityType_Filter_Deprecated(t *testing.T) {
	a := assert.New(t)
	r := SWFActivityType{Name: ptr.String("old-activity"), Status: ptr.String("DEPRECATED")}
	a.Error(r.Filter())
}

func Test_Mock_SWFActivityType_Filter_Registered(t *testing.T) {
	a := assert.New(t)
	r := SWFActivityType{Name: ptr.String("my-activity"), Status: ptr.String("REGISTERED")}
	a.NoError(r.Filter())
}

func Test_Mock_SWFActivityType_Properties(t *testing.T) {
	a := assert.New(t)
	r := SWFActivityType{
		Domain:  ptr.String("my-domain"),
		Name:    ptr.String("my-activity"),
		Version: ptr.String("1.0"),
		Status:  ptr.String("REGISTERED"),
	}
	props := r.Properties()
	a.Equal("my-domain", props.Get("Domain"))
	a.Equal("my-activity", props.Get("Name"))
	a.Equal("1.0", props.Get("Version"))
	a.Equal("REGISTERED", props.Get("Status"))
}

func Test_Mock_SWFActivityType_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("my-activity", (&SWFActivityType{Name: ptr.String("my-activity")}).String())
}
