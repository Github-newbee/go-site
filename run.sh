#!/bin/bash

# Go-Site 项目启动脚本
# 使用方法: ./run.sh [dev|prod|build|clean|wire|migrate]

set -e  # 遇到错误时退出

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目配置
PROJECT_NAME="go-site"
MAIN_FILE="cmd/server/main.go"
WIRE_DIR="cmd/server/wire"
MIGRATION_DIR="cmd/migration"
CONFIG_DEV="config/dev.yml"
CONFIG_PROD="config/prod.yml"

# 帮助信息
help() {
    echo -e "${BLUE}Go-Site 项目启动脚本${NC}"
    echo ""
    echo "使用方法:"
    echo "  ./run.sh [命令]"
    echo ""
    echo "可用命令:"
    echo -e "  ${GREEN}dev${NC}      - 开发模式启动 (默认)"
    echo -e "  ${GREEN}prod${NC}     - 生产模式启动"
    echo -e "  ${GREEN}build${NC}    - 编译项目"
    echo -e "  ${GREEN}clean${NC}    - 清理编译文件"
    echo -e "  ${GREEN}wire${NC}     - 重新生成 Wire 依赖注入代码"
    echo -e "  ${GREEN}migrate${NC}  - 运行数据库迁移"
    echo -e "  ${GREEN}test${NC}     - 运行测试"
    echo -e "  ${GREEN}help${NC}     - 显示帮助信息"
    echo ""
}

# 检查依赖
check_dependencies() {
    echo -e "${YELLOW}检查依赖...${NC}"
    
    # 检查 Go 是否安装
    if ! command -v go &> /dev/null; then
        echo -e "${RED}错误: Go 未安装或不在 PATH 中${NC}"
        exit 1
    fi
    
    # 检查 Wire 是否安装
    if ! command -v wire &> /dev/null; then
        echo -e "${YELLOW}Wire 未安装，正在安装...${NC}"
        go install github.com/google/wire/cmd/wire@latest
    fi
    
    echo -e "${GREEN}依赖检查完成${NC}"
}

# 生成 Wire 代码
generate_wire() {
    echo -e "${YELLOW}生成 Wire 依赖注入代码...${NC}"
    
    if [ -d "$WIRE_DIR" ]; then
        cd $WIRE_DIR
        wire
        cd - > /dev/null
        echo -e "${GREEN}Wire 代码生成完成${NC}"
    else
        echo -e "${RED}错误: Wire 目录不存在: $WIRE_DIR${NC}"
        exit 1
    fi
}

# 编译项目
build_project() {
    echo -e "${YELLOW}编译项目...${NC}"
    
    # 确保 Wire 代码是最新的
    generate_wire
    
    # 编译项目
    go build -o bin/$PROJECT_NAME $MAIN_FILE
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}项目编译成功: bin/$PROJECT_NAME${NC}"
    else
        echo -e "${RED}项目编译失败${NC}"
        exit 1
    fi
}

# 运行数据库迁移
run_migration() {
    echo -e "${YELLOW}运行数据库迁移...${NC}"
    
    if [ -f "$MIGRATION_DIR/main.go" ]; then
        go run $MIGRATION_DIR/*.go -action=up
        echo -e "${GREEN}数据库迁移完成${NC}"
    else
        echo -e "${RED}错误: 迁移文件不存在${NC}"
        exit 1
    fi
}

# 开发模式启动
dev_mode() {
    echo -e "${BLUE}启动开发模式...${NC}"
    
    # 检查配置文件
    if [ ! -f "$CONFIG_DEV" ]; then
        echo -e "${RED}错误: 开发配置文件不存在: $CONFIG_DEV${NC}"
        exit 1
    fi
    
    # 生成 Wire 代码
    generate_wire
    
    # 启动项目
    echo -e "${GREEN}项目启动中... (配置文件: $CONFIG_DEV)${NC}"
    go run $MAIN_FILE -conf=$CONFIG_DEV
}

# 生产模式启动
prod_mode() {
    echo -e "${BLUE}启动生产模式...${NC}"
    
    # 检查配置文件
    if [ ! -f "$CONFIG_PROD" ]; then
        echo -e "${RED}错误: 生产配置文件不存在: $CONFIG_PROD${NC}"
        exit 1
    fi
    
    # 编译项目
    build_project
    
    # 运行编译后的程序
    echo -e "${GREEN}项目启动中... (配置文件: $CONFIG_PROD)${NC}"
    ./bin/$PROJECT_NAME -conf=$CONFIG_PROD
}

# 清理编译文件
clean_build() {
    echo -e "${YELLOW}清理编译文件...${NC}"
    
    if [ -d "bin" ]; then
        rm -rf bin/
    fi
    
    # 清理 Wire 生成的文件
    if [ -f "$WIRE_DIR/wire_gen.go" ]; then
        rm -f $WIRE_DIR/wire_gen.go
    fi
    
    echo -e "${GREEN}清理完成${NC}"
}

# 运行测试
run_tests() {
    echo -e "${YELLOW}运行测试...${NC}"
    go test ./...
    echo -e "${GREEN}测试完成${NC}"
}

# 主函数
main() {
    # 获取命令参数
    COMMAND=${1:-dev}
    
    # 检查是否在项目根目录
    if [ ! -f "go.mod" ]; then
        echo -e "${RED}错误: 请在项目根目录运行此脚本${NC}"
        exit 1
    fi
    
    case $COMMAND in
        "dev")
            check_dependencies
            dev_mode
            ;;
        "prod")
            check_dependencies
            prod_mode
            ;;
        "build")
            check_dependencies
            build_project
            ;;
        "clean")
            clean_build
            ;;
        "wire")
            check_dependencies
            generate_wire
            ;;
        "migrate")
            check_dependencies
            run_migration
            ;;
        "test")
            check_dependencies
            run_tests
            ;;
        "help"|"-h"|"--help")
            help
            ;;
        *)
            echo -e "${RED}错误: 未知命令 '$COMMAND'${NC}"
            echo ""
            help
            exit 1
            ;;
    esac
}

# 运行主函数
main "$@"