package kmgCache

import (
	"testing"
	"time"

	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgTest"
	"github.com/bronze1man/kmg/kmgCmd"
)

func TestFileChangeCacheOneFile(t *testing.T) {
	//可以递归 遍历文件
	//缓存数据文件不存在,没有问题
	//指定的文件不存在,也没有问题
	callLog := make([]string, 32)
	//指定的文件不存在
	kmgFile.MustDeleteFile(getFileChangeCachePath("test_file_change_cache"))
	kmgFile.MustDeleteFile("testFile/notExist")

	MustFileChangeCache("test_file_change_cache", []string{
		"testFile/notExist",
	}, func() {
		callLog[1] = "notExist"
		kmgFile.MustWriteFileWithMkdir("testFile/notExist", []byte("1"))
	})
	kmgTest.Equal(callLog[1], "notExist")

	MustFileChangeCache("test_file_change_cache", []string{
		"testFile/notExist",
	}, func() {
		callLog[2] = "notExist"
	})
	kmgTest.Equal(callLog[2], "")

	time.Sleep(time.Second * 1)

	kmgFile.MustWriteFile("testFile/notExist", []byte("2"))
	MustFileChangeCache("test_file_change_cache", []string{
		"testFile/notExist",
	}, func() {
		callLog[3] = "notExist"
	})
	kmgTest.Equal(callLog[3], "notExist")
}

func TestFileChangeCacheOneDir(t *testing.T) {
	callLog := make([]string, 32)
	//递归可用
	kmgFile.MustDeleteFile(getFileChangeCachePath("test_file_change_cache"))
	kmgFile.MustMkdirAll("testFile/d1/d2")
	kmgFile.MustWriteFile("testFile/d1/d2/f3", []byte("1"))
	MustFileChangeCache("test_file_change_cache", []string{
		"testFile/d1",
	}, func() {
		callLog[3] = "f3"
	})
	kmgTest.Equal(callLog[3], "f3")

	MustFileChangeCache("test_file_change_cache", []string{
		"testFile/d1",
	}, func() {
		callLog[4] = "f3"
	})
	kmgTest.Equal(callLog[4], "")

	time.Sleep(time.Second * 1)
	kmgFile.MustWriteFile("testFile/d1/d2/f3", []byte("2"))
	MustFileChangeCache("test_file_change_cache", []string{
		"testFile/d1",
	}, func() {
		callLog[5] = "f3"
	})
	kmgTest.Equal(callLog[5], "f3")
}

func TestFileMd5ChangeCacheOneDir(t *testing.T) {
	callLog := make([]string, 32)
	//递归可用
	kmgFile.MustDeleteFile(getFileChangeCachePath("test_file_change_cache"))
	kmgFile.MustDelete("testFile/d1")

	kmgFile.MustMkdirAll("testFile/d1/d2")
	kmgFile.MustWriteFile("testFile/d1/d2/f3", []byte("1"))
	MustMd5FileChangeCache("test_file_change_cache", []string{
		"testFile/d1",
	}, func() {
		callLog[3] = "f3"
	})
	kmgTest.Equal(callLog[3], "f3")

	//没有碰过任何东西,缓存有效
	MustMd5FileChangeCache("test_file_change_cache", []string{
		"testFile/d1",
	}, func() {
		callLog[4] = "f3"
	})
	kmgTest.Equal(callLog[4], "")

	//修改文件内容,缓存应该无效
	kmgFile.MustWriteFile("testFile/d1/d2/f3", []byte("2"))
	MustMd5FileChangeCache("test_file_change_cache", []string{
		"testFile/d1",
	}, func() {
		callLog[5] = "f3"
	})
	kmgTest.Equal(callLog[5], "f3")

	//删除文件,缓存应该无效
	kmgFile.MustDelete("testFile/d1/d2/f3")
	MustMd5FileChangeCache("test_file_change_cache", []string{
		"testFile/d1",
	}, func() {
		callLog[6] = "f4"
	})
	kmgTest.Equal(callLog[6], "f4")

	//添加文件,缓存应该无效
	kmgFile.MustWriteFile("testFile/d1/d2/f4", []byte("3"))
	MustMd5FileChangeCache("test_file_change_cache", []string{
		"testFile/d1",
	}, func() {
		callLog[7] = "f4"
	})
	kmgTest.Equal(callLog[7], "f4")

	//读取文件,缓存有效
	kmgFile.MustReadFile("testFile/d1/d2/f4")
	MustMd5FileChangeCache("test_file_change_cache", []string{
		"testFile/d1",
	}, func() {
		callLog[8] = "f4"
	})
	kmgTest.Equal(callLog[8], "")

	//创建目录,缓存有效
	kmgFile.MustMkdir("testFile/d1/d2/f5")
	MustMd5FileChangeCache("test_file_change_cache", []string{
		"testFile/d1",
	}, func() {
		callLog[9] = "f4"
	})
	kmgTest.Equal(callLog[9], "")
}

func TestFileMd5ChangeCacheSymlink(t *testing.T){
	callLog := make([]string, 32)
	//递归可用
	kmgFile.MustDeleteFile(getFileChangeCachePath("test_file_change_cache"))
	kmgFile.MustDelete("testFile")
	kmgFile.MustWriteFileWithMkdir("testFile/d1/d2",[]byte("1"))
	kmgCmd.MustRun("ln -s d1 testFile/d3")
	MustMd5FileChangeCache("test_file_change_cache", []string{
		"testFile",
	}, func() {
		callLog[0] = "f3"
	})
	kmgTest.Equal(callLog[0], "f3")
	MustMd5FileChangeCache("test_file_change_cache", []string{
		"testFile",
	}, func() {
		callLog[1] = "f3"
	})
	kmgTest.Equal(callLog[1], "")
}