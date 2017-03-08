package api_test

import (
	"testing"

	"github.com/oakmail/backend/pkg/api/test"
)

func TestAPI(t *testing.T) {
	api := test.InitAPI()
	defer api.Cleanup()
}
