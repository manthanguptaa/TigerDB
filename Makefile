build:
	go build -o bin/TigerDB

run: build
	./bin/TigerDB

runfollower: build
	./bin/TigerDB --listenaddr :4000 --leaderaddr :3000