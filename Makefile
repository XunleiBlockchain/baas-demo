BUILD_TAGS?=

default: clean build

clean:
	rm -rf baas-demo

build:
	go build -tags '$(BUILD_TAGS)' -o baas-demo .