package pkg

import (
	"os"
	"path"
	"path/filepath"

	"consul-webhook/config"
)

// GetWD 获取当前执行进程所在的目录
func GetWD() string {
	ex, _ := os.Executable()
	dir := filepath.Dir(ex)
	return dir
}

// GetTmpPath 获取repo的缓存目录
func GetTmpPath() string {
	wd := GetWD()
	tmpPath := path.Join(wd, config.GetConfig().Repo.Path)
	return tmpPath
}

// ListToMap 把list转换层map
func ListToMap(keys []string) map[string]int {
	keyMap := make(map[string]int)
	for i, key := range keys {
		keyMap[key] = i
	}
	return keyMap
}
