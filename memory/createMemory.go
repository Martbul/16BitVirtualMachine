package memory

type DataView struct {
	buffer []byte
}

func NewDataView(size int) *DataView {
	return &DataView{buffer: make([]byte, size)}
}

func (dv *DataView) GetBuffer() []byte {
	return dv.buffer
}

func CreateMemory(sizeInBytes int) *DataView {
	return NewDataView(sizeInBytes)
}

// GetUint16 reads a 16-bit unsigned integer (Big-Endian)
func (dv *DataView) GetUint16(offset int) uint16 {
	return uint16(dv.buffer[offset])<<8 | uint16(dv.buffer[offset+1])
}

// SetUint16 writes a 16-bit unsigned integer (Big-Endian)
func (dv *DataView) SetUint16(offset int, value uint16) {
	dv.buffer[offset] = byte(value >> 8) // High byte first
	dv.buffer[offset+1] = byte(value)    // Low byte second
}

// GetUint8 reads an 8-bit unsigned integer at a specific offset
func (dv *DataView) GetUint8(offset int) uint8 {
	return dv.buffer[offset]
}

// SetUint8 writes an 8-bit unsigned integer at a specific offset
func (dv *DataView) SetUint8(offset int, value uint8) {
	dv.buffer[offset] = value
}
