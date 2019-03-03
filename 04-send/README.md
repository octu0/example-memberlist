# send

It is an implementation example that performs simple ping/pong.

## Usage:

`node1` will increment the received `ping` message and send a `pong` message to all nodes.

`node2` periodically sends a `ping` message to all nodes.  
If it received the `pong` message just before, `node2` will increment pong `Value` and send it.

run node1:

```
$ go run c.go node1.go
// => received msg: key=ping value=0
// => send to node2 msg: key=pong value=1
```

run node2:
```
$ go run c.go node2.go --join 192.168.0.25:7947
// => send to node1 msg: key=ping value=0
// => received msg: key=pong value=1
```

### Timeline:

```
                                      +
  node1 says                          | node2 says
+-------------------------------------|-----------------------------------------+
                                      | send to node1 msg: key=ping value=0
  received msg: key=ping value=0      | 
  send to node2 msg: key=pong value=1 |
                                      | received msg: key=pong value=1
                                      | send to node1 msg: key=ping value=2
  received msg: key=ping value=2      | 
  send to node2 msg: key=pong value=3 |
                                      | received msg: key=pong value=3
                                      | send to node1 msg: key=ping value=4
  received msg: key=ping value=4      |
  send to node2 msg: key=pong value=5 |
                                      | received msg: key=pong value=5
                                      |
                                      +
```

## Note:

You can send arbitrary values by using `[]byte` to send data.  
When sending, you can choose to use a [reliable transport](https://godoc.org/github.com/hashicorp/memberlist#Memberlist.SendReliable) or a [best effort transport](https://godoc.org/github.com/hashicorp/memberlist#Memberlist.SendBestEffort).

```
type MyMessage struct {
  Key    string  `json:"key"`
  Value  uint64  `json:"value"`
}
func (m *MyMessage) Bytes() []byte {
  data, err := json.Marshal(m)
  if err != nil {
    return []byte("")
  }
  return data
}

m := &MyMessage{
  Key: "hello world",
  Value: 12345,
}
list.SendReliable(node, m.Bytes())
```

## NextStep:

- RPC (Remote Procedure Call)
- P2P (node communication)
- Proxy / Gateway
