package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/gamelift"
	gamelifttypes "github.com/aws/aws-sdk-go-v2/service/gamelift/types"
)

func Test_Mock_GameLiftAlias_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGameLiftV2Client)
	mockClient.On("ListAliases", mock.Anything, mock.Anything).
		Return(&gamelift.ListAliasesOutput{
			Aliases: []gamelifttypes.Alias{
				{AliasId: ptr.String("alias-12345"), Name: ptr.String("my-alias")},
			},
		}, nil)
	lister := &GameLiftAliasLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGameLiftV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	a.Equal("alias-12345", resources[0].(*GameLiftAlias).String())
	mockClient.AssertExpectations(t)
}

func Test_Mock_GameLiftAlias_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGameLiftV2Client)
	mockClient.On("ListAliases", mock.Anything, mock.Anything).
		Return(&gamelift.ListAliasesOutput{Aliases: []gamelifttypes.Alias{}}, nil)
	lister := &GameLiftAliasLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGameLiftV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_GameLiftAlias_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGameLiftV2Client)
	r := &GameLiftAlias{svc: mockClient, AliasID: ptr.String("alias-12345"), Name: ptr.String("my-alias")}
	mockClient.On("DeleteAlias", mock.Anything, &gamelift.DeleteAliasInput{
		AliasId: r.AliasID,
	}).Return(&gamelift.DeleteAliasOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_GameLiftAlias_Properties(t *testing.T) {
	a := assert.New(t)
	r := GameLiftAlias{AliasID: ptr.String("alias-12345"), Name: ptr.String("my-alias")}
	a.Equal("alias-12345", r.Properties().Get("AliasId"))
	a.Equal("my-alias", r.Properties().Get("Name"))
}

func Test_Mock_GameLiftAlias_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("alias-12345", (&GameLiftAlias{AliasID: ptr.String("alias-12345")}).String())
}
