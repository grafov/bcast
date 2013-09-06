// Broadcast over channels.
package bcast

/*
   bcast package for Go. Broadcasting on a set of channels.

   Copyright © 2013 Alexander I.Grafov <grafov@gmail.com>.
   All rights reserved.
   Use of this source code is governed by a BSD-style
   license that can be found in the LICENSE file.
*/

import (
	// 	h "github.com/emicklei/hopwatch"
	"time"
)

// Internal structure to pack messages together wit info about sender.
type Message struct {
	sender  *chan interface{}
	payload interface{}
}

// Represents member of broadcast group.
type Member struct {
	group *Group            // send messages to others directly to group.In
	In    *chan interface{} // get messages from others to own channel
}

// Represents broadcast group.
type Group struct {
	in  *chan Message       // receive broadcasts from members
	out []*chan interface{} // broadcast messages to members
}

// Create new broadcast group.
func NewGroup() *Group {
	in := make(chan Message)
	return &Group{in: &in}
}

// Broadcast messages received from one group member to others.
// If incoming messages not arrived during `timeout` then function returns.
// Set `timeout` to zero to set unlimited timeout or use Broadcasting().
func (r *Group) Broadcasting(timeout time.Duration) {
	for {
		select {
		case received := <-*r.in:
			switch received.payload.(type) {
			case Member: // unjoining member

				for i, addr := range r.out {
					if addr == received.payload.(Member).In && received.sender == received.payload.(Member).In {
						r.out = append(r.out[:i], r.out[i+1:]...)
						break
					}
				}
			default: // receive payload and broadcast it

				for _, member := range r.out {
					if *received.sender != *member { // not return broadcast to sender

						go func(out *chan interface{}, received *Message) { // non blocking
							*out <- received.payload
						}(member, &received)

					}
				}
			}
		case <-time.After(timeout):
			if timeout > 0 {
				return
			}
		}
	}
}

// Join new member to broadcast.
func (r *Group) Join() *Member {
	out := make(chan interface{})
	r.out = append(r.out, &out)
	return &Member{group: r, In: &out}
}

// Unjoin member from broadcast group.
func (r *Member) Close() {
	*r.group.in <- Message{sender: r.In, payload: *r} // broadcasting of self means member closing
}

// Broadcast Message to others.
func (r *Member) Send(val interface{}) {
	*r.group.in <- Message{sender: r.In, payload: val}
}

// Get broadcast Message.
// As alternative you may get it from `In` channel.
func (r *Member) Recv() interface{} {
	return <-*r.In
}
