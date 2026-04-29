package resources

import (
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
)

func Test_Mock_EC2ManagedPrefixList_Filter_DeleteFailed(t *testing.T) {
	a := assert.New(t)
	r := &EC2ManagedPrefixList{
		PrefixListID: ptr.String("pl-failed"),
		OwnerID:      ptr.String("123456789012"),
		State:        "delete-failed",
		accountID:    ptr.String("123456789012"),
	}
	err := r.Filter()
	a.Error(err)
	a.Contains(err.Error(), "delete failed")
}
