# ENS 实验结果报告
> Yitang Yang

## 一、代码实现
我们分别使用 golang、rust 和 python 实现了 ENS 服务。ENS 服务的重点在于事件的订阅。在使用 golang 实现 ENS 服务程序的时候，我们使用一个字典存储事件与订阅者间的映射关系。字典的键为事件名，字典的值为一保存了 tcp 连接的切片。即字典类型为 `map[string][]net.Conn`。为了避免多个 goroutine 同时修改字典数据，我们还引入了一个互斥锁。

为了实现不同语言间的相互通信，我们参照样例代码中的数据结构，将其作为公共的数据传输单元。其结构如下：
```rust
struct ENSMsgData {
    msg_type: u8,
    topic: [u8; MAX_TOPIC_LENGTH],
    message: [u8; MAX_MESSAGE_LENGTH],
}
```

在使用字节流形式传输之后，我们使用 golang、rust 和 python 中的相关功能将其还原成更易于处理的结构体 `ENSMsg`。并将这一过程封装进 `encode` 和 `decode` 函数中。
```golang
type ENSMsg struct {
	Type    ENSType
	Topic   string
	Message string
}
```

```rust
pub struct ENSMsg {
    pub msg_type: u8,
    pub topic: String,
    pub message: String,
}
```

```python
class ENSMsg:
    def __init__(self, type: int, topic: str = "", message: str = ""):
        self.type = type
        self.topic = topic
        self.message = message
```

之后我们还将数据传输封装进形如 `send_ens_msg` 和 `recv_ens_msg` 的函数中。从而构成了 ENS 服务的不同语言接口。在此基础上我们实现了多种不同语言的客户端。这些客户端之间都能与 ENS 服务程序相互通信，实现跨语言的事件通知。

## 二、实验设计
接下来实验中我们将使用三种语言分别与 golang 实现的 ENS 服务程序建立连接并订阅事件。并使用三种语言实现事件的发布，从而检验实现代码的正确性，理解 ENS 服务的原理。

- 首先我们运行 ENS 服务程序。
- 接着运行三种语言的客户端程序，订阅事件。其中 rust 客户端订阅 rust 和 common 事件；go 客户端订阅 golang 和 common 事件；python 客户端订阅 python 和 common 事件。
- 再之后我们运行另一组三种语言的客户端，分别发布事件。rust 客户端发布 rust 和 common 事件；go 客户端发布 golang 和 common 事件；python 客户端发布 python 和 common 事件。

## 三、实验结果
运行 ENS 服务程序。
```sh
go run ./cmd/goens-server
```

接着分别运行三种语言的客户端程序，订阅各自事件。
**rust**：
```sh
RUST_LOG=info cargo run -- subscribe -t rust -t common
```
```log
[2023-11-27T06:09:33Z INFO  rs_ens] Subscribe topic "rust"
[2023-11-27T06:09:33Z INFO  rs_ens] Subscribe topic "common"
```

**golang**：
```sh
go run ./cmd/goens-consumer -topic golang -topic common
```
```log
2023/11/27 14:09:41 Connecting to ENS Server localhost:4567.
2023/11/27 14:09:41 Subscrib topic "golang".
2023/11/27 14:09:41 Subscrib topic "common".
```

**python**：
```sh
python main.py subscribe -t python -t common
```
```log
27-09-2023 14:09:44:INFO:Connect to localhost:4567 success
27-09-2023 14:09:44:INFO:Subscribe topic "python"
27-09-2023 14:09:44:INFO:Subscribe topic "common"
```

服务程序日志显示三个订阅程序订阅了相关事件
```log
2023/11/27 14:09:30 Event Notification Service running on 0.0.0.0:4567.
2023/11/27 14:09:33 Accept tcp connection from 127.0.0.1:36604.
2023/11/27 14:09:33 127.0.0.1:36604: Subscribe topic rust.
2023/11/27 14:09:33 127.0.0.1:36604: Subscribe topic common.
2023/11/27 14:09:41 Accept tcp connection from 127.0.0.1:36618.
2023/11/27 14:09:41 127.0.0.1:36618: Subscribe topic golang.
2023/11/27 14:09:41 127.0.0.1:36618: Subscribe topic common.
2023/11/27 14:09:44 Accept tcp connection from 127.0.0.1:51634.
2023/11/27 14:09:44 127.0.0.1:51634: Subscribe topic python.
2023/11/27 14:09:44 127.0.0.1:51634: Subscribe topic common.
```

之后在此使用客户端发布事件。首先使用 rust 客户端发布 rust 和 common 事件
```sh
RUST_LOG=info cargo run -- publish -e "rust:i am rust too" -e "common:hello here is rust"
```

在服务程序日志中显示发布了事件
```log
2023/11/27 14:09:52 Accept tcp connection from 127.0.0.1:51648.
2023/11/27 14:09:52 127.0.0.1:51648: Publish event on topic rust, message: "i am rust too".
2023/11/27 14:09:52 127.0.0.1:51648: Publish event on topic common, message: "hello here is rust".
2023/11/27 14:09:52 127.0.0.1:51648: Connection close.
```

rust 订阅进程收到了 rust 和 common 事件的消息
```log
[2023-11-27T06:09:52Z INFO  rs_ens] Receive event, topic: "rust", message: "i am rust too".
[2023-11-27T06:09:52Z INFO  rs_ens] Receive event, topic: "common", message: "hello here is rust".
```

gloang 和 python 订阅进程收到了 common 事件的消息
```log
2023/11/27 14:09:52 Received Event, topic: "common", message: "hello here is rust".
```
```log
27-09-2023 14:09:52:INFO:Receive message on topic "common", message: "hello here is rust"
```

接着使用 golang 客户端发布 gloang 和 common 事件。
```sh
go run ./cmd/goens-publisher -event "golang:i am go too" -event "common:hello here is golang"
```
```log
2023/11/27 14:09:58 Connecting to ENS Server localhost:4567.
2023/11/27 14:09:58 Publish event on topic "golang", message: "i am go too"
2023/11/27 14:09:58 Publish event on topic "common", message: "hello here is golang"
2023/11/27 14:09:58 All events has been published, exit.
```
golang 订阅进程收到了 golang 和 common 事件的消息
```log
2023/11/27 14:09:58 Received Event, topic: "golang", message: "i am go too".
2023/11/27 14:09:58 Received Event, topic: "common", message: "hello here is golang".
```

rust 和 python 订阅进程收到了 common 事件
```log
[2023-11-27T06:09:58Z INFO  rs_ens] Receive event, topic: "common", message: "hello here is golang".
```
```log
27-09-2023 14:09:58:INFO:Receive message on topic "common", message: "hello here is golang"
```

最后再使用 python 客户端发布 python 和 common 事件。
```sh
python main.py publish -e "python:i am python too" -e "common:hello here is python"
```
```log
27-10-2023 14:10:03:INFO:Connect to localhost:4567 success
27-10-2023 14:10:03:INFO:Publish event on topic "python", message: "i am python too"
27-10-2023 14:10:03:INFO:Publish event on topic "common", message: "hello here is python"
27-10-2023 14:10:03:INFO:Publish all events
```

python 订阅进程收到了 python 和 common 事件的消息
```log
27-10-2023 14:10:03:INFO:Receive message on topic "python", message: "i am python too"
27-10-2023 14:10:03:INFO:Receive message on topic "common", message: "hello here is python"
```

rust 和 golang 订阅进程收到了 common 事件
```log
[2023-11-27T06:10:03Z INFO  rs_ens] Receive event, topic: "common", message: "hello here is python".
```
```log
2023/11/27 14:10:03 Received Event, topic: "common", message: "hello here is python".
```

### 完整输出
```log
$ go run ./cmd/goens-server
2023/11/27 14:09:30 Event Notification Service running on 0.0.0.0:4567.
2023/11/27 14:09:33 Accept tcp connection from 127.0.0.1:36604.
2023/11/27 14:09:33 127.0.0.1:36604: Subscribe topic rust.
2023/11/27 14:09:33 127.0.0.1:36604: Subscribe topic common.
2023/11/27 14:09:41 Accept tcp connection from 127.0.0.1:36618.
2023/11/27 14:09:41 127.0.0.1:36618: Subscribe topic golang.
2023/11/27 14:09:41 127.0.0.1:36618: Subscribe topic common.
2023/11/27 14:09:44 Accept tcp connection from 127.0.0.1:51634.
2023/11/27 14:09:44 127.0.0.1:51634: Subscribe topic python.
2023/11/27 14:09:44 127.0.0.1:51634: Subscribe topic common.
2023/11/27 14:09:52 Accept tcp connection from 127.0.0.1:51648.
2023/11/27 14:09:52 127.0.0.1:51648: Publish event on topic rust, message: "i am rust too".
2023/11/27 14:09:52 127.0.0.1:51648: Publish event on topic common, message: "hello here is rust".
2023/11/27 14:09:52 127.0.0.1:51648: Connection close.
2023/11/27 14:09:58 Accept tcp connection from 127.0.0.1:56472.
2023/11/27 14:09:58 127.0.0.1:56472: Publish event on topic golang, message: "i am go too".
2023/11/27 14:09:58 127.0.0.1:56472: Publish event on topic common, message: "hello here is golang".
2023/11/27 14:09:58 127.0.0.1:56472: Connection close.
2023/11/27 14:10:03 Accept tcp connection from 127.0.0.1:38324.
2023/11/27 14:10:03 127.0.0.1:38324: Publish event on topic python, message: "i am python too".
2023/11/27 14:10:03 127.0.0.1:38324: Publish event on topic common, message: "hello here is python".
2023/11/27 14:10:03 127.0.0.1:38324: Connection close.
2023/11/27 14:10:14 127.0.0.1:51634: Connection close.
2023/11/27 14:10:15 127.0.0.1:36618: Connection close.
2023/11/27 14:10:16 127.0.0.1:36604: Connection close.
```

```log
$ RUST_LOG=info cargo run -- subscribe -t rust -t common
[2023-11-27T06:09:33Z INFO  rs_ens] Subscribe topic "rust"
[2023-11-27T06:09:33Z INFO  rs_ens] Subscribe topic "common"
[2023-11-27T06:09:52Z INFO  rs_ens] Receive event, topic: "rust", message: "i am rust too".
[2023-11-27T06:09:52Z INFO  rs_ens] Receive event, topic: "common", message: "hello here is rust".
[2023-11-27T06:09:58Z INFO  rs_ens] Receive event, topic: "common", message: "hello here is golang".
[2023-11-27T06:10:03Z INFO  rs_ens] Receive event, topic: "common", message: "hello here is python".
```

```log
$ go run ./cmd/goens-consumer -topic golang -topic common
2023/11/27 14:09:41 Connecting to ENS Server localhost:4567.
2023/11/27 14:09:41 Subscrib topic "golang".
2023/11/27 14:09:41 Subscrib topic "common".
2023/11/27 14:09:52 Received Event, topic: "common", message: "hello here is rust".
2023/11/27 14:09:58 Received Event, topic: "golang", message: "i am go too".
2023/11/27 14:09:58 Received Event, topic: "common", message: "hello here is golang".
2023/11/27 14:10:03 Received Event, topic: "common", message: "hello here is python".
```

```log
$ python main.py subscribe -t python -t common
27-09-2023 14:09:44:INFO:Connect to localhost:4567 success
27-09-2023 14:09:44:INFO:Subscribe topic "python"
27-09-2023 14:09:44:INFO:Subscribe topic "common"
27-09-2023 14:09:52:INFO:Receive message on topic "common", message: "hello here is rust"
27-09-2023 14:09:58:INFO:Receive message on topic "common", message: "hello here is golang"
27-10-2023 14:10:03:INFO:Receive message on topic "python", message: "i am python too"
27-10-2023 14:10:03:INFO:Receive message on topic "common", message: "hello here is python"
```

```log
$ go run ./cmd/goens-publisher -event "golang:i am go too" -event "common:hello here is golang"
2023/11/27 14:09:58 Connecting to ENS Server localhost:4567.
2023/11/27 14:09:58 Publish event on topic "golang", message: "i am go too"
2023/11/27 14:09:58 Publish event on topic "common", message: "hello here is golang"
2023/11/27 14:09:58 All events has been published, exit.
```

```log
$ python main.py publish -e "python:i am python too" -e "common:hello here is python"
27-10-2023 14:10:03:INFO:Connect to localhost:4567 success
27-10-2023 14:10:03:INFO:Publish event on topic "python", message: "i am python too"
27-10-2023 14:10:03:INFO:Publish event on topic "common", message: "hello here is python"
27-10-2023 14:10:03:INFO:Publish all events
```
