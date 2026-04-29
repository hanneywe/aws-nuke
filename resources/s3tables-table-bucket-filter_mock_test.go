package resources

import (
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
)

func Test_Mock_S3TableBucket_Filter_AWSManaged(t *testing.T) {
	a := assert.New(t)
	r := &S3TableBucket{
		Name:       ptr.String("aws-cloudwatch"),
		BucketType: "aws",
	}
	err := r.Filter()
	a.Error(err)
	a.Contains(err.Error(), "AWS-managed")
}

func Test_Mock_S3TableBucket_Filter_Customer(t *testing.T) {
	a := assert.New(t)
	r := &S3TableBucket{
		Name:       ptr.String("my-bucket"),
		BucketType: "customer",
	}
	a.NoError(r.Filter())
}
