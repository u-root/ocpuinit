// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/u-root/u-root/uroot"
)

var (
	namespace = []uroot.Creator {
		Cmd{Cmd: "dhclient", Args: []string{"-ipv4", "-verbose"}, Background: true, Delay: 10},
		Cmd{Cmd: "wget", Args: []string{"http://100.96.221.129:8080/bzImage"}, Stdout: "bzImage"},
		Cmd{Cmd: "ls", Args: []string{"-l"},},
		Cmd{Cmd: "kexec", Args: []string{"/bzImage"},},
	}
			
	verbose   = flag.Bool("v", false, "print all commands")
	debug     = func(string, ...interface{}) {}
)

type Cmd struct {
	Cmd string
	Args []string
	Background bool
	Delay int
	Stdout string
}

func (c Cmd) Create() error {
	cmd := exec.Command(c.Cmd, c.Args...)
	log.Printf("Run %v", c.String())
	if c.Stdout != "" {
		cmd.Stderr = os.Stdout
		f, err := os.OpenFile(c.Stdout, os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatalf("%v: %v", c.Stdout, err)
		}
		cmd.Stdout = f
		err = cmd.Run()
		//f.Close()
		// That was not enough. It's still staying open. Fuck it.
		syscall.Close(int(f.Fd()))
		return err
	}
	o, err := cmd.CombinedOutput()
	log.Printf("%v", string(o))
	return err
}

func (c Cmd) String() string {
	return fmt.Sprintf("%v %v", c.Cmd, c.Args)
}

func main() {
	flag.Parse()
	log.Printf("Welcome to OCP running NERF and u-root!")

	if *verbose {
		debug = log.Printf
	}

	for _, c := range namespace {
		if err := c.Create(); err != nil {
			log.Printf("Error creating %s: %v", c, err)
			break
		} else {
			log.Printf("Created %v", c)
		}
	}


}
