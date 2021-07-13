package bot

import (
	"fmt"
	"log"
	"os"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

var (
	client    = &http.Client{}
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name: "play",
			Description: "Play a new game of Werewolf",
		},
		{
			Name: "vote",
			Description: "Vote for another player",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type: discordgo.ApplicationCommandOptionUser,
					Name: "for",
					Description: "Who to vote for",
					Required: true,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"play": Play,
		"vote": Vote,
	}
)

type Bot struct {
	owner   *discordgo.User
	session *discordgo.Session
}

func New(auth string) *Bot {
	var err error

	session, err := discordgo.New(auth)
	if err != nil {
		log.Fatal(err)
	}

	b := &Bot{session: session}

	return b
}

func (b *Bot) Start() {
	b.session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Who are you waiting for?")
	})
	err := b.session.Open()
	if err != nil {
		log.Fatal(err)
	}

	if os.Getenv("VERSION") != "" {
		b.session.UpdateGameStatus(0, os.Getenv("VERSION"))
	}

	for _, v := range commands {
		log.Printf("creating slash command %s", v.Name)
		_, err := b.session.ApplicationCommandCreate(b.session.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
	}

	// a handler for "discord sent us an interaction, map it to a handler"
	b.session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func (b *Bot) Close() {
	for _, v := range commands {
		if err := b.session.ApplicationCommandDelete(b.session.State.User.ID, "", v.Name); err != nil {
			log.Printf("error removing command %s: %s", v.Name, err)
		}
	}
	b.session.Close()
	log.Println("But trust our story's end can bring redemption for the pain endured")
}

func Play(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Werewolf coming soon...",
		},
	})
}

func Vote(s *discordgo.Session, i *discordgo.InteractionCreate) {
	from := getNick(i.Member)
	var to string
	toUser := i.ApplicationCommandData().Options[0].UserValue(nil).ID
	toMember, err := s.GuildMember(i.GuildID, toUser)
	if err != nil {
		log.Printf("error getting guild member: %s", err)
		to = i.ApplicationCommandData().Options[0].UserValue(s).Username
	} else {
		to = getNick(toMember)
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("%s voted for %s", from, to),
		},
	})
}

func getNick(m *discordgo.Member) string {
	nick := m.Nick
	if nick == "" {
		nick = m.User.Username
	}
	return nick
}
