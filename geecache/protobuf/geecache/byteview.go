package geecache

type ByteView struct{
	b []byte
}

func (v ByteView) Len() int{
	return len(v.b)
}

func cloneBytes(b []byte) []byte{
	c:=make([]byte,len(b))
	copy(c,b)
	return c
}

func (v ByteView)String()string{
	return string(v.b)
}

func (v ByteView) ByteSlices() []byte{
	return cloneBytes(v.b)
}