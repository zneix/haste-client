DESTDIR ?= /usr/local
install: build
	echo ${DESTDIR}
	install -Dm755 haste ${DESTDIR}/bin/haste

build: haste

haste: main.go
	go build -o haste

run: build
	./haste
