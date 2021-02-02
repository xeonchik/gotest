package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"github.com/tidwall/btree"
	"time"
)

type Item struct {
	Key int
	Value string
	Value2 int
}

type MultiItem struct {
	Value2 int
	Keys []int
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

	rand.Seed(time.Now().UTC().UnixNano())

	for i := 1; i <= 1000000; i++ {
		rnd := rand.Intn(100000 - 10) + 10
		item := &Item{
			Key:   i,
			Value: fmt.Sprintf("test%d", rnd),
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

	rnd := rand.Intn(1000 - 10) + 10
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

	//http.HandleFunc("/push", pushHandler)
	//http.HandleFunc("/list", listHandler)
	//log.Fatal(http.ListenAndServe(":8090", nil))
}

type testStruct struct {
	Id int
	Sectors []int
	Cities []int
	Active bool
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