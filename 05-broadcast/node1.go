package main

import(
  "log"
  "fmt"
  "context"
  "os"
  "gopkg.in/urfave/cli.v1"
  "github.com/hashicorp/memberlist"
)

func action(c *cli.Context) {
  msgCh := make(chan []byte)

  d := new(MyDelegate)
  d.msgCh = msgCh
  d.broadcasts = new(memberlist.TransmitLimitedQueue)

  conf := memberlist.DefaultLocalConfig()
  conf.Name          = "node1"
  conf.BindPort      = 7947 // avoid port confliction
  conf.AdvertisePort = conf.BindPort
  conf.Delegate      = d

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

  run  := true
  for run {
    select {
    case data := <-d.msgCh:
      msg, ok := ParseMyBroadcastMessage(data)
      if ok != true {
        continue
      }
      log.Printf("received broadcast msg: key=%s value=%d", msg.Key, msg.Value)

    case <-stopCtx.Done():
      log.Printf("stop called")
      run = false
    }
  }
  log.Printf("bye.")
}

func main(){
  app         := cli.NewApp()
  app.Action   = action
  app.Run(os.Args)
}
