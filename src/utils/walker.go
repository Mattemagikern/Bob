package utils

import (
	"inc"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Walk() filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		time := info.ModTime()
		name := info.Name()
		var includes []string
		if inc.Sf.Src.MatchString(path) || inc.Sf.Inc.MatchString(path) {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				panic(err)
			}
			tmp := inc.Sf.Inc_pattern.FindAllStringSubmatch(string(data), -1)
			for _, e := range tmp {
				includes = append(includes, e[1])
			}
			inc.File_tree[name] = &inc.File{path, includes, time}
		}
		return nil
	}
}
