package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/omics"
)

// OmicsClient is an interface for the Omics SDK client methods used by all Omics resources.
// It enables mock testing of List and Remove operations.
type OmicsClient interface {
	// Listing
	ListReferenceStores(ctx context.Context, params *omics.ListReferenceStoresInput,
		optFns ...func(*omics.Options)) (*omics.ListReferenceStoresOutput, error)
	ListReferences(ctx context.Context, params *omics.ListReferencesInput,
		optFns ...func(*omics.Options)) (*omics.ListReferencesOutput, error)
	ListRunGroups(ctx context.Context, params *omics.ListRunGroupsInput,
		optFns ...func(*omics.Options)) (*omics.ListRunGroupsOutput, error)
	ListSequenceStores(ctx context.Context, params *omics.ListSequenceStoresInput,
		optFns ...func(*omics.Options)) (*omics.ListSequenceStoresOutput, error)
	ListReadSets(ctx context.Context, params *omics.ListReadSetsInput,
		optFns ...func(*omics.Options)) (*omics.ListReadSetsOutput, error)
	ListWorkflows(ctx context.Context, params *omics.ListWorkflowsInput,
		optFns ...func(*omics.Options)) (*omics.ListWorkflowsOutput, error)

	// Deletion
	DeleteReferenceStore(ctx context.Context, params *omics.DeleteReferenceStoreInput,
		optFns ...func(*omics.Options)) (*omics.DeleteReferenceStoreOutput, error)
	DeleteReference(ctx context.Context, params *omics.DeleteReferenceInput,
		optFns ...func(*omics.Options)) (*omics.DeleteReferenceOutput, error)
	DeleteRunGroup(ctx context.Context, params *omics.DeleteRunGroupInput,
		optFns ...func(*omics.Options)) (*omics.DeleteRunGroupOutput, error)
	DeleteSequenceStore(ctx context.Context, params *omics.DeleteSequenceStoreInput,
		optFns ...func(*omics.Options)) (*omics.DeleteSequenceStoreOutput, error)
	BatchDeleteReadSet(ctx context.Context, params *omics.BatchDeleteReadSetInput,
		optFns ...func(*omics.Options)) (*omics.BatchDeleteReadSetOutput, error)
	DeleteWorkflow(ctx context.Context, params *omics.DeleteWorkflowInput,
		optFns ...func(*omics.Options)) (*omics.DeleteWorkflowOutput, error)
}
