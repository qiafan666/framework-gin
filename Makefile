# 获取当前目录
CUR_DIR := $(shell pwd)
# curd生成代码工具目录
GENERATE_DIR := $(CUR_DIR)/tool/generateCURD
# ws proto目录
WS_PROTO_DIR := $(CUR_DIR)/ws/proto

# 检测操作系统类型
UNAME_S := $(shell uname -s)

# 根据操作系统选择不同的可执行文件
ifeq ($(UNAME_S), Darwin)
    EXECUTABLE := $(GENERATE_DIR)/main
else ifeq ($(UNAME_S), Linux)
    EXECUTABLE := $(GENERATE_DIR)/main
else ifeq ($(UNAME_S), Windows_NT)
    EXECUTABLE := $(GENERATE_DIR)/main.exe
endif

# 根据数据库表生成CURD代码
.PHONY: curd
curd:
	@echo "当前目录: $(CUR_DIR)"
	@echo "执行可执行文件: $(EXECUTABLE)"
	$(EXECUTABLE)


# 生成ws proto代码
.PHONY: ws
ws:
	@echo "当前目录: $(CUR_DIR)"
	@echo "ws proto目录: $(WS_PROTO_DIR)"
	protoc -I=$(WS_PROTO_DIR) --go_out=$(WS_PROTO_DIR) $(WS_PROTO_DIR)/*.proto