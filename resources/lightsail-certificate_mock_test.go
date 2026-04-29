package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	lstypes "github.com/aws/aws-sdk-go-v2/service/lightsail/types"
)

func Test_Mock_LightsailCertificate_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockLightsailClient)

	mockClient.On("GetCertificates", mock.Anything, mock.Anything).
		Return(&lightsail.GetCertificatesOutput{
			Certificates: []lstypes.CertificateSummary{
				{
					CertificateName: ptr.String("my-cert"),
					CertificateArn:  ptr.String("arn:aws:lightsail:us-east-1:123456789012:Certificate/my-cert"),
				},
			},
		}, nil)

	lister := &LightsailCertificateLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testLightsailListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	c := resources[0].(*LightsailCertificate)
	a.Equal("my-cert", *c.CertificateName)
	a.Equal("arn:aws:lightsail:us-east-1:123456789012:Certificate/my-cert", *c.CertificateArn)

	mockClient.AssertExpectations(t)
}

func Test_Mock_LightsailCertificate_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockLightsailClient)

	mockClient.On("GetCertificates", mock.Anything, mock.Anything).
		Return(&lightsail.GetCertificatesOutput{
			Certificates: []lstypes.CertificateSummary{},
		}, nil)

	lister := &LightsailCertificateLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testLightsailListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_LightsailCertificate_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockLightsailClient)

	c := &LightsailCertificate{
		svc:             mockClient,
		CertificateName: ptr.String("my-cert"),
	}

	mockClient.On("DeleteCertificate", mock.Anything, &lightsail.DeleteCertificateInput{
		CertificateName: c.CertificateName,
	}).Return(&lightsail.DeleteCertificateOutput{}, nil)

	a.NoError(c.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_LightsailCertificate_Properties(t *testing.T) {
	a := assert.New(t)

	c := LightsailCertificate{
		CertificateName: ptr.String("my-cert"),
		CertificateArn:  ptr.String("arn:aws:lightsail:us-east-1:123456789012:Certificate/my-cert"),
	}

	props := c.Properties()
	a.Equal("my-cert", props.Get("CertificateName"))
	a.Equal("arn:aws:lightsail:us-east-1:123456789012:Certificate/my-cert", props.Get("CertificateArn"))
}

func Test_Mock_LightsailCertificate_String(t *testing.T) {
	a := assert.New(t)
	c := LightsailCertificate{CertificateName: ptr.String("my-cert")}
	a.Equal("my-cert", c.String())
}
