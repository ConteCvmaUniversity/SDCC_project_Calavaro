package utility

import (
	"bufio"
	"errors"
	"log"
	"net"
	"os"
	"sync"
)

/*
	This package static configure connection port and ip
*/

type Utility int

// Constant value

/*
	Non optimal solution:
	MAXCONNECTION = numberOfPeer + 1 (sequencer)
*/

var (
	Connection = make(chan bool)
	Wg         = new(sync.WaitGroup)
)

type Result_file struct {
	PeerNum int
	Peers   [MAXCONNECTION]string
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// save registration info to reg_node procedure
func (utils *Utility) Save_registration(arg *Info, res *Result_file) error {
	log.Printf("The registration is for %s the ip address:port : %s:%s\n", TypeToString(arg.Type), arg.Address, arg.Port)
	f, err := os.OpenFile(Server_cl_file, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		log.Println(err)
		return errors.New("Impossible to open file")
	}
	/*
		see https://www.joeshaw.org/dont-defer-close-on-writable-files/ for file defer on close
	*/
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(f)

	//save new address on file
	_, err = f.WriteString(TypeToString(arg.Type) + ":" + arg.Address + ":" + arg.Port)
	_, err = f.WriteString("\n")
	err = f.Sync()
	if err != nil {
		return err
	}

	log.Printf("Saved")

	Connection <- true
	Wg.Add(1)
	log.Printf("Waiting other connection")
	Wg.Wait()
	//send back file
	err = prepare_response(res)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func prepare_response(res *Result_file) error {
	res.PeerNum = MAXCONNECTION
	file, err := os.Open(Server_cl_file)
	if err != nil {
		return errors.New("error on open file[prepare_file]")
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	var i int
	for scanner.Scan() {
		line := scanner.Text()
		res.Peers[i] = line
		i++
	}
	if err := scanner.Err(); err != nil {
		return errors.New("error on open file[prepare_file]")
	}
	err = file.Sync()
	if err != nil {
		return errors.New("error on open file[prepare_file]")
	}
	return nil
}
