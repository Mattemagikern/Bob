package utils

import (
	"errors"
	"github.com/Mattemagikern/Bob/inc"
)

func DFS() error {
	for _, r := range inc.Recipes {
		for _, dep := range r.Dependencies {
			ingredient, ok := inc.Recipes[dep]
			if !ok {
				return errors.New("DFS: Missing Recepie " + dep)
			}
			if err := dive(r, ingredient); err != nil {
				return err
			}
		}
	}
	return nil
}

func dive(start *inc.Recipe, latest *inc.Recipe) error {
	for _, v := range latest.Dependencies {
		ingredient, ok := inc.Recipes[v]
		if !ok {
			return errors.New("DFS: Missing Recepie " + v)
		}
		if ingredient == start {
			return errors.New("DFS: Circular Dependency")
		}
		if err := dive(start, ingredient); err != nil {
			return err
		}
	}
	return nil
}
