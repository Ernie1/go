# Go环境与编程

Go是一种开源编程语言，用来构建简单、可靠和高效的软件。

## 下载Go发行版
- [官方二进制发行版](https://golang.org/dl/)

- [源代码安装](https://golang.org/doc/install/source)

- [安装gccgo](https://golang.org/doc/install/gccgo)。

## 安装Go工具
如果要从旧版本的Go升级，必须先[删除现有版本](https://golang.org/doc/install#uninstall)。

对于Linux、macOS和FreeBSD tarballs，[下载](https://golang.org/dl/)后提取到`/usr/local`中，在`/usr/local/go`中创建Go树。例如：
```bash
tar -C /usr/local -xzf go$VERSION.$OS-$ARCH.tar.gz
```
将`/usr/local/go/bin`添加到环境变量。添加
```bash
export PATH=$PATH:/usr/local/go/bin
```
到`/etc/profile`（对于系统范围的安装）或`$home/.profile`并执行
```bash
source $HOME/.profile
```
使其生效。

## 测试安装
通过启动一个工作区并构建一个简单的程序来检查Go是否正确安装。

创建[工作区](https://golang.org/doc/code.html#Workspaces)目录`$home/go`。（如果要使用其他目录，需要[设置GOPATH环境变量](https://golang.org/wiki/SettingGOPATH)。）

在工作区中创建目录`src/hello`，并在该目录中创建一个名为`hello.go`的文件如下：
```go
package main

import "fmt"

func main() {
	fmt.Printf("hello, world\n")
}
```

然后用Go工具构建它:
```bash
$ cd $HOME/go/src/hello
$ go build
```

将在源代码所在的目录中生成一个名为`hello`的可执行文件。执行它：
```bash
$ ./hello
hello, world
```

如果能看到`Hello，World`，说明Go正常工作。

可以运行`go install`将二进制文件安装到工作区的`bin`目录中，或者使用`go clean-i`将其删除。

要了解使用Go工具的基本概念，请阅读[如何编写go代码](https://golang.org/doc/code.html)。