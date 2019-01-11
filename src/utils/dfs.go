package utils

import (
	"errors"
	"fmt"
	"inc"
	"strings"
)

func DFS() (err error) {
	for k, _ := range inc.Recepies {
		visited := make(map[string]bool)
		for _, dep := range inc.Recepies[k].Dependencies {
			visited[dep] = true
			if dive(visited, k, dep) {
				err = errors.New("Circular Dependencie in builder")
				return
			}
		}
	}
	return
}

func dive(visited map[string]bool, start string, latest string) bool {
	visited[latest] = true
	if strings.Compare(latest, "build") == 0 {
		fmt.Println(latest)
		return false
	}
	for _, v := range inc.Recepies[latest].Dependencies {
		if visited[v] == true {
			panic("Circular dependencie!!")
		} else {
			return dive(visited, start, v)
		}
	}
	return false
}
