// +build !custom

package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/sentinel-tools/eventilator/config"
	"github.com/sentinel-tools/eventilator/parser"
)

func PostNotificationEventToSensuJIT(config config.SensuJITConfig, event parser.NotificationEvent) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", config.IP, config.Port))
	if err != nil {
		return err
	}
	hostname, err := GetMyFQDN()
	if err != nil {
		return err
	}
	var se SensuMessage
	spodid := fmt.Sprintf("pod-%s", event.Podname)
	switch event.Eventname {
	case "+odown":
		se = SensuMessage{Name: "redis-master-down",
			Source:      spodid,
			Occurrences: 1,
			ZDSubject:   fmt.Sprintf("Monitoring alert for %s Master Down", event.Podname),
			Description: fmt.Sprintf("Pod %s has a non-responsive Master", event.Podname),
			Host:        hostname,
			Status:      1,
			Output:      fmt.Sprintf("Pod %s has a non-responsive Master", event.Podname),
			ZDTags:      "monitoring redis_master_odown",
		}
	case "-odown":
		se = SensuMessage{Name: "redis-master-down",
			Source:      spodid,
			Occurrences: 1,
			ZDSubject:   fmt.Sprintf("Monitoring alert for %s Master Down", event.Podname),
			Description: fmt.Sprintf("Pod %s has a non-responsive Master", event.Podname),
			Host:        hostname,
			Status:      0,
			Output:      fmt.Sprintf("Pod %s has a non-responsive Master", event.Podname),
			ZDTags:      "monitoring redis_master_odown",
		}
	case "-failover-abort-no-good-slave":
		se = SensuMessage{Name: "sentinel-failover-aborted",
			Source:      spodid,
			Occurrences: 1,
			ZDSubject:   fmt.Sprintf("Monitoring alert for %s Failover Failure", event.Podname),
			Description: fmt.Sprintf("Pod %s Is in Aborted Failover State", event.Podname),
			Host:        hostname,
			Status:      2,
			Output:      fmt.Sprintf("Pod %s has a non-responsive Master and no promotable slave", event.Podname),
			ZDTags:      "monitoring redis_pod_down",
		}
	case "+promoted-slave":
		se = SensuMessage{Name: "redis-master-down",
			Source:      spodid,
			Occurrences: 1,
			ZDSubject:   fmt.Sprintf("Monitoring alert for %s Master Down", event.Podname),
			Description: fmt.Sprintf("Pod %s has a non-responsive Master", event.Podname),
			Host:        hostname,
			Status:      0,
			Output:      fmt.Sprintf("Pod %s has a non-responsive Master", event.Podname),
			ZDTags:      "monitoring redis_master_odown",
		}
	case "+switch-master":
		se = SensuMessage{Name: "sentinel-failover-aborted",
			Source:      spodid,
			Occurrences: 1,
			ZDSubject:   fmt.Sprintf("Monitoring alert for %s Failover Failure", event.Podname),
			Description: fmt.Sprintf("Pod %s Is in Aborted Failover State", event.Podname),
			Host:        hostname,
			Status:      0,
			Output:      fmt.Sprintf("Pod %s has a non-responsive Master and no promotable slave", event.Podname),
			ZDTags:      "monitoring redis_pod_down",
		}
	default:
		log.Printf("I don't handle event %s", event.Eventname)
		return fmt.Errorf("[SensuJIT] I don't handle event %s", event.Eventname)
	}
	msg, err := json.Marshal(se)
	if err != nil {
		return err
	}
	fmt.Fprintf(conn, string(msg))
	return err
}

type SensuMessage struct {
	Name        string `json:"name"`
	Source      string `json:"source"`
	Occurrences int    `json:"occurrences"`
	ZDTags      string `json:"zd_set_tags"`
	ZDSubject   string `json:"zd_subject"`
	Description string `json:"zd_description"`
	Host        string `json:"host"`
	Status      int    `json:"status"`
	Output      string `json:"output"`
}
