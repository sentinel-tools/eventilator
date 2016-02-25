VERSION = $(shell cat .version)


test:
	go test
	go vet


all: eventilator reconfigurator sentinel-scriptify
	@echo all binaries built

release: eventilator reconfigurator sentinel-scriptify
	mkdir -p dist/usr/sbin
	mv eventilator dist/usr/sbin
	mv reconfigurator dist/usr/sbin
	mv sentinel-scriptify dist/usr/sbin
	cd dist && tar -cvzf ../eventilator-${VERSION}.tar.gz usr/ && cd ..
	echo Version=${VERSION}
	ls -lh eventilator-${VERSION}.tar.gz
	@ghr  --username sentinel-tools --token ${GITHUB_TOKEN} ${VERSION} eventilator-${VERSION}.tar.gz

eventilator:
	@echo Building eventilator
	go build 

reconfigurator: 
	@echo Building reconfigurator
	go build -o reconfigurator

sentinel-scriptify:
	@echo Building sentinel-scriptify
	go build -o sentinel-scriptify ./utils/sentinel-scriptify 

clean:
	@echo cleaning up
	rm -rf dist/usr/
