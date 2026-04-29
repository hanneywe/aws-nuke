package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/route53recoverycontrolconfig"
)

// Route53RecoveryControlConfigClient is the interface for the Route53 Recovery Control Config SDK client methods.
type Route53RecoveryControlConfigClient interface {
	ListClusters(ctx context.Context, params *route53recoverycontrolconfig.ListClustersInput,
		optFns ...func(*route53recoverycontrolconfig.Options)) (*route53recoverycontrolconfig.ListClustersOutput, error)
	DeleteCluster(ctx context.Context, params *route53recoverycontrolconfig.DeleteClusterInput,
		optFns ...func(*route53recoverycontrolconfig.Options)) (*route53recoverycontrolconfig.DeleteClusterOutput, error)
	ListControlPanels(ctx context.Context, params *route53recoverycontrolconfig.ListControlPanelsInput,
		optFns ...func(*route53recoverycontrolconfig.Options)) (*route53recoverycontrolconfig.ListControlPanelsOutput, error)
	DeleteControlPanel(ctx context.Context, params *route53recoverycontrolconfig.DeleteControlPanelInput,
		optFns ...func(*route53recoverycontrolconfig.Options)) (*route53recoverycontrolconfig.DeleteControlPanelOutput, error)
	ListRoutingControls(ctx context.Context, params *route53recoverycontrolconfig.ListRoutingControlsInput,
		optFns ...func(*route53recoverycontrolconfig.Options)) (*route53recoverycontrolconfig.ListRoutingControlsOutput, error)
	DeleteRoutingControl(ctx context.Context, params *route53recoverycontrolconfig.DeleteRoutingControlInput,
		optFns ...func(*route53recoverycontrolconfig.Options)) (*route53recoverycontrolconfig.DeleteRoutingControlOutput, error)
}
