.PHONY: clean
clean:
	rm -f ../../app
	rm -f ../../app.key
	rm -f ../../app.toml
	rm -f ../../nacos.toml

.PHONY: build
build:
	go build -ldflags "-s -w" -o ../../app ../main/.

.PHONE: localdev
localdev:
	cp ../../config/localdev/app.key ../../app.key
	cp ../../config/localdev/app.toml ../../app.toml
	cp ../../config/localdev/nacos.toml ../../nacos.toml