@echo off
REM [Windows]本地研发环境，为了方便在IDE下直接启动，将项目启动依赖的文件复制到根目录
REM 设置编码为UTF-8
chcp 65001
echo start

cd ..
cd ..

if exist app.key (echo app.key exists) else (copy .\config\serverconf\app.key app.key)
if exist app.toml (echo app.toml exists) else (copy .\config\serverconf\app.toml app.toml)
if exist nacos.toml (echo app.key exists) else (copy .\config\serverconf\nacos.toml nacos.toml)

echo finish