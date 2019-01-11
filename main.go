package main

import (
	"fmt"
	"os"
	"path/filepath"
	"utils"
)

func main() {
	if err := utils.Parse_builder(); err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	filepath.Walk("../MasterThesis/code", utils.Walk())
	if err := utils.DFS(); err != nil {

	}
	for _, v := range os.Args[1:] {
		err := utils.Execute(v)
		if err != nil {
			fmt.Println(err)
		}
	}
}
