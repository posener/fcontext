bench:
	go test  -benchmem -bench . > results.txt
	$(MAKE) bench-parse

bench-full:
	go test  -benchmem -benchtime=20s -timeout=1h -bench . > results.txt
	$(MAKE) bench-parse

bench-parse:
	cat results.txt | go run ./internal/tools/parsebench/parsebench.go > results.csv
