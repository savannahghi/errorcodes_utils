package errorcodeutil_test

import (
	"testing"

	"github.com/savannahghi/errorcodeutil"
	"github.com/stretchr/testify/assert"
)

func TestModelHasCustomError(t *testing.T) {
	customerrors := errorcodeutil.CustomError{}
	cr := customerrors.Error()
	assert.NotNil(t, cr)
}
