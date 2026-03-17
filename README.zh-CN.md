# gclint

[English](README.md) | 中文

面向工程实践的 Go 风格约束，直接交付为 `gclint`。

`gclint` 用于构建一个名为 `gclint` 的自定义 `golangci-lint` 二进制。

当前仓库只有一个模块插件包：`style`。

## 安装

发布版本号格式为 `YYYY.MM.DD-SHORT_HASH_ID`，例如 `2026.03.17-deadbee`。

发布工作流会基于 `.custom-gcl.yml` 中固定版本对应的官方 `golangci-lint` release 资产，动态生成构建平台矩阵。

### 通过 `install.sh` 安装最新版本

`install.sh` 会自动识别当前平台并下载对应的发布包。

```bash
curl -fsSL https://raw.githubusercontent.com/liasica/gclint/master/install.sh | sh
```

安装到自定义目录：

```bash
curl -fsSL https://raw.githubusercontent.com/liasica/gclint/master/install.sh | GCLINT_INSTALL_DIR="$HOME/.local/bin" sh
```

安装指定版本：

```bash
curl -fsSL https://raw.githubusercontent.com/liasica/gclint/master/install.sh | GCLINT_VERSION="2026.03.17-deadbee" sh
```

需要时也可以手动覆盖平台识别结果：

```bash
curl -fsSL https://raw.githubusercontent.com/liasica/gclint/master/install.sh | \
  GCLINT_OS=linux GCLINT_ARCH=armv7 GCLINT_INSTALL_DIR="$HOME/.local/bin" sh
```

### 手动下载 GitHub Release 资产

Unix-like 平台发布为 `gclint_<version>_<os>_<arch>.tar.gz`，Windows 发布为 `gclint_<version>_<os>_<arch>.zip`。

当前发布平台：

- `darwin`：`amd64`、`arm64`
- `freebsd`：`386`、`amd64`、`arm64`、`armv6`、`armv7`
- `illumos`：`amd64`
- `linux`：`386`、`amd64`、`arm64`、`armv6`、`armv7`、`loong64`、`mips64`、`mips64le`、`ppc64le`、`riscv64`、`s390x`
- `netbsd`：`386`、`amd64`、`arm64`、`armv6`、`armv7`
- `windows`：`386`、`amd64`、`arm64`

### 从源码构建

```bash
make install-lint
make build-lint
install -m 0755 ./.bin/gclint /usr/local/bin/gclint
```

## 快速开始

```bash
make install-lint
make test
make lint
make ci
make clean
```

`make lint` 会根据 `.custom-gcl.yml` 构建 `.bin/gclint`，然后扫描当前仓库。

## 当前规则

- `errshort`：禁止在同一作用域里已声明 `err` 后，再用 `:=` 复用旧 `err`
- `namedreturn`：命名返回值已经赋值后，禁止显式镜像返回
- `chinesekey`：禁止中文 `json` tag key，以及 map literal 中的中文字符串 key
- `layerdep`：禁止低层包导入已配置的高层包
- `varreuse`：用启发式方式检查“语义明确的变量被复用为其他业务对象容器”

## 分层依赖配置

`layerdep` 是配置驱动的规则。

`.golangci.yml` 示例：

```yaml
linters:
  settings:
    custom:
      style:
        type: module
        description: Enforce custom Go style rules with style.
        settings:
          dependency_rules:
            - source: github.com/example/project/internal/repository
              forbidden:
                - github.com/example/project/internal/service
                - github.com/example/project/internal/handler
```

把这些包前缀替换成你项目里的真实分层路径即可。

## 规则覆盖清单

下面的清单将当前仓库实现状态映射到本项目使用的 Go 规范文档。

- [ ] 命名必须贴合真实业务语义，不得随意缩写
- [ ] 单数和复数语义必须准确
- [ ] 不同业务对象不能复用同一个名字
- [ ] 禁止无意义空格、空行和脏格式
- [ ] 不同逻辑块必须清晰分隔，必要时要有说明性注释
- [ ] 注释必须清晰、有逻辑，必要注释默认使用英文
- [x] 同一作用域中已声明 `err` 后，禁止 `:=` 再带上旧 `err`
- [x] 已有明确业务含义的变量，禁止复用为其他业务对象的容器
- [ ] 禁止大范围重复代码
- [x] 低层包禁止依赖高层包，当前通过 `settings.dependency_rules` 配置生效
- [x] JSON key 与固化 map key 禁止使用中文
- [x] 命名返回值赋值后，禁止显式返回镜像结果，必须直接 `return`

## 当前实现边界

- `chinesekey` 当前只检查显式 `json` struct tag 和 map literal 里的显式字符串 key
- `layerdep` 只校验 `settings.dependency_rules` 里列出的包前缀
- `varreuse` 是刻意保持保守的启发式规则，重点看描述性变量名和赋值来源里的稳定语义 token
