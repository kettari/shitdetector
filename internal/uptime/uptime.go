package uptime

import (
	"github.com/hako/durafmt"
	"time"
)

type (
	Uptime struct {
		ID    int64
		Since int64
	}
	Service interface {
		Update() error
		Since() (string, error)
	}
)

func (u Uptime) ToWording() (wording string) {
	t := time.Unix(u.Since, 0)
	return durafmt.Parse(time.Since(t).Round(time.Second)).String()
}
