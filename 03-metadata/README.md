# metadata

## Usage:

`node1` is implemented to refer to the metadata of `memberlist.Node` at regular intervals.

```
$ go run c.go node1.go
// => node1 region: ap-northeast-1, zone: 1a, shard: 100, weight: 0
// :
// :
```

`node2` is implemented to update confidence metadata while joining the cluster.

```
$ go run c.go node2.go --join 192.168.0.25:7947
// => node meta update successful
// :
// :
```

You can see that the information on `node1`, `node2` changes over time.

```
------------------
node1 region: ap-northeast-1, zone: 1a, shard: 100, weight: 0
node2 region: ap-northeast-1, zone: 1c, shard: 100, weight: 0
------------------
node1 region: ap-northeast-1, zone: 1a, shard: 100, weight: 0
node2 region: ap-northeast-1, zone: 1c, shard: 100, weight: 3
------------------
node1 region: ap-northeast-1, zone: 1a, shard: 100, weight: 0
node2 region: ap-northeast-1, zone: 1c, shard: 100, weight: 6
------------------
node1 region: ap-northeast-1, zone: 1a, shard: 100, weight: 0
node2 region: ap-northeast-1, zone: 1c, shard: 100, weight: 9
------------------
node1 region: ap-northeast-1, zone: 1a, shard: 100, weight: 0
node2 region: ap-northeast-1, zone: 1c, shard: 100, weight: 12
------------------
```

## Note:

`meta` is called from `NodeMeta` implemented in [Delegate](https://godoc.org/github.com/hashicorp/memberlist#Delegate).  
In order to implement it with data of `[]byte`, it can have arbitrary value,  
In this implementation [encoding/json](https://golang.org/pkg/encoding/json/) is used, but you can also implement it with [encoding/gob](https://golang.org/pkg/encoding/gob/).

```
type MyDelegate struct {
  meta  MyMetaData
}
func (d *MyDelegate) NodeMeta(limit int) []byte {
  return d.meta.Bytes()
}
func (d *MyDelegate) LocalState(join bool) []byte {
  :
}
func (d *MyDelegate) NotifyMsg(msg []byte) {
  :
}
func (d *MyDelegate) GetBroadcasts(overhead, limit int) [][]byte {
  :
}
func (d *MyDelegate) MergeRemoteState(buf []byte, join bool) {
  :
}

type MyMetaData struct {
  Region   string   `json:"region"`
  Zone     string   `json:"zone"`
  ShardId  uint16   `json:"shard-id"`
  Weight   uint64   `json:"weight"`
}
func (m MyMetaData) Bytes() []byte {
  data, err := json.Marshal(m)
  if err != nil {
    return []byte("")
  }
  return data
}
```

see also - https://github.com/hashicorp/memberlist/blob/master/delegate.go

## NextStep:

- Redundant design
- High availability design
- Load Balancing
