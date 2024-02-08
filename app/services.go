package app

import (
	"crypto/md5"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"consul-webhook/config"
	"consul-webhook/pkg"
)

// RsyncConsul 私有方法 指定tag同步对应仓库配置到consul
func RsyncConsul(tag string) error {
	conf := config.GetConfig()
	prefix := conf.App.Env + "/hxyljk"
	if err := checkout(tag); err != nil {
		log.Printf("clone error %v", err.Error())
	}
	client := getConsulClient()
	keys, _ := client.Keys(prefix, "")
	// 获取仓库中的文件列表
	originKeys, err := getAllFilesKey(prefix)
	if err != nil {
		return err
	}
	if len(originKeys) == 0 {
		log.Println("仓库无此环境配置")
		return nil
	}

	addCh := make(chan string)
	deleteCh := make(chan string)
	var wg sync.WaitGroup
	wg.Add(2)
	// 更新or增加consul key
	go func() {
		defer wg.Done()
		for key := range addCh {
			if err := client.Put(key); err != nil {
				log.Printf("action: update, key: %v, status: failed", key)
				continue
			}
			log.Printf("action: update, key: %v, status: success", key)
		}
	}()
	// 删除consul key
	go func() {
		defer wg.Done()
		for key := range deleteCh {
			if err := client.Delete(key); err != nil {
				log.Printf("action: delete, key: %v, status: failed", key)
				continue
			}
			log.Printf("action: delete, key: %v, status: success", key)
		}
	}()

	// 获取变化的key放入channel中
	tmpMapT := pkg.ListToMap(keys)
	for _, key := range originKeys {
		if _, ok := tmpMapT[key]; !ok {
			addCh <- key
			continue
		}
		if !diff(key, client) {
			addCh <- key
		}
	}
	// 获取需要从consul中删除的key放入channel中
	tmpMapO := pkg.ListToMap(originKeys)
	for _, key := range keys {
		if _, ok := tmpMapO[key]; !ok {
			deleteCh <- key
		}
	}
	close(addCh)
	close(deleteCh)
	wg.Wait()
	return nil
}

// getConsulClient 私有方法获取consul client
func getConsulClient() *pkg.Client {
	//appConfig := config.GetConfig()
	factory, err := pkg.NewFactory()
	if err != nil {
		log.Printf("consul 连接失败: %v", err.Error())
		return nil
	}
	client := factory.NewClient()
	return client
}

// checkout 私有方法 checkout指定tag
func checkout(tag string) error {
	repo := pkg.NewRepo()
	err := repo.Checkout(tag)
	return err
}

// getAllFilesKey 获取仓库所有文件列表
func getAllFilesKey(prefix string) ([]string, error) {
	var originKeys []string
	root := pkg.GetTmpPath()
	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if strings.Contains(filepath.ToSlash(path), prefix) {
			rel, _ := filepath.Rel(root, path)
			originKeys = append(originKeys, filepath.ToSlash(rel))
		}
		return nil
	})
	return originKeys, err
}

// diff 指定key 比较
func diff(key string, client *pkg.Client) bool {
	o, err := os.ReadFile(filepath.Join(pkg.GetTmpPath(), key))
	if err != nil {
		return true
	}
	sumOrigin := md5.Sum(o)
	t, err := client.Get(key)
	if err != nil {
		return true
	}
	sumTarget := md5.Sum(t)
	return sumOrigin == sumTarget
}
