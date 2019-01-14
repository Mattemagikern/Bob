package utils

import (
	"encoding/gob"
	"errors"
	"fmt"
	"inc"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var variables *regexp.Regexp = regexp.MustCompile(`(?m)^(?:\$)?(\S*)\s*(?: =|\+=|\-=)(?:\s*(.*))`)
var recepies *regexp.Regexp = regexp.MustCompile(`(?m)^(\S*): ?(.*)\n((?:\t.*\n?)*)`)
var suffixes *regexp.Regexp = regexp.MustCompile(`(?m)^search\s+\{(.|\n)*^\}`)
var substitute *regexp.Regexp = regexp.MustCompile(`(?s)(.*)\$([^\s]*)`)
var builder *regexp.Regexp = regexp.MustCompile(`(?m)(.*)%(.*):\s?%(.*)$\s((?:\t.*\n?)*)`)

func Parse_builder() (err error) {
	var dat []byte

	if dat, err = ioutil.ReadFile("./BUILDER"); err != nil {
		panic("Could't open BUILDER, exits")
	}

	s := strings.Split(suffixes.FindString(string(dat)), "\n")
	for _, element := range s {
		element := strings.Fields(element)
		switch element[0] {
		case "objects":
			inc.Sf.Objects = regexp.MustCompile(element[1])
		case "src":
			inc.Sf.Src = regexp.MustCompile(element[1])
		case "inc":
			inc.Sf.Inc = regexp.MustCompile(element[1])
		case "include_pattern":
			inc.Sf.Inc_pattern = regexp.MustCompile(element[1])
		}
	}

	for _, v := range recepies.FindAllString(string(dat), -1) {
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
	for _, v := range variables.FindAllString(string(dat), -1) {
		if _, _, err = Substitute(v); err != nil {
			panic("Could not substitute")
		}
	}

	if inc.Sf.Inc == nil || inc.Sf.Src == nil || inc.Sf.Objects == nil || inc.Sf.Inc_pattern == nil {
		err = errors.New("Suffixes: objecs, srcs or, inc were not found")
	}

	return
}

func Substitute(v string) (str string, bo bool, err error) {
	var tmp []string
	var name string
	var expression string
	str = v
	bo = false
	tmp = variables.FindStringSubmatch(v)
	if tmp != nil {
		bo = true
		name = tmp[1]
		expression = tmp[2]
		if strings.Contains(str, "+=") {
			inc.Variables[name].Expression += " " + expression
		} else if strings.Contains(str, "-=") {
			inc.Variables[name].Expression = strings.Trim(inc.Variables[name].Expression, expression)
		} else {
			inc.Variables[name] = &inc.Variable{name, expression}
		}
		if strings.Contains(inc.Variables[name].Expression, "$") {
			for _, v := range strings.Fields(inc.Variables[name].Expression)[1:] {
				if strings.Contains(v, "$") {
					inc.Variables[name].Expression = strings.Replace(inc.Variables[name].Expression, v, substitute.ReplaceAllString(v, "$1"+inc.Variables[substitute.FindStringSubmatch(v)[2]].Expression), -1)
				}
			}
		}
	} else {
		if strings.Contains(str, "$") {
			for _, v := range strings.Fields(str) {
				if strings.Contains(v, "$") {
					str = strings.Trim(strings.Replace(str, v, substitute.ReplaceAllString(v, "${1}"+inc.Variables[substitute.FindStringSubmatch(v)[2]].Expression), -1), " \t")
				}
			}
		}
	}
	return str, bo, nil
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
		fmt.Println(obj)
	}
	return
}
