.PHONY: build_docker run_docker build clean

build:
	go build -o bitcask-server

build_docker:
	docker-compose up -d --build

run_docker:
	docker run -it -d --name bitcask-server -p 8090:8090 -v ./bitcaskDB:/app/bitcaskDB bitcask-server

clean:
	rm -rf bitcask-server

dist_clean: clean
	docker stop bitcask-server
	docker rm bitcask-server
	rm -rf ./bitcaskDB
	rm -rf ./tests/*.json
