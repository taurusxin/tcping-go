# Tcping GO

一个使用 Golang 写的小工具，用于检测TCP端口是否开放

## 使用方法

你可以简单的测试一个域名(或 IP 地址)的某个端口是否开放，默认测试 80 端口

```shell
$ tcping.exe www.qq.com
Using www.qq.com A record: 101.91.42.232
[1] Reply from 101.91.42.232:80: time=14.400ms
[2] Reply from 101.91.42.232:80: time=15.086ms
[3] Reply from 101.91.42.232:80: time=14.557ms
[4] Reply from 101.91.42.232:80: time=11.862ms

Test finished, success 4/4 (100.00%)
min = 11.862ms, max = 15.086ms, avg = 13.976ms
```

也可以指定测试端口号，跟随在域名或 IP 地址后面

```shell
$ tcping.exe www.qq.com 443
Using www.qq.com A record: 101.91.42.232
[1] Reply from 101.91.42.232:443: time=12.348ms
[2] Reply from 101.91.42.232:443: time=10.089ms
[3] Reply from 101.91.42.232:443: time=10.370ms
[4] Reply from 101.91.42.232:443: time=9.131ms

Test finished, success 4/4 (100.00%)
min = 9.131ms, max = 12.348ms, avg = 10.485ms
```

使用 `-6` 指定使用IPv6（仅对域名有效）

```shell
$ tcping.exe www.qq.com -6 
Using www.qq.com AAAA record: 240e:e1:a800:120::36
[1] Reply from [240e:e1:a800:120::36]:80: time=17.576ms
[2] Reply from [240e:e1:a800:120::36]:80: time=14.510ms
[3] Reply from [240e:e1:a800:120::36]:80: time=14.667ms
[4] Reply from [240e:e1:a800:120::36]:80: time=15.449ms

Test finished, success 4/4 (100.00%)
min = 14.510ms, max = 17.576ms, avg = 15.551ms
```

使用 `-c` 指定测试次数，默认为4次

使用 `-s` 指定超时时间，默认为2秒

使用 `-t` 来启用无限次测试模式

使用 `-f` 来启用快速模式，降低两次成功测试之间的间隔时间
