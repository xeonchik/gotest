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

func main() {
	fmt.Println("Hello, World")

	engine.OpenBook()
	engine.PreloadBookPages()

	start := time.Now().UnixNano()

	for i := uint32(0); i < 300; i++ {
		//writeBook(i)
		readBook(i)
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

var records = make(map[uint64]*engine.DataEntity)
var pkIdx = engine.CreatePKIndex()
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
		record := records[it.Key]

		result = append(result, record.Record)
		limit--
		return true
	})

	return result
}

// city IN (cities)
func SelectByCity(cities []int, limit int, offset int) *table.TemporaryDataSet {
	tbl := table.CreateTempTable()

	for _, city := range cities {
		item := citiesIdx.Get(city)
		item.Keys.Tree.Walk(func(itemsWalker []interface{}) {
			for _, itm := range itemsWalker {
				flatItem := itm.(*engine.FlatItem)
				id := flatItem.Value
				entity := records[id]
				tbl.Add(entity.Record)
			}
		})
	}

	return tbl
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

		record := records[it.Key]
		tempTable.Add(record.Record)

		limit--
		return true
	})

	return tempTable
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

		for _, city := range entity.Record.Cities {
			citiesIdx.Add(int(city.Value), entity.Record.ID)
		}

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
