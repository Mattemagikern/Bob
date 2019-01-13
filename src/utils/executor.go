package utils

import (
	"bytes"
	"errors"
	"fmt"
	"inc"
	"os/exec"
	"strings"
	"time"
)

func Execute(recepie string) (err error) {
	if inc.Recepies[recepie] == nil {
		panic("Non valid recepie or ingredient")
	}

	for _, v := range inc.Recepies[recepie].Dependencies {
		if strings.Compare(v, "build") == 0 {
			err = Build()
		} else {
			err = Execute(v)
		}
		if err != nil {
			return
		}
	}

	for indx, str := range inc.Recepies[recepie].Commands {
		inc.Recepies[recepie].Commands[indx] = strings.Trim(str, " \t")
	}
	for _, str := range inc.Recepies[recepie].Commands {
		shell(str)
	}
	return
}

func shell(s string) (err error) {
	var v string
	var boo bool
	var stdout, stderr bytes.Buffer
	v, boo, err = Substitute(s)
	if !boo {
		fmt.Println(v)
		a := strings.Fields(v)
		cmd := exec.Command(a[0], a[1:]...)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		if err = cmd.Run(); err != nil {
			return errors.New(stderr.String())
		}
		stdout.Reset()
		stderr.Reset()
	}
	return
}

func Build() (err error) {
	for k, v := range inc.File_tree {
		if inc.Sf.Src.FindString(v.Path) == "" {
			continue
		}
		file_name := k[:len(k)-len(inc.Build_cmd.Exstensions[2])]
		obj, ok := inc.State[file_name+inc.Build_cmd.Exstensions[1]]
		inc.Variables["<"] = &inc.Variable{"@", v.Path}
		out_path := inc.Build_cmd.Exstensions[0] + file_name + inc.Build_cmd.Exstensions[1]
		inc.Variables["@"] = &inc.Variable{"<", out_path}
		if ok != true {
			for _, j := range inc.Build_cmd.Commands {
				/*TODO: Goroutine shell, if error exit and cancel all other builds?*/
				if err = shell(j); err != nil {
					return err
				}
				inc.State[file_name+inc.Build_cmd.Exstensions[1]] = &inc.Object_file{out_path, inc.Variables["CFLAGS"].Expression, time.Now()}
			}
		} else if v.Timestamp.Sub(obj.Timestamp) != 0 {
			fmt.Println(v.Timestamp.Sub(obj.Timestamp))
			inc.State[file_name+inc.Build_cmd.Exstensions[1]] = &inc.Object_file{out_path, inc.Variables["CFLAGS"].Expression, time.Now()}
		}
	}
	return
}
