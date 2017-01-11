package ringbuf

// RingBuffer2 通过length字段判断环是空还是满
type RingBuffer2 struct {
	data     []interface{}
	size     int
	readPos  int
	writePos int
	length   int
}

func NewRingBuffer2(size int) *RingBuffer2 {
	return &RingBuffer2{
		data: make([]interface{}, size),
		size: size,
	}
}

func (rb *RingBuffer2) Write(v interface{}) error {
	if rb.IsFull() {
		return ErrFull
	}
	rb.data[rb.writePos] = v
	rb.writePos = incr(rb.writePos, rb.size)
	rb.length++
	return nil
}

func (rb *RingBuffer2) Read() (v interface{}, err error) {
	if rb.IsEmpty() {
		err = ErrEmpty
		return
	}
	v = rb.data[rb.readPos]
	rb.readPos = incr(rb.readPos, rb.size)
	rb.length--
	return
}

func (rb *RingBuffer2) IsFull() bool {
	return rb.length == rb.size
}

func (rb *RingBuffer2) IsEmpty() bool {
	return rb.length == 0
}

func (rb *RingBuffer2) Size() int { return rb.size }

var _ RingBuffer = new(RingBuffer2)
