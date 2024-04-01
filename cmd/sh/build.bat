@echo off
REM [Windows] build
REM 设置编码为UTF-8
chcp 65001
echo 开始编译
go env -w GOOS=windows
cd ..
cd main
go build -ldflags "-s -w" -o ../../app.exe .
echo 编译完成
