package engine

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"os"
)

const PageSize = 65536

// Страница содержащая записи DataEntity
type Page struct {
	Number int
	Buffer bytes.Buffer
}

// Представление об одной записе на странице, хранит payload
type DataEntity struct {
	Page   int
	Offset int
	Size   int
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

func (p *Page) ReadDataRecord(offset int64) (*DataRecord, int64) {
	rd := bytes.NewReader(p.Buffer.Bytes())
	_, err := rd.Seek(int64(offset), io.SeekStart)

	if err != nil {
		panic(err)
	}

	decoder := gob.NewDecoder(rd)

	var rec DataRecord
	err = decoder.Decode(&rec)

	if err == io.EOF {
		pos, _ := rd.Seek(0, io.SeekCurrent)
		return nil, pos
	}

	if err != nil {
		panic(err)
	}

	pos, _ := rd.Seek(0, io.SeekCurrent)

	return &rec, pos
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
