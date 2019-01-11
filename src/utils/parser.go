package utils

import (
	"errors"
	"inc"
	"io/ioutil"
	"regexp"
	"strings"
)

func Parse_builder() (err error) {
	var dat []byte
	variables := regexp.MustCompile(`(?m)^([^\s].*)(\W?=)[\s+]?(.*)$`)
	recepies := regexp.MustCompile(`(?m)^(\S*):(.*)(\n(.+))+`)
	suffixes := regexp.MustCompile(`()^suffixes\s+{([^}]*)}`)

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
		if err = substitute(v); err != nil {
			panic("Could not substitute")
		}
	}

	if inc.Sf.Inc == nil || inc.Sf.Src == nil || inc.Sf.Objects == nil || inc.Sf.Inc_pattern == nil {
		err = errors.New("Suffixes: objecs, srcs or, inc were not found")
	}

	return
}

func substitute(v string) error {
	var tmp []string
	var name string
	var expression string
	switch {
	case strings.Contains(v, "+="):
		tmp = strings.Split(v, "+=")
		name = strings.Trim(tmp[0], " \t")
		expression = strings.Trim(tmp[1], " \t")
		inc.Variables[name].Expression += " " + expression

	case strings.Contains(v, "-="):
		tmp = strings.Split(v, "-=")
		name = strings.Trim(tmp[0], " \t")
		expression = strings.Trim(tmp[1], " \t")
		inc.Variables[name].Expression = strings.Trim(inc.Variables[name].Expression, expression)

	case strings.Contains(v, "="):
		tmp = strings.Split(v, "=")
		name = strings.Trim(tmp[0], " \t")
		expression = strings.Trim(tmp[1], " \t")
		inc.Variables[name] = &inc.Variable{name, expression}
	}
	if strings.Contains(inc.Variables[name].Expression, "$") {
		for _, v := range strings.Fields(inc.Variables[name].Expression) {
			if strings.Contains(v, "$") {
				inc.Variables[name].Expression = strings.Replace(inc.Variables[name].Expression, v, inc.Variables[v[1:]].Expression, -1)
			}
		}
	}
	return nil
}
