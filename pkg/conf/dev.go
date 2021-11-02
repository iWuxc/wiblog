// +build !prod

package conf

import (
	"os"
	"path"
	"path/filepath"
)

var workDir = func() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for wd != "" {
		name := filepath.Join(wd, "conf")
		_, err := os.Stat(name)
		if err != nil {
			dir, _ := path.Split(wd)
			wd = path.Clean(dir)
			continue
		}
		return wd
	}
	return ""
}
