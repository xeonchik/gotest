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
	engine.PreloadBookPages()

	start := time.Now().UnixNano()

	for i := uint32(0); i < 1000; i++ {
		//writeBook(i)
		readBook(i)
	}

	timer := (time.Now().UnixNano() - start) / 1000
	log.Printf("init time: %d mcs", timer)

	http.HandleFunc("/select", selectHandler)
	http.HandleFunc("/sort", sortHandler)
	http.HandleFunc("/count", countHandler)
	http.ListenAndServe(":8090", nil)

	engine.Close()
}

var records = make(map[uint64]*engine.DataEntity)
var pkIdx = engine.CreatePKIndex()
var ratingIdx = engine.CreateFloatIndex()

func Sort(limit int, offset int) []*engine.DataRecord {
	result := make([]*engine.DataRecord, 0)

	ratingIdx.Tree.Ascend(nil, func(item interface{}) bool {
		if offset > 0 {
			offset--
			return true // just continue walking
		}
		if limit == 0 {
			return false
		}
		it := item.(*engine.FloatItem)
		record := records[it.Key]

		result = append(result, record.Record)
		limit--
		return true
	})

	return result
}

func Select(limit int, offset int) []*engine.DataRecord {
	result := make([]*engine.DataRecord, 0)

	pkIdx.Tree.Ascend(nil, func(item interface{}) bool {
		if offset > 0 {
			offset--
			return true // just continue walking
		}

		if limit == 0 {
			return false
		}
		it := item.(*engine.PKItem)
		result = append(result, it.Record)
		limit--
		return true
	})

	return result
}

func countHandler(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, "Count items: %d\n", pkIdx.Tree.Len())
}

func sortHandler(rw http.ResponseWriter, req *http.Request) {
	var limit, offset int

	limitParams, ok := req.URL.Query()["limit"]
	if !ok {
		limit = 10
	} else {
		limStr := limitParams[0]
		limit, _ = strconv.Atoi(limStr)
	}

	offsetParams, ok := req.URL.Query()["offset"]
	if !ok {
		offset = 0
	} else {
		offsetStr := offsetParams[0]
		offset, _ = strconv.Atoi(offsetStr)
	}

	start := time.Now().UnixNano()

	result := Sort(limit, offset)

	timer := (time.Now().UnixNano() - start) / 1000

	for i, element := range result {
		fmt.Fprintf(rw, "Entry: %+v %v\n", element, i)
	}

	log.Printf("sort time: %d mcs", timer)
}

func selectHandler(rw http.ResponseWriter, req *http.Request) {
	var limit, offset int

	limitParams, ok := req.URL.Query()["limit"]
	if !ok {
		limit = 10
	} else {
		limStr := limitParams[0]
		limit, _ = strconv.Atoi(limStr)
	}

	offsetParams, ok := req.URL.Query()["offset"]
	if !ok {
		offset = 0
	} else {
		offsetStr := offsetParams[0]
		offset, _ = strconv.Atoi(offsetStr)
	}

	start := time.Now().UnixNano()

	result := Select(limit, offset)

	timer := (time.Now().UnixNano() - start) / 1000

	for i, element := range result {
		fmt.Fprintf(rw, "Entry: %+v %v\n", element, i)
	}

	log.Printf("select time: %d mcs", timer)
}

func readBook(pageNumber uint32) {
	page := engine.GetPage(pageNumber)

	if page == nil {
		return
	}

	var pos int64 = 0

	for i := 0; i < 10000; i++ {
		entity := page.ReadDataRecord(pos)

		if entity == nil {
			break
		}

		pos = pos + int64(entity.Locator.Size)

		// build indexes
		pkIdx.Add(entity.Record, entity.Locator, entity.Record.ID)
		ratingIdx.Add(float64(entity.Record.Sort), entity.Record.ID)

		// add to map
		records[entity.Record.ID] = entity
	}
}

var writePrimary = 1

func writeBook(num uint32) {
	page := &engine.Page{
		Number: int(num),
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

		min := 10.1
		max := 104044.12
		r := min + rand.Float64()*(max-min)
		dataRecord.Sort = float32(r)

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
