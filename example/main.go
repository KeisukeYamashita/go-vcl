package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/KeisukeYamashita/go-vcl/vcl"
)

// Root represents the top level object which is a file
type Root struct {
	ACls []*ACL `vcl:"acl,block"`
}

// ACL are acl blocks
type ACL struct {
	Type      string   `vcl:"type,label"`
	Endpoints []string `vcl:",flat"`
}

func main() {
	dat, err := ioutil.ReadFile("./example/vcl.vcl")
	if err != nil {
		log.Fatal(err)
	}

	r := &Root{}
	if errs := vcl.Decode(dat, r); len(errs) > 0 {
		log.Fatal(errs)
	}

	fmt.Println(r.ACls)
	fmt.Println(r.ACls[0].Endpoints)
}
