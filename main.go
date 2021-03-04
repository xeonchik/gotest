package main

import (
	"encoding/json"
	"fmt"
	"godoc/engine"
	"godoc/table"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var tbl = engine.InitTableSpace("default")

func main() {
	fmt.Println("Hello, World")

	engine.OpenBook()
	engine.PreloadBookPages()

	// create index by sort
	var indexer = func(idx interface{}, record *engine.DataRecord) {
		index := idx.(*engine.FloatIndex)
		index.Add(float64(record.Sort), record.ID)
	}
	var idx = engine.CreateFloatIndex()
	idx.A()

	tbl.AddIndexer("sortIdx", idx, indexer)

	readBook(tbl)

	start := time.Now().UnixNano()

	for i := uint32(0); i < 100; i++ {
		//writeBook(i)
	}

	timer := (time.Now().UnixNano() - start) / 1000
	log.Printf("init time: %d mcs", timer)

	http.HandleFunc("/select", selectHandler)
	http.HandleFunc("/select-cond", selectCondHandler)
	http.HandleFunc("/select-by-city", byCitiesHandler)
	http.HandleFunc("/sort", sortHandler)
	http.HandleFunc("/count", countHandler)

	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		panic(err)
	}

	engine.Close()
}

var ratingIdx = engine.CreateFloatIndex()
var citiesIdx = engine.CreateMulti()

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
		result = append(result, tbl.GetByPK(it.Key))

		limit--
		return true
	})

	return result
}

// city IN (cities)
func SelectByCity(cities []int, limit int, offset int) *table.TemporaryDataSet {
	tempTable := table.CreateTempTable()

	for _, city := range cities {
		item := citiesIdx.Get(city)
		item.Keys.Tree.Walk(func(itemsWalker []interface{}) {
			for _, itm := range itemsWalker {
				flatItem := itm.(*engine.FlatItem)
				id := flatItem.Value
				tempTable.Add(tbl.GetByPK(id))
			}
		})
	}

	return tempTable
}

// id > n, order by sort
func SelectWithConditions(limit int, offset int) *table.TemporaryDataSet {
	tempTable := table.CreateTempTable()
	n := uint64(10000)

	ratingIdx.Tree.Ascend(nil, func(item interface{}) bool {
		if limit == 0 {
			return false
		}

		it := item.(*engine.FloatItem)

		if it.Key <= n {
			return false
		}

		tempTable.Add(tbl.GetByPK(it.Key))

		limit--
		return true
	})

	return tempTable
}

func Select(limit int, offset int) []*engine.DataRecord {
	result := make([]*engine.DataRecord, 0)

	tbl.PrimaryIndex.Tree.Ascend(nil, func(item interface{}) bool {
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
	fmt.Fprintf(rw, "Count items: %d\n", tbl.PrimaryIndex.Tree.Len())
}

func byCitiesHandler(rw http.ResponseWriter, req *http.Request) {
	var limit, offset int
	limit, offset = GetLimitOffsetFromURL(req)

	var cities = []int{30000}

	result := SelectByCity(cities, limit, offset)

	for _, element := range result.Records {
		PrintRecord(rw, element)
	}
}

func sortHandler(rw http.ResponseWriter, req *http.Request) {
	var limit, offset int
	limit, offset = GetLimitOffsetFromURL(req)

	start := time.Now().UnixNano()

	result := Sort(limit, offset)

	timer := (time.Now().UnixNano() - start) / 1000

	for _, element := range result {
		PrintRecord(rw, element)
	}

	log.Printf("sort time: %d mcs", timer)
}

func selectHandler(rw http.ResponseWriter, req *http.Request) {
	var limit, offset int
	limit, offset = GetLimitOffsetFromURL(req)

	start := time.Now().UnixNano()

	result := Select(limit, offset)

	timer := (time.Now().UnixNano() - start) / 1000

	for _, element := range result {
		PrintRecord(rw, element)
	}

	log.Printf("select time: %d mcs", timer)
}

func selectCondHandler(rw http.ResponseWriter, req *http.Request) {
	var limit, offset int
	limit, offset = GetLimitOffsetFromURL(req)

	start := time.Now().UnixNano()

	result := SelectWithConditions(limit, offset)

	timer := (time.Now().UnixNano() - start) / 1000

	for i, element := range result.Records {
		fmt.Fprintf(rw, "Entry: %+v %v\n", element, i)
	}

	log.Printf("select time: %d mcs", timer)
}

func readBook(table *engine.Table) {

	for i := uint32(0); i < engine.GetPagesCount(); i++ {
		page := engine.GetPage(i)
		table.ReadPageRecords(page)
	}

	//	// build indexes
	//	ratingIdx.Add(float64(entity.Record.Sort), entity.Record.ID)
	//
	//	for _, city := range entity.Record.Cities {
	//		citiesIdx.Add(int(city.Value), entity.Record.ID)
	//	}
	//}
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

func PrintRecord(rw http.ResponseWriter, rec *engine.DataRecord) {
	fmt.Fprintf(rw, "Entry: %+v\n", rec)
}

func GetLimitOffsetFromURL(req *http.Request) (int, int) {
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

	return limit, offset
}
