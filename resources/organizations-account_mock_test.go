package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/organizations"
	orgtypes "github.com/aws/aws-sdk-go-v2/service/organizations/types"

	libsettings "github.com/ekristen/libnuke/pkg/settings"
)

func Test_Mock_OrganizationsAccount_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockOrganizationsClient)

	mockClient.On("DescribeOrganization", mock.Anything, mock.Anything).
		Return(&organizations.DescribeOrganizationOutput{
			Organization: &orgtypes.Organization{
				MasterAccountId: ptr.String("111111111111"),
			},
		}, nil)

	mockClient.On("ListAccounts", mock.Anything, mock.Anything).
		Return(&organizations.ListAccountsOutput{
			Accounts: []orgtypes.Account{
				{
					Id:     ptr.String("222222222222"),
					Name:   ptr.String("member-account"),
					Email:  ptr.String("member@example.com"),
					Status: orgtypes.AccountStatusActive,
				},
			},
		}, nil)

	lister := &OrganizationsAccountLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testOrganizationsListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	account := resources[0].(*OrganizationsAccount)
	a.Equal("222222222222", *account.ID)
	a.Equal("member-account", *account.Name)
	a.Equal("member@example.com", *account.Email)
	a.Equal(orgtypes.AccountStatusActive, account.Status)
	a.False(*account.ManagementAccount)
	mockClient.AssertExpectations(t)
}

func Test_Mock_OrganizationsAccount_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockOrganizationsClient)

	mockClient.On("DescribeOrganization", mock.Anything, mock.Anything).
		Return(&organizations.DescribeOrganizationOutput{
			Organization: &orgtypes.Organization{
				MasterAccountId: ptr.String("111111111111"),
			},
		}, nil)

	mockClient.On("ListAccounts", mock.Anything, mock.Anything).
		Return(&organizations.ListAccountsOutput{
			Accounts: []orgtypes.Account{},
		}, nil)

	lister := &OrganizationsAccountLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testOrganizationsListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_OrganizationsAccount_List_ManagementAccount(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockOrganizationsClient)

	mockClient.On("DescribeOrganization", mock.Anything, mock.Anything).
		Return(&organizations.DescribeOrganizationOutput{
			Organization: &orgtypes.Organization{
				MasterAccountId: ptr.String("111111111111"),
			},
		}, nil)

	mockClient.On("ListAccounts", mock.Anything, mock.Anything).
		Return(&organizations.ListAccountsOutput{
			Accounts: []orgtypes.Account{
				{
					Id:     ptr.String("111111111111"),
					Name:   ptr.String("management-account"),
					Email:  ptr.String("mgmt@example.com"),
					Status: orgtypes.AccountStatusActive,
				},
				{
					Id:     ptr.String("222222222222"),
					Name:   ptr.String("member-account"),
					Email:  ptr.String("member@example.com"),
					Status: orgtypes.AccountStatusActive,
				},
			},
		}, nil)

	lister := &OrganizationsAccountLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testOrganizationsListerOpts)
	a.NoError(err)
	a.Len(resources, 2)

	mgmt := resources[0].(*OrganizationsAccount)
	a.True(*mgmt.ManagementAccount)

	member := resources[1].(*OrganizationsAccount)
	a.False(*member.ManagementAccount)
	mockClient.AssertExpectations(t)
}

func Test_Mock_OrganizationsAccount_Remove_WithSetting(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockOrganizationsClient)

	settings := &libsettings.Setting{}
	settings.Set("CloseAccount", true)

	account := &OrganizationsAccount{
		svc:      mockClient,
		ID:       ptr.String("222222222222"),
		settings: settings,
	}

	mockClient.On("CloseAccount", mock.Anything, &organizations.CloseAccountInput{
		AccountId: account.ID,
	}).Return(&organizations.CloseAccountOutput{}, nil)

	a.NoError(account.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_OrganizationsAccount_Remove_WithoutSetting(t *testing.T) {
	a := assert.New(t)

	settings := &libsettings.Setting{}
	settings.Set("CloseAccount", false)

	account := &OrganizationsAccount{
		ID:       ptr.String("222222222222"),
		settings: settings,
	}

	err := account.Remove(context.TODO())
	a.Error(err)
	a.Contains(err.Error(), "CloseAccount setting must be enabled")
}

func Test_Mock_OrganizationsAccount_Filter_Suspended(t *testing.T) {
	a := assert.New(t)
	account := OrganizationsAccount{
		ID:     ptr.String("222222222222"),
		Status: orgtypes.AccountStatusSuspended,
	}
	a.Error(account.Filter())
}

func Test_Mock_OrganizationsAccount_Filter_PendingClosure(t *testing.T) {
	a := assert.New(t)
	account := OrganizationsAccount{
		ID:     ptr.String("222222222222"),
		Status: orgtypes.AccountStatusPendingClosure,
	}
	a.Error(account.Filter())
}

func Test_Mock_OrganizationsAccount_Filter_ManagementAccount(t *testing.T) {
	a := assert.New(t)
	account := OrganizationsAccount{
		ID:                ptr.String("111111111111"),
		Status:            orgtypes.AccountStatusActive,
		ManagementAccount: ptr.Bool(true),
	}
	a.Error(account.Filter())
}

func Test_Mock_OrganizationsAccount_Filter_Active(t *testing.T) {
	a := assert.New(t)
	account := OrganizationsAccount{
		ID:                ptr.String("222222222222"),
		Status:            orgtypes.AccountStatusActive,
		ManagementAccount: ptr.Bool(false),
	}
	a.NoError(account.Filter())
}

func Test_Mock_OrganizationsAccount_Properties(t *testing.T) {
	a := assert.New(t)
	account := OrganizationsAccount{
		ID:                ptr.String("222222222222"),
		Name:              ptr.String("member-account"),
		Email:             ptr.String("member@example.com"),
		Status:            orgtypes.AccountStatusActive,
		ManagementAccount: ptr.Bool(false),
	}
	props := account.Properties()
	a.Equal("222222222222", props.Get("ID"))
	a.Equal("member-account", props.Get("Name"))
	a.Equal("member@example.com", props.Get("Email"))
	a.Equal("ACTIVE", props.Get("Status"))
	a.Equal("false", props.Get("ManagementAccount"))
}

func Test_Mock_OrganizationsAccount_String(t *testing.T) {
	a := assert.New(t)
	account := OrganizationsAccount{
		ID: ptr.String("222222222222"),
	}
	a.Equal("222222222222", account.String())
}
