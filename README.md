bcast package for Go
====================

Broadcasting on a set of channels in Go. Go channels offer different usage patterns but not ready to use broadcast pattern.
This library solves the problem in direct way.

Usage example:

			import (
						 "github.com/grafov/bcast"
			)

			group := NewGroup() // create broadcast group

			member1 := group.Join() // joined member1
			member2 := group.Join() // joined member2
			member1.Send("test message") // broadcast message
			// `bcast.Send` use interface{} type then any values may be broadcasted

			val := member1.Recv() // value received

Of course it all has sense for members in different goroutines. You need pass `group` object to goroutine and call group.Join() in it.
See more examples in a test suit `bcast_test.go`.

License
-------

Library licensed under GPLv3. See LICENSE.

Project status
--------------

[![Is maintained?](http://stillmaintained.com/grafov/bcast.png)](http://stillmaintained.com/grafov/bcast)
