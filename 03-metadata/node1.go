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
    Zone:   "1a",
    ShardId: 100,
    Weight:  0,
  }

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

  tick := time.NewTicker(3 * time.Second)
  run  := true
  for run {
    select {
    case <-tick.C:
      for _, node := range list.Members() {
        meta, ok := ParseMyMetaData(node.Meta)
        if ok != true {
          continue
        }

        log.Printf(
          "%s region: %s, zone: %s, shard: %d, weight: %d",
          node.Name,
          meta.Region,
          meta.Zone,
          meta.ShardId,
          meta.Weight,
        )
      }
      log.Printf("------------------")

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
  app.Run(os.Args)
}
