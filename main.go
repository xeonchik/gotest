package main

import (
	"encoding/json"
	"fmt"
	"godoc/engine"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func main() {
	fmt.Println("Hello, World")

	engine.OpenBook()

	for i := 0; i < 100; i++ {
		//writeBook(i)
		//readBook(i)
	}

	//engine.FlushIndexToDisk(pkIdx, "primary.idx")
	idx, err := engine.ReadIndexFromDisk("primary.idx")

	if err != nil {
		panic(err)
	}

	start := time.Now().UnixNano()

	idx.Tree.Ascend(&engine.PKItem{
		PrimaryKey: 1,
	}, func(item interface{}) bool {
		it := item.(*engine.PKItem)

		it.Record = engine.ReadRecordByLocator(&it.Locator)
		it.Locator.Loaded = true

		return true
	})

	idx.Tree.Ascend(&engine.PKItem{
		PrimaryKey: 1,
	}, func(item interface{}) bool {
		it := item.(*engine.PKItem)
		_ = engine.ReadRecordByLocator(&it.Locator)
		return true
	})

	timer := (time.Now().UnixNano() - start) / 1000
	fmt.Printf("Result time: %d mcs", timer)
	engine.Close()
}

var pkIdx = engine.CreatePKIndex()

//var multiIdx = engine.CreateMulti()

func readBook(num int) {
	page := engine.ReadPage(num)

	var pos int64 = 0

	for i := 0; i < 10000; i++ {
		dataRecord, locator := page.ReadDataRecord(pos)

		if dataRecord == nil {
			break
		}

		pos = pos + int64(locator.Size)

		// build indexes
		pkIdx.Add(dataRecord, *locator, dataRecord.ID)
	}
}

var writePrimary = 1

func writeBook(num int) {
	page := &engine.Page{
		Number: num,
	}

	rand.Seed(time.Now().UTC().UnixNano())

	for i := 0; i < 1220; i++ {
		if page.Free() < 1024 { //fixme: replace to record.Size() comparison
			fmt.Println(page.Free())
			break
		}

		dataRecord := &engine.DataRecord{
			ID:     uint64(writePrimary),
			Data:   strconv.Itoa(i),
			Active: false,
		}

		for ci := 0; ci < 10; ci++ {
			city := &engine.City{Value: int32(rand.Intn(100000-1000) + 1000)}
			dataRecord.Cities = append(dataRecord.Cities, city)
		}

		writePrimary++

		recordOffset := page.PlaceRecord(dataRecord)
		fmt.Printf("datarecord #%d offset %d\n", writePrimary, recordOffset)
	}

	engine.OpenBook()
	engine.WritePage(page)
}

type testStruct struct {
	Id      int
	Sectors []int
	Cities  []int
	Active  bool
}

func pushHandler(rw http.ResponseWriter, req *http.Request) {

	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		panic(err)
	}

	var t testStruct
	err = json.Unmarshal(body, &t)
	if err != nil {
		panic(err)
	}
	log.Println(t.Sectors)
}

func listHandler(rw http.ResponseWriter, req *http.Request) {

}
