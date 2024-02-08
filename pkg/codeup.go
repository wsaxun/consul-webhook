package pkg

import (
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// AliCode 阿里云 codeup webhook推送的消息格式
type AliCode struct {
	TotalCommitsCount int    `json:"total_commits_count"`
	UserEmail         string `json:"user_email"`
	Before            string `json:"before"`
	UserExternUID     string `json:"user_extern_uid"`
	UserName          string `json:"user_name"`
	CheckoutSha       string `json:"checkout_sha"`
	Repository        struct {
		Name            string `json:"name"`
		Description     string `json:"description"`
		VisibilityLevel int    `json:"visibility_level"`
		GitHTTPURL      string `json:"git_http_url"`
		GitSSHURL       string `json:"git_ssh_url"`
		URL             string `json:"url"`
		Homepage        string `json:"homepage"`
	} `json:"repository"`
	ObjectKind string `json:"object_kind"`
	Ref        string `json:"ref"`
	ProjectID  int    `json:"project_id"`
	UserID     int    `json:"user_id"`
	Commits    []struct {
		Author struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"author"`
		ID        string    `json:"id"`
		Message   string    `json:"message"`
		URL       string    `json:"url"`
		Timestamp time.Time `json:"timestamp"`
	} `json:"commits"`
	After    string `json:"after"`
	AliyunPk string `json:"aliyun_pk"`
}

// GetAliCode 转换成struct
func GetAliCode(c *gin.Context) (AliCode, error) {
	var body AliCode
	err := errors.New("invalid event")
	header := c.GetHeader("X-Codeup-Event")
	if header != "Tag Push Hook" {
		return body, err
	}
	err = c.ShouldBind(&body)
	return body, err
}

// GetCurrentTag 根据传递过来的消息获取tag
func GetCurrentTag(code AliCode) string {
	return strings.ReplaceAll(code.Ref, "refs/tags/", "")
}

// IsDeleteTag 通过 total_commits_count 字段判断是否删除tag事件
func IsDeleteTag(code AliCode) bool {
	return code.TotalCommitsCount <= 0
}
