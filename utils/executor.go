package utils

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Mattemagikern/Bob/parser"
	"log"
	"os"
	"os/exec"
)

func Execute(recipe string) error {
	if parser.Recipes[recipe] == nil {
		log.Println("Invalid recipe: " + recipe)
		return errors.New("Execute: Invalid Recipe " + recipe)
	}

	for _, v := range parser.Recipes[recipe].Ingredients {
		if err := Execute(v); err != nil {
			log.Println(err)
			return err
		}
	}
	parser.Store["@"] = recipe
	for _, str := range parser.Recipes[recipe].Cmds {
		if !parser.Update_variables(str) {
			if err := shell(str); err != nil {
				return err
			}
		}
	}
	return nil
}

func shell(str string) error {
	s := parser.Substitute(str)
	var b bytes.Buffer
	cmd := exec.Command("bash", "-c", s)
	cmd.Stderr = os.Stderr
	cmd.Stdout = &b
	if err := cmd.Run(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return err
		}
		return err
	}
	if b.Len() > 0 {
		os.Stdout.Write(b.Bytes())
		return nil
	}
	fmt.Println(str)
	return nil
}
