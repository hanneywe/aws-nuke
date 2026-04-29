package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/lightsail"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const LightsailCertificateResource = "LightsailCertificate"

func init() {
	registry.Register(&registry.Registration{
		Name:     LightsailCertificateResource,
		Scope:    nuke.Account,
		Resource: &LightsailCertificate{},
		Lister:   &LightsailCertificateLister{},
	})
}

type LightsailCertificateLister struct {
	svc LightsailClient
}

func (l *LightsailCertificateLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = lightsail.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	resp, err := svc.GetCertificates(ctx, &lightsail.GetCertificatesInput{})
	if err != nil {
		return nil, err
	}

	for _, cert := range resp.Certificates {
		resources = append(resources, &LightsailCertificate{
			svc:             svc,
			CertificateName: cert.CertificateName,
			CertificateArn:  cert.CertificateArn,
		})
	}

	return resources, nil
}

type LightsailCertificate struct {
	svc             LightsailClient
	CertificateName *string
	CertificateArn  *string
}

func (r *LightsailCertificate) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteCertificate(ctx, &lightsail.DeleteCertificateInput{
		CertificateName: r.CertificateName,
	})
	return err
}

func (r *LightsailCertificate) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *LightsailCertificate) String() string {
	return *r.CertificateName
}
