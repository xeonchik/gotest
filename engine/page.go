package engine

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

const PageSize = 65536

// Страница содержащая записи DataRowLocator
type Page struct {
	Number int
	Buffer bytes.Buffer
	Dirty  bool
}

// Представление об одной записе на странице, хранит payload
type DataRowLocator struct {
	Page   int
	Offset int64
	Size   int
	Loaded bool
}

func (p *Page) Free() int {
	return PageSize - p.Buffer.Len()
}

func (p *Page) PlaceRecord(record *DataRecord) int {
	data, err := record.MarshalBinary()
	if err != nil {
		panic(err)
	}

	length := len(data)
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(length))

	p.Buffer.Write(buf)
	p.Buffer.Write(data)

	return p.Buffer.Len()
}

func ReadRecordByLocator(locator *DataRowLocator) *DataRecord {
	if locator.Loaded {
		return nil
	}

	_ = ReadPage(locator.Page)

	offset := int64(locator.Page*PageSize) + locator.Offset
	_, err := pageFile.Seek(offset, io.SeekStart)

	if err != nil {
		panic(err)
	}

	buf := make([]byte, 4)
	pageFile.Seek(offset, io.SeekStart)
	_, err = pageFile.Read(buf)

	if err != nil {
		panic(err)
	}

	length := binary.BigEndian.Uint32(buf)
	recordBuffer := make([]byte, length)
	_, err = pageFile.Read(recordBuffer)
	if err != nil {
		panic(err)
	}
	record := &DataRecord{}
	err = record.UnmarshalBinary(recordBuffer)
	if err == io.EOF {
		return nil
	}

	locator.Loaded = true

	return record
}

func (p *Page) ReadDataRecord(offset int64) (*DataRecord, *DataRowLocator) {

	buf := make([]byte, 4)
	_, err := p.Buffer.Read(buf)

	if err == io.EOF {
		return nil, nil
	}

	if err != nil {
		panic(err)
	}

	length := binary.BigEndian.Uint32(buf)
	recordBuffer := make([]byte, length)
	p.Buffer.Read(recordBuffer)

	record := &DataRecord{}
	n, err := record.Unmarshal(recordBuffer)

	if err == io.EOF {
		return nil, nil
	}

	if err != nil {
		panic(err)
	}

	locator := &DataRowLocator{
		Page:   p.Number,
		Offset: offset,
		Size:   n + 4,
		Loaded: true,
	}

	return record, locator
}

var pageFile *os.File = nil

func OpenBook() {
	if pageFile == nil {
		fo, err := os.OpenFile("book.dat", os.O_RDWR, 0644)
		if err != nil {
			panic(err)
		}
		pageFile = fo
	}
}

func Close() {
	if pageFile != nil {
		pageFile.Close()
	}
}

func WritePage(page *Page) {
	if pageFile == nil {
		panic("book.dat is not opened.")
	}

	if page.Buffer.Len() > PageSize {
		panic("page oversize")
	}

	offset := PageSize * page.Number

	length, err := pageFile.WriteAt(page.Buffer.Bytes(), int64(offset))
	if err != nil {
		panic(err)
	}
	fmt.Println("written: ", length)
}

var pages []*Page

func ReadPage(num int) *Page {
	if pages[num] != nil {
		return pages[num]
	}

	if pageFile == nil {
		panic("book.dat is not opened.")
	}

	buf := make([]byte, PageSize)
	offset := num * PageSize
	pageFile.Seek(int64(offset), io.SeekStart)
	_, err := pageFile.Read(buf)

	if err != nil {
		panic(err)
	}

	buffer := &bytes.Buffer{}
	buffer.Write(buf)

	pages[num] = &Page{
		Number: num,
		Buffer: *buffer,
	}

	return pages[num]
}
