#!/bin/bash

# 判断平台
OS=$(uname -s)

# 设置可执行文件名
EXECUTABLE="server"

# 编译 Go 程序
if [ "$OS" == "Linux" ]; then
    GOOS=linux GOARCH=amd64 go build -o $EXECUTABLE cmd/server/main.go
elif [[ "$OS" == "MINGW64_NT"* ]] || [[ "$OS" == "MINGW32_NT"* ]]; then
    GOOS=windows GOARCH=amd64 go build -o $EXECUTABLE.exe cmd/server/main.go
else
    echo "Unsupported OS: $OS"
    exit 1
fi

# 检查编译是否成功
if [ $? -ne 0 ]; then
    echo "Build failed"
    exit 1
fi

# 执行可执行文件
if [[ "$OS" == "MINGW64_NT"* ]] || [[ "$OS" == "MINGW32_NT"* ]]; then
    start ./$EXECUTABLE.exe -conf ./config/prod.yml
    pause
else
    ./$EXECUTABLE -conf ./config/prod.yml
fi