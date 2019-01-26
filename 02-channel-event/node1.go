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
  conf := memberlist.DefaultLocalConfig()
  conf.Name          = "node1"
  conf.BindPort      = 7947 // avoid port confliction
  conf.AdvertisePort = conf.BindPort
  conf.Events        = new(MyEventDelegate)

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
      devt := conf.Events.(*MyEventDelegate)
      if devt == nil {
        log.Printf("consistent isnt initialized")
        continue
      }
      log.Printf("current node size: %d", devt.consistent.Size())

      keys := []string{"hello", "world"}
      for _, key := range keys {
        node, ok := devt.consistent.GetNode(key)
        if ok == true {
          log.Printf("node1 search %s => %s", key, node)
        } else {
          log.Printf("no node available")
        }
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
  app.Run(os.Args)
}
