package resources

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/s3control"
	s3controltypes "github.com/aws/aws-sdk-go-v2/service/s3control/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const S3PublicAccessBlockResource = "S3PublicAccessBlock"

func init() {
	registry.Register(&registry.Registration{
		Name:     S3PublicAccessBlockResource,
		Scope:    nuke.Account,
		Resource: &S3PublicAccessBlock{},
		Lister:   &S3PublicAccessBlockLister{},
	})
}

type S3PublicAccessBlockLister struct {
	svc S3AccountClient
}

func (l *S3PublicAccessBlockLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = s3control.NewFromConfig(*opts.Config)
	}

	resp, err := svc.GetPublicAccessBlock(ctx, &s3control.GetPublicAccessBlockInput{
		AccountId: opts.AccountID,
	})
	if err != nil {
		var noPAB *s3controltypes.NoSuchPublicAccessBlockConfiguration
		if errors.As(err, &noPAB) {
			return nil, nil
		}
		return nil, err
	}

	if resp.PublicAccessBlockConfiguration == nil {
		return nil, nil
	}

	cfg := resp.PublicAccessBlockConfiguration
	return []resource.Resource{
		&S3PublicAccessBlock{
			svc:                   svc,
			accountID:             opts.AccountID,
			BlockPublicAcls:       cfg.BlockPublicAcls,
			BlockPublicPolicy:     cfg.BlockPublicPolicy,
			IgnorePublicAcls:      cfg.IgnorePublicAcls,
			RestrictPublicBuckets: cfg.RestrictPublicBuckets,
		},
	}, nil
}

type S3PublicAccessBlock struct {
	svc                   S3AccountClient
	accountID             *string
	BlockPublicAcls       *bool
	BlockPublicPolicy     *bool
	IgnorePublicAcls      *bool
	RestrictPublicBuckets *bool
}

func (r *S3PublicAccessBlock) Remove(ctx context.Context) error {
	_, err := r.svc.DeletePublicAccessBlock(ctx, &s3control.DeletePublicAccessBlockInput{
		AccountId: r.accountID,
	})
	return err
}

func (r *S3PublicAccessBlock) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *S3PublicAccessBlock) String() string {
	return "S3PublicAccessBlock"
}
