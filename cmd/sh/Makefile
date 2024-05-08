.PHONY: clean
clean:
	rm -f ../../app
	rm -f ../../app.log
	rm -f ../../request.log
	rm -f ../../app.key
	rm -f ../../app.toml
	rm -f ../../nacos.toml

.PHONE: localdev
localdev:
	cp ../../config/localdev/app.key ../../app.key
	cp ../../config/serverconf/app.toml ../../app.toml
	cp ../../config/localdev/nacos.toml ../../nacos.toml


CONF_DIR = ../../config
TARGET_WIN_DIR = ../../doc/download/windows-amd64
TARGET_WIN_ZIP = $(TARGET_WIN_DIR)/windows-amd64.zip


.PHONY: windows_build
windows_build:
	go env -w CGO_ENABLED=0 GOOS=windows GOARCH=amd64
	go build -ldflags "-s -w" -o $(TARGET_WIN_DIR)/app.exe ../main/.

.PHONY: windows_copy
windows_copy:
	cp $(CONF_DIR)/serverconf/app.toml $(TARGET_WIN_DIR)/app.toml
	cp $(CONF_DIR)/localdev/nacos.toml $(TARGET_WIN_DIR)/nacos.toml
	cp $(CONF_DIR)/storeconf/mysqlconf/otp.sql $(TARGET_WIN_DIR)/mysql-otp.sql
	cp $(CONF_DIR)/storeconf/pgsqlconf/otp.sql $(TARGET_WIN_DIR)/pgsql-otp.sql

.PHONY: windows_zip
zip_win:
	zip -q -r -f $(TARGET_WIN_ZIP) $(TARGET_WIN_DIR)


.PHONY: build
build: windows_build windows_copy windows_zip