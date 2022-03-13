package handlers

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/masagatech/nav-vts/app/models"
	"github.com/masagatech/nav-vts/app/queue"
	"github.com/masagatech/nav-vts/app/servers"
	"github.com/masagatech/nav-vts/app/shared"
	"github.com/masagatech/nav-vts/app/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type Handlers struct {
	ampq *queue.RMQ
}

//Terminal connection persists
func (h *Handlers) addClient(conn net.Conn, client string, allwspd int, vtsid int,
	PushClients []string, LLoc []float64, LAC int, LACC int, LTime time.Time) shared.ClientsMod {
	shared.Locker.Lock()
	defer shared.Locker.Unlock()
	ipaddress := conn.RemoteAddr().String()
	fmt.Println("Client Added  ", ipaddress)
	shared.Clients[ipaddress] = shared.ClientsMod{Con: conn, Imei: client, Lstm: LTime,
		Allwspd: allwspd, VtsID: vtsid, PClients: PushClients, Loc: LLoc, AC: LAC, Acc: LACC}
	return shared.Clients[ipaddress]
}

func (h *Handlers) RemoveClient(conn net.Conn) {
	shared.Locker.Lock()
	defer shared.Locker.Unlock()
	ipaddress := conn.RemoteAddr().String()
	delete(shared.Clients, ipaddress)
}

func (h *Handlers) getClient(conn net.Conn) (client shared.ClientsMod) {
	ipaddress := conn.RemoteAddr().String()
	return shared.Clients[ipaddress]
}

func (h *Handlers) setClient(conn net.Conn, client shared.ClientsMod) {
	shared.Locker.Lock()
	defer shared.Locker.Unlock()

	ipaddress := conn.RemoteAddr().String()
	shared.Clients[ipaddress] = client
}

func NewInit(tcpServer *servers.Server, ampq *queue.RMQ) {

	h := &Handlers{
		ampq: ampq,
	}

	tcpServer.OnNewClient(func(c *servers.Client) {
		fmt.Println("On new client connect")
		fmt.Println(c.Conn().RemoteAddr().String())

	})

	tcpServer.OnNewMessage(func(c *servers.Client, _data []byte) {
		if !(_data[0] == 0x78 || _data[0] == 0x79) {

			// if _data[0] == 0x5b && _data[len(_data)-1] == 0x5d {
			// 	fmt.Printf("%02x\n", bytes.Trim(_data, "\x00"), " watch detected")
			// } else {
			fmt.Println("invalid command", _data[len(_data)-1])
			// }

		} else if _data[0] == 0x78 { //check for 78 commands
			if _data[3] == 0x01 {
				//fmt.Println("login request")
				h.registerDevice(_data, c.Conn()) // Registration Response
			} else if _data[3] == 0x13 {
				h.heartBeat(_data, c.Conn()) // HeartBeat Response
			} else if (_data[3] == 0x22) || (_data[3] == 0x12) {
				h.locationDt(_data, c.Conn()) //location data
			} else if _data[3] == 0x15 {
				//	commandReply(_data, lendata, connection) //location data
			}
		}

	})
	tcpServer.OnClientConnectionClosed(func(c *servers.Client, err error) {
		fmt.Println(err)

	})

}

func (h *Handlers) registerDevice(_data []byte, connection net.Conn) {
	reply := []byte{0x78, 0x78, 0x05, 0x01}          //assign reply variable
	serial := _data[12:14]                           //get crc from data
	_crxCRC := append([]byte{0x05, 0x01}, serial...) // create crc string
	_crxCRCF := utils.GetCrc16(_crxCRC)              // get computed crc in variable
	_crxCRCF = append(serial, _crxCRCF...)           // append final crc and reply data
	reply = append(reply, _crxCRCF...)               // append final crc and reply data
	reply = append(reply, 0x0D, 0x0A)                //EOF
	//get imei number
	_imei := fmt.Sprintf("%x", _data[4:12])[1:16] //getting imei number
	fmt.Println(_imei)

	connection.Write(reply)
	h.addClient(connection, _imei, 0, 0, nil, nil, 1, 1, time.Now())
	fmt.Println("send to queue")

	h.ampq.PublishOnQueue("", "test", "device registered")

	//h.ampq.PublishOnQueue()
}

//getting heart beat
func (h *Handlers) heartBeat(_data []byte, connection net.Conn) {
	_clnt := h.getClient(connection)
	if _clnt.Imei == "" {
		return
	}

	reply := []byte{0x78, 0x78, 0x05, 0x13} //assign reply variable
	serial := _data[7:9]                    //get crc from data

	_crxCRC := append([]byte{0x05, 0x13}, serial...) // create crc string
	_crxCRCF := utils.GetCrc16(_crxCRC)              // get computed crc in variable
	_crxCRCF = append(serial, _crxCRCF...)           // append final crc and reply data
	reply = append(reply, _crxCRCF...)               // append final crc and reply data
	reply = append(reply, 0x0D, 0x0A)                //EOF
	//Client get by ipaddress

	//extract data from received data
	//fmt.Println(_data[4:5])
	_prd := fmt.Sprintf("%08b", _data[4:5])
	_prd = _prd[1 : len(_prd)-1]
	btrt := "BTRY"
	fmt.Println(_prd)
	data := models.HertBt{
		Acttm:  time.Now(),
		Actvt:  "hrtbt",
		Sertm:  time.Now(),
		Speed:  0,
		Imei:   _clnt.Imei,
		Flag:   "inprog",
		Appvr:  "1.0",
		Vhid:   _clnt.Imei,
		Btr:    utils.Batryper(int(_data[5])),
		Btrst:  btrt,
		Alm:    (_prd[2:3] + _prd[3:4] + _prd[4:5]), //100: SOS,011: Low Battery Alarm,010: Power Cut Alarm,001: Shock Alarm,000: Normal
		Gsmsig: utils.Networkper(int(_data[6])),     //0x00: no signal,0x01: extremely weak signal,0x02: very weak signal,0x03: good signal,0x04: strong signal
	}

	data.Oe, _ = strconv.Atoi(_prd[0:1])   //1: oil and electricity disconnected, 0: gas oil and electricity
	data.Gp, _ = strconv.Atoi(_prd[1:2])   //1: GPS tracking is on,0: GPS tracking is off
	data.Chrg, _ = strconv.Atoi(_prd[5:6]) //1: Charge On,0: Charge Off
	data.Acc, _ = strconv.Atoi(_prd[6:7])  //1: ACC high,0: ACC Low
	data.Df, _ = strconv.Atoi(_prd[7:8])   //1: Defense Activated,0: Defense Deactivated
	//
	// fmt.Println(data.Acc)

	_clnt.Lstm = data.Sertm
	_clnt.Acc = data.Acc
	otherdata := bson.M{
		"actvt":    "loc",
		"sertm":    time.Now(),
		"imei":     _clnt.Imei,
		"alwspeed": _clnt.Allwspd,
		"isp":      false,
		"flag":     "acc",
		"acc":      data.Acc,
		"appvr":    "1.0",
		"loc":      _clnt.Loc,
		"bearing":  0,
		"speed":    0,
		"vhid":     _clnt.Imei,
	}
	// if _clnt.Acc != data.Acc {

	// 	// go fcm.SendACCAlertTotopic(_clnt.Imei, data.Acc)
	// }

	// if data.Chrg == 1 {
	// 	data.Btrst = "CHRG"
	// }

	fmt.Println(otherdata["sertm"])

	connection.Write(reply)

	//ampq.PublishOnQueue("test", string(reply))
	h.ampq.PublishOnQueue("", "test", otherdata)

}

func (h *Handlers) locationDt(_data []byte, connection net.Conn) {

	// _clnt := getClient(connection)
	// if _clnt.Imei == "" {
	// 	return
	// }
	// _dt := "20" + fmt.Sprintf("%d-%d-%d %d:%d:%d", _data[4], _data[5], _data[6], _data[7], _data[8], _data[9]) //conver to Date
	// //fmt.Println(_dt)
	// crs := fmt.Sprintf("%x", _data[10:11])            //Quantity of GPS	information	satellites
	// _stlt, _ := strconv.ParseInt("0x0"+crs[1:], 0, 8) //satlites
	// //extract data from received data
	// var _lat float64
	// var _lon float64

	// _lat = float64(binary.BigEndian.Uint32(_data[11:15])) / (30000 * 60) //Lattitude
	// _lon = float64(binary.BigEndian.Uint32(_data[15:19])) / (30000 * 60) //Longitude

	// _courus := fmt.Sprintf("%016b", binary.BigEndian.Uint16(_data[20:22]))
	// _bearing, _ := strconv.ParseInt(_courus[6:], 2, 64) // get bearing
	// point := []float64{utils.ToFixed(_lon, 6), utils.ToFixed(_lat, 6)}
	//fmt.Println(_clnt.allwspd)

	// data := bson.M{
	// 	"gpstm":    _dt,
	// 	"actvt":    "loc",
	// 	"sertm":    time.Now(),
	// 	"imei":     _clnt.Imei,
	// 	"alwspeed": _clnt.Allwspd,
	// 	"isp":      false,
	// 	"flag":     "inprog",
	// 	"appvr":    "1.0",
	// 	"acc":      _clnt.Acc,
	// 	"sat":      _stlt,
	// 	"loc":      point,
	// 	"postyp":   _courus[2:3],
	// 	"bearing":  _bearing,
	// 	"speed":    _data[19],
	// 	"vhid":     _clnt.Imei}

	// _clnt.Lstm = time.Now()
	// _clnt.Loc = point

	// if _clnt.Allwspd > 0 {

	// 	crspeed := int(_data[19])

	// 	//fmt.Println(crspeed)
	// 	if crspeed > _clnt.Allwspd {
	// 		// speed voilence
	// 		//fmt.Println(int(_data[19]))
	// 		go fcm.SendSpeedAlertTotopic(_clnt.Imei, crspeed)
	// 		data["lstspd"] = crspeed
	// 		data["lstspdtm"] = time.Now()
	// 		data["isp"] = true
	// 	}
	// }

	// setClient(connection, _clnt)
	// //need to call mongo db
	// models.UpdateData(data, _clnt.Imei, "loc", nil)
	// checkGeofence(point, _data[19], _clnt.Imei)

	// //send tp 3rd party client
	// datashare.SendToClient(_clnt)
}
