package core

type GuildProvider struct {
	Save func(guild *Guild)
	Load func() map[string]*Guild
}
