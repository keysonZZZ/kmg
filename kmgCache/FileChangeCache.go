package kmgCache

import (
	"os"
	"path/filepath"

	"github.com/bronze1man/kmg/encoding/kmgGob"
	"github.com/bronze1man/kmg/kmgConfig"
	"github.com/bronze1man/kmg/kmgCrypto"
	"github.com/bronze1man/kmg/kmgFile"
)

//这个测试要用到.
func getFileChangeCachePath(key string) string {
	// 此处key长度不可控,所以用md5.
	return filepath.Join(kmgConfig.DefaultEnv().TmpPath, "FileChangeCache", kmgCrypto.Md5HexFromString(key))
}

// key 表示一组缓存，pathList 即这组缓存相关的文件
// 比较容易用错的情况：
// key := 1
// pathList := []string{"/a.txt","/b.txt"}
// a.txt 或者 b.txt 中任意有文件发生变化，这个 key 对应的缓存都会被更新
// 如果要对单个文件进行缓存控制，那么应该一个文件一个 key
func MustMd5FileChangeCache(key string, pathList []string, f func()) {
	// 此处需要考虑,
	//   用户新添加了一个文件
	//   用户删除了一个文件
	//   用户编辑了一个文件
	//   用户在目录里面添加了一个文件
	//   用在在目录里面删除了一个文件
	//读取文件修改时间缓存信息
	toChange := false
	cacheInfo := map[string]string{}
	cacheFilePath := getFileChangeCachePath(key)
	err := kmgGob.ReadFile(cacheFilePath, &cacheInfo)
	if err != nil {
		//忽略缓存读取的任何错误
		cacheInfo = map[string]string{}
	}
	hasReadFileMap := map[string]bool{}
	for _, path := range pathList {
		statList, err := kmgFile.GetAllFileAndDirectoryStat(path)
		if err != nil {
			if os.IsNotExist(err) {
				toChange = true
				//fmt.Printf("[MustFileChangeCache] path:[%s] not exist\n", path)
				break
			}
			panic(err)
		}

		for _, stat := range statList {
			if stat.Fi.IsDir() {
				continue
			}
			hasReadFileMap[stat.FullPath] = true
			if stat.Fi.Mode()&os.ModeSymlink == os.ModeSymlink {
				if "symlink_"+kmgFile.MustReadSymbolLink(stat.FullPath)!=cacheInfo[stat.FullPath]{
					toChange = true
					break
				}
				continue
			}
			if kmgCrypto.MustMd5File(stat.FullPath) != cacheInfo[stat.FullPath] {
				toChange = true
				//fmt.Printf("[MustMd5FileChangeCache] path:[%s] mod md5 not match save[%s] file[%s]\n", stat.FullPath,
				//	cacheMd5, kmgCrypto.MustMd5File(stat.FullPath))
				break
			}
		}
		if toChange {
			break
		}
	}
	//删除一个已经存储在缓存列表里面的文件,是一个修改.
	for fullPath := range cacheInfo {
		if !hasReadFileMap[fullPath] {
			toChange = true
			break
		}
	}
	if !toChange {
		return
	}
	f()
	cacheInfo = map[string]string{}
	for _, path := range pathList {
		statList, err := kmgFile.GetAllFileAndDirectoryStat(path)
		if err != nil {
			panic(err)
		}
		for _, stat := range statList {
			if stat.Fi.IsDir() {
				continue
			}
			if stat.Fi.Mode()&os.ModeSymlink == os.ModeSymlink {
				linkToPath:=kmgFile.MustReadSymbolLink(stat.FullPath)
				cacheInfo[stat.FullPath] = "symlink_"+linkToPath
				continue
			}
			cacheInfo[stat.FullPath] = kmgCrypto.MustMd5File(stat.FullPath)
		}
	}
	kmgFile.MustMkdirForFile(cacheFilePath)
	kmgGob.MustWriteFile(cacheFilePath, cacheInfo)
	//保存文件缓存信息
	return
}
