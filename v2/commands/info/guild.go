package info

import (
	bot "github.com/ubergeek77/uberbot/v2/core"
)

var guildInfoCmd = bot.CreateCommandInfo("guild", "displays guild info", true, bot.Utility)

func guildInfo(ctx *bot.CmdContext) {
	response := bot.NewResponse(ctx, true, false, 0)
	// set author fields instead of using title
	response.PrependAuthor(0, ctx.Guild.Name, "", ctx.Guild.IconURL())
	// append footer data
	response.AppendFooter(0, ctx.Guild.ID, "", true)
	response.Send(true, "", "", 0)
}
func init() {
	bot.AddCommand(guildInfoCmd, guildInfo)
	bot.AddSlashCommand(guildInfoCmd)
}
