package parser

import (
	"log"
	"testing"
)

func TestSubsitutePath(t *testing.T) {
	Store["DESTDIR"] = "/mnt"
	Store["CC"] = "go build"
	Store["@"] = "test"
	input := "$DESTDIR/etc/stuff"
	str := Substitute(input)
	if str != "/mnt/etc/stuff" {
		t.Errorf("error:%s\n", str)
	}
	input = "$CC -o $DESTDIR/usr/bin/$@"
	str = Substitute(input)
	if str != "go build -o /mnt/usr/bin/test" {
		t.Errorf("error:%s\n", str)
	}

}

func TestSubstitute(t *testing.T) {
	Store["CC"] = "clang"
	Store["CFLAGS"] = "-Wall"

	s := "export GOPATH=$(echo pass); $CC -o main main.c"
	str := Substitute(s)
	if str != "export GOPATH=pass; clang -o main main.c" {
		t.Errorf("error: exp:export GOPATH=pass; clang -o main main.c [[ acctual: %s\n", str)
	}

	str = Substitute("$CC hello $(echo stuff)")
	if str != "clang hello stuff" {
		t.Errorf("error:%s\n", str)
	}

	s = "export HELLO=$(echo pass); $CC $CFLAGS -o server main.c"
	str = Substitute(s)
	if str != "export HELLO=pass; clang -Wall -o server main.c" {
		t.Errorf("error:%s\n", str)
	}
}

func TestParseBuilder(t *testing.T) {
	builder := `CC      = gcc
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
	$CFLAGS += -DDEBUG`

	ParseBuilder(builder)
	if Store["CC"] != "gcc" {
		t.Errorf("Builder parsing failed")
	}
}
