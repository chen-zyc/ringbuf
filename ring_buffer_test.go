package ringbuf

import (
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	bufferSize = 16
)

func testCorrectness(t *testing.T, rb RingBuffer) {
	require.Equal(t, bufferSize, rb.Size())

	require.Equal(t, true, rb.IsEmpty())
	require.Equal(t, false, rb.IsFull())

	// 写一个，读一个
	require.Empty(t, rb.Write(1), "write 1")
	require.Equal(t, false, rb.IsEmpty())
	require.Equal(t, false, rb.IsFull())

	d, err := rb.Read()
	require.Empty(t, err)
	require.Equal(t, 1, d.(int))
	require.Equal(t, true, rb.IsEmpty())
	require.Equal(t, false, rb.IsFull())

	// 写满
	for i := 0; i < bufferSize; i++ {
		require.Empty(t, rb.Write(i))
	}
	require.Equal(t, false, rb.IsEmpty())
	require.Equal(t, true, rb.IsFull())
	require.Equal(t, ErrFull, rb.Write(-1))

	// 全读
	for i := 0; i < bufferSize; i++ {
		d, err = rb.Read()
		require.Empty(t, err)
		require.Equal(t, i, d.(int))
	}
	require.Equal(t, true, rb.IsEmpty())
	require.Equal(t, false, rb.IsFull())
	d, err = rb.Read()
	require.Equal(t, ErrEmpty, err)
	require.Empty(t, d)

	// 再写和读
	for i := 1; i <= 5; i++ {
		require.Empty(t, rb.Write(i))
	}
	require.Equal(t, false, rb.IsEmpty())
	require.Equal(t, false, rb.IsFull())
	for i := 1; i <= 5; i++ {
		d, err = rb.Read()
		require.Empty(t, err)
		require.Equal(t, i, d.(int))
	}
	require.Equal(t, true, rb.IsEmpty())
	require.Equal(t, false, rb.IsFull())
}

func benchRB(b *testing.B, rb RingBuffer) {
	for i := 0; i < b.N; i++ {
		rb.Write(i)
		rb.Read()
	}
}

func benchPBC(b *testing.B, rb RingBuffer) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rb.Write(1)
			rb.Read()
		}
	})
}

func TestRB_1(t *testing.T) {
	testCorrectness(t, NewRingBuffer1(bufferSize))
}

func TestRB_2(t *testing.T) {
	testCorrectness(t, NewRingBuffer2(bufferSize))
}

func TestRB_3(t *testing.T) {
	testCorrectness(t, NewRingBuffer3(bufferSize))
}

func TestRB_4(t *testing.T) {
	testCorrectness(t, NewRingBuffer4(bufferSize))
}

func BenchmarkRB_1(b *testing.B) {
	benchRB(b, NewRingBuffer1(bufferSize))
}

func BenchmarkRB_2(b *testing.B) {
	benchRB(b, NewRingBuffer2(bufferSize))
}

func BenchmarkRB_3(b *testing.B) {
	benchRB(b, NewRingBuffer3(bufferSize))
}

func BenchmarkRB_4(b *testing.B) {
	benchRB(b, NewRingBuffer4(bufferSize))
}

func BenchmarkRBL_1(b *testing.B) {
	benchRB(b, NewRingBufferLock(NewRingBuffer1(bufferSize)))
}

func BenchmarkRBL_2(b *testing.B) {
	benchRB(b, NewRingBufferLock(NewRingBuffer2(bufferSize)))
}

func BenchmarkRBL_3(b *testing.B) {
	benchRB(b, NewRingBufferLock(NewRingBuffer3(bufferSize)))
}

func BenchmarkRBL_4(b *testing.B) {
	benchRB(b, NewRingBufferLock(NewRingBuffer4(bufferSize)))
}

func BenchmarkCorrRBL_1(b *testing.B) {
	benchPBC(b, NewRingBufferLock(NewRingBuffer1(bufferSize)))
}

func BenchmarkCorrRBL_2(b *testing.B) {
	benchPBC(b, NewRingBufferLock(NewRingBuffer2(bufferSize)))
}

func BenchmarkCorrRBL_3(b *testing.B) {
	benchPBC(b, NewRingBufferLock(NewRingBuffer3(bufferSize)))
}

func BenchmarkCorrRBL_4(b *testing.B) {
	benchPBC(b, NewRingBufferLock(NewRingBuffer4(bufferSize)))
}
