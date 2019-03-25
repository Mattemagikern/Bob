package utils

import (
	"errors"
	"github.com/Mattemagikern/Bob/inc"
)

func DFS() error {
	for k := range inc.Recipes {
		if _, ok := inc.Recipes[k]; !ok {
			return errors.New("DFS: Missing Recepie " + k)
		}
		visited := make(map[string]bool)
		for _, dep := range inc.Recipes[k].Dependencies {
			if _, ok := inc.Recipes[dep]; !ok {
				return errors.New("DFS: Missing Recepie " + dep)
			}
			visited[dep] = true
			if dive(visited, k, dep) {
				return errors.New("DFS: Circular Dependency in builder")
			}
		}
	}
	return nil
}

func dive(visited map[string]bool, start string, latest string) bool {
	visited[latest] = true
	if latest == "build" {
		return false
	}
	for _, v := range inc.Recipes[latest].Dependencies {
		if visited[v] == true {
			return true
		} else {
			return dive(visited, start, v)
		}
	}
	return false
}
