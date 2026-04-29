package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/redshift"
	redshifttypes "github.com/aws/aws-sdk-go-v2/service/redshift/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

var testRedshiftListerOpts = &nuke.ListerOpts{
	Logger: logrus.WithField("test", true),
}

func Test_Mock_RedshiftHsmClientCertificate_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockRedshiftClient)

	mockClient.On("DescribeHsmClientCertificates", mock.Anything, mock.Anything).
		Return(&redshift.DescribeHsmClientCertificatesOutput{
			HsmClientCertificates: []redshifttypes.HsmClientCertificate{
				{
					HsmClientCertificateIdentifier: ptr.String("my-hsm-cert"),
				},
			},
		}, nil)

	lister := &RedshiftHsmClientCertificateLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testRedshiftListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	cert := resources[0].(*RedshiftHsmClientCertificate)
	a.Equal("my-hsm-cert", *cert.HsmClientCertificateIdentifier)

	mockClient.AssertExpectations(t)
}

func Test_Mock_RedshiftHsmClientCertificate_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockRedshiftClient)

	mockClient.On("DescribeHsmClientCertificates", mock.Anything, mock.Anything).
		Return(&redshift.DescribeHsmClientCertificatesOutput{
			HsmClientCertificates: []redshifttypes.HsmClientCertificate{},
		}, nil)

	lister := &RedshiftHsmClientCertificateLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testRedshiftListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_RedshiftHsmClientCertificate_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockRedshiftClient)

	cert := &RedshiftHsmClientCertificate{
		svc:                            mockClient,
		HsmClientCertificateIdentifier: ptr.String("my-hsm-cert"),
	}

	mockClient.On("DeleteHsmClientCertificate", mock.Anything, &redshift.DeleteHsmClientCertificateInput{
		HsmClientCertificateIdentifier: cert.HsmClientCertificateIdentifier,
	}).Return(&redshift.DeleteHsmClientCertificateOutput{}, nil)

	a.NoError(cert.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_RedshiftHsmClientCertificate_Properties(t *testing.T) {
	a := assert.New(t)

	cert := RedshiftHsmClientCertificate{
		HsmClientCertificateIdentifier: ptr.String("my-hsm-cert"),
	}

	props := cert.Properties()
	a.Equal("my-hsm-cert", props.Get("HsmClientCertificateIdentifier"))
}

func Test_Mock_RedshiftHsmClientCertificate_String(t *testing.T) {
	a := assert.New(t)
	cert := RedshiftHsmClientCertificate{HsmClientCertificateIdentifier: ptr.String("my-hsm-cert")}
	a.Equal("my-hsm-cert", cert.String())
}
