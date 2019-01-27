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
  d := new(MyDelegate)
  d.meta = MyMetaData{
    Region: "ap-northeast-1",
    Zone:   "1c",
    ShardId: 100,
    Weight:  0,
  }

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

  tick := time.NewTicker(1 * time.Second)
  run  := true
  for run {
    select {
    case <-tick.C:
      d.meta.Weight = d.meta.Weight + 1

      if err := list.UpdateNode(1 * time.Second); err != nil {
        log.Printf("node meta update failure")
      } else {
        log.Printf("node meta update successful")
      }

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
