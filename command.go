package main

/*
 * command.go
 * Handles commands
 * By J. Stuart McMurray
 * Created 20150815
 * Last Modified 20150815
 */

/*
 * ntsh -- The "Nice Try" shell
 * version 0.0.1, August 15, 2015
 *
 * Copyright (C) 2015 Stuart McMurray and Josiah Hamilton
 *
 * This software is provided 'as-is', without any express or implied
 * warranty.  In no event will the authors be held liable for any damages
 * arising from the use of this software.
 *
 * Permission is granted to anyone to use this software for any purpose,
 * including commercial applications, and to alter it and redistribute it
 * freely, subject to the following restrictions:
 *
 * 1. The origin of this software must not be misrepresented; you must not
 *    claim that you wrote the original software. If you use this software
 *    in a product, an acknowledgment in the product documentation would be
 *    appreciated but is not required.
 * 2. Altered source versions must be plainly marked as such, and must not be
 *    misrepresented as being the original software.
 * 3. This notice may not be removed or altered from any source distribution.
 *
 * Stuart McMurray      Josiah Hamilton
 * kd5pbo@gmail.com     dev.x.josiah@mamber.net
 */

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

/* Command is the function prototype for a defined command.  It will be called
when the user types a command starting with c.  a are the remaining arguments.
Output should be sent to out, which will be closed when the command returns. */
type Command func(c string, a []string, out io.Writer)

/* Commands holds the set of defined commands */
var (
	commands  map[string]Command
	commandsL sync.Mutex
)

/* Register registrs f to be called when the user executes c */
func Register(c string, f Command) {
	commandsL.Lock()
	defer commandsL.Unlock()
	/* Make sure we actually have a map */
	if nil == commands {
		commands = make(map[string]Command)
	}
	/* Don't double-add */
	if _, ok := commands[c]; ok {
		panic(fmt.Sprintf("%v already defined", c))
	}
	/* Register command */
	commands[c] = f
}

/* run parses a line and runs the appropriate command.  Source should be
something that identifies the caller, like an IP address.  Ding should be a
"\a" for a bell every command, or the empty string. */
func run(line, source, ding string) error {
	commandsL.Lock()
	defer commandsL.Unlock()
	/* Make sure we actually have a command */
	if 0 == len(line) {
		return nil
	}
	/* Split into fields */
	a := strings.Fields(line)
	if 0 == len(a) {
		return nil
	}
	/* Get the function to call */
	f, ok := commands[a[0]]
	if !ok {
		fmt.Printf("Nice Try!\n")
		log.Printf(
			"%v!: Unable to find command %v",
			source,
			strconv.Quote(line),
		)
		return nil
	}
	/* Comms channel */
	pr, pw := io.Pipe()
	defer pr.Close()
	defer pw.Close()
	go io.Copy(os.Stdout, pr)

	/* Start the command */
	log.Printf("%v%v: %v", ding, source, strconv.Quote(line))

	f(a[0], a[1:], pw)
	return nil
}
