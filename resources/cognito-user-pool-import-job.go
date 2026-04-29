package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	cognitotypes "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const CognitoUserPoolImportJobResource = "CognitoUserPoolImportJob"

func init() {
	registry.Register(&registry.Registration{
		Name:     CognitoUserPoolImportJobResource,
		Scope:    nuke.Account,
		Resource: &CognitoUserPoolImportJob{},
		Lister:   &CognitoUserPoolImportJobLister{},
	})
}

type CognitoUserPoolImportJobLister struct {
	svc CognitoClient
}

func (l *CognitoUserPoolImportJobLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = cognitoidentityprovider.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	// Step 1: List all user pools
	poolParams := &cognitoidentityprovider.ListUserPoolsInput{
		MaxResults: aws.Int32(60),
	}
	for {
		poolOutput, err := svc.ListUserPools(ctx, poolParams)
		if err != nil {
			return nil, err
		}

		// Step 2: For each pool, list import jobs
		for _, pool := range poolOutput.UserPools {
			jobParams := &cognitoidentityprovider.ListUserImportJobsInput{
				UserPoolId: pool.Id,
				MaxResults: aws.Int32(60),
			}

			jobOutput, err := svc.ListUserImportJobs(ctx, jobParams)
			if err != nil {
				return nil, err
			}

			for _, job := range jobOutput.UserImportJobs {
				jobStatus := string(job.Status)
				resources = append(resources, &CognitoUserPoolImportJob{
					svc:        svc,
					JobID:      job.JobId,
					JobName:    job.JobName,
					UserPoolID: pool.Id,
					Status:     &jobStatus,
				})
			}
		}

		if poolOutput.NextToken == nil {
			break
		}
		poolParams.NextToken = poolOutput.NextToken
	}

	return resources, nil
}

type CognitoUserPoolImportJob struct {
	svc        CognitoClient
	JobID      *string `property:"name=JobId"`
	JobName    *string
	UserPoolID *string `property:"name=UserPoolId"`
	Status     *string
}

func (r *CognitoUserPoolImportJob) Remove(ctx context.Context) error {
	_, err := r.svc.StopUserImportJob(ctx, &cognitoidentityprovider.StopUserImportJobInput{
		UserPoolId: r.UserPoolID,
		JobId:      r.JobID,
	})
	return err
}

func (r *CognitoUserPoolImportJob) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *CognitoUserPoolImportJob) String() string {
	return *r.JobID
}

func (r *CognitoUserPoolImportJob) Filter() error {
	if r.Status == nil {
		return nil
	}
	switch cognitotypes.UserImportJobStatusType(*r.Status) {
	case cognitotypes.UserImportJobStatusTypeSucceeded,
		cognitotypes.UserImportJobStatusTypeFailed,
		cognitotypes.UserImportJobStatusTypeExpired,
		cognitotypes.UserImportJobStatusTypeStopped:
		return fmt.Errorf("already %s", *r.Status)
	}
	return nil
}
