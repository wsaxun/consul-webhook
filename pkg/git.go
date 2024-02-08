package pkg

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"

	"consul-webhook/config"
)

type Repo struct {
	Path        string
	GitRepo     string
	UserName    string
	Password    string
	UpdateFiles []string
	DeleteFiles []string
	RepoObj     *git.Repository
}

// Clone 清空缓存目录后执行clone操作
func (r *Repo) Clone() error {
	repoPath := GetTmpPath()
	_ = os.RemoveAll(repoPath)
	cloneRepo, err := git.PlainClone(repoPath, false, &git.CloneOptions{
		URL: r.GitRepo,
		Auth: &http.BasicAuth{
			Username: r.UserName,
			Password: r.Password,
		},
	})
	r.RepoObj = cloneRepo
	return err
}

// Checkout 克隆代码后checkout到指定tag
func (r *Repo) Checkout(tag string) error {
	err := r.Clone()
	if err != nil {
		log.Panicf("clone error: %v\n", err.Error())
		return err
	}
	worktree, _ := r.RepoObj.Worktree()
	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewTagReferenceName(tag),
	})
	return err
}

// GetOneCommitFiles 根据commitID 获取变化的文件
// 判断文件状态的依据是 全部都是 - 符号则为该文件被删除，  全部都是 + 符号则为该文件是新增的
// 带 + 和 - 号则为此文件编辑过
func (r *Repo) GetOneCommitFiles(commitID string) ([]string, []string, error) {
	if r.RepoObj == nil {
		err := r.Clone()
		if err != nil {
			log.Panicf("clone error: %v\n", err.Error())
			return []string{}, []string{}, err
		}
	}
	c := r.getCommitIter(commitID)
	err := c.ForEach(func(commit *object.Commit) error {
		f, _ := commit.Stats()
		if commit.ID().String() == commitID {
			l := strings.Split(f.String(), "|")
			item := strings.TrimSpace(l[0])
			status := l[1]
			if strings.Contains(status, "+") {
				r.UpdateFiles = append(r.UpdateFiles, item)
			} else {
				r.DeleteFiles = append(r.DeleteFiles, item)
			}
		}
		return nil
	})
	return r.UpdateFiles, r.DeleteFiles, err
}

// getCommitIter Log方法无法获取特定commitid的历史，故使用指定时间来减少获取数据条目数
func (r *Repo) getCommitIter(commitID string) object.CommitIter {
	const util = "-720h"
	now := time.Now()
	duration, _ := time.ParseDuration(util)
	old := now.Add(duration)
	c, _ := r.RepoObj.Log(&git.LogOptions{From: plumbing.NewHash(commitID),
		Order: git.LogOrderCommitterTime,
		Since: &old,
	})
	return c
}

// NewRepo 创建repo对象
func NewRepo() *Repo {
	consulConfig, err := GetConsulConfig()
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	repo := Repo{
		Path:     consulConfig.Repo.Path,
		GitRepo:  consulConfig.Repo.Url,
		UserName: consulConfig.Repo.Username,
		Password: consulConfig.Repo.Password,
	}
	return &repo
}

func InitRepo() {
	repo := NewRepo()
	localConf := config.GetConfig()
	localConf.Repo = &config.Repo{
		Password: repo.Password,
		Username: repo.UserName,
		Url:      repo.GitRepo,
		Path:     repo.Path,
	}
}
