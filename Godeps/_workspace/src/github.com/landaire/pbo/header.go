//go:generate stringer -type=Flag

package pbo

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"time"
)

type Flag uint32

const (
	Uncompressed Flag = 0x00000000
	Packed       Flag = 0x43707273
	ProductEntry Flag = 0x56657273 // resistance/elite/arma
)

type FileEntry struct {
	Name          string
	Flag          Flag
	UnpackedSize  uint32
	ReservedField uint32
	Timestamp     time.Time
	DataBlockSize uint32
	dataOffset    int64
	pbo           *Pbo
}

func (f FileEntry) IsNull() bool {
	return f.Name == "" && f.Flag == 0 &&
		f.UnpackedSize == 0 && f.ReservedField == 0 &&
		f.Timestamp.Unix() == 0 && f.DataBlockSize == 0
}

// Implement the io.Reader interface
func (f FileEntry) Read(p []byte) (n int, err error) {
	file := f.pbo.file

	offset, err := file.Seek(0, os.SEEK_CUR)
	if offset > f.dataOffset+int64(f.DataBlockSize) {
		return 0, io.EOF
	}

	if uint32(len(p)) > f.DataBlockSize {
		return file.Read(p[:f.DataBlockSize])
	}

	return file.Read(p)
}

func (f *FileEntry) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case os.SEEK_SET:
		if offset < 0 || offset > int64(f.DataBlockSize) {
			return 0, os.ErrInvalid
		}

		f.pbo.file.Seek(f.dataOffset+offset, whence)
		return offset, nil

	// TODO: Support other whence values
	default:
		return 0, os.ErrInvalid
	}
}

// Gets the length of the entry block
func (f FileEntry) EntrySize() int {
	// Length of the name includes a null terminator
	// the 4 * 5 is 4 bytes (uint32) for the 5 fields
	return len(f.Name) + 1 + (4 * 5)
}

func (f FileEntry) String() string {
	return fmt.Sprintf(
		"Name: %s\nFlag: 0x%08X (%s)\nOriginal Size: %d\nReserved: %d\nTimestamp: %d (%s)\nData Size: %d",
		f.Name,
		uint32(f.Flag), f.Flag.String(),
		f.UnpackedSize,
		f.ReservedField,
		f.Timestamp.Unix(), f.Timestamp,
		f.DataBlockSize,
	)
}

type HeaderExtension struct {
	FileEntry
	ExtendedFields map[string]string
}

func (e *HeaderExtension) ReadExtendedFields(r *bufio.Reader) {
	maybeNullByte, _ := r.Peek(1)
	if maybeNullByte[0] != 0 {
		e.ExtendedFields = make(map[string]string)
	}

	for maybeNullByte[0] != 0 {
		key, _ := r.ReadString(0)
		value, _ := r.ReadString(0)

		e.ExtendedFields[key[:len(key)-1]] = value[:len(value)-1]

		maybeNullByte, _ = r.Peek(1)
	}

	r.ReadByte()
}

// Gets the length of the entry block
func (f HeaderExtension) EntrySize() int {
	// Length of the name includes a null terminator
	// the 4 * 5 is 4 bytes (uint32) for the 5 fields
	baseSize := f.FileEntry.EntrySize()

	for key, val := range f.ExtendedFields {
		// + 2 for the null terminator for each key/val
		baseSize += len(key) + len(val) + 2
	}

	// There's a null terminator at the end of the block
	return baseSize + 1
}

func readEntry(r *bufio.Reader) FileEntry {
	entry := FileEntry{}
	entry.Name, _ = r.ReadString(0)
	entry.Name = entry.Name[:len(entry.Name)-1]

	var timestamp uint32
	fields := []interface{}{
		&entry.Flag,
		&entry.UnpackedSize,
		&entry.ReservedField,
		&timestamp,
		&entry.DataBlockSize,
	}

	for _, field := range fields {
		// Ignore errors -- swag
		binary.Read(r, binary.LittleEndian, field)
	}

	entry.Timestamp = time.Unix(int64(timestamp), 0)

	return entry
}
