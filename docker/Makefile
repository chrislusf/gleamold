BINARY = gleamold
SOURCE_DIR = github.com/chrislusf/gleamold/distributed/gleamold

all: gleamold

.PHONY : clean deps docker run_docker

clean:
	go clean $(SOURCE_DIR)
	rm -f $(BINARY)

deps:
	go get $(GO_FLAGS) -d $(SOURCE_DIR)

gleamold: deps
	GOOS=linux GOARCH=amd64 go build -o $(BINARY) $(SOURCE_DIR)

docker: gleamold
	cp ../script/MessagePack.lua .
	docker build -t chrislusf/gleamold .
	rm MessagePack.lua
	rm -f $(BINARY)

run_docker: docker
	docker-compose up
