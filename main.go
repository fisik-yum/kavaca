/*
   kavaca- a discord bot that acts as a mail forwarder
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

// this is a modified version of ping pong from discordgo examples

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
	token          string
	ownerID        string
	prefix         string
	defaultChannel string
	bindings       []binding
)

func init() {
	_, err := os.Stat("bindings.json")
	if os.IsNotExist(err) {
		os.Create("bindings.json")
	}
	load_bindings()
	fmt.Println(bindings)
	save_bindings()
	fmt.Println(bindings)
	read_config()
	fmt.Println(bindings)
}

func main() {

	dg, err := discordgo.New("Bot " + token)
	check(err)
	dg.AddHandler(messageCreate)
	dg.ShouldReconnectOnError = true
	dg.Identify.Intents = 12800
	err = dg.Open()
	check(err)
	if defaultChannel == "" {
		dc, err := dg.UserChannelCreate(ownerID)
		check(err)
		defaultChannel = dc.ID
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	save_bindings()
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if s.State.User.ID == m.Author.ID {
		return
	}
	if m.Author.ID != ownerID {
		if find_bind(m.Author.ID) == "" {
			create_bind(m.Author.ID, defaultChannel)
		}
		cID := find_bind(m.Author.ID)
		s.ChannelMessageSend(cID, ("**" + m.Author.Username + "<" + m.Author.ID + ">" + ":** " + m.Content))
		return
	}

	//command handling, modular prefixes
	cmd := trim_index(m.Content, 0)
	if m.Author.ID == ownerID && strings.HasPrefix(cmd, prefix) {
		if strings.HasPrefix(cmd, prefix+"bind") { //bind command. add functions to modify binds, otherwise useless.
			bindUS := trim_index(m.Content, 1)
			bindCID := trim_index(m.Content, 2)
			if find_bind(bindUS) == "" {
				create_bind(bindUS, bindCID)
			}
			modify_bind(bindUS, bindCID)
			s.ChannelMessageSend(m.ChannelID, "rebound user "+bindUS)
			return
		}
		if strings.HasPrefix(cmd, prefix+"listbinds") {
			fmt.Println(bindings)
			save_bindings()
			return
		}
		if strings.HasPrefix(cmd, prefix+"reset") {
			bindUS := trim_index(m.Content, 1)
			if find_bind(bindUS) == "" {
				create_bind(bindUS, defaultChannel)
			}
			modify_bind(bindUS, defaultChannel)
			s.ChannelMessageSend(m.ChannelID, "reset bind for user "+bindUS)
			return
		}
		if strings.HasPrefix(cmd, prefix+"savebinds") {
			save_bindings()
		}
		if strings.HasPrefix(cmd, prefix+"info") {
			s.ChannelMessageSend(m.ChannelID, ("` kavaca built with discordgo " + discordgo.VERSION + "`"))
			return
		}
	}
}
