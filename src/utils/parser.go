package utils

import (
	"errors"
	"inc"
	"io/ioutil"
	"regexp"
	"strings"
)

var variables *regexp.Regexp = regexp.MustCompile(`(?m)^(?:\$)?(\S*)\s*(?: =|\+=|\-=)(?:\s*(.*))`)
var recepies *regexp.Regexp = regexp.MustCompile(`(?m)^(\S*):(.*)(\n(.+))+`)
var suffixes *regexp.Regexp = regexp.MustCompile(`(?m)^search\s+\{(.|\n)*^\}`)

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
		lines := strings.Split(v, "\n")
		tmp := strings.Split(lines[0], ":")
		name := tmp[0]
		lines = lines[1:]
		ingredients := strings.Fields(tmp[1])
		inc.Recepies[name] = &inc.Recepie{name, ingredients, lines}
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
					inc.Variables[name].Expression = strings.Replace(inc.Variables[name].Expression, v, inc.Variables[v[1:]].Expression, -1)
				}
			}
		}
	} else {
		if strings.Contains(str, "$") {
			for _, v := range strings.Fields(str)[1:] {
				if strings.Contains(v, "$") {
					str = strings.Replace(str, v, inc.Variables[v[1:]].Expression, -1)
				}
			}
		}
	}
	return str, bo, nil
}
