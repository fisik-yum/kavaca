/*
   Ghatam- a discord bot that acts as a mail forwarder
   Copyright (C) 2021  fisik_yum

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

//read (and write) configuration files
package main

import (
	"encoding/json"
	"io/ioutil"
)

func readRecipient() owner { // main config file for user
	f, err := ioutil.ReadFile("config.json")
	check(err)
	var userData owner
	err = json.Unmarshal([]byte(f), &userData)
	check(err)
	return userData
}

func load_bindings(bindingmap []binding) { //yet to be implemented

}
