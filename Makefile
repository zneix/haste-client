DESTDIR ?= "/usr/local"
default: build

install: build
	echo ${DESTDIR}
	install -m 755 -t "${DESTDIR}/bin" -D haste

build: haste

haste: main.go
	go build -o haste

run: build
	./haste
