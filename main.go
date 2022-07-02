// Buanu, an experimental command-line tool for cryptography.
// Copyright (C) 2022 Brian Zhu
//
// This program is free software: you can redistribute it
// and/or modify it under the terms of the GNU General Public
// License as published by the Free Software Foundation, either
// version 3 of the License, or (at your option) any later
// version.
//
// This program is distributed in the hope that it will be
// useful, but WITHOUT ANY WARRANTY; without even the implied
// warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR
// PURPOSE.  See the GNU General Public License for more
// details.
//
// You should have received a copy of the GNU General Public
// License along with this program.  If not, see
// <https://www.gnu.org/licenses/>.
package main

import (
	"buanu/cmdsum"
	_ "embed"
	"fmt"
	"log"
	"os"
)

//go:embed help.txt
var help string

func main() {
	if len(os.Args) < 2 {
		fmt.Print(help)
		os.Exit(1)
	}
	subcommand := os.Args[1]
	switch subcommand {
	case "sum":
		cmdsum.Run()
	default:
		log.Fatalf("%q is not a valid subcommand.", subcommand)
	}
}
