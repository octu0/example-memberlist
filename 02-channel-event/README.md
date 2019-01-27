# channel event

In this example, use [serialx/hashring](https://github.com/serialx/hashring) as an implementation of Consistent hashing.  

## Usage:

Both `node1`, `node2`, and `node3` are searching `KEY` by consistent hasing using the node information of memberlist.  
Because it is consistent hashing, you can see that searching for keys does not significantly affect the increase / decrease of nodes.


run node1.

```
$ go run c.go node1.go
// => current node size: 1
// => node1 search hello => 192.168.0.25:7947
// => node1 search world => 192.168.0.25:7947
```

run node2.

```
$ go run c.go node2.go --join 192.168.0.25:7947
// => join 192.168.0.25:7947
// => current node size: 2
// => node2 search foo => 192.168.0.25:7947
// => node2 search bar => 192.168.0.25:7947
```

run node3.

```
$ go run c.go node3.go --join 192.168.0.25:7947
// => join 192.168.0.25:7948
// => join 192.168.0.25:7947
// => current node size: 3
// => node3 search foo => 192.168.0.25:7949
// => node3 search world => 192.168.0.25:7948
```

### Timeline:

```
+                                         +                                      +
|                                         |                                      |
| node1 says                              |node2 says                            |node3 says
+-------------------------------------------------------------------------------------------------------------------------+
| current node size: 1                    |                                      |
| node1 search hello => 192.168.0.25:7947 |                                      |
| node1 search world => 192.168.0.25:7947 |                                      |
+-------------------------------------------------------------------------------------------------------------------------+
|                                         |join 192.168.0.25:7947                |
|                                         |current node size: 2                  |
|                                         |node2 search foo => 192.168.0.25:7947 |
|                                         |node2 search bar => 192.168.0.25:7947 |
+-------------------------------------------------------------------------------------------------------------------------+
| join 192.168.0.25:7948                  |                                      |
| current node size: 2                    |current node size: 2                  |
| node1 search hello => 192.168.0.25:7948 |node2 search foo => 192.168.0.25:7947 |
| node1 search world => 192.168.0.25:7948 |node2 search bar => 192.168.0.25:7947 |
+-------------------------------------------------------------------------------------------------------------------------+
|                                         |                                      |join 192.168.0.25:7947
|                                         |                                      |join 192.168.0.25:7948
|                                         |                                      |current node size: 3
|                                         |                                      |node3 search foo => 192.168.0.25:7949
|                                         |                                      |node3 search world => 192.168.0.25:7948
+-------------------------------------------------------------------------------------------------------------------------+
|current node size: 3                     |current node size: 3                  |
|node1 search hello => 192.168.0.25:7949  |node2 search foo => 192.168.0.25:7949 |current node size: 3
|node1 search world => 192.168.0.25:7948  |node2 search bar => 192.168.0.25:7947 |node3 search foo => 192.168.0.25:7949
|                                         |                                      |node3 search world => 192.168.0.25:7948
+-------------------------------------------------------------------------------------------------------------------------+

finally:
hello => 192.168.0.25:7949 (moves 3 times)
world => 192.168.0.25:7948 (moves 2 times)
foo   => 192.168.0.25:7949 (moves 2 times)
bar   => 192.168.0.25:7947 (moves 1 times)
```

## Note:

EventDelegate to receive event of join/leave/update of memberlist nodes.  
implement the method of [EventDelegate](https://godoc.org/github.com/hashicorp/memberlist#EventDelegate) Interface and assign it to [Config#Events](https://godoc.org/github.com/hashicorp/memberlist#Config)

```
type MyEventDelegate struct {}
func (d *MyEventDelegate) NotifyJoin(node *memberlist.Node) {
  // join event
}
func (d *MyEventDelegate) NotifyLeave(node *memberlist.Node) {
  // leave event
}
func (d *MyEventDelegate) NotifyUpdate(node *memberlist.Node) {
  // update event
}

conf := memberlist.DefaultLocalConfig()
conf.Events = new(MyEventDelegate)
list, err  := memberlist.Create(conf)
if err != nil {
	log.Fatal(err)
}
```

see also - https://github.com/hashicorp/memberlist/blob/master/event_delegate.go

## NextStep:

- Distributed cache (like memcached)
- Service discovery
- Automatic Configuration / Provision
