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
	fmt.Fprintf(ctx, "Entry: %+v\n", rec)
}

func countHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Count items: %d\n", tbl.PrimaryIndex.Tree.Len())
}

func byCitiesHandler(ctx *fasthttp.RequestCtx) {
	var limit, offset int
	limit, offset = GetLimitOffsetFromURL(ctx)

	var cities = []int{30000}

	result := SelectByCity(cities, limit, offset)

	for _, element := range result.Records {
		PrintRecord(ctx, element)
	}
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

	for i, element := range result.Records {
		fmt.Fprintf(ctx, "Entry: %+v %v\n", element, i)
	}

	log.Printf("select time: %d mcs", timer)
}
