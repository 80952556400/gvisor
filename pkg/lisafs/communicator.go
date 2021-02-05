// Copyright 2021 The gVisor Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lisafs

// Communicator is a server side utility which represents exactly how the
// server is communicating with the client.
type Communicator interface {
	// PayloadBuf returns a slice to the payload section of its internal buffer
	// where the message can be marshalled. The handlers should use this to
	// populate the payload buffer with the message.
	//
	// The payload buffer contents *should* be preserved across calls with
	// different sizes. Note that this is not a guarantee, because a compromised
	// owner of a "shared" payload buffer can tamper with its contents anytime,
	// even when it's not its turn to do so.
	PayloadBuf(size uint32) []byte

	// SndRcvMessage sends message m. The caller must have populated PayloadBuf()
	// with payloadLen bytes. The caller expects to receive wantFDs FDs.
	// Any recieved FDs must be accessible via ReleaseFDs(). It returns the
	// response message along with the response payload length.
	SndRcvMessage(m MID, payloadLen uint32, wantFDs uint8) (MID, uint32, error)

	// AddDonationFD starts tracking fd. The next call to ReleaseFDs will include
	// fd in the order it was added. Communicator takes ownership of fd.
	AddDonationFD(fd int)

	// ReleaseFDs returns the accumulated FDs and stops tracking them. The
	// ownership of the FDs is transferred to the caller.
	ReleaseFDs() []int
}

// fdTracker is a partial implementation of Communicator. It can be embedded in
// Communicator implementations to keep track of FD donations.
type fdTracker struct {
	fds []int
}

// AddDonationFD implements Communicator.AddDonationFD.
func (d *fdTracker) AddDonationFD(fd int) {
	d.fds = append(d.fds, fd)
}

// ReleaseFDs implements Communicator.ReleaseFDs.
func (d *fdTracker) ReleaseFDs() []int {
	ret := d.fds
	d.fds = d.fds[:0]
	return ret
}
