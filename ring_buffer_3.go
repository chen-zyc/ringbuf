package ringbuf

// RingBuffer3 使用镜像指示位判断环是满还是空.
// size 必须是2的次幂
// 参考 <https://zh.wikipedia.org/wiki/%E7%92%B0%E5%BD%A2%E7%B7%A9%E8%A1%9D%E5%8D%80>
type RingBuffer3 struct {
	data              []interface{}
	size              int
	readPos, writePos int
}

func NewRingBuffer3(size int) *RingBuffer3 {
	// TODO 判断size是否是2的次幂
	return &RingBuffer3{
		data: make([]interface{}, size),
		size: size,
	}
}

func (rb *RingBuffer3) Write(v interface{}) error {
	if rb.IsFull() {
		return ErrFull
	}
	rb.data[rb.writePos&(rb.size-1)] = v // writePos % size

	rb.writePos = rb.incr(rb.writePos)
	return nil
}

func (rb *RingBuffer3) Read() (v interface{}, err error) {
	if rb.IsEmpty() {
		err = ErrEmpty
		return
	}
	v = rb.data[rb.readPos&(rb.size-1)]
	rb.readPos = rb.incr(rb.readPos)
	return
}

func (rb *RingBuffer3) IsFull() bool {
	// readPos 和 writePos 是否相差 size
	return rb.writePos == (rb.readPos ^ rb.size)
}

func (rb *RingBuffer3) IsEmpty() bool {
	return rb.readPos == rb.writePos
}

func (rb *RingBuffer3) Size() int { return rb.size }

func (rb *RingBuffer3) incr(i int) int {
	return (i + 1) & (2*rb.size - 1)
}

var _ RingBuffer = new(RingBuffer3)
