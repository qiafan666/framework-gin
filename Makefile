# 获取当前目录
CUR_DIR := $(shell pwd)
# curd生成代码工具目录
GENERATE_DIR := $(CUR_DIR)/tool/gormt
# ws proto目录
WS_PROTO_DIR := $(CUR_DIR)/ws/proto
# goose mysql连接信息
GOOSE_MYSQL_DSN := "root:admin123@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
# 检测操作系统类型
UNAME_S := $(shell uname -s)

# 根据操作系统选择不同的可执行文件
ifeq ($(UNAME_S), Darwin)
    EXECUTABLE := $(GENERATE_DIR)/gormt
else ifeq ($(UNAME_S), Linux)
    EXECUTABLE := $(GENERATE_DIR)/gormt
else ifeq ($(UNAME_S), Windows_NT)
    EXECUTABLE := $(GENERATE_DIR)/gormt.exe
endif

.PHONY: init
init:
	cd $(GOPATH)/bin; \
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest
	go install github.com/go-swagger/go-swagger/cmd/swagger@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest

# 根据数据库表生成CURD代码
.PHONY: curd
curd:
	@echo "当前目录: $(CUR_DIR)"
	@echo "执行可执行文件: $(EXECUTABLE)"
	$(EXECUTABLE)


# 生成ws proto代码
.PHONY: proto
proto:
	@echo "当前目录: $(CUR_DIR)"
	@echo "ws proto目录: $(WS_PROTO_DIR)"
	protoc -I=$(WS_PROTO_DIR) --go_out=$(WS_PROTO_DIR) $(WS_PROTO_DIR)/*.proto

goose_create:
	@goose -dir ./migrations -table goose_version create default sql
goose_up:
	@goose -dir ./migrations -table goose_version mysql $(GOOSE_MYSQL_DSN) up
goose_one:
	@goose -dir ./migrations -table goose_version mysql $(GOOSE_MYSQL_DSN) up-by-one
goose_version:
	@goose -dir ./migrations -table goose_version mysql $(GOOSE_MYSQL_DSN) version
goose_status:
	@goose -dir ./migrations -table goose_version mysql $(GOOSE_MYSQL_DSN) status
goose_down:
	@goose -dir ./migrations -table goose_version mysql $(GOOSE_MYSQL_DSN) down
