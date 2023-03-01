
package msgsniper

import bot "github.com/ubergeek77/uberbot/core"

var sniperAddInfo = bot.CreateCommandInfo("add", "adds to the sniper", false, bot.Moderation)
var sniperRemoveInfo = bot.CreateCommandInfo("remove", "removes from the sniper", false, bot.Moderation)

func subCommandAdd(ctx *bot.Context) {
	response := bot.NewResponse(ctx, false, true)
	if ctx.Args["option"].StringValue() == "" {
		response.Send(false, "Sniper", "Invalid Syntax, you did not provide a word to add!")
		return
	}
	ctx.Guild.AddWord(ctx.Args["option"].StringValue())
	response.Send(true, "Sniper", "Added Word: "+ctx.Args["option"].StringValue()+" to the sniper list")
	return
}

func subCommandRemove(ctx *bot.Context) {
	response := bot.NewResponse(ctx, false, true)
	if ctx.Args["option"].StringValue() == "" {
		response.Send(false, "Sniper", "Invalid Syntax, you did not provide a word to remove!")
	}
	ctx.Guild.RemoveWord(ctx.Args["option"].StringValue())
	response.Send(true, "Sniper", "Remove Word: "+ctx.Args["option"].StringValue()+" from the sniper list")
}

func init() {
	sniperAddInfo.SetParent(false, "sniper")
	sniperRemoveInfo.SetParent(false, "sniper")
	sniperAddInfo.AddArg("option", bot.String, bot.ArgOption, "string to add", true, "")
	sniperRemoveInfo.AddArg("option", bot.String, bot.ArgOption, "string to remove", true, "")
	bot.AddChildCommand(sniperAddInfo, subCommandAdd)
	bot.AddChildCommand(sniperRemoveInfo, subCommandRemove)
}
