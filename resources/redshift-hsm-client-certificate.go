package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/redshift"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const RedshiftHsmClientCertificateResource = "RedshiftHsmClientCertificate"

func init() {
	registry.Register(&registry.Registration{
		Name:     RedshiftHsmClientCertificateResource,
		Scope:    nuke.Account,
		Resource: &RedshiftHsmClientCertificate{},
		Lister:   &RedshiftHsmClientCertificateLister{},
	})
}

type RedshiftHsmClientCertificateLister struct {
	svc RedshiftClient
}

func (l *RedshiftHsmClientCertificateLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = redshift.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := redshift.NewDescribeHsmClientCertificatesPaginator(svc, &redshift.DescribeHsmClientCertificatesInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.HsmClientCertificates {
			resources = append(resources, &RedshiftHsmClientCertificate{
				svc:                            svc,
				HsmClientCertificateIdentifier: item.HsmClientCertificateIdentifier,
			})
		}
	}

	return resources, nil
}

type RedshiftHsmClientCertificate struct {
	svc                            RedshiftClient
	HsmClientCertificateIdentifier *string
}

func (r *RedshiftHsmClientCertificate) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteHsmClientCertificate(ctx, &redshift.DeleteHsmClientCertificateInput{
		HsmClientCertificateIdentifier: r.HsmClientCertificateIdentifier,
	})
	return err
}

func (r *RedshiftHsmClientCertificate) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *RedshiftHsmClientCertificate) String() string {
	return *r.HsmClientCertificateIdentifier
}
