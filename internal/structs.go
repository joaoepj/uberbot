package internal

import "time"

type Mute struct {
	ID         string
	ExpireDate time.Time
	Length     string
	UseRole    bool
}
