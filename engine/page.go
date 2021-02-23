package engine

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"os"
)

const PageSize = 65536

// Страница содержащая записи DataRowLocator
type Page struct {
	Number int
	Buffer bytes.Buffer
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

func (p *Page) PlaceRecord(record DataRecord) int {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)

	err := encoder.Encode(record)
	if err != nil {
		panic(err)
	}

	offset := p.Buffer.Len()
	p.Buffer.Write(buffer.Bytes())
	return offset
}

func (p *Page) ReadDataRecord(offset int64) (*DataRecord, *DataRowLocator) {
	rd := bytes.NewReader(p.Buffer.Bytes())
	_, err := rd.Seek(int64(offset), io.SeekStart)

	if err != nil {
		panic(err)
	}

	decoder := gob.NewDecoder(rd)

	rec := &DataRecord{}

	err = decoder.Decode(&rec)

	if err == io.EOF {
		return nil, nil
	}

	if err != nil {
		return nil, nil
	}

	pos, _ := rd.Seek(0, io.SeekCurrent)
	size := pos - offset

	locator := &DataRowLocator{
		Page:   p.Number,
		Offset: pos,
		Size:   int(size),
		Loaded: true,
	}

	return rec, locator
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

func ReadPage(num int) *Page {

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
	}
}
