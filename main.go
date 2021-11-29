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

package main

// this is a modified version of the "ping pong example in discordgo examples"

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// variables found in config.json, which needs to exist
var (
	token    string
	ownerID  string
	prefix   string
	bindings []binding
)

func init() {
	_, err := os.Stat("bindings.json")
	if os.IsNotExist(err) {
		os.Create("bindings.json")
	}

	token = readRecipient().Token
	ownerID = readRecipient().ID
	prefix = readRecipient().Prefix
}

func main() {

	dg, err := discordgo.New("Bot " + token)
	check(err)
	dg.AddHandler(messageCreate)

	dg.Identify.Intents = 12800
	err = dg.Open()
	check(err)

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if s.State.User.ID == m.Author.ID {
		return
	}
	if find_bind(m.Author.ID) == "" && m.Author.ID != ownerID {
		nc, err := s.UserChannelCreate(ownerID)
		check(err)
		create_bind(m.Author.ID, nc.ID)
	}
	if m.Author.ID != ownerID { // basic forwarding //wil have to implement message redirects
		cID := find_bind(m.Author.ID)
		s.ChannelMessageSend(cID, ("**" + m.Author.Username + "<" + m.Author.ID + ">" + ":** " + m.Content))
		return
	}

	//command handling, modular prefixes
	cmd := trim_index(m.Content, 0)
	if m.Author.ID == ownerID && strings.HasPrefix(cmd, prefix) {
		s.ChannelMessageSend(m.ChannelID, ("` ghatam built with discordgo " + discordgo.VERSION + "`"))
		if strings.HasPrefix(cmd, prefix+"bind") { //bind command. add functions to modify binds, otherwise useless.
			bindID := trim_index(m.Content, 1)
			if find_bind(bindID) == "" { // this can obviously be shrunk
				nc, err := s.UserChannelCreate(ownerID)
				check(err)
				create_bind(m.Author.ID, nc.ID)
			}
			modify_bind(bindID, m.ChannelID)
			s.ChannelMessageSend(m.ChannelID, "Created binding rule for user "+bindID)
			fmt.Println(bindings)
		}
		if strings.HasPrefix(cmd, prefix+"listbinds") { //bind command
			fmt.Println(bindings)
		}
	}
}
