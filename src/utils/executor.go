package utils

import (
	"errors"
	"fmt"
	"inc"
	"os"
	"os/exec"
	"parser"
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
				err = errors.New("utils: Execute: " + err.Error())
				return
			}
		} else {
			if err = Execute(v); err != nil {
				err = errors.New("utils: Execute: " + err.Error())
				return
			}
		}
	}

	for _, str := range inc.Recepies[recepie].Commands {
		if err = shell(str); err != nil {
			err = errors.New("utils: Execute: " + err.Error())
			return err
		}
	}
	return
}
func shell(s string) (err error) {
	var str string
	if tmp := parser.Variables.FindStringSubmatch(s); tmp != nil {
		str, err = parser.Substitute(tmp[3])
		parser.Update_vars(tmp[1], tmp[2], str)
		return nil
	}
	str, err = parser.Substitute(s)
	a := strings.Fields(str)
	cmd := exec.Command(a[0], a[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			fmt.Println(exitError)
			return errors.New("shell: " + err.Error())
		}
	}
	return nil
}

func Build() (err error) {
	errs := make(chan error, 10)
	var visited int
	for k, v := range inc.File_tree {
		if inc.Sf.Src.FindString(v.Path) == "" {
			continue
		}
		visited++
		file_name := k[:len(k)-len(inc.Build_cmd.Exstensions[2])]
		object_name := file_name + inc.Build_cmd.Exstensions[1]
		out_path := inc.Build_cmd.Exstensions[0] + object_name
		if !strings.Contains(inc.Variables["Objects"].Expression, out_path) {
			inc.Variables["Objects"].Expression += " " + out_path
		}
		obj := inc.State[object_name]
		go build(v, &out_path, obj, &inc.Variables["CFLAGS"].Expression, errs)
		inc.State[object_name] = &inc.Object_file{out_path, inc.Variables["CFLAGS"].Expression, time.Now()}
	}
	for visited != 0 {
		err = <-errs
		if err != nil {
			err = errors.New("Build: " + err.Error())
			return
		}
		visited--
	}
	return
}

func build(v *inc.File, out_path *string, obj *inc.Object_file, flags *string, errors chan error) {
	if obj == nil || obj.Timestamp.Sub(v.Timestamp) < 0 || obj.Flags != *flags || check_dependencies(v, obj) {
		for _, j := range inc.Build_cmd.Commands {
			j = strings.Replace(j, "$@", *out_path, -1)
			j = strings.Replace(j, "$<", v.Path, -1)
			if err := shell(j); err != nil {
				errors <- err
			}
		}
	}
	errors <- nil
}

func check_dependencies(v *inc.File, obj *inc.Object_file) bool {
	for _, dep := range v.Inc {
		include, ok := inc.File_tree[dep]
		if ok && obj.Timestamp.Sub(include.Timestamp) < 0 {
			return true
		}
		if ok && check_dependencies(include, obj) {
			return true
		}
	}
	return false
}
