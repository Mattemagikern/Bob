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
var recipes *regexp.Regexp = regexp.MustCompile(`(?m)^(\S*): ?(.*)\n((?:\t.*\n?)*)`)
var builder *regexp.Regexp = regexp.MustCompile(`(?m)(.*)%(.*):\s?%(.*)$\s((?:\t.*\n?)*)`)
var test *regexp.Regexp = regexp.MustCompile(`(?:\$\(.*\)|[^\s]\S*)`)
var wow *regexp.Regexp = regexp.MustCompile(`\$\((.*)\)`)

func Parse_builder(file string) (err error) {
	for _, v := range recipes.FindAllString(file, -1) {
		var tmp []string
		if tmp = builder.FindStringSubmatch(v); tmp != nil {
			cmds := strings.Split(tmp[4], "\n")
			cmds = cmds[:len(cmds)-1]
			inc.Build_cmd = &inc.Build{"build", tmp[1:4], cmds}
		} else {
			tmp = recipes.FindStringSubmatch(v)
			name := tmp[1]
			ingredients := strings.Fields(tmp[2])
			cmds := strings.Split(tmp[3], "\n")
			cmds = cmds[:len(cmds)-1]
			for i, v := range cmds {
				cmds[i] = strings.Trim(v, "\t")
			}
			inc.Recipes[name] = &inc.Recipe{name, ingredients, cmds}
			if _, ok := inc.Recipes["default"]; !ok {
				inc.Recipes["default"] = inc.Recipes[name]
			}
		}
	}

	for _, v := range Variables.FindAllString(file, -1) {
		tmp := Variables.FindStringSubmatch(v)
		ama, err := Substitute(tmp[3])
		if err != nil {
			return err
		}
		Update_vars(tmp[1], tmp[2], ama)
	}
	if inc.Variables["src"] == nil {
		fmt.Println("Missing regex pattern for src files, may not work as you intend")
	}

	if inc.Variables["src"] != nil {
		inc.Sf.Src = regexp.MustCompile(inc.Variables["src"].Expression)
	} else {
		inc.Sf.Src = regexp.MustCompile(`$^`)
	}

	if inc.Variables["inc"] != nil {
		inc.Sf.Inc = regexp.MustCompile(inc.Variables["inc"].Expression)
	} else {
		inc.Sf.Inc = regexp.MustCompile(`$^`)
	}
	if inc.Variables["inc_pattern"] != nil {
		inc.Sf.Inc_pattern = regexp.MustCompile(inc.Variables["inc_pattern"].Expression)
	} else {
		inc.Sf.Inc_pattern = regexp.MustCompile(`$^`)
	}

	return
}

func Substitute(v string) (str string, err error) {
	if !strings.Contains(v, "$") {
		return v, nil
	}
	for _, value := range test.FindAllString(v, -1) {
		indx := strings.Index(value, "$")
		switch {
		case indx == -1:
			str += value
		case indx == 0 && value[indx+1] != '(':
			elm, ok := inc.Variables[value[1:]]
			if !ok {
				return "", errors.New("Unknown variable")
			}
			str += elm.Expression
		default:
			str += value[:indx]
			if value[indx+1] == '(' {
				i := strings.Index(value, ")")
				if i == -1 {
					return "", errors.New("Malformed variable")
				}
				cmd := value[indx+2 : i]
				out, err := exec.Command("bash", "-c", cmd).Output()
				if err != nil {
					return "", errors.New("Failed substitute command")
				}
				str += string(out[:len(out)-1])
				str += value[i+1:]
			}
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
