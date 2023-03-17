package gcache

//ByteView 为lru.Cache实际存储的类型
//b 存储实际的值 entry.Value，并且实现了Value接口
type ByteView struct {
	b []byte
}

// Len 返回b 所占内存大小，也为了满足Value接口
func (v ByteView) Len() int {
	return len(v.b)
}

// ByteSlice 返回一个拷贝的值，防止缓存被外部代码修改.
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

// String 返回b的string类型
func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
