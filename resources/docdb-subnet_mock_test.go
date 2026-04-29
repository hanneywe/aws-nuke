package resources

import (
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
)

func Test_Mock_DocDBSubnetGroup_Filter_Default(t *testing.T) {
	a := assert.New(t)
	r := DocDBSubnetGroup{Name: ptr.String("default")}
	err := r.Filter()
	a.Error(err)
	a.Contains(err.Error(), "cannot delete default subnet group")
}

func Test_Mock_DocDBSubnetGroup_Filter_Custom(t *testing.T) {
	a := assert.New(t)
	r := DocDBSubnetGroup{Name: ptr.String("my-custom-subnet-group")}
	a.NoError(r.Filter())
}

func Test_Mock_DocDBSubnetGroup_Properties(t *testing.T) {
	a := assert.New(t)
	r := DocDBSubnetGroup{Name: ptr.String("my-subnet-group")}
	a.Equal("my-subnet-group", r.Properties().Get("Name"))
}

func Test_Mock_DocDBSubnetGroup_String(t *testing.T) {
	a := assert.New(t)
	r := DocDBSubnetGroup{Name: ptr.String("my-subnet-group")}
	a.Equal("my-subnet-group", r.String())
}
