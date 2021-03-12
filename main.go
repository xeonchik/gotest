package main

import (
	"fmt"
	"godoc/engine"
	"godoc/table"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"sort"
	"strconv"
	"time"
)

var tbl = engine.InitTableSpace("default")

func main() {
	go func() {
		err := http.ListenAndServe("localhost:6060", nil)
		if err == nil {
			log.Println("Started debug at localhost:6060")
		}
	}()

	fmt.Println("Hello! Starting GoDoc!")
	start := time.Now().UnixNano()

	engine.OpenBook()
	engine.PreloadBookPages()

	// create index by sort
	var indexer = func(idx interface{}, record *engine.DataRecord) {
		index := idx.(*engine.FloatIndex)
		index.Add(float64(record.Sort), record.ID)
	}
	err := tbl.AddIndexer("sortIdx", engine.CreateFloatIndex(), indexer)
	if err != nil {
		panic(err)
	}

	// index citiesIdx
	var indexerCities = func(idx interface{}, record *engine.DataRecord) {
		index := idx.(*engine.MultiIndex)

		for _, city := range record.Cities {
			index.Add(int(city.Value), record.ID)
		}
	}
	err = tbl.AddIndexer("citiesIdx", engine.CreateMulti(), indexerCities)

	if err != nil {
		panic(err)
	}

	readBook(tbl)

	for i := uint32(0); i < 100; i++ {
		//writeBook(i)
	}

	timer := (time.Now().UnixNano() - start) / 1000
	log.Printf("init time for %d entities: %d mcs", tbl.PrimaryIndex.Tree.Len(), timer)

	StartServer()
	engine.Close()
}

var ratingIdx = engine.CreateFloatIndex()

func Sort(limit int, offset int) []*engine.DataRecord {
	result := make([]*engine.DataRecord, 0)

	idx := tbl.Indexes["sortIdx"]
	index := idx.Index.(*engine.FloatIndex)

	index.Tree.Ascend(nil, func(item interface{}) bool {
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

// SelectByCity city IN (cities)
func SelectByCity(city int, limit int, offset int) *table.TemporaryDataSet {
	idx := tbl.Indexes["citiesIdx"].Index
	index := idx.(*engine.MultiIndex)

	// item of multi
	mItem := index.Get(city)
	leng := len(mItem.Keys)

	idxSort := tbl.Indexes["sortIdx"].Index
	sortIndex := idxSort.(*engine.FloatIndex)

	tempTable := table.CreateTempTable(0)

	sortIndex.Tree.Ascend(nil, func(item interface{}) bool {
		if limit == 0 {
			return false
		}

		it := item.(*engine.FloatItem)
		key := int(it.Key)

		i := sort.Search(leng, func(i int) bool {
			return mItem.Keys[i] >= key
		})

		if i < leng && key == mItem.Keys[i] {
			if offset > 0 {
				offset--
				return true
			}

			pk := tbl.PK(uint64(key))
			tempTable.AddPK(pk)
			limit--
		}
		return true
	})

	return tempTable
}

// SelectWithConditions / SelectWithConditions id > n, order by sort
func SelectWithConditions(limit int, offset int) *table.TemporaryDataSet {
	tempTable := table.CreateTempTable(ratingIdx.Tree.Len())
	n := uint64(10000)

	idx := tbl.Indexes["sortIdx"]
	index := idx.Index.(*engine.FloatIndex)

	index.Tree.Ascend(nil, func(item interface{}) bool {
		if limit == 0 {
			return false
		}

		it := item.(*engine.FloatItem)

		if it.Key <= n {
			return true
		}

		tempTable.AddPK(tbl.PK(it.Key))

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
			city := &engine.City{Value: int32(rand.Intn(1200-1000) + 1000)}
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
