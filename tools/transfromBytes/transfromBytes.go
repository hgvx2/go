package transfromBytes

import (
	"bytes"
	"encoding/binary"
)

func BytesToUint16(b []byte) uint16 {
	byteBuf := bytes.NewBuffer(b)
	var x uint16
	binary.Read(byteBuf, binary.BigEndian, &x)
	return x
}

func Uint16ToBytes(n uint16) []byte {
	byteBuf := bytes.NewBuffer([]byte{})
	binary.Write(byteBuf, binary.BigEndian, n)
	return byteBuf.Bytes()
}

//整形转换成字节
func Int32ToBytes(n int32) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}
//字节转换成整形
func BytesToInt32(b []byte) int32 {
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return x
}
