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
var suffixes *regexp.Regexp = regexp.MustCompile(`(?m)^search\s+\{(.|\n)*^\}`)
var substitute *regexp.Regexp = regexp.MustCompile(`(?s)(.*)\$([^\s]*)`)
var builder *regexp.Regexp = regexp.MustCompile(`(?m)(.*)%(.*):\s?%(.*)$\s((?:\t.*\n?)*)`)
var cmds *regexp.Regexp = regexp.MustCompile(`(?:^\s?([^\s]*)\s?)=\s?\$\((.*)\)`)
var test *regexp.Regexp = regexp.MustCompile(`(?:\$\(.*\)|[^\s]\S*)`)
var wow *regexp.Regexp = regexp.MustCompile(`\$\((.*)\)`)

func Parse_builder(file string) (err error) {

	s := strings.Split(suffixes.FindString(file), "\n")
	for _, element := range s {
		element := strings.Fields(element)
		switch element[0] {
		case "src":
			inc.Sf.Src = regexp.MustCompile(element[1])
		case "inc":
			inc.Sf.Inc = regexp.MustCompile(element[1])
		case "include_pattern":
			inc.Sf.Inc_pattern = regexp.MustCompile(element[1])
		}
	}

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
			inc.Recepies[name] = &inc.Recepie{name, ingredients, cmds}
		}
	}

	for _, v := range Variables.FindAllString(file, -1) {
		tmp := Variables.FindStringSubmatch(v)
		ama, _ := Substitute(tmp[3])
		Update_vars(tmp[1], tmp[2], ama)
	}

	if inc.Sf.Inc == nil || inc.Sf.Src == nil || inc.Sf.Inc_pattern == nil {
		err = errors.New("Suffixes: objecs, srcs or, inc were not found")
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
			str += string(tmps)

		case value[0] == '$' && ok:
			str += elm.Expression

		default:
			str += value
		}
		str += " "
	}
	str = strings.Trim(str, " ")

	fmt.Println(str)
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
		err = errors.New("Update vars wrong input")
	}
	return
}

func Parse_state() (err error) {
	var f *os.File
	if f, err = os.OpenFile(".state", os.O_RDONLY, 0644); err != nil {
		return err
	}
	dec := gob.NewDecoder(f)
	for err != io.EOF {
		var obj inc.Object_file
		err = dec.Decode(&obj)
		inc.State[filepath.Base(obj.Path)] = &obj
	}
	err = nil
	return
}
