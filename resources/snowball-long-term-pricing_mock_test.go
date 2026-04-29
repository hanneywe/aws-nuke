package resources

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/snowball"
	snowballtypes "github.com/aws/aws-sdk-go-v2/service/snowball/types"
	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

func Test_Mock_SnowballLongTermPricing_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSnowballClient)

	now := time.Now()
	mockClient.On("ListLongTermPricing", mock.Anything, mock.Anything).
		Return(&snowball.ListLongTermPricingOutput{
			LongTermPricingEntries: []snowballtypes.LongTermPricingListEntry{
				{
					LongTermPricingId:          ptr.String("ltp-123"),
					LongTermPricingStatus:      ptr.String("Active"),
					IsLongTermPricingAutoRenew: ptr.Bool(true),
					LongTermPricingStartDate:   &now,
					LongTermPricingEndDate:     &now,
				},
				{
					LongTermPricingId:          ptr.String("ltp-456"),
					LongTermPricingStatus:      ptr.String("Active"),
					IsLongTermPricingAutoRenew: ptr.Bool(false),
				},
			},
		}, nil)

	lister := &SnowballLongTermPricingLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), &nuke.ListerOpts{})
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*SnowballLongTermPricing)
	a.Equal("ltp-123", *r.LongTermPricingID)
	a.Equal("Active", *r.Status)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SnowballLongTermPricing_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSnowballClient)

	mockClient.On("ListLongTermPricing", mock.Anything, mock.Anything).
		Return(&snowball.ListLongTermPricingOutput{
			LongTermPricingEntries: []snowballtypes.LongTermPricingListEntry{},
		}, nil)

	lister := &SnowballLongTermPricingLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), &nuke.ListerOpts{})
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SnowballLongTermPricing_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSnowballClient)

	r := &SnowballLongTermPricing{
		svc:               mockClient,
		LongTermPricingID: ptr.String("ltp-123"),
	}

	mockClient.On("UpdateLongTermPricing", mock.Anything, &snowball.UpdateLongTermPricingInput{
		LongTermPricingId:          r.LongTermPricingID,
		IsLongTermPricingAutoRenew: aws.Bool(false),
	}).Return(&snowball.UpdateLongTermPricingOutput{}, nil)

	err := r.Remove(context.TODO())
	a.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SnowballLongTermPricing_Properties(t *testing.T) {
	a := assert.New(t)

	r := SnowballLongTermPricing{
		LongTermPricingID: ptr.String("ltp-123"),
		Status:            ptr.String("Active"),
	}

	props := r.Properties()
	a.Equal("ltp-123", props.Get("LongTermPricingID"))
	a.Equal("Active", props.Get("Status"))
}

func Test_Mock_SnowballLongTermPricing_String(t *testing.T) {
	a := assert.New(t)
	r := SnowballLongTermPricing{
		LongTermPricingID: ptr.String("ltp-123"),
	}
	a.Equal("ltp-123", r.String())
}
