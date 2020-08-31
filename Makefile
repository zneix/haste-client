build:
	go build
	mv ./haste-client ./haste

run: build
	./haste