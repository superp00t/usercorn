.DEFAULT_GOAL := build
.PHONY: get test cov bench deps usercorn

all: usercorn

clean:
	rm -f usercorn

build: get all

# dependency targets
DEST = $(shell mkdir -p deps/build; cd deps && pwd)
FIXRPATH := touch
LIBEXT := so
OS := $(shell uname -s)
ARCH := $(shell uname -m)

ifeq "$(OS)" "Darwin"
	LIBEXT = dylib
	FIXRPATH = @install_name_tool \
		-add_rpath @executable_path/lib \
		-add_rpath @executable_path/deps/lib \
		-change libunicorn.dylib @rpath/libunicorn.dylib \
		-change libunicorn.1.dylib @rpath/libunicorn.1.dylib \
		-change libunicorn.2.dylib @rpath/libunicorn.2.dylib \
		-change libcapstone.dylib @rpath/libcapstone.dylib \
		-change libcapstone.3.dylib @rpath/libcapstone.3.dylib \
		-change libcapstone.4.dylib @rpath/libcapstone.4.dylib \
		-change libkeystone.dylib @rpath/libkeystone.dylib \
		-change libkeystone.0.dylib @rpath/libkeystone.0.dylib \
		-change libkeystone.1.dylib @rpath/libkeystone.1.dylib
endif

# figure out if we can download Go

deps/lib/libunicorn.1.$(LIBEXT):
	cd deps/build && \
	git clone https://github.com/unicorn-engine/unicorn.git && git --git-dir unicorn fetch; \
	cd unicorn && git clean -fdx && git reset --hard origin/master && cmake -DCMAKE_INSTALL_PREFIX=$(DEST) . && \
	make && make install

deps/lib/libcapstone.3.$(LIBEXT):
	cd deps/build && \
	git clone https://github.com/aquynh/capstone.git && git --git-dir capstone pull; \
	cd capstone && git clean -fdx && git reset --hard origin/master; \
	mkdir build && cd build && cmake -DCAPSTONE_BUILD_STATIC=OFF -DCMAKE_INSTALL_PREFIX=$(DEST) -DCMAKE_BUILD_TYPE=RELEASE .. && \
	make -j2 install

deps/lib/libkeystone.0.$(LIBEXT):
	cd deps/build && \
	git clone https://github.com/keystone-engine/keystone.git && git --git-dir keystone pull; \
	cd keystone; git clean -fdx && git reset --hard origin/master; mkdir build && cd build && \
	cmake -DCMAKE_INSTALL_PREFIX=$(DEST) -DCMAKE_BUILD_TYPE=Release -DBUILD_SHARED_LIBS=ON -DLLVM_TARGETS_TO_BUILD="all" -G "Unix Makefiles" .. && \
	make -j2 install

deps: deps/lib/libunicorn.1.$(LIBEXT) deps/lib/libcapstone.3.$(LIBEXT) deps/lib/libkeystone.0.$(LIBEXT) deps/$(GODIR)

# Go executable targets

export CGO_CFLAGS = -I$(DEST)/include
export CGO_LDFLAGS = -L$(DEST)/lib

GOBUILD := go build
SHELL := env LD_LIBRARY_PATH=$(LD_LIBRARY_PATH):$(DEST)/lib DYLD_LIBRARY_PATH=$(DYLD_LIBRARY_PATH):$(DEST)/lib PATH=$(PATH) /bin/bash

# TODO: more DRY?
usercorn:
	rm -f usercorn
	$(GOBUILD) -o usercorn github.com/superp00t/usercorn/go/cmd/main
	$(FIXRPATH) usercorn

