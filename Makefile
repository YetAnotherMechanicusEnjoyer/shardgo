CMD_FOLDER	=	cmd
MAIN_FOLDERS	=	cli	\
								server

build: $(MAIN_FOLDERS)

$(MAIN_FOLDERS):
	go build -o bin/$@ ./$(CMD_FOLDER)/$@

run:
	go run ./$(CMD_FOLDER)/server

test:
	go test -v ./...

clean:
	rm -rf bin/

docker:
	docker build -t shardgo .

.PHONY: build run test clean docker
