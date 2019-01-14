package utils

import (
	"fmt"
	"inc"
	"os"
	"os/exec"
	"strings"
	"time"
)

func Execute(recepie string) (err error) {
	if inc.Recepies[recepie] == nil {
		panic("Non valid recepie or ingredient")
	}

	for _, v := range inc.Recepies[recepie].Dependencies {
		if v == "build" {
			if err = Build(); err != nil {
				return
			}
		} else {
			if err = Execute(v); err != nil {
				return
			}
		}
	}

	for indx, str := range inc.Recepies[recepie].Commands {
		inc.Recepies[recepie].Commands[indx] = strings.Trim(str, " \t")
	}
	for _, str := range inc.Recepies[recepie].Commands {
		if err = shell(str); err != nil {
		}

	}
	return
}

func shell(s string) (err error) {
	var v string
	var boo bool
	v, boo, _ = Substitute(s)
	if !boo {
		fmt.Println(v)
		a := strings.Fields(v)
		cmd := exec.Command(a[0], a[1:]...)
		cmd.Stderr = os.Stderr
		if err = cmd.Run(); err != nil {
			return err.(*exec.ExitError)
		}
	}
	return nil
}

func Build() (err error) {
	errors := make(chan error)
	str := make(chan string)
	for k, v := range inc.File_tree {
		if inc.Sf.Src.FindString(v.Path) == "" {
			continue
		}
		go build(k, v, str, errors)
	}
	for range inc.File_tree {
		obj := <-str
		inc.Variables["Objects"].Expression += " " + obj
		if err = <-errors; err != nil {
			return
		}
	}
	return
}

func build(k string, v *inc.File, objects chan string, errors chan error) {
	file_name := k[:len(k)-len(inc.Build_cmd.Exstensions[2])]
	obj, ok := inc.State[file_name+inc.Build_cmd.Exstensions[1]]
	inc.Variables["<"] = &inc.Variable{"@", v.Path}
	out_path := inc.Build_cmd.Exstensions[0] + file_name + inc.Build_cmd.Exstensions[1]
	inc.Variables["@"] = &inc.Variable{"<", out_path}
	objects <- out_path
	if ok != true {
		inc.State[file_name+inc.Build_cmd.Exstensions[1]] = &inc.Object_file{out_path, inc.Variables["CFLAGS"].Expression, time.Now()}
		for _, j := range inc.Build_cmd.Commands {
			errors <- shell(j)
		}
	} else if obj.Timestamp.Sub(v.Timestamp) < 0 {
		for _, j := range inc.Build_cmd.Commands {
			errors <- shell(j)
		}
	}
	errors <- nil
}
