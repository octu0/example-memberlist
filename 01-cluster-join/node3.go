package main

import(
  "log"
  "os"
  "time"
  "gopkg.in/urfave/cli.v1"
  "github.com/hashicorp/memberlist"
)

func action(c *cli.Context) {
  conf := memberlist.DefaultLocalConfig()
  conf.Name = "node3"
  conf.BindPort = 7948 // avoid port confliction
  conf.AdvertisePort = conf.BindPort

  list, err := memberlist.Create(conf)
  if err != nil {
    log.Fatal(err)
  }

  local := list.LocalNode()
  log.Printf("node3 at %s:%d", local.Addr.To4().String(), local.Port)

  join := c.String("join")
  log.Printf("cluster join to %s", join)

  if _, err := list.Join([]string{join}); err != nil {
    log.Fatal(err)
  }

  for _, member := range list.Members() {
    log.Printf("Member: %s(%s:%d)", member.Name, member.Addr.To4().String(), member.Port)
  }

  log.Printf("wait 10 sec")
  time.Sleep(10 * time.Second)

  log.Printf("leave cluster")
  list.Leave(5 * time.Second)
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
