# cluster join/leave

## Usage:

`node1` is the node to start first, since you will see the address to join the cluster, please make a note.

```
$ go run node1.go
// => node1 at 192.168.0.25:7946
```

`node2` is a node that only joins the cluster and finishes.  
When you start up using the address of `node1` above, the member information participating in the cluster is displayed.

```
$ go run node2.go --join 192.168.0.25:7946
// => Member: node1(192.168.0.25:7946)
// => Member: node2(192.168.0.25:7947)
```

`node3` is a node that waits for a certain period of time after joining the cluster,  
by exchanging the boot order of `node2` and `node3`, you can see that participating members are changing.

```
# session1
$ go run node3.go --join 192.168.0.25:7946
// => Member: node1(192.168.0.25:7946)
// => Member: node3(192.168.0.25:7948)

# session2
$ go run node2.go --join 192.168.0.25:7946
// => Member: node1(192.168.0.25:7946)
// => Member: node2(192.168.0.25:7947)
// => Member: node3(192.168.0.25:7948)
```

## Note:

Make sure `memberlist.Config#Name` (used as the node name) is unique.  
Also, make sure `BindPort` and `AdvertisePort` are unique addresses.  
It is used to communicate between nodes.

## Configure:

this example using [DefaultLocalConfig](https://godoc.org/github.com/hashicorp/memberlist#DefaultLocalConfig) because it is intended for testing in the local environment.  
If you want to test over the LAN/WAN network please set the appropriate Timeout / Interval.

```
conf := memberlist.DefaultLocalConfig()
conf.Name = "node1"
conf.BindAddr = "192.168.0.1"
conf.BindPort = 7901
conf.AdvertiseAddr = "192.168.0.1"
conf.AdvertisePort = 7901

list, err := memberlist.Create(conf)
if err != nil {
  log.Fatal(err)
}
```

see also - https://github.com/hashicorp/memberlist/blob/master/config.go
