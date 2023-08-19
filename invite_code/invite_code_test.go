package invite_code

import (
	"fmt"
	"testing"

	"gotest.tools/assert"
)

func TestIdToCode(t *testing.T) {
	id := 5
	invteCode := DefaultInviteCodeHandler.IdToCode(id)
	fmt.Println(invteCode)
}

func TestCodeToId(t *testing.T) {
	code := "2A99UCYP"
	id := DefaultInviteCodeHandler.CodeToId(code)
	assert.Equal(t, 5, id)
	fmt.Println(id)
}
