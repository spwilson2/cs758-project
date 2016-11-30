all: build

build: main.go
	go build main.go

clean:
	rm -f *.txt main
