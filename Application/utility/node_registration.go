package utility

import (
	"container/list"
	"errors"
	"log"
	"net/rpc"
	"strconv"
)

/*
	Set info
*/
func setInfo(info *Info, nodeType NodeType, port int) error {
	info.Type = nodeType
	info.Address = GetLocalIP()
	if info.Address == "" {
		return errors.New("Impossible to find local ip")
	}

	info.Port = strconv.Itoa(port)
	return nil
}

/*
	Registration function for peer
*/
func Registration(peers *list.List, nodeType NodeType, port int) {

	var info Info
	var res Result_file

	addr := Server_addr + ":" + strconv.Itoa(Server_port)
	// Try to connect to addr
	server, err := rpc.Dial("tcp", addr)
	if err != nil {
		log.Fatal("Error in dialing: ", err)
	}
	defer server.Close()

	//set info to send
	err = setInfo(&info, nodeType, port)
	if err != nil {
		log.Fatal("Error on setInfo: ", err)
	}

	//call procedure
	log.Printf("Call to registration node")
	err = server.Call("Utility.Save_registration", &info, &res)
	if err != nil {
		log.Fatal("Error save_registration procedure: ", err)
	}

	//check result
	for e := 0; e < res.PeerNum; e++ {
		var item Info
		var tmp string
		tmp, item.Address, item.Port = ParseLine(res.Peers[e], ":")
		item.Type = StringToType(tmp)
		item.ID = e
		peers.PushBack(item)

	}

}
