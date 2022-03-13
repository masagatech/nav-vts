package shared

import (
	"net"
	"sync"
	"time"
)

var Locker sync.Mutex
var Clients = make(map[string]ClientsMod)

type ClientsMod struct {
	Con      net.Conn
	Imei     string
	Lstm     time.Time
	Allwspd  int
	VtsID    int
	Speed    int
	Loc      []float64
	Acc      int
	AC       int
	PClients []string
}
