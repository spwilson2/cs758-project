.PHONY:all
all: ../generated/bmain ../generated/amain

../generated/%.go:../src/main.go
	python -c 'import build; print build.buildProject()'

../generated/%:../generated/%.go
	go build -o $@ $<
