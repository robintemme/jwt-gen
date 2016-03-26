BINARY=jwt

SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go' -d 1)

GOOS=linux
GOARCH=arm
GOARM=7

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(SOURCES)
	go build -o ${BINARY} ${SOURCES}

.PHONY: clean
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

cross:
	make clean
	GOOS=${GOOS} GOARCH=${GOARCH} GOARM=${GOARM} make

