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

	for i := 0; i < 10; i++ {
		//writeBook(i)
		readBook(i)
	}

	//pkIdx.Print()

	fmt.Println(pkIdx.Get(200))
	multiItem := multiIdx.Get(412)

	fmt.Println("Items in multi: ", multiIdx.GetTree().Len())

	multiItem.Keys.Tree.Ascend(nil, func(item interface{}) bool {
		it := item.(*engine.FlatItem)
		fmt.Println("PK: ", it.Value)
		return true
	})

	engine.Close()
	//readBook(0)
}

var pkIdx = engine.CreatePKIndex("123")
var multiIdx = engine.CreateMulti()
var idxMultiCities = engine.CreateMulti()

func readBook(num int) {
	page := engine.ReadPage(num)

	var pos int64 = 0

	for i := 0; i < 10000; i++ {
		_rec, _pos := page.ReadDataRecord(pos)
		pos = _pos

		if _rec == nil {
			break
		}

		pkIdx.Add(_rec, _rec.Primary)
		val, _ := strconv.Atoi(_rec.Data)

		multiIdx.Add(val, _rec.Primary)

		idxMultiCities.AddArray(_rec.Cities, _rec.Primary)

		fmt.Println(_rec)
	}

	//engine.FlushIndexToDisk(pkIdx, "primary.idx")
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
