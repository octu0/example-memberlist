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

  e := new(MyEventDelegate)
  e.Num = 0

  d := new(MyDelegate)
  d.msgCh = msgCh
  d.broadcasts = new(memberlist.TransmitLimitedQueue)
  d.broadcasts.NumNodes = func() int {
    log.Printf("broadcast nodes = %d", e.Num)
    return e.Num
  }
  d.broadcasts.RetransmitMult = 1

  conf := memberlist.DefaultLocalConfig()
  conf.Name          = "node3"
  conf.BindPort      = 7949 // avoid port confliction
  conf.AdvertisePort = conf.BindPort
  conf.Events        = e
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
  e.Num = list.NumMembers()

  stopCtx, cancel := context.WithCancel(context.TODO())
  go wait_signal(cancel)

  i    := uint64(0)
  tick := time.NewTicker(3 * time.Second)
  run  := true
  for run {
    select {
    case <-tick.C:
      m := MyBroadcastMessage{
        Key: "I am node3",
        Value: i,
      }

      log.Printf("send broadcast msg: key=%s value=%d", m.Key, m.Value)
      d.broadcasts.QueueBroadcast(m)

    case data := <-d.msgCh:
      msg, ok := ParseMyBroadcastMessage(data)
      if ok != true {
        continue
      }

      log.Printf("received broadcast msg: key=%s value=%d", msg.Key, msg.Value)
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

