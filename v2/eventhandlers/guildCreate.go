package eventhandlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/ubergeek77/uberbot/v2/core"
)

func guildCreate(s *discordgo.Session, evt *discordgo.GuildCreate) {
	core.Log.Infof("found guild %s (%s) with %d members", evt.Name, evt.ID, len(evt.Members))
	// this has to be done since we don't manage our own state.
	// the guild we get from this event isn't updated, idk why it's a pointer
	g, err := s.State.Guild(evt.ID)
	if err != nil {
		core.Log.Errorf("unable to find guild %s (%s). maybe race condition?")
		return
	}
	if core.GuildExists(g.ID) {
		core.Log.Infof("guild already exists in memory already, are we reconnecting?")
		core.Log.Infof("figuring out if we should reregister...")
	} else {
		core.AddGuild(g)
	}
}

func init() {
	core.AddHandler(guildCreate)
}
