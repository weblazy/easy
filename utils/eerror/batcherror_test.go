package eerror

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	err1 = "first error"
	err2 = "second error"
)

func TestBatchErrorNil(t *testing.T) {
	var batch BatchError
	assert.Nil(t, batch)
}

func TestBatchErrorOneError(t *testing.T) {
	var batch BatchError
	batch = append(batch, errors.New(err1))
	assert.NotNil(t, batch)
	assert.Equal(t, err1, batch.Error())
}

func TestBatchErrorWithErrors(t *testing.T) {
	var batch BatchError
	batch = append(batch, errors.New(err1))
	batch = append(batch, errors.New(err2))
	assert.NotNil(t, batch)
	assert.Equal(t, fmt.Sprintf("%s\n%s", err1, err2), batch.Error())
}
