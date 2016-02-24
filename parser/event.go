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
	Podname       string
	Role          string
	State         string
	Eventname     string
	OldMasterIP   string
	OldMasterPort int
	NewMasterIP   string
	NewMasterPort int
}

// NotificationEvent contains the various options passe din for the different
// notification events. These notification script options are event type specific.
// For details see the switch sequence in the processor
type NotificationEvent struct {
	Podname       string
	Role          string
	Eventname     string
	IP            string
	Leader        string
	Port          int
	Extra         string
	Quorum        int
	Votes         int
	OldMasterIP   string
	OldMasterPort int
	NewMasterIP   string
	NewMasterPort int
	Epoch         int
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
	case "-sdown":
		//master pod1 127.0.0.1 6502
		ne.Role = args[0]
		ne.Podname = args[1]
		ne.IP = args[2]
		ne.Port, err = strconv.Atoi(args[3])
	case "+failover-state-wait-promotion":
		//slave 127.0.0.1:6379 127.0.0.1 6379 @ tp1 127.0.0.1 7000
		ne.Role = args[0]
		ne.IP = args[2]
		ne.Port, err = strconv.Atoi(args[3])
		ne.Podname = args[5]
	case "+sdown", "-odown", "+failover-state-select-slave":
		//+sdown slave 127.0.0.1:6379 127.0.0.1 6379 @ tp1 127.0.0.1 7000
		ne.Role = args[0]
		ne.Podname = args[1]
		ne.IP = args[2]
		ne.Port, err = strconv.Atoi(args[3])
		ne.Extra = strings.Join(args[4:], " ")
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
	default:
		return ne, fmt.Errorf("Don't know how to handle event '%s' yet. Its args are %v", event, os.Args[1:])
	}
	return
}
