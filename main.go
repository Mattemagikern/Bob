package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"inc"
	"io/ioutil"
	"os"
	"parser"
	"path/filepath"
	"strings"
	"utils"
)

func main() {
	var index int = 1

	if dat, err := ioutil.ReadFile("./BUILDER"); err == nil {
		if err := parser.Parse_builder(string(dat)); err != nil {
			fmt.Println("main: ", err.Error())
			os.Exit(1)
		}
	} else {
		fmt.Println("main: ", err.Error())
		os.Exit(1)
	}

	filepath.Walk("../MasterThesis/code", Walk())
	if err := utils.DFS(); err != nil {
		fmt.Println("main: " + err.Error())
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		os.Exit(1)
	}

	if strings.Compare(os.Args[1], "clean") == 0 {
		if _, err := os.Create(".state"); err != nil {
			fmt.Println("main: " + err.Error())
			os.Exit(1)
		}
		index = 2
	}

	if err := parser.Parse_state(); err != nil {
		fmt.Println(err)
	}

	for _, v := range os.Args[index:] {
		err := utils.Execute(v)
		if err != nil {
			fmt.Println("main: " + err.Error())
		}
	}
	if f, err := os.Create(".state"); err == nil {
		var buffer bytes.Buffer
		var enc *gob.Encoder = gob.NewEncoder(&buffer)
		for _, v := range inc.State {
			if err := enc.Encode(v); err != nil {
				fmt.Println("Main: Error encoding state, exits")
				os.Exit(1)
			}
			if _, err := f.Write(buffer.Bytes()); err != nil {
				fmt.Println("main: " + err.Error())
				os.Exit(1)
			}
			buffer.Reset()
		}
	} else {
		fmt.Println("Main: " + err.Error())
		os.Exit(1)
	}

}

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
			inc.File_tree[name] = &inc.File{Path: path, Inc: includes, Timestamp: time}
		}
		return nil
	}
}
