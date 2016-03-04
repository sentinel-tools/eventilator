package parser

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ReconfigurationEvent contains the leader's role, the event state, new and
// old master IP information
// client-reconfig-script options:
// <master-name> <role> <state> <from-ip> <from-port> <to-ip> <to-port>
type ReconfigurationEvent struct {
	Podname       string `json:"podname"`
	Role          string `json:"-"`
	State         string `json:"state"`
	Eventname     string `json:"event"`
	OldMasterIP   string `json:"old_master_ip"`
	OldMasterPort int    `json:"old_master_port"`
	NewMasterIP   string `json:"new_master_ip"`
	NewMasterPort int    `json:"new_master_port"`
}

// NotificationEvent contains the various options passe din for the different
// notification events. These notification script options are event type specific.
// For details see the switch sequence in the processor
type NotificationEvent struct {
	Podname       string `json:"podname,omitempty"`
	Role          string `json:"role,omitempty"`
	Eventname     string `json:"event,omitempty"`
	IP            string `json:"ip,omitempty"`
	Leader        string `json:"leader,omitempty"`
	Port          int    `json:"port,omitempty"`
	Extra         string `json:"extra,omitempty"`
	Quorum        int    `json:"quorum,omitempty"`
	Votes         int    `json:"votes,omitempty"`
	OldMasterIP   string `json:"old_master_ip,omitempty"`
	OldMasterPort int    `json:"old_master_port,omitempty"`
	NewMasterIP   string `json:"new_master_ip,omitempty"`
	NewMasterPort int    `json:"new_master_port,omitempty"`
	Epoch         int    `json:"epoch,omitempty"`
}

func ParseReconfiguration(args []string) (re ReconfigurationEvent, err error) {
	re.Podname = args[0]
	re.Role = args[1]
	re.State = args[2]
	re.OldMasterIP = args[3]
	re.OldMasterPort, err = strconv.Atoi(args[4])
	re.NewMasterIP = args[5]
	re.NewMasterPort, err = strconv.Atoi(args[6])
	return
}

func ParseNotification(event string, args []string) (ne NotificationEvent, err error) {
	ne.Eventname = event
	switch event {
	case "+vote-for-leader":
		//runid epoch
		ne.Leader = args[0]
		ne.Epoch, err = strconv.Atoi(args[1])
	case "+new-epoch":
		//16
		ne.Epoch, err = strconv.Atoi(args[0])
	case "+odown":
		ne.Role = args[0]
		ne.Podname = args[1]
		ne.IP = args[2]
		ne.Port, err = strconv.Atoi(args[3])
		qdata := strings.Split(args[5], "/")
		ne.Quorum, err = strconv.Atoi(qdata[1])
		ne.Votes, err = strconv.Atoi(qdata[0])
	case "+slave":
		//slave 127.0.0.1:7000 127.0.0.1 7000 @ tp1 127.0.0.1 6379
		ne.Role = args[0]
		ne.IP = args[2]
		ne.Port, err = strconv.Atoi(args[3])
		ne.Podname = args[5]
	case "+failover-state-wait-promotion":
		//slave 127.0.0.1:6379 127.0.0.1 6379 @ tp1 127.0.0.1 7000
		ne.Role = args[0]
		ne.IP = args[2]
		ne.Port, err = strconv.Atoi(args[3])
		ne.Podname = args[5]
	case "-sdown", "+sdown", "-odown", "+failover-state-select-slave":
		//master PODID MASTERIP MASTERORT
		//slave SLAVENAME SLAVEIP SLAVEPORT @ PODNAME MASTERIP MASTERPORT
		//sentinel SENTINELNAME SENTINELIP SENTINELPORT @ PODNAME MASTERIP MASTERPORT
		ne.Role = args[0]
		switch ne.Role {
		case "master":
			ne.Podname = args[1]
			ne.IP = args[2]
			ne.Port, err = strconv.Atoi(args[3])
		case "slave":
			ne.Podname = args[5]
			ne.IP = args[2]
			ne.Port, err = strconv.Atoi(args[3])
			ne.Extra = strings.Join(args[5:], " ")
		case "sentinel":
			ne.Podname = args[5]
			ne.IP = args[2]
			ne.Port, err = strconv.Atoi(args[3])
			ne.Extra = strings.Join(args[5:], " ")
		}
	case "+try-failover":
		ne.Role = args[0]
		ne.Podname = args[1]
		ne.IP = args[2]
		ne.Port, err = strconv.Atoi(args[3])
	case "+selected-slave":
		//slave 127.0.0.1:7000 127.0.0.1 7000 @ tp1 127.0.0.1 6379
		ne.Role = args[0]
		ne.Podname = args[1]
		ne.IP = args[2]
		ne.Port, _ = strconv.Atoi(args[3])
		ne.Podname = args[5]
	case "+switch-master":
		//tp1 127.0.0.1 6379 127.0.0.1 7000
		ne.Podname = args[0]
		ne.OldMasterIP = args[1]
		ne.OldMasterPort, _ = strconv.Atoi(args[2])
		ne.IP = args[3]
		ne.Port, _ = strconv.Atoi(args[4])
		ne.NewMasterIP = args[3]
		ne.NewMasterPort, _ = strconv.Atoi(args[4])
	case "+promoted-slave", "+convert-to-slave":
		//slave 127.0.0.1:6379 127.0.0.1 6379 @ tp1 127.0.0.1 7000
		ne.Role = args[0]
		ne.IP = args[2]
		ne.Port, err = strconv.Atoi(args[3])
		ne.Podname = args[5]

		ne.OldMasterIP = args[6]
		ne.OldMasterPort, _ = strconv.Atoi(args[7])
		ne.NewMasterIP = ne.IP
		ne.NewMasterPort = ne.Port
	case "+failover-end", "+failover-state-reconf-slaves", "+elected-leader", "-failover-abort-no-good-slave", "-failover-abort-not-elected":
		//master tp1 127.0.0.1 6379
		ne.Role = args[0]
		ne.Podname = args[1]
		ne.OldMasterIP = args[2]
		ne.OldMasterPort, _ = strconv.Atoi(args[3])
	case "+monitor":
		//master tp1 127.0.0.1 7000 quorum 1
		ne.Role = args[0]
		ne.Podname = args[1]
		ne.IP = args[2]
		ne.Port, _ = strconv.Atoi(args[3])
	case "+set":
		// master PODNAME MASTERIP MASTERPORT DIRECTIVE VALUE
		ne.Role = args[0]
		ne.Podname = args[1]
		ne.IP = args[2]
		ne.Port, _ = strconv.Atoi(args[3])
		ne.Extra = fmt.Sprintf("%s=%s", args[4], args[5])
	default:
		return ne, fmt.Errorf("Don't know how to handle event '%s' yet. Its args are %v", event, os.Args[1:])
	}
	return
}
