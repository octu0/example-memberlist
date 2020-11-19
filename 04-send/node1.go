package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/memberlist"
	"gopkg.in/urfave/cli.v1"
)

func action(c *cli.Context) {
	msgCh := make(chan []byte)

	d := new(MyDelegate)
	d.msgCh = msgCh

	conf := memberlist.DefaultLocalConfig()
	conf.Name = "node1"
	conf.BindPort = 7947 // avoid port confliction
	conf.AdvertisePort = conf.BindPort
	conf.Delegate = d

	list, err := memberlist.Create(conf)
	if err != nil {
		log.Fatal(err)
	}

	local := list.LocalNode()
	list.Join([]string{
		fmt.Sprintf("%s:%d", local.Addr.To4().String(), local.Port),
	})

	stopCtx, cancel := context.WithCancel(context.TODO())
	go wait_signal(cancel)

	run := true
	for run {
		select {
		case data := <-d.msgCh:
			msg, ok := ParseMyMessage(data)
			if ok != true {
				continue
			}

			log.Printf("received msg: key=%s value=%d", msg.Key, msg.Value)

			if msg.Key == "ping" {
				m := new(MyMessage)
				m.Key = "pong"
				m.Value = msg.Value + 1

				// pong to all
				for _, node := range list.Members() {
					if node.Name == conf.Name {
						continue // skip self
					}
					log.Printf("send to %s msg: key=%s value=%d", node.Name, m.Key, m.Value)
					list.SendReliable(node, m.Bytes())
				}
			}

		case <-stopCtx.Done():
			log.Printf("stop called")
			run = false
		}
	}
	log.Printf("bye.")
}

func main() {
	app := cli.NewApp()
	app.Action = action
	app.Run(os.Args)
}
