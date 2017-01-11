package ringbuf

import (
	"errors"
	"sync"
)

type BufferType int

const (
	ONE_MORE       BufferType = 1 << iota // 多使用一个元素
	FIELD                                 // 多使用一个字段
	POWER_2                               // 环大小是2的次幂
	VIRTUAL                               // 使用虚拟空间，POWER_2 的普通版
	WITH_LOCK                             // 加锁
	WITH_OVERWRITE                        // TODO 可覆盖写
)

var (
	ErrFull  = errors.New("ring buffer is full")
	ErrEmpty = errors.New("ring buffer is empty")
)

type RingBuffer interface {
	Write(v interface{}) error
	Read() (interface{}, error)
	Size() int
	IsFull() bool
	IsEmpty() bool
}

func NewRingBuffer(size int, typ BufferType) (rb RingBuffer) {
	switch {
	case typ&ONE_MORE != 0:
		rb = NewRingBuffer1(size)
	case typ&FIELD != 0:
		rb = NewRingBuffer2(size)
	case typ%POWER_2 != 0:
		rb = NewRingBuffer3(size)
	case typ%VIRTUAL != 0:
		rb = NewRingBuffer4(size)
	default:
		rb = NewRingBuffer4(size)
	}
	if typ % WITH_LOCK {
		rb = NewRingBufferLock(rb)
	}
	return
}

type RingBufferLock struct {
	rb RingBuffer
	sync.Mutex
}

func NewRingBufferLock(rb RingBuffer) *RingBufferLock {
	return &RingBufferLock{
		rb: rb,
	}
}

func (rbl *RingBufferLock) Write(v interface{}) error {
	rbl.Lock()
	err := rbl.rb.Write(v)
	rbl.Unlock()
	return err
}

func (rbl *RingBufferLock) Read() (v interface{}, err error) {
	rbl.Lock()
	v, err = rbl.rb.Read()
	rbl.Unlock()
	return
}

func (rbl *RingBufferLock) Size() int {
	// 不需要加锁
	return rbl.rb.Size()
}

func (rbl *RingBufferLock) IsFull() bool {
	rbl.Lock()
	full := rbl.rb.IsFull()
	rbl.Unlock()
	return full
}

func (rbl *RingBufferLock) IsEmpty() bool {
	rbl.Lock()
	empty := rbl.rb.IsEmpty()
	rbl.Unlock()
	return empty
}

func incr(i, n int) int {
	i++
	if i >= n {
		return 0
	}
	return i
}
