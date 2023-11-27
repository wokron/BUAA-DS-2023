# (go|rs|py)-ens
使用 go 实现的 ENS（Event Notification Service）服务器和客户端。以及 rust 和 python 实现的 ENS 客户端。

## go-ens
实现了四个命令行程序：`server`、`client`、`consumer` 和 `publisher`。可以使用如下命令构建

```sh
$ go build ./cmd/goens-<program name>
```

运行 `<program> -h` 可以得到参数说明，这里不再详述。

## rs-ens
实现了一个命令行程序，可以订阅事件并从服务器接收消息，或者向服务器发布事件。可以使用如下命令构建

```sh
$ cargo build --release
```

程序的使用方法如下所示，首先需要指定服务器的地址和端口号，默认为 `localhost` 和 `3456`。
```text
Usage: rs-ens [OPTIONS] <COMMAND>

Commands:
  subscribe  
  publish    
  help       Print this message or the help of the given subcommand(s)

Options:
  -s, --server <SERVER>  [default: localhost]
  -p, --port <PORT>      [default: 3456]
  -h, --help             Print help
```

之后对于订阅和发布，程序分别定义了两个子命令 `subscribe` 和 `publish`。
```sh
$ rs-ens -s <server> -p <port> subscribe
$ rs-ens -s <server> -p <port> publish
```

其中 `subscribe` 的使用方法如下，需要使用 `-t` 指定想要订阅的事件名。如 `rs-ens subscribe -t tech -t math` 将订阅 `tech` 和 `math` 两个事件。
```text
Usage: rs-ens subscribe [OPTIONS]

Options:
  -t, --topic <TOPIC>  
  -h, --help           Print help
```

`publish` 的使用方法如下，使用 `-e` 指定要发布的事件。其中事件的格式为 `<topic>:<message>`。如 `rs-ens publish -t tech:chatgpt -t math:convex` 将向 `tech` 事件发送一条内容为 `chatgpt` 的消息；向 `math` 事件发送一条内容为 `convex` 的消息。
```
Usage: rs-ens publish [OPTIONS]

Options:
  -e, --event <EVENT>  
  -h, --help           Print help
```

> 注意 rs-ens 的默认日志等级为 error。因此想要查看 info 信息还需在命令前加上 `RUST_LOG=info`。

## py-ens
和 rs-ens 类似，同样实现了一个命令行程序，可以订阅事件并从服务器接收消息，或者向服务器发布事件。

直接使用 python 运行 main.py 即可。
```sh
$ python main.py
```

py-ens 的使用方法和 rs-ens 类似。需要设定服务器地址和端口号，拥有 `subscribe` 和 `publish` 子命令。
```text
usage: main.py [-h] [--server SERVER] [--port PORT] {publish,subscribe} ...

positional arguments:
  {publish,subscribe}
    publish             publish events
    subscribe           subscribe topics and wait for receiveing

options:
  -h, --help            show this help message and exit
  --server SERVER, -s SERVER
  --port PORT, -p PORT
```

subscribe 子命令的用法与 rs-ens 相同
```text
usage: main.py subscribe [-h] [--topic TOPIC]

options:
  -h, --help            show this help message and exit
  --topic TOPIC, -t TOPIC
                        topic to subscribe
```

publish 子命令的用法也与 rs-ens 相同
```text
usage: main.py publish [-h] [--event [EVENT ...]]

options:
  -h, --help            show this help message and exit
  --event [EVENT ...], -e [EVENT ...] event to publish
```
