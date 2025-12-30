package data

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
)



func ReadData(file *os.File) (map[string]string, error) {
	pathMap := make(map[string]string)

	fileInfo, err := file.Stat()
	if err != nil {
		return pathMap, fmt.Errorf("Error getting file info: %w", err)
	}

	// take file size in bytes and make a buffer of that size
	size := fileInfo.Size()
	buff := make([]byte, size)

	// read entire file into memory
	_, err = io.ReadFull(file, buff)
	if err != nil {
		return pathMap, fmt.Errorf("Error reading file into buffer: %w", err)
	}

	// key length integer should always fit in 1 byte
	var keyLen uint8
	// value length integer should always fit in 2 bytes
	var valLen uint16
	// sliding pointer to navigate buffer
	var offset int

	// iterate through buffer and deserialize
	for offset < len(buff) {

		// read length of key, use length to read in key, adjust offset
		// simple type conversion since length is only 1 byte and not a []byte
		keyLen = buff[offset]
		offset++
		kl := int(keyLen)
		keyBytes := buff[offset : offset+kl]
		offset += kl

		// read length of value, use length to read in value, adjust offset
		// length contained in 2 bytes, need to convert []byte to a uint16 value
		valLen = binary.LittleEndian.Uint16(buff[offset : offset+2])
		offset += 2
		vl := int(valLen)
		valBytes := buff[offset : offset+vl]
		offset += vl
		// add key-value to map
		pathMap[string(keyBytes)] = string(valBytes)

	}

	return pathMap, nil

}

func EnsureData(filepath string) (*os.File, error) {

	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, errors.New(fmt.Sprint("Error opening file: ", err))
	}

	return file, nil

}

func DataUpdate(hashmap map[string]string, file *os.File) error {

	bufSize := 0

	for key, val := range hashmap {
		// key length is stored using 1 byte value length is stored using 2 bytes
		bufSize += 3 + len(key) + len(val)
	}

	b := make([]byte, 0, bufSize)
	buffer := bytes.NewBuffer(b)

	var valLen [2]byte //reusable valLen array of 2 bytes
	for key, val := range hashmap {
		// insert key length and key into buffer
		keyBytes := []byte(key)
		keyLen := byte(uint8(len(keyBytes)))

		buffer.WriteByte(keyLen)
		buffer.Write(keyBytes)
		
		// insert value length and value into buffer
		valBytes := []byte(val)
		binary.LittleEndian.PutUint16(valLen[:], uint16(len(valBytes)))
		
		buffer.Write(valLen[:])
		buffer.Write(valBytes)
	}

	err := file.Truncate(0)
	if err != nil {
		return fmt.Errorf("Error truncating file: %w", err)
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("Error seeking to beginning of file: %w", err)
	}

	_, err = file.Write(buffer.Bytes())
	if err != nil {
		return fmt.Errorf("Error writing contents of buffer to file: %w", err)
	}

	return nil
}
