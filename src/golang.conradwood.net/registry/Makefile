.PHONY:	client server
all:
	make client server
server:
	cd server && go install ${LDFLAGS} registry-server.go `ls -1 *.go | grep -v registry-server.go`
client:
	cd client && go install ${LDFLAGS} registry-client.go

