package engine

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

const PageSize = 65536

// Page represent a one piece of the large store, that contains a DataRecord(s)
// limited by the PageSize
type Page struct {
	Number int
	Buffer bytes.Buffer
	Dirty  bool
}

// DataRowLocator is a pointer of stored data record
type DataRowLocator struct {
	Page   int
	Offset int64
	Size   int
	Loaded bool
}

// DataEntity is a combination of record and it pointer (locator)
type DataEntity struct {
	Record  *DataRecord
	Locator DataRowLocator
}

var pagesMap = make(map[uint32]*Page)

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

func (p *Page) ReadDataRecord(offset int64) *DataEntity {

	buf := make([]byte, 4)
	_, err := p.Buffer.Read(buf)

	if err == io.EOF {
		return nil
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
		return nil
	}

	if err != nil {
		panic(err)
	}

	locator := DataRowLocator{
		Page:   p.Number,
		Offset: offset,
		Size:   n + 4,
		Loaded: true,
	}

	entity := &DataEntity{
		Record:  record,
		Locator: locator,
	}

	return entity
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
		err := pageFile.Close()
		if err != nil {
			panic(err)
		}
		pageFile = nil
	}
}

/// Flushed all dirty pages to disk
func Flush() {

}

func PreloadBookPages() {
	pagesCount := GetPagesCount()
	for i := 0; uint32(i) < pagesCount; i++ {
		pagesMap[uint32(i)] = _readPage(i)
	}
}

func GetPagesCount() uint32 {
	if pageFile == nil {
		panic("book.dat is not opened.")
	}

	info, err := pageFile.Stat()

	if err != nil {
		panic(err)
	}

	return uint32(info.Size() / PageSize)
}

func GetPage(num uint32) *Page {
	page, ok := pagesMap[num]
	if !ok {
		return nil
	}
	return page
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

func _readPage(num int) *Page {
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

	return &Page{
		Number: num,
		Buffer: *buffer,
		Dirty:  false,
	}
}
