package bcast

/*
   bcast package for Go. Broadcasting on a set of channels.

   Copyright © 2013 Alexander I.Grafov <grafov@gmail.com>.
   All rights reserved.
   Use of this source code is governed by a BSD-style
   license that can be found in the LICENSE file.
*/

import (
	"testing"
	"time"
	"sync"
)

// Create new broadcast group.
// Join two members.
func TestNewGroupAndJoin(t *testing.T) {
	group := NewGroup()
	member1 := group.Join()
	member2 := group.Join()
	if member1.group != member2.group {
		t.Fatal("group for these members must be same")
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
		t.Fatal("incorrect length of `out` slice")
	}

	go group.Broadcasting(2 * time.Second)

	member1.Close()
	if len(group.out) > 1 || group.out[0] != member2.In {
		t.Fatal("unjoin member does not work")
	}
}

type Adder struct {
	count int
	l sync.Mutex
}

func (a *Adder) Add(c int) {
	a.l.Lock()
	a.count += c
	a.l.Unlock()
}

// Create new broadcast group.
// Join 12 members.
// Broadcast one integer from each member.
func TestBroadcast(t *testing.T) {
	var valcount Adder
	valcount.count = 0

	group := NewGroup()

	for i := 1; i <= 12; i++ {
		go func(i int) {
			m := group.Join()
			m.Send(i)

			for {
				val := m.Recv()
				if val.(int) == i {
					t.Fatal("sent value was received by sender")
				}
				valcount.Add(1)
			}
		}(i)
	}

	group.Broadcasting(100 * time.Millisecond)
	if valcount.count != 12*12-12 { // number of channels * number of messages - number of channels
		t.Fatal("not all messages broadcasted")
	}
}

// Create new broadcast group.
// Join 12 members.
// Broadcast one integer from each member.
func TestBroadcastFromMember(t *testing.T) {
	var valcount int

	group := NewGroup()

	for i := 1; i <= 12; i++ {
		go func(i int, group *Group) {
			m := group.Join()
			m.Send(i)

			for {
				val := m.Recv()
				if val.(int) == i {
					t.Fatal("sent value was received by sender")
				}
				valcount++
			}
		}(i, group)
	}

	group.Broadcasting(100 * time.Millisecond)
	if valcount != 12*12-12 { // number of channels * number of messages - number of channels
		t.Fatal("not all messages broadcasted")
	}
}

// Create new broadcast group.
// Join 12 members.
// Make group broadcast to all group members.
func TestGroupBroadcast(t *testing.T) {
	var valcount int

	group := NewGroup()

	for i := 1; i <= 12; i++ {
		go func(i int, group *Group) {
			m := group.Join()
			if val := m.Recv(); val != "group message" {
				t.Fatal("incorrect message received")
			}
			valcount++
		}(i, group)
	}

	go group.Broadcasting(200 * time.Millisecond)
	group.Send("group message")
	time.Sleep(150 * time.Millisecond)

	if valcount != 12 {
		t.Fatal("not all messages broadcasted")
	}
}

// Create new broadcast group.
// Join 128 members.
// Make group broadcast to all members.
func TestBroadcastOnLargeNumberOfMembers(t *testing.T) {
	const max = 128
	var valcount int

	group := NewGroup()
	for i := 1; i <= max; i++ {
		go func(i int, group *Group) {
			m := group.Join() // here is the problem
			m.Send(i)
			for {
				if val := m.Recv(); val.(int) == i {
					t.Fatal("sent value should not received by sender")
				}
				valcount++
			}
		}(i, group)
	}

	group.Broadcasting(1 * time.Second)
	if valcount != max*max-max { // number of channels * number of messages - number of channels
		t.Fatalf("not all messages broadcasted (%d/%d)", valcount, max*max-max)
	}
}
