package utils

import (
	"errors"
	"inc"
	"strings"
)

func DFS() (err error) {
	for k := range inc.Recepies {
		visited := make(map[string]bool)
		for _, dep := range inc.Recepies[k].Dependencies {
			visited[dep] = true
			if dive(visited, k, dep) {
				err = errors.New("DFS: Circular Dependencie in builder")
				return
			}
		}
	}
	return
}

func dive(visited map[string]bool, start string, latest string) bool {
	visited[latest] = true
	if strings.Compare(latest, "build") == 0 {
		return false
	}
	for _, v := range inc.Recepies[latest].Dependencies {
		if visited[v] == true {
			return true
		} else {
			return dive(visited, start, v)
		}
	}
	return false
}
