package v1

import (
	"testing"

	"github.com/ebusiness/go-disney/utils"
)

// go test v1/index*.go -v
func TestIndex(t *testing.T) {

	utils.CreaterTestForHTTP(t, "/test", "/test", index, nil)
}
