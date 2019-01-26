package main

import(
  "log"
  "os"
  "os/signal"
  "syscall"
  "gopkg.in/urfave/cli.v1"
  "github.com/hashicorp/memberlist"
)

func action(c *cli.Context) {
  conf := memberlist.DefaultLocalConfig()
  conf.Name = "node1"

  list, err := memberlist.Create(conf)
  if err != nil {
    log.Fatal(err)
  }

  local := list.LocalNode()
  log.Printf("node1 at %s:%d", local.Addr.To4().String(), local.Port)

  log.Printf("wait for other member connections")
  wait_signal()
}

func wait_signal(){
  signal_chan := make(chan os.Signal, 2)
  signal.Notify(signal_chan, syscall.SIGTERM)
  signal.Notify(signal_chan, syscall.SIGINT)
  for {
    select {
    case s := <-signal_chan:
      log.Printf("signal %s happen", s.String())
      return
    }
  }
}

func main(){
  app         := cli.NewApp()
  app.Action   = action
  app.Run(os.Args)
}
