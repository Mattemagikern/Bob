package utils

import (
	"fmt"
	"inc"
	"os"
	"os/exec"
	"regexp"
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
			return err
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
	obj_chan := make(chan *inc.Object_file)
	names := make(chan string)
	for k, v := range inc.File_tree {
		if inc.Sf.Src.FindString(v.Path) == "" {
			continue
		}
		file_name := k[:len(k)-len(inc.Build_cmd.Exstensions[2])]
		out_path := inc.Build_cmd.Exstensions[0] + file_name + inc.Build_cmd.Exstensions[1]
		inc.Variables["Objects"].Expression += " " + inc.Build_cmd.Exstensions[0] + file_name + inc.Build_cmd.Exstensions[1]
		obj := inc.State[file_name+inc.Build_cmd.Exstensions[1]]
		go build(v, out_path, obj, file_name, names, obj_chan, errors)
	}
	for range inc.File_tree {
		obj_file := <-obj_chan
		inc.State[<-names] = obj_file
		err = <-errors
		if err != nil {
			return
		}
	}
	return
}

var name *regexp.Regexp = regexp.MustCompile(`($@)`)
var path *regexp.Regexp = regexp.MustCompile(`($<)`)

func build(v *inc.File, out_path string, obj *inc.Object_file, file_name string, names chan string, obj_chan chan *inc.Object_file, errors chan error) {
	if obj == nil || obj.Timestamp.Sub(v.Timestamp) < 0 {
		obj_chan <- &inc.Object_file{out_path, inc.Variables["CFLAGS"].Expression, time.Now()}
		names <- file_name
		for _, j := range inc.Build_cmd.Commands {
			j = strings.Replace(j, "$@", out_path, -1)
			j = strings.Replace(j, "$<", v.Path, -1)
			if err := shell(j); err != nil {
				fmt.Println(err)
				errors <- err
				return
			}
		}
		errors <- nil
		return
	}
	names <- ""
	obj_chan <- nil
	errors <- nil
}
