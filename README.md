<h1 style="text-align:center"> 
ICS_Y86
<br>
</h1>

<h4 style="text-align:center">
A Y86 simulator project for Fudan ICS 2022 Autumn
</h4>

## Feature

### CPU设计

CPU后端基于Golang，分为三个组件

- `Instruction`：存储指令，包含机器码、指令地址、汇编码、注释信息。用于解析机器码为`ICode`, `IFun`, `RegA`, `RegB`, `ValC`, `ValJmp`信息。
- `Device`：设备的抽象，包含`PC`, `Reg`, `CC`, `Stat`, `Mem`，用于模拟设备底层，封装了读取内存、写入内存、Push、Pop、条件码检查、以及计算操作。

- `Controller`：用于管理整体程序的运行，包含一个`Device`和一个`map[uint64]Instruction`，包括了文件解析、程序运行控制等模块。

运行流程是，构造`Constroller`对象，和内嵌的`Device`对象，初始化；调用`Parser`函数，从`io.Stdin`中读取文件信息，并且解析成`Instruction`；再调用`Run`函数，对指令解释执行。

执行的策略是：先判断当前状态码是否为`AOK`，如果不是就直接退出循环；再根据`PC`值读取当前指令，然后获取`ICode`, `IFun`, `RegA`, `RegB`, `ValC`, `ValJmp`；对`iCode`进行枚举，内部再对`IFun`进行枚举，如果匹配，就按照指令要求进行读、写、运算、设置条件码等操作，如果读写了非法的内存地址，则设置条件码为`ADR`；如果没有匹配，就设置为状态码为`INS`。每条指令执行完成后，对当前`Device`进行`json`解析，程序运行结束后打包输出。

在内存的读写方面，`Memory`是一个单字节数组，封装了`Read`，`Write`等操作。`Read`操作时根据内存地址获取对应的八个字节的值（小端法），`Write`操作则是将内存地址写入到地址对应的八个字节中。为了节省时间，在每次写入的时候记录了最大写入内存的大小，最后`json`扫描并打包的过程可以更快。

### 前端设计

To be done.

## Usage

### Stage 1

#### prerequisite

Go 1.19

#### build and run

```shell
cd ./Stage1/ICS_Y86_Backend
go build .
./ICS_Y86_Backend < hello.yo > hello.json # for cmd and bash
Get-Content hello.yo | .\ICS_Y86_Backend.exe > hello.json # for powershell
```

### Stage 2

#### run

## License

MIT License

Copyright (c) 2022-present ck ct fkx