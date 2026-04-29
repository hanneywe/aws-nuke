package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	cognitotypes "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

func Test_Mock_CognitoUserPoolImportJob_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockCognitoClient)

	mockClient.On("ListUserPools", mock.Anything, mock.Anything).
		Return(&cognitoidentityprovider.ListUserPoolsOutput{
			UserPools: []cognitotypes.UserPoolDescriptionType{
				{Id: ptr.String("us-east-1_pool1")},
			},
		}, nil)

	mockClient.On("ListUserImportJobs", mock.Anything, mock.Anything).
		Return(&cognitoidentityprovider.ListUserImportJobsOutput{
			UserImportJobs: []cognitotypes.UserImportJobType{
				{
					JobId:   ptr.String("import-job-123"),
					JobName: ptr.String("test-import"),
					Status:  cognitotypes.UserImportJobStatusTypeInProgress,
				},
			},
		}, nil)

	lister := &CognitoUserPoolImportJobLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCognitoListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	job := resources[0].(*CognitoUserPoolImportJob)
	assertions.Equal("import-job-123", *job.JobID)
	assertions.Equal("test-import", *job.JobName)
	assertions.Equal("us-east-1_pool1", *job.UserPoolID)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CognitoUserPoolImportJob_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockCognitoClient)

	mockClient.On("ListUserPools", mock.Anything, mock.Anything).
		Return(&cognitoidentityprovider.ListUserPoolsOutput{
			UserPools: []cognitotypes.UserPoolDescriptionType{
				{Id: ptr.String("us-east-1_pool1")},
			},
		}, nil)

	mockClient.On("ListUserImportJobs", mock.Anything, mock.Anything).
		Return(&cognitoidentityprovider.ListUserImportJobsOutput{
			UserImportJobs: []cognitotypes.UserImportJobType{},
		}, nil)

	lister := &CognitoUserPoolImportJobLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCognitoListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CognitoUserPoolImportJob_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockCognitoClient)

	job := &CognitoUserPoolImportJob{
		svc:        mockClient,
		JobID:      ptr.String("import-job-123"),
		JobName:    ptr.String("test-import"),
		UserPoolID: ptr.String("us-east-1_pool1"),
		Status:     ptr.String("InProgress"),
	}

	mockClient.On("StopUserImportJob", mock.Anything, &cognitoidentityprovider.StopUserImportJobInput{
		UserPoolId: job.UserPoolID,
		JobId:      job.JobID,
	}).Return(&cognitoidentityprovider.StopUserImportJobOutput{}, nil)

	err := job.Remove(context.TODO())
	assertions.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CognitoUserPoolImportJob_Properties(t *testing.T) {
	assertions := assert.New(t)

	job := CognitoUserPoolImportJob{
		JobID:      ptr.String("import-job-123"),
		JobName:    ptr.String("test-import"),
		UserPoolID: ptr.String("us-east-1_pool1"),
		Status:     ptr.String("InProgress"),
	}

	properties := job.Properties()
	assertions.Equal("import-job-123", properties.Get("JobId"))
	assertions.Equal("test-import", properties.Get("JobName"))
	assertions.Equal("us-east-1_pool1", properties.Get("UserPoolId"))
	assertions.Equal("InProgress", properties.Get("Status"))
}

func Test_Mock_CognitoUserPoolImportJob_String(t *testing.T) {
	assertions := assert.New(t)
	job := CognitoUserPoolImportJob{JobID: ptr.String("import-job-123")}
	assertions.Equal("import-job-123", job.String())
}

func Test_Mock_CognitoUserPoolImportJob_Filter(t *testing.T) {
	assertions := assert.New(t)

	// Succeeded should be filtered
	succeededJob := CognitoUserPoolImportJob{Status: ptr.String("Succeeded")}
	assertions.Error(succeededJob.Filter())

	// Failed should be filtered
	failedJob := CognitoUserPoolImportJob{Status: ptr.String("Failed")}
	assertions.Error(failedJob.Filter())

	// Expired should be filtered
	expiredJob := CognitoUserPoolImportJob{Status: ptr.String("Expired")}
	assertions.Error(expiredJob.Filter())

	// Stopped should be filtered
	stoppedJob := CognitoUserPoolImportJob{Status: ptr.String("Stopped")}
	assertions.Error(stoppedJob.Filter())

	// InProgress should NOT be filtered
	inProgressJob := CognitoUserPoolImportJob{Status: ptr.String("InProgress")}
	assertions.NoError(inProgressJob.Filter())
}
