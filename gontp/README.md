# gontp
用 go 实现的 ntp 服务器和客户端。只实现了 ntp 的基本逻辑，**十分不完善**。

## 构建
> 脚本均位于 `/script`。

任意位置执行
```sh
source setvars.sh
```

执行构建
```sh
bash build.sh
```

可以在 `/out` 文件夹下找到构建得到的客户端和服务器程序 `gontp-client` 和 `gontp-server`。

## 命令行使用
命令行后添加 `-h` 可以查看命令行参数：
```text
Usage of gontp-client:
  -port int
        the port of ntp service (default 123)
  -result
        whether to show the rrt and offset (default true)
  -server string
        the server of ntp (default "ntp.ntsc.ac.cn")
  -settime
        whether to modify the system time
  -timeout int
        request timeout sec (default 10)
```
```text
Usage of gontp-server:
  -host string
        the host name of ntp service (default "0.0.0.0")
  -port int
        the port of ntp service (default 123)
```
