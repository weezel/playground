include ../../Makefile

name		?= huhhuh

show_mem_profile:
	go tool pprof -http=: mem.prof

show_cpu_profile:
	go tool pprof -http=: cpu.prof

build: _build
