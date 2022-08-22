package component

import "github.com/robfig/cron"

type Cron struct {
	Named
	*cron.Cron
}

func NewCron() *Cron {
	return &Cron{Named: componentCron, Cron: cron.New()}
}
