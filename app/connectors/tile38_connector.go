package connectors

import (
	"strconv"

	"github.com/masagatech/nav-vts/app/models"
	"github.com/xjem/t38c"
)

func NewTile38(config *models.Config) *t38c.Client {
	client, err := t38c.New(config.Tile38.Host+":"+strconv.Itoa(config.Tile38.Port), t38c.Debug)
	if err != nil {
		panic(err)
	}
	return client
}
