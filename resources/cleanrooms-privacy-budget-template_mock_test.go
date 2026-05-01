package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/cleanrooms"
	cleanroomstypes "github.com/aws/aws-sdk-go-v2/service/cleanrooms/types"
)

func Test_Mock_CleanRoomsPrivacyBudgetTemplate_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCleanRoomsClient)

	mockClient.On("ListMemberships", mock.Anything, mock.Anything).
		Return(&cleanrooms.ListMembershipsOutput{
			MembershipSummaries: []cleanroomstypes.MembershipSummary{
				{Id: ptr.String("test-membershipidentifier")},
			},
		}, nil)

	mockClient.On("ListPrivacyBudgetTemplates", mock.Anything, mock.Anything).
		Return(&cleanrooms.ListPrivacyBudgetTemplatesOutput{
			PrivacyBudgetTemplateSummaries: []cleanroomstypes.PrivacyBudgetTemplateSummary{
				{Id: ptr.String("test-privacybudgettemplateidentifier"), CollaborationId: ptr.String("test-collaborationid")},
			},
		}, nil)

	lister := &CleanRoomsPrivacyBudgetTemplateLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCleanRoomsListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*CleanRoomsPrivacyBudgetTemplate)
	a.Equal("test-membershipidentifier", *r.MembershipIdentifier)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CleanRoomsPrivacyBudgetTemplate_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCleanRoomsClient)

	mockClient.On("ListMemberships", mock.Anything, mock.Anything).
		Return(&cleanrooms.ListMembershipsOutput{
			MembershipSummaries: []cleanroomstypes.MembershipSummary{},
		}, nil)

	lister := &CleanRoomsPrivacyBudgetTemplateLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCleanRoomsListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CleanRoomsPrivacyBudgetTemplate_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCleanRoomsClient)

	r := &CleanRoomsPrivacyBudgetTemplate{
		svc:                             mockClient,
		MembershipIdentifier:            ptr.String("test-membershipidentifier"),
		PrivacyBudgetTemplateIdentifier: ptr.String("test-privacybudgettemplateidentifier"),
	}

	mockClient.On("DeletePrivacyBudgetTemplate", mock.Anything,
		&cleanrooms.DeletePrivacyBudgetTemplateInput{
			MembershipIdentifier:            r.MembershipIdentifier,
			PrivacyBudgetTemplateIdentifier: r.PrivacyBudgetTemplateIdentifier,
		}).Return(&cleanrooms.DeletePrivacyBudgetTemplateOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_CleanRoomsPrivacyBudgetTemplate_Properties(t *testing.T) {
	a := assert.New(t)
	r := &CleanRoomsPrivacyBudgetTemplate{
		MembershipIdentifier:            ptr.String("test-membershipidentifier"),
		PrivacyBudgetTemplateIdentifier: ptr.String("test-privacybudgettemplateidentifier"),
		CollaborationID:                 ptr.String("test-collaborationid"),
	}
	props := r.Properties()
	a.Equal("test-membershipidentifier", props.Get("MembershipIdentifier"))
	a.Equal("test-privacybudgettemplateidentifier", props.Get("PrivacyBudgetTemplateIdentifier"))
	a.Equal("test-collaborationid", props.Get("CollaborationID"))
}

func Test_Mock_CleanRoomsPrivacyBudgetTemplate_String(t *testing.T) {
	a := assert.New(t)
	r := &CleanRoomsPrivacyBudgetTemplate{
		PrivacyBudgetTemplateIdentifier: ptr.String("test-privacybudgettemplateidentifier"),
	}
	a.Equal("test-privacybudgettemplateidentifier", r.String())
}
