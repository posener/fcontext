package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

var line = regexp.MustCompile(`^Benchmark\/([^\/]+)\/([^\/]+)\/(.*)-4\s+(\S+)\s+(\S+) ns\/op\s+(\S+) B\/op\s+(\S+) allocs\/op`)

type Benchmark struct {
	Cycles   int
	NsOp     float64
	BytesOp  float64
	AllocsOp float64
}

func main() {
	results := make(map[string]map[string]Benchmark)
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		m := line.FindStringSubmatch(s.Text())
		if len(m) == 0 {
			log.Println("Skip", s.Text())
			continue
		}
		var (
			ctx      = m[1]
			fn       = m[2]
			test     = m[3]
			cycles   = m[4]
			nsOp     = m[5]
			bytesOp  = m[6]
			allocsOp = m[7]
		)

		var (
			benchmark Benchmark
			err       error
		)
		benchmark.Cycles, err = strconv.Atoi(cycles)
		panicOnErr(err)
		benchmark.AllocsOp, err = strconv.ParseFloat(allocsOp, 64)
		panicOnErr(err)
		benchmark.NsOp, err = strconv.ParseFloat(nsOp, 64)
		panicOnErr(err)
		benchmark.BytesOp, err = strconv.ParseFloat(bytesOp, 64)
		panicOnErr(err)

		test = fn + "/" + test

		if results[test] == nil {
			results[test] = make(map[string]Benchmark)
		}
		results[test][ctx] = benchmark
	}

	fmt.Println("test,ns/op ratio,B/op ratio,allocs/op ratio")
	for test, compare := range results {
		fmt.Printf("%s,%.2f,%.2f,%.2f\n",
			test,
			saveDivide(compare["stdctx"].NsOp, compare["fcontext"].NsOp),
			saveDivide(compare["stdctx"].BytesOp, compare["fcontext"].BytesOp),
			saveDivide(compare["stdctx"].AllocsOp, compare["fcontext"].AllocsOp),
		)
	}
}

func saveDivide(a, b float64) float64 {
	if b == 0 {
		return 0
	}
	return a / b
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
