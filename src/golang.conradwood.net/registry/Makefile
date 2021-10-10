.PHONY:	client server
server:
	cd server && go install registry-server.go `ls -1 *.go | grep -v registry-server.go`
client:
	cd client && go install registry-client.go

all:
	make client server
