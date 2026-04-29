package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/mailmanager"
)

// MailManagerClient is the interface for the SES Mail Manager SDK client methods.
type MailManagerClient interface {
	ListAddonInstances(ctx context.Context, params *mailmanager.ListAddonInstancesInput,
		optFns ...func(*mailmanager.Options)) (*mailmanager.ListAddonInstancesOutput, error)
	DeleteAddonInstance(ctx context.Context, params *mailmanager.DeleteAddonInstanceInput,
		optFns ...func(*mailmanager.Options)) (*mailmanager.DeleteAddonInstanceOutput, error)
	ListAddonSubscriptions(ctx context.Context, params *mailmanager.ListAddonSubscriptionsInput,
		optFns ...func(*mailmanager.Options)) (*mailmanager.ListAddonSubscriptionsOutput, error)
	DeleteAddonSubscription(ctx context.Context, params *mailmanager.DeleteAddonSubscriptionInput,
		optFns ...func(*mailmanager.Options)) (*mailmanager.DeleteAddonSubscriptionOutput, error)
	ListAddressLists(ctx context.Context, params *mailmanager.ListAddressListsInput,
		optFns ...func(*mailmanager.Options)) (*mailmanager.ListAddressListsOutput, error)
	DeleteAddressList(ctx context.Context, params *mailmanager.DeleteAddressListInput,
		optFns ...func(*mailmanager.Options)) (*mailmanager.DeleteAddressListOutput, error)
	ListArchives(ctx context.Context, params *mailmanager.ListArchivesInput,
		optFns ...func(*mailmanager.Options)) (*mailmanager.ListArchivesOutput, error)
	DeleteArchive(ctx context.Context, params *mailmanager.DeleteArchiveInput,
		optFns ...func(*mailmanager.Options)) (*mailmanager.DeleteArchiveOutput, error)
	ListRelays(ctx context.Context, params *mailmanager.ListRelaysInput,
		optFns ...func(*mailmanager.Options)) (*mailmanager.ListRelaysOutput, error)
	DeleteRelay(ctx context.Context, params *mailmanager.DeleteRelayInput,
		optFns ...func(*mailmanager.Options)) (*mailmanager.DeleteRelayOutput, error)
	ListRuleSets(ctx context.Context, params *mailmanager.ListRuleSetsInput,
		optFns ...func(*mailmanager.Options)) (*mailmanager.ListRuleSetsOutput, error)
	DeleteRuleSet(ctx context.Context, params *mailmanager.DeleteRuleSetInput,
		optFns ...func(*mailmanager.Options)) (*mailmanager.DeleteRuleSetOutput, error)
}
