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
	for _, v := range os.Args[index:] {
		err := utils.Execute(v)
		if err != nil {
			fmt.Println(err)
		}
	}
	if _, err := os.Create(".state"); err == nil {
		var state bytes.Buffer
		enc := gob.NewEncoder(&state)
		dec := gob.NewDecoder(&state)
		for _, v := range inc.State {
			if err := enc.Encode(v); err != nil {
				fmt.Println("Error encoding state, exits")
				os.Exit(1)
			}
			fmt.Println(state.Bytes())
			var object inc.Object_file

			if err := dec.Decode(&object); err != nil {
				fmt.Println("decode error:", err)
			}
			fmt.Printf("%s: {%s,%s}\n", object.Path, object.Flags, object.Timestamp)
		}
	} else {
		fmt.Println("Couldn't open/create state, exits")
		fmt.Println(err)
	}

}
