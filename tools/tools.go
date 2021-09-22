package tools

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"regexp"
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

	for _, fi := range fileInfos {
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

var (
	regexpBrackets = regexp.MustCompile(`<[\S\s]+?>`)
	regexpEnter    = regexp.MustCompile(`\s+`)
)

// IgnoreHtmlTag 去掉 html tag
func IgnoreHtmlTag(src string) string {
	// 去除所有尖括号内的HTML代码
	src = regexpBrackets.ReplaceAllString(src, "")
	// 去除换行符
	return regexpEnter.ReplaceAllString(src, "")
}
