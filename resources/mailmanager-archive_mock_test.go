package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/mailmanager"
	mailmanagertypes "github.com/aws/aws-sdk-go-v2/service/mailmanager/types"
)

func Test_Mock_MailManagerArchive_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMailManagerClient)
	mockClient.On("ListArchives", mock.Anything, mock.Anything).
		Return(&mailmanager.ListArchivesOutput{
			Archives: []mailmanagertypes.Archive{
				{
					ArchiveId:    ptr.String("ar-12345"),
					ArchiveName:  ptr.String("my-archive"),
					ArchiveState: mailmanagertypes.ArchiveStateActive,
				},
			},
		}, nil)
	lister := &MailManagerArchiveLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMailManagerListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	ar := resources[0].(*MailManagerArchive)
	a.Equal("ar-12345", *ar.ArchiveID)
	a.Equal("my-archive", *ar.ArchiveName)
	mockClient.AssertExpectations(t)
}

func Test_Mock_MailManagerArchive_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMailManagerClient)
	mockClient.On("ListArchives", mock.Anything, mock.Anything).
		Return(&mailmanager.ListArchivesOutput{Archives: []mailmanagertypes.Archive{}}, nil)
	lister := &MailManagerArchiveLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMailManagerListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_MailManagerArchive_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMailManagerClient)
	ar := &MailManagerArchive{svc: mockClient, ArchiveID: ptr.String("ar-12345")}
	mockClient.On("DeleteArchive", mock.Anything, &mailmanager.DeleteArchiveInput{ArchiveId: ar.ArchiveID}).
		Return(&mailmanager.DeleteArchiveOutput{}, nil)
	a.NoError(ar.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_MailManagerArchive_Filter_PendingDeletion(t *testing.T) {
	a := assert.New(t)
	ar := MailManagerArchive{
		ArchiveID:    ptr.String("ar-12345"),
		ArchiveName:  ptr.String("my-archive"),
		ArchiveState: mailmanagertypes.ArchiveStatePendingDeletion,
	}
	a.Error(ar.Filter())
	a.Contains(ar.Filter().Error(), "already pending deletion")
}

func Test_Mock_MailManagerArchive_Filter_Active(t *testing.T) {
	a := assert.New(t)
	ar := MailManagerArchive{
		ArchiveID:    ptr.String("ar-12345"),
		ArchiveName:  ptr.String("my-archive"),
		ArchiveState: mailmanagertypes.ArchiveStateActive,
	}
	a.NoError(ar.Filter())
}

func Test_Mock_MailManagerArchive_Properties(t *testing.T) {
	a := assert.New(t)
	ar := MailManagerArchive{ArchiveID: ptr.String("ar-12345"), ArchiveName: ptr.String("my-archive")}
	a.Equal("my-archive", ar.Properties().Get("ArchiveName"))
}

func Test_Mock_MailManagerArchive_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("my-archive", (&MailManagerArchive{ArchiveName: ptr.String("my-archive")}).String())
}
