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

	filepath.Walk("..", utils.Walk())
	/*
		fmt.Println("recepies")
		for k, v := range inc.Recepies {
			fmt.Println(k, v)
		}
		fmt.Println("variables")
		for k, v := range inc.Variables {
			fmt.Println(k, v)
		}
		fmt.Println("file-tree")
		for k, v := range inc.File_tree {
			fmt.Println(k, v, v.Timestamp)
		}
	*/
	if err := utils.DFS(); err != nil {

	}
	for _, v := range os.Args[1:] {
		err := utils.Execute(v)
		if err != nil {
			fmt.Println("error in command exec")
			fmt.Println(err)
		}
	}
}
