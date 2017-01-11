package ringbuf

type RingBuffer4 struct {
	data              []interface{}
	size              int
	readPos, writePos int
}

func NewRingBuffer4(size int) *RingBuffer4 {
	return &RingBuffer4{
		data: make([]interface{}, size),
		size: size,
	}
}

func (rb *RingBuffer4) Write(v interface{}) error {
	if rb.IsFull() {
		return ErrFull
	}
	rb.data[rb.mod(rb.writePos, rb.size)] = v // writePos % size

	rb.writePos = rb.incr(rb.writePos)
	return nil
}

func (rb *RingBuffer4) Read() (v interface{}, err error) {
	if rb.IsEmpty() {
		err = ErrEmpty
		return
	}
	v = rb.data[rb.mod(rb.readPos, rb.size)]
	rb.readPos = rb.incr(rb.readPos)
	return
}

func (rb *RingBuffer4) IsFull() bool {
	// readPos 和 writePos 是否相差 size
	return rb.mod(rb.readPos+rb.size, 2*rb.size) == rb.writePos
}

func (rb *RingBuffer4) IsEmpty() bool {
	return rb.readPos == rb.writePos
}

func (rb *RingBuffer4) Size() int { return rb.size }

func (rb *RingBuffer4) incr(i int) int {
	i++
	if i >= 2*rb.size {
		return 0
	}
	return i
}

func (rb *RingBuffer4) mod(i, n int) int {
	if i >= n { // 本应该是 `for i >= n` 的，但是本例中 i 不会大于 2*n
		return i - n
	}
	return i
}

var _ RingBuffer = new(RingBuffer4)
