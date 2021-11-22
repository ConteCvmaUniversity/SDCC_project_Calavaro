package utility

const (
	MAXCONNECTION int    = 4 //MAXPEERS + 1 sequencer
	MAXPEERS      int    = 3
	Server_port   int    = 4321
	Server_addr   string = "10.10.1.50"
	// Server_addr string = "localhost" //if running outside docker
	Seq_addr          string = "10.10.1.51"
	Client_port       int    = 2345
	Server_cl_file    string = "/tmp/clients.txt"
	Peer_msg_seq_file string = "/tmp/messageSeq.txt"
	Peer_msg_sca_file string = "/tmp/messageSca.txt"
	Peer_msg_vec_file string = "/tmp/messageVec.txt"
	Launch_Test       bool   = false //launch all peer in test mode
	Clean_Test_Dir    bool   = true
	Test_Dir          string = "/go/src/app/peer/files/"
)
