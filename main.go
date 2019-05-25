package main

import (
	"github.com/Mattemagikern/Bob/parser"
	"github.com/Mattemagikern/Bob/utils"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	var jobs = []string{}
	dat, err := ioutil.ReadFile("./Blueprint")
	if err != nil {
		log.Fatal(err)
	}
	parser.ParseBuilder(string(dat))
	if err := utils.DFS(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	jobs = append(jobs, os.Args[1:]...)

	for _, v := range jobs {
		err := utils.Execute(v)
		if err != nil {
			log.Fatal(err, v)
		}
	}
}
