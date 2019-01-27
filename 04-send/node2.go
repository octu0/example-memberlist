package main

import(
  "log"
  "fmt"
  "time"
  "context"
  "os"
  "gopkg.in/urfave/cli.v1"
  "github.com/hashicorp/memberlist"
)

func action(c *cli.Context) {
  msgCh := make(chan []byte)

  d := new(MyDelegate)
  d.msgCh = msgCh

  conf := memberlist.DefaultLocalConfig()
  conf.Name          = "node2"
  conf.BindPort      = 7948 // avoid port confliction
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
  join := c.String("join")
  log.Printf("cluster join to %s", join)

  if _, err := list.Join([]string{join}); err != nil {
    log.Fatal(err)
  }

  stopCtx, cancel := context.WithCancel(context.TODO())
  go wait_signal(cancel)

  i    := uint64(0)
  tick := time.NewTicker(3 * time.Second)
  run  := true
  for run {
    select {
    case <-tick.C:
      m := new(MyMessage)
      m.Key = "ping"
      m.Value = i

      // ping to all
      for _, node := range list.Members() {
        if node.Name == conf.Name {
          continue // skip self
        }
        log.Printf("send to %s msg: key=%s value=%d", node.Name, m.Key, m.Value)
        list.SendReliable(node, m.Bytes())
      }

    case data := <-d.msgCh:
      msg, ok := ParseMyMessage(data)
      if ok != true {
        continue
      }

      log.Printf("received msg: key=%s value=%d", msg.Key, msg.Value)
      i = msg.Value + 1

    case <-stopCtx.Done():
      log.Printf("stop called")
      run = false
    }
  }
  tick.Stop()
  log.Printf("bye.")
}

func main(){
  app         := cli.NewApp()
  app.Action   = action
  app.Flags    = []cli.Flag{
    cli.StringFlag{
      Name: "join, j",
      Usage: "cluster join address",
      Value: "127.0.0.1:xxxx",
      EnvVar: "JOIN_ADDR",
    },
  }
  app.Run(os.Args)
}
