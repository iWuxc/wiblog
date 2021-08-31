package tools

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"path"
)

//EncryptPassword encrypt password
func EncryptPassword(name, passwd string) string {
	salt := "%$@w*<>"
	h := sha256.New()
	io.WriteString(h, salt)
	io.WriteString(h, name)
	io.WriteString(h, passwd)
	return fmt.Sprintf("%x", h.Sum(nil))
}

//ReadDirFiles 读取目录
func ReadDirFiles(dir string, filter func(name string) bool) (files []string) {
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}

	for _, fi := range fileInfos  {
		if filter(fi.Name()) {
			continue
		}
		if fi.IsDir() {
			files = append(files, ReadDirFiles(path.Join(dir, fi.Name()), filter)...)
			continue
		}
		files = append(files, path.Join(dir, fi.Name()))
	}

	return files
}
