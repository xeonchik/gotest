package main

import (
	"encoding/json"
	"fmt"
	"gotest/engine"
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

	//idx.Print()

	fmt.Println(idx.Get(4021).Locator)
	fmt.Println(idx.Get(40021).Locator)

	start := time.Now().UnixNano()

	acc := make([]int, 100000)

	idx.Tree.Ascend(nil, func(item interface{}) bool {
		it := item.(*engine.PKItem)
		acc = append(acc, it.PrimaryKey)

		////record := engine.ReadRecordByLocator(it.Locator)
		//
		//if record == nil {
		//	return true
		//}

		//fmt.Println(record.Primary)
		return true
	})

	timer := (time.Now().UnixNano() - start) / 1000
	fmt.Printf("Result time: %d mcs", timer)

	engine.Close()
}

var pkIdx = engine.CreatePKIndex()
var multiIdx = engine.CreateMulti()
var idxMultiCities = engine.CreateMulti()

func readBook(num int) {
	page := engine.ReadPage(num)

	var pos int64 = 0

	for i := 0; i < 10000; i++ {
		_rec, locator := page.ReadDataRecord(pos)

		if _rec == nil {
			break
		}

		pos = locator.Offset

		// build indexes
		pkIdx.Add(_rec, *locator, _rec.Primary)
		val, _ := strconv.Atoi(_rec.Data)

		multiIdx.Add(val, _rec.Primary)
		idxMultiCities.AddArray(_rec.Cities, _rec.Primary)

		//fmt.Println(_rec)
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

		var cities = make([]int, 10)

		for ci := 0; ci < 10; ci++ {
			cities[ci] = rand.Intn(100000-1000) + 1000
		}

		dataRecord := engine.DataRecord{
			Primary: writePrimary,
			Data:    strconv.Itoa(i),
			Cities:  cities,
			Active:  false,
		}

		writePrimary++

		recordOffset := page.PlaceRecord(dataRecord)
		fmt.Printf("datarecord #%d offset %d\n", i, recordOffset)
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
