package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type roleMap map[string]string

var guildMap = make(map[string]roleMap)

func main() {
	//Get auth token
	authToken := os.Getenv("DISCORD_BOT_TOKEN")
	if authToken == "" {
		log.Fatal("No auth token")
	}
	discord, err := discordgo.New("Bot " + authToken)
	if err != nil {
		log.Fatal("Bot creation failed")
	}

	discord.AddHandler(messageCreate)
	discord.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	err = discord.Open()
	if err != nil {
		log.Fatal("Could not open a connection")
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	discord.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	var avocadoRole string
	var baguetteRole string
	var ekksdeeRole string
	//Ignore own messages
	if m.Author.ID == s.State.User.ID {
		return
	}
	//Get avocado role ids
	if _, ok := guildMap[m.GuildID]; ok {
		avocadoRole = guildMap[m.GuildID]["avocado"]
		baguetteRole = guildMap[m.GuildID]["baguette"]
		ekksdeeRole = guildMap[m.GuildID]["ekks dee"]
	} else {
		roles, err := s.GuildRoles(m.GuildID)
		rmap := roleMap{}
		if err != nil {
			log.Println("Error retrieving roles")
		}
		for _, role := range roles {
			if strings.ToLower(role.Name) == "avocado" {
				rmap["avocado"] = role.ID
			}
			if strings.ToLower(role.Name) == "baguette" {
				rmap["baguette"] = role.ID
			}
			if strings.ToLower(role.Name) == "ekks dee" {
				rmap["ekks dee"] = role.ID
			}
		}
		guildMap[m.GuildID] = rmap
		avocadoRole = guildMap[m.GuildID]["avocado"]
		baguetteRole = guildMap[m.GuildID]["baguette"]
		ekksdeeRole = guildMap[m.GuildID]["ekks dee"]
	}
	// Avocado user
	if strings.HasPrefix(m.Content, "!avocado") || strings.HasPrefix(m.Content, "!🥑") {
		s.GuildMemberRoleAdd(m.GuildID, m.Mentions[0].ID, avocadoRole)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Avocado'd %s uwu", m.Mentions[0].Username))
	}
	// Unavocado user
	if strings.HasPrefix(m.Content, "!unavocado") || strings.HasPrefix(m.Content, "!❌🥑") {
		s.GuildMemberRoleRemove(m.GuildID, m.Mentions[0].ID, avocadoRole)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Unavocado'd %s uwu", m.Mentions[0].Username))
	}
	// Baguette user
	if strings.HasPrefix(m.Content, "!baguette") {
		fmt.Println(baguetteRole)
		s.GuildMemberRoleAdd(m.GuildID, m.Mentions[0].ID, baguetteRole)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Baguette'd %s uwu", m.Mentions[0].Username))
	}
	// Unbaguette user
	if strings.HasPrefix(m.Content, "!unbaguette") {
		s.GuildMemberRoleRemove(m.GuildID, m.Mentions[0].ID, baguetteRole)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Unbaguette'd %s uwu", m.Mentions[0].Username))
	}
	// XD user
	if strings.HasPrefix(m.Content, "!XD") || strings.HasPrefix(m.Content, "!😆") {
		fmt.Println(ekksdeeRole)
		s.GuildMemberRoleAdd(m.GuildID, m.Mentions[0].ID, ekksdeeRole)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("😆'd %s uwu", m.Mentions[0].Username))
	}
	// UnXD user
	if strings.HasPrefix(m.Content, "!unXD") || strings.HasPrefix(m.Content, "!❌😆") {
		s.GuildMemberRoleRemove(m.GuildID, m.Mentions[0].ID, ekksdeeRole)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Un😆'd %s uwu", m.Mentions[0].Username))
	}
}
