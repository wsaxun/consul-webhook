package pkg

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetAliCode(t *testing.T) {
	context, _ := gin.CreateTestContext(httptest.NewRecorder())
	context.Request, _ = http.NewRequest(
		"POST",
		"/",
		strings.NewReader(`{
		"object_kind": "tag_push",
		"ref": "refs/tags/v1.0.0"
	}`),
	)
	context.Request.Header.Set("X-Codeup-Event", "Tag Push Hook")
	context.Request.Header.Set("Content-Type", "application/json; charset=utf-8")
	mockCode := AliCode{
		Ref: "refs/tags/v1.0.0",
	}

	code, err := GetAliCode(context)
	if err != nil {
		t.Errorf("GetAliCode returned an error: %v", err)
	}
	if code.Ref != mockCode.Ref {
		t.Errorf("GetAliCode returned incorrect AliCode object: got %v, want %v", code, mockCode)
	}
}

func TestGetCurrentTag(t *testing.T) {
	mockCode := AliCode{
		Ref: "refs/tags/v1.0.0",
	}
	tag := GetCurrentTag(mockCode)
	if tag != "v1.0.0" {
		t.Errorf("GetCurrentTag returned an error: %v", tag)
	}
}
