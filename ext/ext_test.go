package ext

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateKeyValue(t *testing.T) {
	assert.Equal(t, X_NOTEFRED+"-foo:bar", CreateKeyValue("foo", "bar"))
}
