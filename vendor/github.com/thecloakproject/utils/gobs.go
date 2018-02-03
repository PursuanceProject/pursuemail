package utils

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
)

func GobEncode(w io.Writer, data interface{}) error {
	err := gob.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Errorf("Error encoding '%+v': %v", data, err)
	}
	return nil
}

func GobDecode(data []byte, structure interface{}) error {
	err := gob.NewDecoder(bytes.NewReader(data)).Decode(structure)
	if err != nil {
		return fmt.Errorf("Error decoding '%+v': %v", structure, err)
	}
	return nil
}
