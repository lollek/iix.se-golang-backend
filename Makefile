GO	=	go
SRC	=	$(wilcard *.go)

all:	build

build: $(SRC)
	$(GO) build
	@ctags -R .

docker: $(SRC)
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o iix.se-golang-backend .
	sudo docker build -t loll3k/iix.se-golang-backend .
