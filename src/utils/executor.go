package utils

import (
	"bytes"
	"fmt"
	"inc"
	"os/exec"
	"strings"
)

func Execute(recepie string) (err error) {
	if inc.Recepies[recepie] == nil {
		panic("Non valid recepie or ingredient")
	}
	var stdout, stderr bytes.Buffer

	for _, v := range inc.Recepies[recepie].Dependencies {
		err = Execute(v)
	}

	for indx, str := range inc.Recepies[recepie].Commands {
		inc.Recepies[recepie].Commands[indx] = strings.Trim(str, " \t")
	}
	for _, str := range inc.Recepies[recepie].Commands {
		var name, expression string
		var tmp []string
		switch {
		case strings.Contains(str, "+="):
			tmp = strings.Split(str, "+=")
			name = strings.Trim(tmp[0][1:], " ")
			expression = strings.Trim(tmp[1], " ")
			inc.Variables[name].Expression += " " + expression

		case strings.Contains(str, "-="):
			tmp = strings.Split(str, "-=")
			name = strings.Trim(tmp[0][1:], " ")
			expression = strings.Trim(tmp[1], " ")
			inc.Variables[name].Expression = strings.Trim(inc.Variables[name].Expression, expression)

		case strings.Contains(str, "="):
			tmp = strings.Split(str, "=")
			name = strings.Trim(tmp[0][1:], " ")
			expression = strings.Trim(tmp[1], " ")
			inc.Variables[name] = &inc.Variable{name, expression}
		default:
			if strings.Contains(str, "$") {
				for _, v := range strings.Fields(str) {
					if strings.Contains(v, "$") {
						str = strings.Replace(str, v, inc.Variables[v[1:]].Expression, -1)
					}
				}
			}
			tmp = strings.SplitN(str, " ", 2)
			cmd := exec.Command(strings.Trim(tmp[0], " "), strings.Trim(tmp[1], " "))
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			if err = cmd.Run(); err != nil {
				fmt.Println(stderr.String())
				panic("command failed!")
			}
		}
		fmt.Printf("%s", stdout.String())
		stdout.Reset()
		stderr.Reset()
	}
	return
}
