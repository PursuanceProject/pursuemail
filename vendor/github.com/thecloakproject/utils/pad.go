package utils

func PadBytes(data []byte, blockSize int) []byte {
	// Add padding (originally for correctness, now for simplicity)
	for len(data) % blockSize != 0 {
		data = append(data, 0x0)
	}
	return data
}

func PadString(data string, multiple int) string {
	for len(data) % multiple != 0 {
		data += string(0x0)
	}
	return data
}
