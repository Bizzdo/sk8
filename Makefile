IMAGE=sk8

build: $(IMAGE)

$(IMAGE): *.go Makefile
	go build -o ./$(IMAGE) .

clean:
	rm ./$(IMAGE)

fix:
	go get -v

install: $(IMAGE)
	cp $(IMAGE) $(GOPATH)/bin