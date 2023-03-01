package buttons

import (
	"github.com/bwmarrin/discordgo"
	bot "github.com/ubergeek77/uberbot/core"
)

var buttonsInfo = bot.CreateCommandInfo(
	"buttons",
	"description",
	true,
	bot.Utility)

func buttons(ctx *bot.Context) {
	// Message components are buttons, selects.
	// Ephemeral means that when ran as a slash command, the message will only be shown to the user
	// NewResponse(ctx *bot.CmdContext, messageComponents bool, ephemeral bool)
	response := bot.NewResponse(ctx, true, false)

	arg := ctx.Args["buttonname"].StringValue()
	response.PrependField("My button name is", arg, false)
	response.AppendField("Lets use", "buttons everywhere", false)

	response.AppendButton("Append Button", discordgo.PrimaryButton, "", "buttons:1", 0)
	response.AppendButton("Button's are cool!", discordgo.SecondaryButton, "", "buttons:2", 0)
	response.Send(true, "Lets test buttons", "")
}

func init() {
	buttonsInfo.AddArg("buttonname", bot.String, bot.ArgOption, "button name", true, "")
	bot.AddCommand(buttonsInfo, buttons)
	bot.AddSlashCommand(buttonsInfo)
}
