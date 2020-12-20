package main

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func main() {
	discord, err := discordgo.New("Bot " + "<authentication token>")
	if err != nil {
		log.Fatal("oops")
	}
	fmt.Println(discord.LogLevel)
}
