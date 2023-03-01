package deleteabusechecker

import (
	"encoding/json"
	"time"

	"github.com/bwmarrin/discordgo"
	tlog "github.com/ubergeek77/tinylog"
	bot "github.com/ubergeek77/uberbot/core"
)

var log = tlog.NewTaggedLogger("SnipeCassian", tlog.NewColor("38;5;81"))

type deleteInfo struct {
	DeleteCount  int64 `json:"deleteCount"`
	CheckExpires int64 `json:"checkExpires"`
}

type userList map[string]deleteInfo

func handleDeletedMessage(s *discordgo.Session, messageDelete *discordgo.MessageDelete) {
	// If the message is not cached, do nothing
	if messageDelete.BeforeDelete == nil {
		return
	}

	// Get the guild ID
	guild := bot.Guilds[messageDelete.BeforeDelete.GuildID]

	// Get the author object
	author := messageDelete.BeforeDelete.Author

	// Ignore bot admins
	if bot.IsAdmin(messageDelete.BeforeDelete.Author.ID) {
		return
	}

	// Ignore moderators
	if guild.IsMod(messageDelete.BeforeDelete.Author.ID) {
		return
	}

	// Only apply to Cassian
	if messageDelete.BeforeDelete.Author.ID != "468164718045954048" {
		return
	}

	// Ignore everything if the mute role is unset or invalid
	muteRoleId := guild.Info.MuteRoleId
	if muteRoleId == "" {
		return
	}

	_, err := guild.GetRole(muteRoleId)
	if err != nil {
		return
	}

	// Get a guild map for this data, or initialize one if necessary
	var list userList
	rawData, err := guild.GetMap("deleteAbuse")
	if err != nil {
		list = userList{}
	} else {
		// Convert the interface to the userInfo type via json tags
		b, err := json.Marshal(rawData)
		if err != nil {
			log.Errorf("Failed to marshal guild storage data")
			return
		}
		list = userList{}
		err = json.Unmarshal(b, &list)
		if err != nil {
			log.Errorf("Failed to unmarshal guild storage data")
			return
		}
	}

	// Check if the current author has a record in the map, make a new one if not
	var userInfo deleteInfo
	if exist, ok := list[author.ID]; ok {
		userInfo = exist
	} else {
		userInfo.DeleteCount = 0
		userInfo.CheckExpires = time.Now().Unix() + 60
	}

	// Increment the counter
	userInfo.DeleteCount++

	// If the check has expired, set the counter to 0 and add a new expiry
	// If the time hasn't expired, check the message count and act accordingly
	if time.Now().Unix() > userInfo.CheckExpires {
		userInfo.DeleteCount = 0
		userInfo.CheckExpires = time.Now().Unix() + 60
	} else if userInfo.DeleteCount >= 5 {
		// Start a new response
		response := bot.NewResponse(&bot.Context{
			Guild: guild,
			Cmd:   bot.CommandInfo{},
			Args:  nil,
			Message: &discordgo.Message{
				ChannelID: messageDelete.ChannelID,
				Author:    bot.Session.State.User,
			},
		}, false, false)

		// Try muting the member
		err := guild.Mute(author.ID, 300)
		if err != nil {
			log.Errorf("Unable to mute delete abuser %s: %s", author.Username, err)
			response.Send(false, "Failed to mute member", "An error occurred when muting a member abusing the delete feature")
		} else {
			// Reset deleted messages
			userInfo.DeleteCount = 0
			response.PrependField("Target Member", author.Mention(), false)
			response.PrependField("Mute duration", "5 Minutes", false)
			response.Send(true, "Someone deleted messages too fast!", "")
		}
	}

	// Save the user data to the list
	list[author.ID] = userInfo

	// Convert the list back to a map[string]interface{}
	interfaceList := make(map[string]interface{})

	for key, value := range list {
		interfaceList[key] = value
	}

	// Save the list
	guild.StoreMap("deleteAbuse", interfaceList)
}

func init() {
	bot.AddHandler(handleDeletedMessage)
}
