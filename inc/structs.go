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

type Recipe struct {
	Name         string
	Dependencies []string
	Commands     []string
}
type Build struct {
	Name       string
	Extensions []string
	Commands   []string
}

type Variable struct {
	Name       string
	Expression string
}

var File_tree map[string]*File = make(map[string]*File)
var Inc_tree map[string]*File = make(map[string]*File)
var State map[string]*Object_file = make(map[string]*Object_file)
var Recipes map[string]*Recipe = make(map[string]*Recipe)
var Sf Suffix
var Build_cmd *Build

var Variables = map[string]*Variable{
	"Objects": {"Objects", ""},
	"CFLAGS":  {"CFLAGS", ""},
}
