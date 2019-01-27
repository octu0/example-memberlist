package main

import(
  "log"
  "fmt"
  "os"
  "os/signal"
  "syscall"
  "context"
  "github.com/hashicorp/memberlist"
  "github.com/serialx/hashring"
)

type MyEventDelegate struct {
  consistent *hashring.HashRing
}
func (d *MyEventDelegate) NotifyJoin(node *memberlist.Node) {
  hostPort := fmt.Sprintf("%s:%d", node.Addr.To4().String(), node.Port)
  log.Printf("join %s", hostPort)
  if d.consistent == nil {
    d.consistent = hashring.New([]string{hostPort})
  } else {
    d.consistent = d.consistent.AddNode(hostPort)
  }
}
func (d *MyEventDelegate) NotifyLeave(node *memberlist.Node) {
  hostPort := fmt.Sprintf("%s:%d", node.Addr.To4().String(), node.Port)
  log.Printf("leave %s", hostPort)
  if d.consistent != nil {
    d.consistent = d.consistent.RemoveNode(hostPort)
  }
}
func (d *MyEventDelegate) NotifyUpdate(node *memberlist.Node) {
  // skip
}

func wait_signal(cancel context.CancelFunc){
  signal_chan := make(chan os.Signal, 1)
  signal.Notify(signal_chan, syscall.SIGINT)
  for {
    select {
    case s := <-signal_chan:
      log.Printf("signal %s happen", s.String())
      cancel()
    }
  }
}
