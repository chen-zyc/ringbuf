package ringbuf

// RingBuffer1 总是保持一个单元为空，以判断环是满还是空
type RingBuffer1 struct {
	data     []interface{}
	size     int
	readPos  int
	writePos int
}

func NewRingBuffer1(size int) *RingBuffer1 {
	return &RingBuffer1{
		data: make([]interface{}, size+1), // 多余的一个做哨兵
		size: size + 1,
	}
}

func (rb *RingBuffer1) Write(v interface{}) error {
	if rb.IsFull() {
		return ErrFull
	}
	rb.data[rb.writePos] = v
	rb.writePos = incr(rb.writePos, rb.size)
	return nil
}

func (rb *RingBuffer1) Read() (v interface{}, err error) {
	if rb.IsEmpty() {
		err = ErrEmpty
		return
	}
	v = rb.data[rb.readPos]
	rb.readPos = incr(rb.readPos, rb.size)
	return
}

func (rb *RingBuffer1) IsFull() bool {
	return incr(rb.writePos, rb.size) == rb.readPos
}

func (rb *RingBuffer1) IsEmpty() bool {
	return rb.writePos == rb.readPos
}

func (rb *RingBuffer1) Size() int { return rb.size - 1 }

var _ RingBuffer = new(RingBuffer1)
