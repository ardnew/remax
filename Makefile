# configuration
output  := dist
tarflag := -z
tarext  := tar.gz

# supported platforms
arch    := 386 amd64 arm arm64
os      := linux

# system commands
rm      := rm -rf
go      := go
mkdir   := mkdir -p
mv      := mv
cp      := cp -r
tar     := tar $(tarflag) -cf

# build identity
ident   := $(shell $(go) run . -version)
bin     := $(word 1, $(ident))
version := $(word 3, $(ident))
license := LICENSE
readme  := README.md

# target lists
target  := $(patsubst %,$(bin)-$(version)-%,$(arch))
dist    := $(patsubst %,$(output)/%.$(tarext),$(target))

.PHONY: all
all: $(dist)

.PHONY: clean
clean:
	$(rm) "$(output)"

$(arch):
	GOOS=$(os) GOARCH=$(@) $(go) build

$(target): $(bin)-$(version)-%: %
	@[ -d "$(output)/$(@)" ] || $(mkdir) "$(output)/$(@)"
	@$(mv) "$(bin)" "$(output)/$(@)"
	@$(cp) "$(license)" "$(readme)" "$(output)/$(@)"

$(dist): $(output)/%.$(tarext): %
	$(tar) "$(@)" -C "$(output)" "$(<)"

