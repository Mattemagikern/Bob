package utils

import (
	"fmt"
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
		switch {
		case inc.Sf.Objects.MatchString(path):
			fmt.Printf("*.o visited: %s\n", path)
			return nil
		case inc.Sf.Inc.MatchString(path) || inc.Sf.Src.MatchString(path):
			data, err := ioutil.ReadFile(path)
			if err != nil {
				panic(err)
			}
			tmp := inc.Sf.Inc_pattern.FindAllStringSubmatch(string(data), -1)
			for _, e := range tmp {
				includes = append(includes, e[1])
			}
		default:
			return nil
		}
		inc.File_tree[name] = inc.File{path, includes, time}
		return nil
	}
}
