package eventhandlers

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/ubergeek77/uberbot/v2/core"
	"github.com/ubergeek77/uberbot/v2/workers"
)

func readyEventHandler(s *discordgo.Session, evt *discordgo.Ready) {
	// Update presence and set a worker to update presence
	UpdatePresence()
	// update presence every 12 hours
	core.WorkerManager.AddWorker("presence", workers.Worker{Duration: "0 */12 * * *", WorkerFunc: UpdatePresence})
	// Update slash commands, if not bypassed
	if os.Getenv("BYPASS_SLASH_REG") != "true" {
		core.RegisterSlashCommands()
	}
	// Add all registered workers
	if core.WorkerManager.IsRunning != true {
		core.WorkerManager.AddWorkers()
		core.WorkerManager.Start()
	}
}

func UpdatePresence() {
	shard := core.Session.ShardID
	err := core.Session.UpdateStatusComplex(discordgo.UpdateStatusData{
		Status: string(discordgo.StatusDoNotDisturb),
		Activities: []*discordgo.Activity{
			{
				//nolint:gosec // we don't need cryptographically good random numbers
				Name: fmt.Sprintf("MMBN (%s | %d)", core.VERSION, shard),
				Type: discordgo.ActivityTypeStreaming,
				URL:  "https://twitch.tv/ubergeek77",
			},
		},
		AFK: false,
	})
	if err != nil {
		core.Log.Error("unable to update presence")
		core.Log.Error(err.Error())
	}
}

func init() {
	core.AddHandler(readyEventHandler)
}
