package kmgFile

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func IsDotFile(path string) bool {
	if path == "./" {
		return false
	}
	base := filepath.Base(path)
	if strings.HasPrefix(base, ".") {
		return true
	}
	return false
}

func GetFileBaseWithoutExt(p string) string {
	return filepath.Base(p[:len(p)-len(filepath.Ext(p))])
}

func WriteFile(path string, content []byte) (err error) {
	return ioutil.WriteFile(path, content, os.FileMode(0777))
}
func MustWriteFile(path string, content []byte) {
	err := ioutil.WriteFile(path, content, os.FileMode(0777))
	if err != nil {
		panic(err)
	}
}
func ReadFileAll(path string) (content []byte, err error) {
	return ioutil.ReadFile(path)
}

func Mkdir(path string) (err error) {
	return os.MkdirAll(path, os.FileMode(0777))
}

func MustMkdirAll(dirname string) {
	err := os.MkdirAll(dirname, os.FileMode(0777))
	if err != nil {
		panic(err)
	}
}

//保证一个文件的路径可以写入
func MkdirForFile(path string) (err error) {
	path = filepath.Dir(path)
	return os.MkdirAll(path, os.FileMode(0777))
}

func AppendFile(path string, content []byte) (err error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0777))
	if err != nil {
		return
	}
	defer f.Close()
	_, err = f.Write(content)
	return
}

func FileExist(path string) (exist bool, err error) {
	_, err = os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, err
}

//from http://stackoverflow.com/a/13027975/1586797
func RemoveExtFromFilePath(path string) string {
	return path[:len(path)-len(filepath.Ext(path))]
}

//just some Knowledge,you can direct call ioutil.ReadDir
func ReadDir(dirname string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dirname)
}

//copy file
// * override dst file if it exist,
// * mkdir if base dir not exist
//from http://stackoverflow.com/a/21067803/1586797
func CopyFile(src, dst string) (err error) {
    in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("[CopyFile] openSrc err[%s]", err.Error())
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(filepath.Dir(dst), os.FileMode(0777))
			if err != nil {
				return err
			}
			out, err = os.Create(dst)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("[CopyFile] createDst err[%s]", err.Error())
		}
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	//why this?
	//err = out.Sync()
	return
}
