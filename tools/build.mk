TARGETS:= amain bmain
TARGETS:= $(addprefix ../generated/,$(TARGETS))
GEN_SRCS:= $(addsuffix .go,$(TARGETS))

.PHONY:all
all: gen build

.PHONY:gen
gen: $(GEN_SRCS)
../generated/%.go:../src/main.go
	python -c 'import build; print build.generate("$(notdir $@)")'

.PHONY:build
../generated/%:../generated/%.go
	python -c 'import build; print build.build("$(notdir $@)")'
