package parser

import (
	"fmt"
	"inc"
	"testing"
)

func TestSubstitute(t *testing.T) {
	inc.Variables = map[string]*inc.Variable{
		"CC":     {"CC", "clang"},
		"CFLAGS": {"CFLAGS", "-Wall"},
	}
	str, err := Substitute("$CC hello $(echo stuff)")
	if err != nil {
		fmt.Println(str)
		t.Errorf(err.Error())
	}
	fmt.Println(str)

	s := "export HELLO=$(pwd); $CC $CFLAGS -o server ../MasterThesis/code/test/server.c"
	str, err = Substitute(s)
	fmt.Println(str)
	if err != nil {
		t.Errorf("error:%s\n", err.Error())
	}
	fmt.Println(str)
	s = "export GOOPATH=$(pwd); $CC -o bin/master master"
	str, err = Substitute(s)
	fmt.Println(str)
	if err != nil {
		t.Errorf("error:%s\n", err.Error())
	}
	fmt.Println(str)
}

func TestUpdate_vars(t *testing.T) {
	inc.Variables = map[string]*inc.Variable{
		"CC":     {"CC", "clang"},
		"CFLAGS": {"CFLAGS", "-Wall"},
	}
	expected := "clang NopeTown"
	name := "CC"
	delimiter := "+="
	expression := "NopeTown"
	Update_vars(name, delimiter, expression)
	if inc.Variables["CC"].Expression != expected {
		t.Errorf("Result:%s, Expected:%s\n", inc.Variables["CC"].Expression, expected)
	}
}

func TestParseBuilder(t *testing.T) {
	builder := `search {
	inc .*\/inc\/(.*)\.h$
	src .*\/src\/(.*)\.c$
	include_pattern (?m)(?:^#include[\s]*)[<|"](.*)[>|"]
}

CC      = gcc
dep     =  -I../MasterThesis/code/inc -I../MasterThesis/code/inc/pkg -I../MasterThesis/code/test -I../MasterThesis/code/test/testvector
CFLAGS   = -Wall -Wextra -std=c99 -Wno-format -Wno-parentheses -Wno-empty-body -DROHC_UDP -pthread
$CFLAGS  += $dep

all: debug build server client
	echo $LA

server:
	export HELLO=$(pwd); $CC $CFLAGS -o server $Objects ../MasterThesis/code/test/server.c

client:
	$CC $CFLAGS -o client $Objects ../MasterThesis/code/test/client.c

debug:
	$LA = $(ls -la)
	$CFLAGS += -DDEBUG

../MasterThesis/code/obj/%.o:%.c
	$CC $CFLAGS -c -o $@ $<`
	Parse_builder(builder)
	if inc.Variables["CC"].Expression != "gcc" {
		t.Errorf("Builder parsing failed")
	}
}
