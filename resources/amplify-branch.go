package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/amplify"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const AmplifyBranchResource = "AmplifyBranch"

func init() {
	registry.Register(&registry.Registration{
		Name:     AmplifyBranchResource,
		Scope:    nuke.Account,
		Resource: &AmplifyBranch{},
		Lister:   &AmplifyBranchLister{},
	})
}

type AmplifyBranchLister struct {
	svc AmplifyClient
}

func (l *AmplifyBranchLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = amplify.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	appPaginator := amplify.NewListAppsPaginator(svc, &amplify.ListAppsInput{})
	for appPaginator.HasMorePages() {
		appResp, err := appPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for i := range appResp.Apps {
			branchPaginator := amplify.NewListBranchesPaginator(svc, &amplify.ListBranchesInput{
				AppId: appResp.Apps[i].AppId,
			})
			for branchPaginator.HasMorePages() {
				branchResp, err := branchPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}
				for j := range branchResp.Branches {
					resources = append(resources, &AmplifyBranch{
						svc:        svc,
						BranchName: branchResp.Branches[j].BranchName,
						BranchArn:  branchResp.Branches[j].BranchArn,
						AppID:      appResp.Apps[i].AppId,
					})
				}
			}
		}
	}
	return resources, nil
}

type AmplifyBranch struct {
	svc        AmplifyClient
	BranchName *string
	BranchArn  *string
	AppID      *string `property:"name=AppId"`
}

func (r *AmplifyBranch) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteBranch(ctx, &amplify.DeleteBranchInput{
		AppId:      r.AppID,
		BranchName: r.BranchName,
	})
	return err
}

func (r *AmplifyBranch) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *AmplifyBranch) String() string {
	return *r.BranchName
}
