package bcast

/*
   bcast package for Go. Broadcasting on a set of channels.
   Copyright © 2013 Alexander I.Grafov <grafov@gmail.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see {http://www.gnu.org/licenses/}.
*/

import (
	// h "github.com/emicklei/hopwatch"
	"testing"
	"time"
)

// Create new broadcast group.
// Join two members.
func TestNewGroupAndJoin(t *testing.T) {
	group := NewGroup()
	member1 := group.Join()
	member2 := group.Join()
	if member1.group != member2.group {
		panic("group for these members must be same")
	}
}

// Create new broadcast group.
// Join two members.
// Unjoin first member.
func TestUnjoin(t *testing.T) {
	group := NewGroup()
	member1 := group.Join()
	member2 := group.Join()
	if len(group.out) != 2 {
		panic("incorrect length of `out` slice")
	}

	go group.Broadcasting(0)

	member1.Close()
	time.Sleep(1 * time.Second) // because Broadcasting executed concurrently
	if len(group.out) > 1 || group.out[0] != member2.In {
		panic("unjoin member does not work")
	}
}

// Create new broadcast group.
// Join 12 members.
// Broadcast one integer from each member.
func TestBroadcast(t *testing.T) {
	var valcount int

	group := NewGroup()

	for i := 1; i <= 12; i++ {
		go func(i int, group *Group) {
			m := group.Join()
			m.Send(i)

			for {
				val := m.Recv()
				if val.(int) == i {
					panic("sent value was received by sender")
				}
				valcount++
			}
		}(i, group)
	}

	group.Broadcasting(100 * time.Millisecond)
	if valcount != 12*12-12 { // number of channels * number of messages - number of channels
		panic("not all messages broadcasted")
	}
}

// Create new broadcast group.
// Join 512 members.
// Broadcast one integer from each member.
func TestBroadcastOnLargeNumberOfMembers(t *testing.T) {
	var valcount int

	group := NewGroup()
	for i := 1; i <= 512; i++ {
		go func(i int, group *Group) {
			m := group.Join()
			m.Send(i)
			for {
				val := m.Recv()
				if val.(int) == i {
					panic("sent value was received by sender")
				}
				valcount++
			}
		}(i, group)
	}
	group.Broadcasting(100 * time.Millisecond)
	if valcount != 512*512-512 { // number of channels * number of messages - number of channels
		panic("not all messages broadcasted")
	}
}
