# gRPC实践
## 简介 
gRPC的功能是让客户端应用像调用一个本地对象一样直接调用另一台机器上的服务端应用的方法，使分布式应用和服务更容易创建。与许多RPC系统一样，gRPC基于定义服务的思想，指定可以使用参数和返回类型远程调用的方法。在服务端，服务器实现此接口并运行gRPC服务器来处理客户端调用。在客户端，客户机有一个stub（在某些语言中仅称为客户机），它提供与服务器相同的方法。

![](https://grpc.io/img/landing-2.svg)

gRPC客户机和服务器可以在各种环境中运行和通信，可以用任何GRPC支持的语言编写。例如，可以轻松地在Java中用GO、Python或Ruby中的客户端创建gRPC服务器。
## 安装
以macOS为例。
```bash
brew install autoconf automake libtool
brew tap grpc/grpc
brew install --with-plugins grpc
```
编译器插件protoc-gen-go会安装在`$GOBIN`，默认是`$GOPATH/bin`，需要添加到环境变量中使其能被协议编译器protoc找到。
```bash
$ export PATH=$PATH:$GOPATH/bin
```
## 搭建示例
如果安装完成，可打开示例路径。
```bash
$ cd $GOPATH/src/google.golang.org/grpc/examples/helloworld
```
gRPC服务被定义在一个`.proto`文件中。
- helloworld.proto
    ```js
    // Copyright 2015 gRPC authors.
    //
    // Licensed under the Apache License, Version 2.0 (the "License");
    // you may not use this file except in compliance with the License.
    // You may obtain a copy of the License at
    //
    //     http://www.apache.org/licenses/LICENSE-2.0
    //
    // Unless required by applicable law or agreed to in writing, software
    // distributed under the License is distributed on an "AS IS" BASIS,
    // WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    // See the License for the specific language governing permissions and
    // limitations under the License.

    syntax = "proto3";

    option java_multiple_files = true;
    option java_package = "io.grpc.examples.helloworld";
    option java_outer_classname = "HelloWorldProto";

    package helloworld;

    // The greeting service definition.
    service Greeter {
    // Sends a greeting
    rpc SayHello (HelloRequest) returns (HelloReply) {}
    }

    // The request message containing the user's name.
    message HelloRequest {
    string name = 1;
    }

    // The response message containing the greetings
    message HelloReply {
    string message = 1;
    }
    ```
protoc编译器可以把它编译成一个相应的`.pb.go`文件。
```bash
$ protoc -I helloworld/ helloworld/helloworld.proto --go_out=plugins=grpc:helloworld
```
- helloworld.pb.go
    ```go
    // Code generated by protoc-gen-go. DO NOT EDIT.
    // source: helloworld.proto

    package helloworld

    import (
        context "context"
        fmt "fmt"
        proto "github.com/golang/protobuf/proto"
        grpc "google.golang.org/grpc"
        math "math"
    )

    // Reference imports to suppress errors if they are not otherwise used.
    var _ = proto.Marshal
    var _ = fmt.Errorf
    var _ = math.Inf

    // This is a compile-time assertion to ensure that this generated file
    // is compatible with the proto package it is being compiled against.
    // A compilation error at this line likely means your copy of the
    // proto package needs to be updated.
    const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

    // The request message containing the user's name.
    type HelloRequest struct {
        Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
        XXX_NoUnkeyedLiteral struct{} `json:"-"`
        XXX_unrecognized     []byte   `json:"-"`
        XXX_sizecache        int32    `json:"-"`
    }

    func (m *HelloRequest) Reset()         { *m = HelloRequest{} }
    func (m *HelloRequest) String() string { return proto.CompactTextString(m) }
    func (*HelloRequest) ProtoMessage()    {}
    func (*HelloRequest) Descriptor() ([]byte, []int) {
        return fileDescriptor_17b8c58d586b62f2, []int{0}
    }

    func (m *HelloRequest) XXX_Unmarshal(b []byte) error {
        return xxx_messageInfo_HelloRequest.Unmarshal(m, b)
    }
    func (m *HelloRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
        return xxx_messageInfo_HelloRequest.Marshal(b, m, deterministic)
    }
    func (m *HelloRequest) XXX_Merge(src proto.Message) {
        xxx_messageInfo_HelloRequest.Merge(m, src)
    }
    func (m *HelloRequest) XXX_Size() int {
        return xxx_messageInfo_HelloRequest.Size(m)
    }
    func (m *HelloRequest) XXX_DiscardUnknown() {
        xxx_messageInfo_HelloRequest.DiscardUnknown(m)
    }

    var xxx_messageInfo_HelloRequest proto.InternalMessageInfo

    func (m *HelloRequest) GetName() string {
        if m != nil {
            return m.Name
        }
        return ""
    }

    // The response message containing the greetings
    type HelloReply struct {
        Message              string   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
        XXX_NoUnkeyedLiteral struct{} `json:"-"`
        XXX_unrecognized     []byte   `json:"-"`
        XXX_sizecache        int32    `json:"-"`
    }

    func (m *HelloReply) Reset()         { *m = HelloReply{} }
    func (m *HelloReply) String() string { return proto.CompactTextString(m) }
    func (*HelloReply) ProtoMessage()    {}
    func (*HelloReply) Descriptor() ([]byte, []int) {
        return fileDescriptor_17b8c58d586b62f2, []int{1}
    }

    func (m *HelloReply) XXX_Unmarshal(b []byte) error {
        return xxx_messageInfo_HelloReply.Unmarshal(m, b)
    }
    func (m *HelloReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
        return xxx_messageInfo_HelloReply.Marshal(b, m, deterministic)
    }
    func (m *HelloReply) XXX_Merge(src proto.Message) {
        xxx_messageInfo_HelloReply.Merge(m, src)
    }
    func (m *HelloReply) XXX_Size() int {
        return xxx_messageInfo_HelloReply.Size(m)
    }
    func (m *HelloReply) XXX_DiscardUnknown() {
        xxx_messageInfo_HelloReply.DiscardUnknown(m)
    }

    var xxx_messageInfo_HelloReply proto.InternalMessageInfo

    func (m *HelloReply) GetMessage() string {
        if m != nil {
            return m.Message
        }
        return ""
    }

    func init() {
        proto.RegisterType((*HelloRequest)(nil), "helloworld.HelloRequest")
        proto.RegisterType((*HelloReply)(nil), "helloworld.HelloReply")
    }

    func init() { proto.RegisterFile("helloworld.proto", fileDescriptor_17b8c58d586b62f2) }

    var fileDescriptor_17b8c58d586b62f2 = []byte{
        // 175 bytes of a gzipped FileDescriptorProto
        0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0xc8, 0x48, 0xcd, 0xc9,
        0xc9, 0x2f, 0xcf, 0x2f, 0xca, 0x49, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x42, 0x88,
        0x28, 0x29, 0x71, 0xf1, 0x78, 0x80, 0x78, 0x41, 0xa9, 0x85, 0xa5, 0xa9, 0xc5, 0x25, 0x42, 0x42,
        0x5c, 0x2c, 0x79, 0x89, 0xb9, 0xa9, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x60, 0xb6, 0x92,
        0x1a, 0x17, 0x17, 0x54, 0x4d, 0x41, 0x4e, 0xa5, 0x90, 0x04, 0x17, 0x7b, 0x6e, 0x6a, 0x71, 0x71,
        0x62, 0x3a, 0x4c, 0x11, 0x8c, 0x6b, 0xe4, 0xc9, 0xc5, 0xee, 0x5e, 0x94, 0x9a, 0x5a, 0x92, 0x5a,
        0x24, 0x64, 0xc7, 0xc5, 0x11, 0x9c, 0x58, 0x09, 0xd6, 0x25, 0x24, 0xa1, 0x87, 0xe4, 0x02, 0x64,
        0xcb, 0xa4, 0xc4, 0xb0, 0xc8, 0x14, 0xe4, 0x54, 0x2a, 0x31, 0x38, 0x19, 0x70, 0x49, 0x67, 0xe6,
        0xeb, 0xa5, 0x17, 0x15, 0x24, 0xeb, 0xa5, 0x56, 0x24, 0xe6, 0x16, 0xe4, 0xa4, 0x16, 0x23, 0xa9,
        0x75, 0xe2, 0x07, 0x2b, 0x0e, 0x07, 0xb1, 0x03, 0x40, 0x5e, 0x0a, 0x60, 0x4c, 0x62, 0x03, 0xfb,
        0xcd, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0x0f, 0xb7, 0xcd, 0xf2, 0xef, 0x00, 0x00, 0x00,
    }

    // Reference imports to suppress errors if they are not otherwise used.
    var _ context.Context
    var _ grpc.ClientConn

    // This is a compile-time assertion to ensure that this generated file
    // is compatible with the grpc package it is being compiled against.
    const _ = grpc.SupportPackageIsVersion4

    // GreeterClient is the client API for Greeter service.
    //
    // For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
    type GreeterClient interface {
        // Sends a greeting
        SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error)
    }

    type greeterClient struct {
        cc *grpc.ClientConn
    }

    func NewGreeterClient(cc *grpc.ClientConn) GreeterClient {
        return &greeterClient{cc}
    }

    func (c *greeterClient) SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error) {
        out := new(HelloReply)
        err := c.cc.Invoke(ctx, "/helloworld.Greeter/SayHello", in, out, opts...)
        if err != nil {
            return nil, err
        }
        return out, nil
    }

    // GreeterServer is the server API for Greeter service.
    type GreeterServer interface {
        // Sends a greeting
        SayHello(context.Context, *HelloRequest) (*HelloReply, error)
    }

    func RegisterGreeterServer(s *grpc.Server, srv GreeterServer) {
        s.RegisterService(&_Greeter_serviceDesc, srv)
    }

    func _Greeter_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
        in := new(HelloRequest)
        if err := dec(in); err != nil {
            return nil, err
        }
        if interceptor == nil {
            return srv.(GreeterServer).SayHello(ctx, in)
        }
        info := &grpc.UnaryServerInfo{
            Server:     srv,
            FullMethod: "/helloworld.Greeter/SayHello",
        }
        handler := func(ctx context.Context, req interface{}) (interface{}, error) {
            return srv.(GreeterServer).SayHello(ctx, req.(*HelloRequest))
        }
        return interceptor(ctx, in, info, handler)
    }

    var _Greeter_serviceDesc = grpc.ServiceDesc{
        ServiceName: "helloworld.Greeter",
        HandlerType: (*GreeterServer)(nil),
        Methods: []grpc.MethodDesc{
            {
                MethodName: "SayHello",
                Handler:    _Greeter_SayHello_Handler,
            },
        },
        Streams:  []grpc.StreamDesc{},
        Metadata: "helloworld.proto",
    }
    ```
此文件包含生成的客户端和服务器代码，包括填充、序列化和检索`HelloRequest`和`HelloReply`消息的类型。
## 尝试
编译和执行示例的客户端和服务器代码。
```bash
$ go run greeter_server/main.go
```
打开一个新的终端。
```bash
$ go run greeter_client/main.go
```
在客户端可以看到
```bash
2019/01/16 17:45:15 Greeting: Hello world
```
在服务端可以看到
```bash
2019/01/16 17:45:15 Received: world
```
## 更新一个gRPC服务
给这个gRPC服务的服务端增加一个方法给客户端调用。
编辑`helloworld/helloworld.proto`，增添一个`SayHelloAgain`的方法。
```js
// The greeting service definition.
service Greeter {
  // Sends a greeting
  rpc SayHello (HelloRequest) returns (HelloReply) {}
  // Sends another greeting
  rpc SayHelloAgain (HelloRequest) returns (HelloReply) {}
}

// The request message containing the user's name.
message HelloRequest {
  string name = 1;
}

// The response message containing the greetings
message HelloReply {
  string message = 1;
}
```
用protoc编译器把它编译成一个相应的`.pb.go`文件。
```bash
$ protoc -I helloworld/ helloworld/helloworld.proto --go_out=plugins=grpc:helloworld
```

更新服务端`greeter_server/main.go`代码，添加函数
```go
func (s *server) SayHelloAgain(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
        return &pb.HelloReply{Message: "Hello again " + in.Name}, nil
}
```
更新客户端`greeter_client/main.go`代码，添加函数
```go
r, err = c.SayHelloAgain(ctx, &pb.HelloRequest{Name: name})
if err != nil {
        log.Fatalf("could not greet: %v", err)
}
log.Printf("Greeting: %s", r.Message)
```
编译和执行示例的客户端和服务器代码。
```bash
$ go run greeter_server/main.go
```
打开一个新的终端。
```bash
$ go run greeter_client/main.go
```
在客户端可以看到
```bash
2019/01/16 20:35:49 Greeting: Hello world
2019/01/16 20:35:49 Greeting: Hello again world
```
在服务端可以看到
```bash
2019/01/16 20:35:49 Received: world
```