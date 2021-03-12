package main

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"godoc/engine"
	"log"
	"time"
)

func selectHandler(ctx *fasthttp.RequestCtx) {

	var limit, offset int
	limit, offset = GetLimitOffsetFromURL(ctx)

	start := time.Now().UnixNano()

	result := Select(limit, offset)

	timer := (time.Now().UnixNano() - start) / 1000

	for _, element := range result {
		PrintRecord(ctx, element)
	}

	log.Printf("select time: %d mcs", timer)

	fmt.Fprintf(ctx, "this is the first part of body\n")
}

func barHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "this is the first part of body\n")
}

func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/select":
		selectHandler(ctx)
	case "/select-cond":
		selectCondHandler(ctx)
	case "/select-by-city":
		byCitiesHandler(ctx)
	case "/sort":
		sortHandler(ctx)
	case "/count":
		countHandler(ctx)
	case "/stats":
		statsHandler(ctx)
	default:
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
	}
}

func StartServer() {
	log.Println("Listening HTTP at *:8090")
	fasthttp.ListenAndServe(":8090", fastHTTPHandler)
}

func GetLimitOffsetFromURL(ctx *fasthttp.RequestCtx) (int, int) {
	var limit, offset int

	limit, _ = ctx.URI().QueryArgs().GetUint("limit")
	if limit < 0 {
		limit = 1000
	}

	offset, _ = ctx.URI().QueryArgs().GetUint("offset")
	if offset < 0 {
		offset = 0
	}

	return limit, offset
}

func PrintRecord(ctx *fasthttp.RequestCtx, rec *engine.DataRecord) {
	var cities []int32

	for _, city := range rec.Cities {
		cities = append(cities, city.Value)
	}

	fmt.Fprintf(ctx, "Entry: ID: %+v Cities: %+v Sort: %+v Active: %v\n", rec.ID, cities, rec.Sort, rec.Active)
}

func countHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Count items: %d\n", tbl.PrimaryIndex.Tree.Len())
}

func statsHandler(ctx *fasthttp.RequestCtx) {
	for name, idx := range tbl.Indexes {
		var index engine.IndexType = idx.Index.(engine.IndexType)
		fmt.Fprintf(ctx, "Index %s size is %d bytes\n", name, index.GetSize())
	}
}

func byCitiesHandler(ctx *fasthttp.RequestCtx) {
	var limit, offset int
	limit, offset = GetLimitOffsetFromURL(ctx)

	city, err := ctx.QueryArgs().GetUint("city")

	if err != nil {
		ctx.Error("Undefined parameter 'city'.", fasthttp.StatusBadRequest)
		return
	}

	start := time.Now().UnixNano()

	result := SelectByCity(city, limit, offset)

	timer := (time.Now().UnixNano() - start) / 1000

	for {
		pk := result.Read()
		if pk == nil {
			break
		}
		PrintRecord(ctx, pk.Record)
	}

	log.Printf("by cities time: %d mcs", timer)
}

func sortHandler(ctx *fasthttp.RequestCtx) {
	var limit, offset int
	limit, offset = GetLimitOffsetFromURL(ctx)

	start := time.Now().UnixNano()

	result := Sort(limit, offset)

	timer := (time.Now().UnixNano() - start) / 1000

	for _, element := range result {
		PrintRecord(ctx, element)
	}

	log.Printf("sort time: %d mcs", timer)
}

func selectCondHandler(ctx *fasthttp.RequestCtx) {
	var limit, offset int
	limit, offset = GetLimitOffsetFromURL(ctx)

	start := time.Now().UnixNano()

	result := SelectWithConditions(limit, offset)

	timer := (time.Now().UnixNano() - start) / 1000

	for {
		pk := result.Read()
		if pk == nil {
			break
		}
		PrintRecord(ctx, pk.Record)
	}

	log.Printf("select time: %d mcs", timer)
}
