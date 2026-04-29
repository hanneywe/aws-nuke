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

func Test_Mock_SWFDomain_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSWFClient)
	mockClient.On("ListDomains", mock.Anything, mock.Anything).
		Return(&swf.ListDomainsOutput{
			DomainInfos: []swftypes.DomainInfo{
				{Name: ptr.String("my-domain"), Status: swftypes.RegistrationStatusRegistered},
			},
		}, nil)
	lister := &SWFDomainLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSWFListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	d := resources[0].(*SWFDomain)
	a.Equal("my-domain", *d.Name)
	a.Equal("REGISTERED", *d.Status)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SWFDomain_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSWFClient)
	mockClient.On("ListDomains", mock.Anything, mock.Anything).
		Return(&swf.ListDomainsOutput{DomainInfos: []swftypes.DomainInfo{}}, nil)
	lister := &SWFDomainLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSWFListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SWFDomain_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSWFClient)
	d := &SWFDomain{svc: mockClient, Name: ptr.String("my-domain")}
	mockClient.On("DeprecateDomain", mock.Anything, &swf.DeprecateDomainInput{Name: d.Name}).
		Return(&swf.DeprecateDomainOutput{}, nil)
	a.NoError(d.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_SWFDomain_Filter_Deprecated(t *testing.T) {
	a := assert.New(t)
	d := SWFDomain{Name: ptr.String("old-domain"), Status: ptr.String("DEPRECATED")}
	a.Error(d.Filter())
}

func Test_Mock_SWFDomain_Filter_Registered(t *testing.T) {
	a := assert.New(t)
	d := SWFDomain{Name: ptr.String("my-domain"), Status: ptr.String("REGISTERED")}
	a.NoError(d.Filter())
}

func Test_Mock_SWFDomain_Properties(t *testing.T) {
	a := assert.New(t)
	d := SWFDomain{Name: ptr.String("my-domain"), Status: ptr.String("REGISTERED")}
	a.Equal("my-domain", d.Properties().Get("Name"))
	a.Equal("REGISTERED", d.Properties().Get("Status"))
}

func Test_Mock_SWFDomain_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("my-domain", (&SWFDomain{Name: ptr.String("my-domain")}).String())
}
