package main

import (
	"log"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/therealbill/libredis/client"
)

var (
	pods []string
	app  *cli.App
	rc   *client.Redis
	path string
)

func main() {
	app := cli.NewApp()
	app.Name = "sentinel-scriptify"
	app.Usage = "Comb through a Sentinel and add eventilator and/or reconfigurator to all pods it manages"
	app.EnableBashCompletion = true
	author := cli.Author{Name: "Bill Anderson", Email: "therealbill@me.com"}
	app.Authors = append(app.Authors, author)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "sentineladdress, s",
			Value: "127.0.0.1",
			Usage: "Address of the sentinel",
		},
		cli.IntFlag{
			Name:  "sentinelport, p",
			Value: 26379,
			Usage: "Port of the sentinel",
		},
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "Be verbose in logging output",
		},
		cli.BoolFlag{
			Name:  "both, b",
			Usage: "Add both eventilator and registrator",
		},
		cli.BoolFlag{
			Name:  "eventilator, e",
			Usage: "Add eventilator",
		},
		cli.BoolFlag{
			Name:  "reconfigurator, r",
			Usage: "Add reconfigurator",
		},
		cli.StringFlag{
			Name:  "commandpath, c",
			Value: "/usr/sbin",
			Usage: "Directory where eventilator/reconfigurator lives",
		},
	}
	app.Action = updatePods
	app.Before = getPodList
	app.Run(os.Args)
}

func eventilate(pod string) error {
	return rc.SentinelSetString(pod, "notification-script", path+"/eventilator")
}

func reconfigurate(pod string) error {
	return rc.SentinelSetString(pod, "client-reconfig-script", path+"/reconfigurator")
}

func updatePods(c *cli.Context) {
	path = c.String("commandpath")
	verbose := c.Bool("verbose")
	if verbose {
		log.Printf("Updating %d pods", len(pods))
	}
	both := c.Bool("both")
	eventilator := false
	reconfigurator := false
	if both {
		eventilator = true
		reconfigurator = true
	} else {
		eventilator = c.Bool("eventilator")
		reconfigurator = c.Bool("reconfigurator")
	}
	for _, pod := range pods {
		if eventilator {
			if verbose {
				log.Printf("Eventilating %s", pod)
			}
			err := eventilate(pod)
			if err != nil {
				if strings.Contains(err.Error(), "seems non existing or non executable") {
					log.Fatalf("eventilator is not found in %s or is not executable, aborting entire run now.", path)
				}
				log.Printf("Error while eventilating '%s': %v", pod, err)
			}
		}
		if reconfigurator {
			if verbose {
				log.Printf("Reconfigurating %s", pod)
			}
			err := reconfigurate(pod)
			if err != nil {
				if strings.Contains(err.Error(), "seems non existing or non executable") {
					log.Fatalf("reconfigurator is not found in %s or is not executable, aborting entire run now.", path)
				}
				log.Printf("Error while reconfigurating '%s': %v", pod, err)
			}
		}
	}
}

func getPodList(c *cli.Context) (err error) {
	saddr := c.String("sentineladdress")
	sport := c.Int("sentinelport")

	rc, err = client.Dial(saddr, sport)
	if err != nil {
		log.Fatalf("Unable to connect to sentinel at %s:%d. Err: %v", saddr, sport, err)
		return err
	}
	pods, err = rc.Role()
	return err
}
