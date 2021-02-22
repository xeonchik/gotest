package engine

import (
	"fmt"
	"github.com/tidwall/btree"
)

type FlatIndex struct {
	Tree *btree.BTree
}

type FlatItem struct {
	Value int
}

func byFlatVal(a, b interface{}) bool {
	i1, i2 := a.(*FlatItem), b.(*FlatItem)
	return i1.Value < i2.Value
}

type MultiIndex struct {
	tree *btree.BTree
}

type MultiItem struct {
	Keys     *FlatIndex
	IdxValue int
}

func byIdxVal(a, b interface{}) bool {
	i1, i2 := a.(*MultiItem), b.(*MultiItem)
	return i1.IdxValue < i2.IdxValue
}

type PKIndex struct {
	tree *btree.BTree
	name string
}

type PKItem struct {
	value *DataRecord
	key   int
}

func byKey(a, b interface{}) bool {
	i1, i2 := a.(*PKItem), b.(*PKItem)
	return i1.key < i2.key
}

/// indexValue - value
/// key - PK Link
func (idx *MultiIndex) Add(indexValue int, key int) {
	// check that exists
	item := idx.tree.Get(&MultiItem{
		IdxValue: indexValue,
	})

	if item == nil {
		item := &MultiItem{
			Keys:     &FlatIndex{Tree: btree.New(byFlatVal)},
			IdxValue: indexValue,
		}
		item.Keys.Tree.Set(&FlatItem{Value: key})
		idx.tree.Set(item)
	} else {
		it := item.(*MultiItem)
		flatItem := &FlatItem{Value: key}
		it.Keys.Tree.Set(flatItem)
	}
}

func (idx *MultiIndex) GetTree() *btree.BTree {
	return idx.tree
}

func (idx *MultiIndex) Print() {
	point := &MultiItem{
		IdxValue: 0,
	}

	idx.tree.Ascend(point, func(item interface{}) bool {
		it := item.(*MultiItem)
		fmt.Println("Value: ", it.IdxValue)
		it.Keys.Tree.Ascend(nil, func(item2 interface{}) bool {
			fmt.Print(item2.(*FlatItem).Value, " ")
			return true
		})
		return true
	})
}

func (idx *PKIndex) Get(key int) *DataRecord {
	StatsObj.Hits += 1
	return idx.tree.Get(&PKItem{
		key: key,
	}).(*PKItem).value
}

func (idx *MultiIndex) Get(key int) *MultiItem {
	StatsObj.Hits += 1
	return idx.tree.Get(&MultiItem{
		IdxValue: key,
	}).(*MultiItem)
}

func (idx *PKIndex) Add(record *DataRecord, key int) {
	item := PKItem{
		value: record,
		key:   key,
	}

	if idx.tree.Get(&item) != nil {
		panic("PK already exists")
	}

	idx.tree.Set(&item)
}

func (idx *PKIndex) Print() {
	point := &PKItem{
		key: 0,
	}

	idx.tree.Ascend(point, func(item interface{}) bool {
		it := item.(*PKItem)
		fmt.Println(it.value.Primary)
		return true
	})
}

func CreatePKIndex(name string) *PKIndex {
	idx := &PKIndex{
		tree: btree.New(byKey),
		name: name,
	}
	return idx
}

func CreateMulti() *MultiIndex {
	return &MultiIndex{
		tree: btree.New(byIdxVal),
	}
}

func (idx *MultiIndex) AddArray(arr []int, key int) {
	for _, element := range arr {
		idx.Add(element, key)
	}
}
