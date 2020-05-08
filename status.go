package theilliminationgame

import (
	"github.com/maisiesadler/theilliminationgame/models"
)

func (g *Game) isRunning() bool {
	return g.db.State == models.StateRunning
}

func (g *Game) isFinished() bool {
	return g.db.State == models.StateFinished || g.db.State == models.StateArchived
}
