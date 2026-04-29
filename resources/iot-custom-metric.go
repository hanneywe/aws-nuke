package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const IoTCustomMetricResource = "IoTCustomMetric"

func init() {
	registry.Register(&registry.Registration{
		Name:     IoTCustomMetricResource,
		Scope:    nuke.Account,
		Resource: &IoTCustomMetric{},
		Lister:   &IoTCustomMetricLister{},
	})
}

type IoTCustomMetricLister struct {
	svc IoTClient
}

func (l *IoTCustomMetricLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = iot.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	params := &iot.ListCustomMetricsInput{
		MaxResults: aws.Int32(100),
	}

	for {
		resp, err := svc.ListCustomMetrics(ctx, params)
		if err != nil {
			return nil, err
		}

		for _, name := range resp.MetricNames {
			resources = append(resources, &IoTCustomMetric{
				svc:        svc,
				MetricName: aws.String(name),
			})
		}

		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}

	return resources, nil
}

type IoTCustomMetric struct {
	svc        IoTClient
	MetricName *string
}

func (r *IoTCustomMetric) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteCustomMetric(ctx, &iot.DeleteCustomMetricInput{
		MetricName: r.MetricName,
	})
	return err
}

func (r *IoTCustomMetric) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *IoTCustomMetric) String() string {
	return *r.MetricName
}
