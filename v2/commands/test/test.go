package test

import (
	"fmt"
	"runtime"

	"github.com/bwmarrin/discordgo"
	bot "github.com/ubergeek77/uberbot/v2/core"
)

var buttonsInfo = bot.CreateCommandInfo(
	"test",
	"description",
	true,
	bot.Utility)

func buttons(ctx *bot.CmdContext) {
	response := bot.NewResponse(ctx, true, false, 1)

	response.PrependField(0, "Buttons are", "very cool!", false)
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	response.PrependField(0, "Info:", fmt.Sprintf("```go version: %s,\nos: %s,\nmemory: %dmb```", runtime.Version(), runtime.GOOS, mem.Sys/1000000), false)
	response.AppendButton("Append Button", discordgo.PrimaryButton, "", "buttons:1", 0)
	response.AppendButton("Button's are cool!", discordgo.SecondaryButton, "", "buttons:2", 0)
	response.Send(true, "Buttons", "idk", 0)
}

func handleButtons(ctx *bot.InteractionCtx) {
	content := "Currently testing customid " + ctx.MessageComponentData().CustomID
	ctx.Message.Embeds[0].Description = content
	err := ctx.Session.InteractionRespond(ctx.Interaction, &discordgo.InteractionResponse{
		// Buttons also may update the message which they was attached to.
		// Or may just acknowledge (InteractionResponseDeferredMessageUpdate) that the event was received and not update the message.
		// To update it later you need to use interaction response edit endpoint.
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			TTS:     false,
			Content: "the game",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		bot.Log.Errorf("error while responding to interaction")
		bot.Log.Error(err.Error())
		return
	}
	_, err = ctx.Session.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Embeds:  ctx.Message.Embeds,
		ID:      ctx.Message.ID,
		Channel: ctx.Message.ChannelID,
	})

	if err != nil {
		bot.Log.Errorf("error while responding to interaction")
		bot.Log.Error(err.Error())
		return
	}
}

func init() {
	bot.AddCommand(buttonsInfo, buttons)
	bot.AddInteractHandler(&bot.InteractionInfo{
		Id: "modal1:1",
	}, handleButtons)
	bot.AddInteractHandler(&bot.InteractionInfo{
		Id: "buttons:2",
	}, handleButtons)
	bot.AddSlashCommand(buttonsInfo)
}
