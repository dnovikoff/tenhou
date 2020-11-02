protoc_version := 3.13.0
protoc_dir := gobin/protoc/v$(protoc_version)
protoc_bin := $(protoc_dir)/bin/protoc
protoc_gen_go := gobin/protoc-gen-go

ifeq ($(OS),Windows_NT)
    protoc_bin := $(protoc_bin).exe
    ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
        protocarch := win64
    else
        protocarch := win32
    endif
else
    uname_s := $(shell uname)
    uname_m := $(shell uname -m)
    ifeq ($(uname_s),Linux)
        pos := linux
    else ifeq ($(uname_s),Darwin)
        pos := osx
    else
        $(error Unknown os $(uname_s))
    endif
    ifeq ($(uname_m),x86_64)
       protocarch := $(pos)-x86_64
    else
       protocarch := $(pos)-x86_32
    endif
endif

protoc_include := $(protoc_dir)/include
protoc_zip_name := protoc-$(protoc_version)-$(protocarch).zip
protoc_url := https://github.com/google/protobuf/releases/download/v$(protoc_version)/$(protoc_zip_name)
protoc_zip_output := gobin/$(protoc_zip_name)

$(protoc_bin):
	rm -rf $(protoc_dir)
	mkdir -p $(protoc_dir)
	curl --retry 5 -L $(protoc_url) -o $(protoc_zip_output)
	unzip $(protoc_zip_output) -d $(protoc_dir)
	rm $(protoc_zip_output)
	touch $(protoc_bin) # override time from zip

# Ignore timestamp
protoc-bin: | $(protoc_bin)

$(protoc_gen_go): $(protoc_bin)
	go build -mod vendor -o $@ ./vendor/github.com/golang/protobuf/protoc-gen-go

protoc_cmd := $(protoc_bin) -I $(protoc_include) -I ./proto
protoc_go_cmd := $(protoc_cmd) --plugin=protoc-gen-go=$(protoc_gen_go)
