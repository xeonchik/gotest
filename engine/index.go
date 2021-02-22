package engine

import (
	"fmt"
	"github.com/tidwall/btree"
)

type FlatIndex struct {
	tree *btree.BTree
}

type FlatItem struct {
	value int
}

func byFlatVal(a, b interface{}) bool {
	i1, i2 := a.(*FlatItem), b.(*FlatItem)
	return i1.value < i2.value
}

type MultiIndex struct {
	tree *btree.BTree
}

type MultiItem struct {
	keys     *FlatIndex
	idxValue int
}

func byIdxVal(a, b interface{}) bool {
	i1, i2 := a.(*MultiItem), b.(*MultiItem)
	return i1.idxValue < i2.idxValue
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

func (idx *MultiIndex) Add(indexValue int, key int) {
	// check that exists
	item := idx.tree.Get(&MultiItem{
		idxValue: indexValue,
	}).(*MultiItem)

	fmt.Println(item)
}

func (idx *PKIndex) Get(key int) *DataRecord {
	return idx.tree.Get(&PKItem{
		key: key,
	}).(*PKItem).value
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
