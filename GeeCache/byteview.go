package geecache

// 定义数据存储类型

type ByteView struct {
	// 这里使用[]byte，作为实际数据存储的原因是：byte可以支持任意类型的数据，作为二进制的方式进行，其实现的Len函数是为了满足之前cache的value类型。
	b []byte
}

func (v ByteView) Len() int {
	return len(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

// 使用深拷贝，使内部缓存不被发现，而且可以输出

func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func (v ByteView) String() string {
	return string(v.b)
}
