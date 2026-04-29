package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	snstypes "github.com/aws/aws-sdk-go-v2/service/sns/types"
)

func Test_Mock_SNSDataProtectionPolicy_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSNSV2Client)

	mockClient.On("ListTopics", mock.Anything, mock.Anything).
		Return(&sns.ListTopicsOutput{
			Topics: []snstypes.Topic{
				{TopicArn: ptr.String("arn:aws:sns:us-east-1:123456789012:my-topic")},
			},
		}, nil)

	mockClient.On("GetDataProtectionPolicy", mock.Anything, &sns.GetDataProtectionPolicyInput{
		ResourceArn: ptr.String("arn:aws:sns:us-east-1:123456789012:my-topic"),
	}).Return(&sns.GetDataProtectionPolicyOutput{
		DataProtectionPolicy: ptr.String(`{"Name":"policy"}`),
	}, nil)

	lister := &SNSDataProtectionPolicyLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSNSV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*SNSDataProtectionPolicy)
	a.Equal("arn:aws:sns:us-east-1:123456789012:my-topic", *r.TopicArn)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SNSDataProtectionPolicy_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSNSV2Client)

	mockClient.On("ListTopics", mock.Anything, mock.Anything).
		Return(&sns.ListTopicsOutput{
			Topics: []snstypes.Topic{},
		}, nil)

	lister := &SNSDataProtectionPolicyLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSNSV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SNSDataProtectionPolicy_List_NoPolicy(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSNSV2Client)

	mockClient.On("ListTopics", mock.Anything, mock.Anything).
		Return(&sns.ListTopicsOutput{
			Topics: []snstypes.Topic{
				{TopicArn: ptr.String("arn:aws:sns:us-east-1:123456789012:no-policy-topic")},
			},
		}, nil)

	mockClient.On("GetDataProtectionPolicy", mock.Anything, &sns.GetDataProtectionPolicyInput{
		ResourceArn: ptr.String("arn:aws:sns:us-east-1:123456789012:no-policy-topic"),
	}).Return(&sns.GetDataProtectionPolicyOutput{
		DataProtectionPolicy: ptr.String(""),
	}, nil)

	lister := &SNSDataProtectionPolicyLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSNSV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SNSDataProtectionPolicy_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSNSV2Client)

	r := &SNSDataProtectionPolicy{
		svc:      mockClient,
		TopicArn: ptr.String("arn:aws:sns:us-east-1:123456789012:my-topic"),
	}

	mockClient.On("PutDataProtectionPolicy", mock.Anything,
		&sns.PutDataProtectionPolicyInput{
			ResourceArn:          r.TopicArn,
			DataProtectionPolicy: ptr.String(""),
		}).Return(&sns.PutDataProtectionPolicyOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_SNSDataProtectionPolicy_Properties(t *testing.T) {
	a := assert.New(t)
	r := &SNSDataProtectionPolicy{
		TopicArn: ptr.String("arn:aws:sns:us-east-1:123456789012:my-topic"),
	}
	props := r.Properties()
	a.Equal("arn:aws:sns:us-east-1:123456789012:my-topic", props.Get("TopicArn"))
}

func Test_Mock_SNSDataProtectionPolicy_String(t *testing.T) {
	a := assert.New(t)
	r := &SNSDataProtectionPolicy{
		TopicArn: ptr.String("arn:aws:sns:us-east-1:123456789012:my-topic"),
	}
	a.Equal("arn:aws:sns:us-east-1:123456789012:my-topic", r.String())
}
