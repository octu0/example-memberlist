# broadcast

Broadcasting data to other nodes than self node.

## Usage:

`node1` is receiving broadcast messages of all nodes.

`node2` and `node3` are sending their messages on periodic interval.

run node1:

```
$ go run c.go node1.go
// => received broadcast msg: key=I am node2 value=0
// => received broadcast msg: key=I am node3 value=1
```

run node2:
```
$ go run c.go node2.go
// => info: broadcast nodes = 2
// => send broadcast msg: key=I am node2 value=0
// => received broadcast msg: key=I am node3 value=0
```

run node3:
```
$ go run c.go node3.go
// => info: broadcast nodes = 3
// => received broadcast msg: key=I am node2 value=2
// => send broadcast msg: key=I am node3 value=3
```

### Timeline:

```
                                               +                                                +
  node1 says                                   | node2 says                                     | node3 says
+----------------------------------------------|------------------------------------------------|-------------------------------------------------+
                                               | send broadcast msg: key=I am node2 value=0     |
received broadcast msg: key=I am node2 value=0 |                                                | received broadcast msg: key=I am node2 value=0
                                               |                                                | send broadcast msg: key=I am node3 value=1
received broadcast msg: key=I am node3 value=1 | received broadcast msg: key=I am node3 value=1 |
                                               | send broadcast msg: key=I am node2 value=2     |
received broadcast msg: key=I am node2 value=2 |                                                | received broadcast msg: key=I am node2 value=2
```

## Note:

In [TransmitLimitedQueue](https://godoc.org/github.com/hashicorp/memberlist#TransmitLimitedQueue), the one with lower transmit count reaches, so the order may be incorrect.
And beware of [NumNodes()](https://godoc.org/github.com/hashicorp/memberlist#TransmitLimitedQueue), you have to correctly set the number of nodes [Broadcast messages](https://godoc.org/github.com/hashicorp/memberlist#Broadcast) may not reach to node.

```
type MyEventDelegate struct {
  Num int
}
func (d *MyEventDelegate) NotifyJoin(node *memberlist.Node) {
  d.Num += 1
}
func (d *MyEventDelegate) NotifyLeave(node *memberlist.Node) {
  d.Num -= 1
}
:
:

type MyDelegate struct {
  msgCh      chan []byte
  broadcasts *memberlist.TransmitLimitedQueue
}
func (d *MyDelegate) NotifyMsg(msg []byte) {
  d.msgCh <- msg
}
func (d *MyDelegate) GetBroadcasts(overhead, limit int) [][]byte {
  return d.broadcasts.GetBroadcasts(overhead, limit)
}
:
:

// main
e := new(MyEventDelegate)
e.Num = 0

d := new(MyDelegate)
d.broadcasts = new(memberlist.TransmitLimitedQueue)
d.broadcasts.NumNodes = func() int {
  return e.Num
}
:
:

mlist, _ := memberlist.Create(conf)
mlist.Join([]string{join})

// init num
e.Num = mlist.NumMembers()
```

## NextStep:

- Distributed Lock Manager
- Distributed Cache
- Deployment Manager
