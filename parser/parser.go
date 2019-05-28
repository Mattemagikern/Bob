package parser

import (
	"log"
	"os/exec"
	"regexp"
	"strings"
)

type Recipe struct {
	Name        string
	Ingredients []string
	Cmds        []string
}

var Store = map[string]string{"@": "RecipieName", "var": "hello"}
var Recipes = map[string]*Recipe{}

func Init_variable(input string) bool {
	indx := strings.Index(input, "=")
	if indx == -1 {
		return false
	}
	variableName := strings.TrimSpace(input[:indx])
	expression := strings.TrimSpace(input[indx+1:])
	expression = shell(expression)
	Store[variableName] = expression
	return true
}

func Update_variables(input string) bool {
	indx := strings.IndexAny(input, "+=-")
	if indx == -1 {
		return false
	}
	i := strings.IndexAny(input[:indx], "$")
	if i == -1 || (input[indx+1] != '=' && input[indx] != '=') {
		return false
	}

	variableName := strings.TrimSpace(input[i+1 : indx])
	switch {
	case input[indx:indx+2] == "+=":
		expression := input[indx+2:]
		expression = shell(expression)
		Store[variableName] += expression
	case input[indx:indx+2] == "-=":
		expression := input[indx+2:]
		expression = shell(expression)
		Store[variableName] = strings.Trim(Store[variableName], expression)
	case input[indx] == '=':
		expression := input[indx+1:]
		expression = shell(expression)
		Store[variableName] = expression
	}
	return true
}

/* Used in EXEC command */
func Substitute(input string) string {
	str := ""
	for {
		indx := strings.IndexAny(input, "$/ ")
		switch {
		case indx == -1 || indx == len(input)-1:
			return str + input
		case input[indx:indx+2] == "$(":
			if i := strings.Index(input, ")"); i != -1 {
				str += input[:indx]
				out := shell(input[indx : i+1])
				str += out
				input = input[i+1:]
				continue
			}
			log.Fatal("Subsitute Error: malformed shell command:", input)
		case input[indx] == '$':
			next := strings.IndexAny(input[indx+1:], "$/.= ")
			if next != -1 {
				str += Store[input[indx+1:next+1]]
				input = input[next+1:]
			} else {
				str += Store[input[indx+1:]]
				return str
			}
		default:
			//a space
			str += input[:indx+1]
			input = input[indx+1:]
		}
	}
}

func shell(str string) string {
	str = strings.TrimSpace(str)
	i := strings.LastIndex(str, ")")
	if i == -1 || str[:2] != "$(" {
		return str
	}
	cmd := str[2:i]
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.Fatal(str, err)
	}
	if string(out)[len(out)-1] == '\n' {
		return string(out[:len(out)-1])
	}
	return string(out)
}

func ParseBuilder(builder string) {
	var r *regexp.Regexp = regexp.MustCompile(`(?m)^(\S*)\s?:\s?(.*)$`)
	lines := strings.Split(builder, "\n")
	for i := 0; i < len(lines); i++ {
		if r.Match([]byte(lines[i])) {
			tmp := strings.Split(lines[i], ":")
			name := strings.TrimSpace(tmp[0])
			tmp = strings.Split(strings.TrimSpace(tmp[1]), " ")
			ingredients := []string{}
			if tmp[0] != "" {
				ingredients = tmp
			}
			i++
			cmds := []string{}
			for ; i < len(lines) && len(lines[i]) > 1 && lines[i][0] == '\t'; i++ {
				cmds = append(cmds, strings.TrimSpace(lines[i]))
			}
			Recipes[name] = &Recipe{name, ingredients, cmds}
		} else {
			if !Update_variables(lines[i]) {
				Init_variable(lines[i])
			}
		}
	}
}

/*
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
*/
