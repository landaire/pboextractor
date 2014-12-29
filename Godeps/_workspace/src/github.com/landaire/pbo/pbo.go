package pbo

import (
	"bufio"
	"os"
)

type Pbo struct {
	file            *os.File
	HeaderExtension *HeaderExtension
	Entries         []FileEntry
	dataOffset      int64
}

// Reads the file given by path and returns
// a Pbo pointer and err != nil if no errors occurred
func NewPbo(path string) (*Pbo, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	pbo := Pbo{
		file: file,
	}

	// Create a new buffered reader
	reader := bufio.NewReader(file)

	for {
		entry := readEntry(reader)
		if entry.Flag == ProductEntry {
			extension := HeaderExtension{
				FileEntry: entry,
			}

			extension.ReadExtendedFields(reader)
			pbo.HeaderExtension = &extension

			pbo.dataOffset += int64(extension.EntrySize())

			continue
		}

		pbo.dataOffset += int64(entry.EntrySize())

		if entry.IsNull() {
			break
		}

		entry.pbo = &pbo
		pbo.Entries = append(pbo.Entries, entry)
	}

	// Loop through all of our entries and set their data offset
	baseOffset := pbo.dataOffset
	for i := range pbo.Entries {
		entry := &pbo.Entries[i]
		// If the block is compressed, use the compressed size. If it's not, use the actual size
		entry.dataOffset = baseOffset
		baseOffset += int64(entry.DataBlockSize)
	}

	return &pbo, nil
}
