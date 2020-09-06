package framework

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

type DiscordService interface {
	GetSession() *discordgo.Session
	SetSession(session *discordgo.Session)
	GetGuild() *discordgo.Guild
	GetChannel() *discordgo.Channel
	GetUser() *discordgo.User
	GetMessage() *discordgo.MessageCreate
	GetEmbedLayout() *discordgo.MessageEmbed
	GetArgs() []string
	SetArgs(args []string)

	RetrieveUser(id string) (*discordgo.User, error)
	Reply(content string) *discordgo.Message
	Embed(content *discordgo.MessageEmbed) *discordgo.Message
	Help(content string) *discordgo.Message
}

type Discord struct {
	Session            *discordgo.Session
	Guild              *discordgo.Guild
	TextChannel        *discordgo.Channel
	User               *discordgo.User
	Message            *discordgo.MessageCreate
	EmbedMessageLayout *discordgo.MessageEmbed
	Args               []string
}

func NewDiscord(session *discordgo.Session, guild *discordgo.Guild, textChannel *discordgo.Channel,
	user *discordgo.User, message *discordgo.MessageCreate) DiscordService {
	return &Discord{
		Session:     session,
		Guild:       guild,
		TextChannel: textChannel,
		User:        user,
		Message:     message,
		EmbedMessageLayout: &discordgo.MessageEmbed{
			Title: "Svachan",
			Color: 261306,
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Svachan v2.0",
			},
		},
	}
}

func (d Discord) GetSession() *discordgo.Session {
	return d.Session
}

func (d *Discord) SetSession(session *discordgo.Session) {
	d.Session = session
}

func (d Discord) GetGuild() *discordgo.Guild {
	return d.Guild
}

func (d Discord) GetUser() *discordgo.User {
	return d.User
}

func (d Discord) GetChannel() *discordgo.Channel {
	return d.TextChannel
}

func (d Discord) GetMessage() *discordgo.MessageCreate {
	return d.Message
}

func (d Discord) GetEmbedLayout() *discordgo.MessageEmbed {
	return d.EmbedMessageLayout
}

func (d Discord) GetArgs() []string {
	return d.Args
}

func (d *Discord) SetArgs(args []string) {
	d.Args = args
}

func (d Discord) RetrieveUser(id string) (*discordgo.User, error) {
	return d.Session.User(id)
}

func (d Discord) Reply(content string) *discordgo.Message {
	if content == "" {
		return nil
	}

	msg, err := d.Session.ChannelMessageSend(d.TextChannel.ID, content)
	if err != nil {
		fmt.Println("Error whilst sending message,", err)
		return nil
	}
	return msg
}

func (d Discord) Embed(content *discordgo.MessageEmbed) *discordgo.Message {
	msg, err := d.Session.ChannelMessageSendEmbed(d.TextChannel.ID, content)
	if err != nil {
		fmt.Println("Error whilst sending message,", err)
		return nil
	}
	return msg
}

func (d Discord) Help(content string) *discordgo.Message {
	if content == "" {
		return nil
	}

	help := d.EmbedMessageLayout
	help.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL: "https://gbf.wiki/images/9/9b/20205.jpg",
	}
	help.Description = content

	return d.Embed(help)
}
