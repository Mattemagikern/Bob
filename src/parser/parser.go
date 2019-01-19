package parser

import (
	"encoding/gob"
	"errors"
	"fmt"
	"inc"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var Variables *regexp.Regexp = regexp.MustCompile(`(?m)^\$?(\S*)\s*(=|\+=|\-=)(?:\s*(.*))`)
var recepies *regexp.Regexp = regexp.MustCompile(`(?m)^(\S*): ?(.*)\n((?:\t.*\n?)*)`)
var builder *regexp.Regexp = regexp.MustCompile(`(?m)(.*)%(.*):\s?%(.*)$\s((?:\t.*\n?)*)`)
var test *regexp.Regexp = regexp.MustCompile(`(?:\$\(.*\)|[^\s]\S*)`)
var wow *regexp.Regexp = regexp.MustCompile(`\$\((.*)\)`)

func Parse_builder(file string) (err error) {
	for _, v := range recepies.FindAllString(file, -1) {
		var tmp []string
		if tmp = builder.FindStringSubmatch(v); tmp != nil {
			cmds := strings.Split(tmp[4], "\n")
			cmds = cmds[:len(cmds)-1]
			inc.Build_cmd = &inc.Build{"build", tmp[1:4], cmds}
		} else {
			tmp = recepies.FindStringSubmatch(v)
			name := tmp[1]
			ingredients := strings.Fields(tmp[2])
			cmds := strings.Split(tmp[3], "\n")
			cmds = cmds[:len(cmds)-1]
			for i, v := range cmds {
				cmds[i] = strings.Trim(v, "\t")
			}
			inc.Recepies[name] = &inc.Recepie{name, ingredients, cmds}
		}
	}

	for _, v := range Variables.FindAllString(file, -1) {
		tmp := Variables.FindStringSubmatch(v)
		ama, _ := Substitute(tmp[3])
		Update_vars(tmp[1], tmp[2], ama)
	}
	if inc.Variables["src"] == nil {
		fmt.Println("Missing regex pattern for src files, exits")
		os.Exit(1)
	}
	inc.Sf.Src = regexp.MustCompile(inc.Variables["src"].Expression)

	if inc.Variables["inc"] != nil {
		inc.Sf.Inc = regexp.MustCompile(inc.Variables["inc"].Expression)
	} else {
		inc.Sf.Inc = regexp.MustCompile(`$^`)
	}
	if inc.Variables["inc_pattern"] != nil {
		inc.Sf.Inc_pattern = regexp.MustCompile(inc.Variables["inc_pattern"].Expression)
	} else {
		inc.Sf.Inc_pattern = regexp.MustCompile(`$^`)
		err = errors.New("parser: Parse_builder: Missing inc or inc_pattern, This will not garuantee the correctness of your object files.")
	}

	return
}

func Substitute(v string) (str string, err error) {
	if !strings.Contains(v, "$") {
		return v, nil
	}
	var tmps []byte
	for _, value := range test.FindAllString(v, -1) {
		elm, ok := inc.Variables[value[1:]]
		switch {
		case strings.Contains(value, "$("):
			cmd := wow.FindStringSubmatch(v)
			cmd = strings.Fields(cmd[1])
			tmps, err = exec.Command(cmd[0], cmd[1:]...).Output()
			if err != nil {
				err = errors.New("parser: Substitute: " + err.Error())
			}
			str += string(tmps)

		case value[0] == '$' && ok:
			str += elm.Expression

		default:
			str += value
		}
		str += " "
	}
	str = strings.Trim(str, " ")
	return
}

func Update_vars(name string, delimiter string, str string) (err error) {
	switch {
	case delimiter == "+=":
		inc.Variables[name].Expression += " " + str
	case delimiter == "-=":
		inc.Variables[name].Expression = strings.Trim(inc.Variables[name].Expression, str)
	case delimiter == "=":
		inc.Variables[name] = &inc.Variable{name, str}
	default:
		err = errors.New("parser: Uptade_vars: Update vars wrong input")
	}
	switch {
	case name == "inc":
		inc.Sf.Inc = regexp.MustCompile(inc.Variables[name].Expression)
	case name == "src":
		inc.Sf.Src = regexp.MustCompile(inc.Variables[name].Expression)
	case name == "inc_pattern":
		inc.Sf.Inc_pattern = regexp.MustCompile(inc.Variables[name].Expression)
	}

	return
}

func Parse_state() (err error) {
	var f *os.File
	if f, err = os.OpenFile(".state", os.O_RDONLY, 0644); err != nil {
		err = errors.New("Parse_state" + err.Error())
		return
	}
	dec := gob.NewDecoder(f)
	for err != io.EOF {
		var obj inc.Object_file
		err = dec.Decode(&obj)
		inc.State[filepath.Base(obj.Path)] = &obj
	}
	return nil
}
