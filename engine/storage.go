package engine

import (
	"bufio"
	"encoding/gob"
	"github.com/golang/protobuf/proto"
	pb "gotest/proto/dist/proto"
	"io/ioutil"
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

func ReadIndexFromDisk(name string) (*PKIndex, error) {
	in, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}

	indexStore := &pb.PKIndexStore{}
	err = proto.Unmarshal(in, indexStore)
	if err != nil {
		return nil, err
	}

	pkIdx := CreatePKIndex()

	for _, item := range indexStore.Items {
		record := &DataRecord{
			Primary: int(item.Primary),
		}
		locator := &DataRowLocator{
			Page:   int(item.PageNumber),
			Offset: item.Offset,
			Size:   int(item.Size),
			Loaded: false,
		}
		pkIdx.Load(record, *locator, record.Primary)
	}

	return pkIdx, nil
}

func FlushIndexToDisk(index *PKIndex, name string) {
	indexStore := &pb.PKIndexStore{}

	index.Tree.Ascend(nil, func(it interface{}) bool {
		item := it.(*PKItem)

		indexItem := &pb.PKIndexItem{
			Primary:    int32(item.PrimaryKey),
			PageNumber: int32(item.Locator.Page),
			Offset:     item.Locator.Offset,
			Size:       int32(item.Locator.Size),
		}
		indexStore.Items = append(indexStore.Items, indexItem)
		return true
	})

	out, err := proto.Marshal(indexStore)

	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(name, out, 0660)

	if err != nil {
		panic(err)
	}
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
