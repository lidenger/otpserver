@echo off
REM [Windows] 清除临时文件
REM 设置编码为UTF-8
chcp 65001
echo start clear up tmp file
cd ..
cd ..

if exist app.exe (del /f app.exe)
if exist app.log (del /f app.log)
if exist request.log (del /f request.log)
if exist app.key (del /f app.key)
if exist app.toml (del /f app.toml)

if exist nacos.toml (del /f nacos.toml)

echo clear up tmp file finish


