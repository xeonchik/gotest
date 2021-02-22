package engine

import (
	"bufio"
	"encoding/gob"
	"log"
	"os"
)

type DataRecord struct {
	Primary int
	Data    string
	Sectors []int
	Cities  []int
	Active  bool
}

var reader *bufio.Reader = nil

func ReadFromStorage() (*DataRecord, error) {
	if reader == nil {
		fo, err := os.Open("data")
		if err != nil {
			panic(err)
		}

		reader = bufio.NewReader(fo)
	}
	decoder := gob.NewDecoder(reader)

	var rec DataRecord
	err := decoder.Decode(&rec)
	if err != nil {
		log.Println("decode error:", err)
	}

	return &rec, err
}

func WriteToStorage(record *DataRecord) {
	fo, err := os.OpenFile("data", os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	//st, err := os.Stat("data")
	//fo.Seek(0, os.O_APPEND)

	defer fo.Close()

	writer := bufio.NewWriter(fo)
	encoder := gob.NewEncoder(writer)

	err = encoder.Encode(record)

	if err != nil {
		log.Println("encode error:", err)
	}

	writer.Flush()
}
