package inc

import (
	"regexp"
	"time"
)

type Suffix struct {
	Objects     *regexp.Regexp
	Inc         *regexp.Regexp
	Src         *regexp.Regexp
	Inc_pattern *regexp.Regexp
}

type Object_file struct {
	Path      string
	Flags     string
	Timestamp time.Time
}

type File struct {
	Path      string
	Inc       []string
	Timestamp time.Time
}

type Recepie struct {
	Name         string
	Dependencies []string
	Commands     []string
}

type Variable struct {
	Name       string
	Expression string
}

var File_tree map[string]File = make(map[string]File)
var Recepies map[string]*Recepie = make(map[string]*Recepie)
var Variables map[string]*Variable = make(map[string]*Variable)
var Sf Suffix
