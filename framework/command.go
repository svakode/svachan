package framework

import "github.com/bwmarrin/discordgo"

type (
	Command func(*Context) *discordgo.Message

	CmdMap map[string]Command

	CommandHandler struct {
		cmds CmdMap
	}
)

func NewCommandHandler() *CommandHandler {
	return &CommandHandler{make(CmdMap)}
}

func (handler CommandHandler) Get(name string) (*Command, bool) {
	cmd, found := handler.cmds[name]

	return &cmd, found
}

func (handler CommandHandler) Register(name string, command Command) {
	handler.cmds[name] = command
}
