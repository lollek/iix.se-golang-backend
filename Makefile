GO	=	go
SRC	=	$(wilcard *.go)

all:	build

build: $(SRC)
	$(GO) build
	@ctags -R .
.PHONY: build

docker: $(SRC)
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o iix.se-golang-backend .
	sudo docker build -t loll3k/iix.se-golang-backend .
.PHONY: docker

clean:
	$(RM) iix.se-golang-backend
.PHONY: clean
