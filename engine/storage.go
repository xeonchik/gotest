package engine

import (
	"github.com/golang/protobuf/proto"
	pb "godoc/proto/dist/proto"
	"io/ioutil"
)

func ReadIndexFromDisk(name string) (*PKIndex, error) {
	in, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}

	indexStore := &pb.PKIndexStore{}
	err = proto.Unmarshal(in, indexStore)
	if err != nil {
		return nil, err
	}

	pkIdx := CreatePKIndex()

	for _, item := range indexStore.Items {
		record := &DataRecord{
			ID: uint64(item.Primary),
		}
		locator := &DataRowLocator{
			Page:   int(item.PageNumber),
			Offset: item.Offset,
			Size:   int(item.Size),
			Loaded: false,
		}
		pkIdx.Load(record, *locator, record.ID)
	}

	return pkIdx, nil
}

func FlushIndexToDisk(index *PKIndex, name string) {
	indexStore := &pb.PKIndexStore{}

	index.Tree.Ascend(nil, func(it interface{}) bool {
		item := it.(*PKItem)

		indexItem := &pb.PKIndexItem{
			Primary:    int32(item.PrimaryKey),
			PageNumber: int32(item.Locator.Page),
			Offset:     item.Locator.Offset,
			Size:       int32(item.Locator.Size),
		}
		indexStore.Items = append(indexStore.Items, indexItem)
		return true
	})

	out, err := proto.Marshal(indexStore)

	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(name, out, 0660)

	if err != nil {
		panic(err)
	}
}
