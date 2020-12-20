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

type RoleMap map[string]string

var GuildMap = make(map[string]RoleMap)

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
	var modocadoRole string
	//Ignore own messages
	if m.Author.ID == s.State.User.ID {
		return
	}
	if _, ok := GuildMap[m.GuildID]; ok {
		avocadoRole = GuildMap[m.GuildID]["avocado"]
		modocadoRole = GuildMap[m.GuildID]["modocado"]
	} else {
		roles, err := s.GuildRoles(m.GuildID)
		roleMap := RoleMap{}
		if err != nil {
			log.Println("Error retrieving roles")
		}
		for _, role := range roles {
			if strings.ToLower(role.Name) == "avocado" {
				roleMap["avocado"] = role.ID
			}
			if strings.ToLower(role.Name) == "modocado" {
				roleMap["modocado"] = role.ID
			}
		}
		GuildMap[m.GuildID] = roleMap
		avocadoRole = GuildMap[m.GuildID]["avocado"]
		modocadoRole = GuildMap[m.GuildID]["modocado"]
	}
	// Avocado user
	if strings.HasPrefix(m.Content, "!avocado") {
		s.GuildMemberRoleAdd(m.GuildID, m.Mentions[0].ID, avocadoRole)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Avocado'd %s uwu", m.Mentions[0].Mention()))
	}
	// Unavocado user
	if strings.HasPrefix(m.Content, "!unavocado") {
		s.GuildMemberRoleRemove(m.GuildID, m.Mentions[0].ID, avocadoRole)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Unavocado'd %s uwu", m.Mentions[0].Mention()))
	}
	// Modocado user
	if strings.HasPrefix(m.Content, "!modocado") {
		fmt.Println(modocadoRole)
		s.GuildMemberRoleAdd(m.GuildID, m.Mentions[0].ID, modocadoRole)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Avocado'd %s uwu", m.Mentions[0].Mention()))
	}
	// Unmodocado user
	if strings.HasPrefix(m.Content, "!unmodocado") {
		s.GuildMemberRoleRemove(m.GuildID, m.Mentions[0].ID, modocadoRole)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Unavocado'd %s uwu", m.Mentions[0].Mention()))
	}
}
