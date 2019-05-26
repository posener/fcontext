bench:
	go test  -benchmem -bench .

bench-full:
	go test  -benchmem -benchtime=20s -bench .
