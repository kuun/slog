// Copyright 2013 Google Inc. All Rights Reserved.
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

// Package buffer is a buffer pool used by slog, most code is from
// github/golang/glog
package buffer

import (
	"bytes"
	"sync"
)

// buffer holds a byte Buffer for reuse. The zero value is ready for use.
type Buffer struct {
	bytes.Buffer
	Tmp  [24]byte // temporary byte array for creating headers.
	next *Buffer
}

// freeList is a list of byte buffers, maintained under freeListMu.
var freeList *Buffer

// current buffers count in free list
var freeListLen int

const (
	// free list's max length
	maxFreeCount = 256
	// cached buffer's max size in free list
	cachedBufferMaxSize = 256
)

// freeListMu maintains the free list. It is separate from the main mutex
// so buffers can be grabbed and printed to without holding the main lock,
// for better parallelization.
var freeListMu sync.Mutex

// Some custom tiny helper functions to print the log header efficiently.
const digits = "0123456789"

// twoDigits formats a zero-prefixed two-digit integer at buf.tmp[i].
func (buf *Buffer) TwoDigits(i, d int) {
	buf.Tmp[i+1] = digits[d%10]
	d /= 10
	buf.Tmp[i] = digits[d%10]
}

// nDigits formats an n-digit integer at buf.tmp[i],
// padding with pad on the left.
// It assumes d >= 0.
func (buf *Buffer) NDigits(n, i, d int, pad byte) {
	j := n - 1
	for ; j >= 0 && d > 0; j-- {
		buf.Tmp[i+j] = digits[d%10]
		d /= 10
	}
	for ; j >= 0; j-- {
		buf.Tmp[i+j] = pad
	}
}

// someDigits formats a zero-prefixed variable-width integer at buf.tmp[i].
func (buf *Buffer) SomeDigits(i, d int) int {
	// Print into the top, then copy down. We know there's space for at least
	// a 10-digit number.
	j := len(buf.Tmp)
	for {
		j--
		buf.Tmp[j] = digits[d%10]
		d /= 10
		if d == 0 {
			break
		}
	}
	return copy(buf.Tmp[i:], buf.Tmp[j:])
}

// getBuffer returns a new, ready-to-use buffer.
func GetBuffer() *Buffer {
	freeListMu.Lock()
	b := freeList
	if b != nil {
		freeList = b.next
		freeListLen -= 1
	}
	freeListMu.Unlock()
	if b == nil {
		b = new(Buffer)
	} else {
		b.next = nil
		b.Reset()
	}
	return b
}

// putBuffer returns a buffer to the free list.
func PutBuffer(b *Buffer) {
	if b.Len() >= cachedBufferMaxSize {
		// Let big buffers die a natural death.
		return
	}
	freeListMu.Lock()
	// If buffer pool is full, let the buffer die a natural death.
	if freeListLen < maxFreeCount {
		b.next = freeList
		freeList = b
	}
	freeListMu.Unlock()
}
