package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"inc"
	"os"
	"path/filepath"
	"strings"
	"utils"
)

func main() {
	var index int = 1
	if err := utils.Parse_builder(); err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	filepath.Walk("../MasterThesis/code", utils.Walk())
	if err := utils.DFS(); err != nil {

	}
	if strings.Compare(os.Args[1], "clean") == 0 {
		fmt.Println("clean!")
		index = 2
		if fi, err := os.Create(".state"); err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			fi.Close()
		}
	}
	if err := utils.Parse_state(); err != nil {
		fmt.Println(err)
	}

	for _, v := range os.Args[index:] {
		err := utils.Execute(v)
		if err != nil {
			fmt.Println(err)
		}
	}
	if f, err := os.Create(".state"); err == nil {
		var buffer bytes.Buffer
		var enc *gob.Encoder = gob.NewEncoder(&buffer)
		for _, v := range inc.State {
			if err := enc.Encode(v); err != nil {
				fmt.Println("Error encoding state, exits")
				os.Exit(1)
			}
			if _, err := f.Write(buffer.Bytes()); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(v)
			buffer.Reset()
		}
	} else {
		fmt.Println("Couldn't create .state, exits")
		fmt.Println(err)
	}

}
