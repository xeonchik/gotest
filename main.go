package main

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/btree"
	"godoc/engine"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Item struct {
	Key    int
	Value  string
	Value2 int
}

type MultiItem struct {
	Value2 int
	Keys   []int
}

func byKeys(a, b interface{}) bool {
	i1, i2 := a.(*Item), b.(*Item)
	return i1.Key < i2.Key
}

func byValue2(a, b interface{}) bool {
	i1, i2 := a.(*MultiItem), b.(*MultiItem)
	return i1.Value2 < i2.Value2
}

var primary = btree.New(byKeys)
var valuesIdx = btree.New(byValue2)

func main() {
	fmt.Println("Hello, World")

	engine.OpenBook()

	for i := 0; i < 100; i++ {
		readBook(i)
	}

	//pkIdx.Print()

	fmt.Println(pkIdx.Get(200))

	engine.Close()
	//readBook(0)
}

var pkIdx = engine.CreatePKIndex("123")

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

func IndexTest() {
	for i := 1; i <= 1000; i++ {
		rnd := rand.Intn(1000-10) + 10
		item := &Item{
			Key:    i,
			Value:  fmt.Sprintf("test%d", rnd),
			Value2: rnd,
		}
		primary.Set(item)

		val2 := rnd
		valItem := valuesIdx.Get(&MultiItem{
			Value2: val2,
		})

		var keys []int

		if valItem != nil {
			multi := valItem.(*MultiItem)
			keys = append(multi.Keys, i)
			multi.Keys = keys
			valuesIdx.Set(multi)
		} else {
			keys = append(keys, i)
			item := &MultiItem{
				Value2: val2,
				Keys:   keys,
			}
			valuesIdx.Set(item)
		}
	}

	primary.Descend(&Item{
		Key: 10,
	}, func(item interface{}) bool {
		it := item.(*Item)
		fmt.Printf("%d %s\n", it.Key, it.Value)
		return true
	})

	fmt.Println("Length: ", primary.Len())
	fmt.Println("Length Value2 idx: ", valuesIdx.Len())

	rnd := rand.Intn(1000-10) + 10
	item := primary.Get(&Item{Key: rnd})
	if item != nil {
		it := item.(*Item)
		fmt.Println(it.Value)
	}

	valItem := valuesIdx.Get(&MultiItem{
		Value2: 11,
	})
	if valItem != nil {
		it := valItem.(*MultiItem)
		fmt.Println(it.Keys)
	}
}
