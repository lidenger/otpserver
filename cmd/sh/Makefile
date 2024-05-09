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
DOWNLOAD_DIR = ../../doc/download

TARGET_WIN_DIR = $(DOWNLOAD_DIR)/windows-amd64
TARGET_WIN_ZIP = $(DOWNLOAD_DIR)/windows-amd64.zip

.PHONY: windows_build
windows_build:
	go env -w CGO_ENABLED=0 GOOS=windows GOARCH=amd64
	go build -ldflags "-s -w" -o $(TARGET_WIN_DIR)/app.exe ../main/.

.PHONY: windows_conf
windows_conf:
	cp $(CONF_DIR)/serverconf/app.toml $(TARGET_WIN_DIR)/app.toml
	cp $(CONF_DIR)/localdev/nacos.toml $(TARGET_WIN_DIR)/nacos.toml
	cp $(CONF_DIR)/storeconf/mysqlconf/otp.sql $(TARGET_WIN_DIR)/mysql-otp.sql
	cp $(CONF_DIR)/storeconf/pgsqlconf/otp.sql $(TARGET_WIN_DIR)/pgsql-otp.sql

.PHONY: windows_zip
windows_zip:
	rm -rf $(TARGET_WIN_ZIP)
	zip -q -r $(TARGET_WIN_ZIP) $(TARGET_WIN_DIR)
	rm -rf $(TARGET_WIN_DIR)


TARGET_LINUX_DIR = $(DOWNLOAD_DIR)/linux-amd64
TARGET_LINUX_ZIP = $(DOWNLOAD_DIR)/linux-amd64.zip

.PHONY: linux_build
linux_build:
	go env -w CGO_ENABLED=0 GOOS=linux GOARCH=amd64
	go build -ldflags "-s -w" -o $(TARGET_LINUX_DIR)/app ../main/.

.PHONY: linux_conf
linux_conf:
	cp $(CONF_DIR)/serverconf/app.toml $(TARGET_LINUX_DIR)/app.toml
	cp $(CONF_DIR)/localdev/nacos.toml $(TARGET_LINUX_DIR)/nacos.toml
	cp $(CONF_DIR)/storeconf/mysqlconf/otp.sql $(TARGET_LINUX_DIR)/mysql-otp.sql
	cp $(CONF_DIR)/storeconf/pgsqlconf/otp.sql $(TARGET_LINUX_DIR)/pgsql-otp.sql

.PHONY: linux_zip
linux_zip:
	rm -rf $(TARGET_LINUX_ZIP)
	zip -q -r $(TARGET_LINUX_ZIP) $(TARGET_LINUX_DIR)
	rm -rf $(TARGET_LINUX_DIR)



TARGET_MACOS_ARM_DIR = $(DOWNLOAD_DIR)/macos-arm64
TARGET_MACOS_ARM_ZIP = $(DOWNLOAD_DIR)/macos-arm64.zip

.PHONY: macos_arm64_build
macos_arm64_build:
	go env -w CGO_ENABLED=0 GOOS=darwin GOARCH=arm64
	go build -ldflags "-s -w" -o $(TARGET_MACOS_ARM_DIR)/app ../main/.

.PHONY: macos_arm64_conf
macos_arm64_conf:
	cp $(CONF_DIR)/serverconf/app.toml $(TARGET_MACOS_ARM_DIR)/app.toml
	cp $(CONF_DIR)/localdev/nacos.toml $(TARGET_MACOS_ARM_DIR)/nacos.toml
	cp $(CONF_DIR)/storeconf/mysqlconf/otp.sql $(TARGET_MACOS_ARM_DIR)/mysql-otp.sql
	cp $(CONF_DIR)/storeconf/pgsqlconf/otp.sql $(TARGET_MACOS_ARM_DIR)/pgsql-otp.sql

.PHONY: macos_arm64_zip
macos_arm64_zip:
	rm -rf $(TARGET_MACOS_ARM_ZIP)
	zip -q -r $(TARGET_MACOS_ARM_ZIP) $(TARGET_MACOS_ARM_DIR)
	rm -rf $(TARGET_MACOS_ARM_DIR)


TARGET_MACOS_AMD_DIR = $(DOWNLOAD_DIR)/macos-amd64
TARGET_MACOS_AMD_ZIP = $(DOWNLOAD_DIR)/macos-amd64.zip

.PHONY: macos_amd64_build
macos_amd64_build:
	go env -w CGO_ENABLED=0 GOOS=darwin GOARCH=amd64
	go build -ldflags "-s -w" -o $(TARGET_MACOS_AMD_DIR)/app ../main/.

.PHONY: macos_amd64_conf
macos_amd64_conf:
	cp $(CONF_DIR)/serverconf/app.toml $(TARGET_MACOS_AMD_DIR)/app.toml
	cp $(CONF_DIR)/localdev/nacos.toml $(TARGET_MACOS_AMD_DIR)/nacos.toml
	cp $(CONF_DIR)/storeconf/mysqlconf/otp.sql $(TARGET_MACOS_AMD_DIR)/mysql-otp.sql
	cp $(CONF_DIR)/storeconf/pgsqlconf/otp.sql $(TARGET_MACOS_AMD_DIR)/pgsql-otp.sql

.PHONY: macos_amd64_zip
macos_amd64_zip:
	rm -rf $(TARGET_MACOS_AMD_ZIP)
	zip -q -r $(TARGET_MACOS_AMD_ZIP) $(TARGET_MACOS_AMD_DIR)
	rm -rf $(TARGET_MACOS_AMD_DIR)


.PHONY: build
build: windows_build windows_conf windows_zip \
	   linux_build linux_conf linux_zip \
	   macos_arm64_build macos_arm64_conf macos_arm64_zip \
	   macos_amd64_build macos_amd64_conf macos_amd64_zip