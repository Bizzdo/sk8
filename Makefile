IMAGE=sk8

build: $(IMAGE)

$(IMAGE): src/*.go Makefile
	go build -o $(IMAGE) ./src

clean:
	rm ./$(IMAGE)

fix:
	go get -v

install: $(IMAGE)
	mkdir -p ~/bin
	cp $(IMAGE) ~/bin