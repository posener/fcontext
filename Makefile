bench:
	go test  -benchmem -bench . > results.txt
	$(MAKE) parse-bench

bench-full:
	go test  -benchmem -benchtime=20s -bench . > results.txt
	$(MAKE) parse-bench

bench-parse:
	cat results.txt | go run ./tools/parsebench/parsebench.go > results.csv
