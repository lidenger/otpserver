@echo off
REM [Windows] build
REM 设置编码为UTF-8
chcp 65001

REM windows amd64
echo 开始编译 [windows amd64]
go env -w CGO_ENABLED=0 GOOS=windows GOARCH=amd64
cd ..
cd main
go build -ldflags "-s -w" -o ../../doc/download/windows-amd64/app.exe .
cd ..
cd ..
copy .\config\serverconf\app.toml .\doc\download\windows-amd64\app.toml
copy .\config\localdev\nacos.toml .\doc\download\windows-amd64\nacos.toml
if exist .\doc\download\windows-amd64.zip (del /f .\doc\download\windows-amd64.zip)
powershell Compress-Archive -Path .\doc\download\windows-amd64 -DestinationPath .\doc\download\windows-amd64.zip
rmdir /s /q .\doc\download\windows-amd64
echo 编译完成  [windows amd64]

REM linux amd64
echo 开始编译 [linux amd64]
go env -w CGO_ENABLED=0 GOOS=linux GOARCH=amd64
cd cmd
cd main
go build -ldflags "-s -w" -o ../../doc/download/linux-amd64/app .
cd ..
cd ..
copy .\config\serverconf\app.toml .\doc\download\linux-amd64\app.toml
copy .\config\localdev\nacos.toml .\doc\download\linux-amd64\nacos.toml
if exist .\doc\download\linux-amd64.zip (del /f .\doc\download\linux-amd64.zip)
powershell Compress-Archive -Path .\doc\download\linux-amd64 -DestinationPath .\doc\download\linux-amd64.zip
rmdir /s /q .\doc\download\linux-amd64
echo 编译完成  [linux amd64]

REM macOS amd64 [macOS]
echo 开始编译 [macOS amd64]
go env -w CGO_ENABLED=0 GOOS=darwin GOARCH=amd64
cd cmd
cd main
go build -ldflags "-s -w" -o ../../doc/download/macos-amd64/app .
cd ..
cd ..
copy .\config\serverconf\app.toml .\doc\download\macos-amd64\app.toml
copy .\config\localdev\nacos.toml .\doc\download\macos-amd64\nacos.toml
if exist .\doc\download\macos-amd64.zip (del /f .\doc\download\macos-amd64.zip)
powershell Compress-Archive -Path .\doc\download\macos-amd64 -DestinationPath .\doc\download\macos-amd64.zip
rmdir /s /q .\doc\download\macos-amd64
echo 编译完成 [macOS amd64]


REM darwin arm64 [macOS M系列芯片]
echo 开始编译 [macOS M系列芯片]
go env -w CGO_ENABLED=0 GOOS=darwin GOARCH=arm64
cd cmd
cd main
go build -ldflags "-s -w" -o ../../doc/download/macos-arm64/app .
cd ..
cd ..
copy .\config\serverconf\app.toml .\doc\download\macos-arm64\app.toml
copy .\config\localdev\nacos.toml .\doc\download\macos-arm64\nacos.toml
if exist .\doc\download\macos-arm64.zip (del /f .\doc\download\macos-arm64.zip)
powershell Compress-Archive -Path .\doc\download\macos-arm64 -DestinationPath .\doc\download\macos-arm64.zip
rmdir /s /q .\doc\download\macos-arm64
echo 编译完成 [macOS M系列芯片]
