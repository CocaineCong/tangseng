package se

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// BinaryWrite any -> bytes.Buffer
func BinaryWrite(buf *bytes.Buffer, v any) (err error) {
	size := binary.Size(v)
	if size <= 0 {
		return fmt.Errorf("encodePostings binary.Size err,size: %v", size)
	}

	return binary.Write(buf, binary.LittleEndian, v)
}
