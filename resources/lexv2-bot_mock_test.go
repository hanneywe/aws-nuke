package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/lexmodelsv2"
	lexmodelsv2types "github.com/aws/aws-sdk-go-v2/service/lexmodelsv2/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

var testLexV2ListerOpts = &nuke.ListerOpts{
	Logger: logrus.WithField("test", true),
}

func Test_Mock_LexV2Bot_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockLexV2Client)

	mockClient.
		On("ListBots", mock.Anything, mock.Anything).
		Return(
			&lexmodelsv2.ListBotsOutput{
				BotSummaries: []lexmodelsv2types.BotSummary{
					{
						BotId:   ptr.String("bot-12345"),
						BotName: ptr.String("test-bot"),
					},
				},
			}, nil,
		)

	lister := &LexV2BotLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testLexV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	bot := resources[0].(*LexV2Bot)
	assertions.Equal("bot-12345", *bot.BotID)
	assertions.Equal("test-bot", *bot.BotName)

	mockClient.AssertExpectations(t)
}

func Test_Mock_LexV2Bot_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockLexV2Client)

	mockClient.
		On("ListBots", mock.Anything, mock.Anything).
		Return(
			&lexmodelsv2.ListBotsOutput{
				BotSummaries: []lexmodelsv2types.BotSummary{},
			}, nil,
		)

	lister := &LexV2BotLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testLexV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_LexV2Bot_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockLexV2Client)

	bot := &LexV2Bot{
		svc:     mockClient,
		BotID:   ptr.String("bot-12345"),
		BotName: ptr.String("test-bot"),
	}

	mockClient.
		On(
			"DeleteBot",
			mock.Anything,
			&lexmodelsv2.DeleteBotInput{
				BotId:                  bot.BotID,
				SkipResourceInUseCheck: true,
			},
		).
		Return(&lexmodelsv2.DeleteBotOutput{}, nil)

	err := bot.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_LexV2Bot_Properties(t *testing.T) {
	assertions := assert.New(t)

	bot := LexV2Bot{
		BotID:   ptr.String("bot-12345"),
		BotName: ptr.String("test-bot"),
	}

	properties := bot.Properties()

	assertions.Equal("bot-12345", properties.Get("BotId"))
	assertions.Equal("test-bot", properties.Get("BotName"))
}

func Test_Mock_LexV2Bot_String(t *testing.T) {
	assertions := assert.New(t)

	bot := LexV2Bot{
		BotID:   ptr.String("bot-12345"),
		BotName: ptr.String("test-bot"),
	}

	assertions.Equal("test-bot", bot.String())
}
