# Tcping GO

一个使用go写的小工具，用于检测TCP端口是否开放。

## 使用方法

```shell
tcping -h

Usage of tcping:
  -c, --count int          测试次数，默认为4次 (default 4)
  -h, --help               显示帮助信息
  -t, --infinite           无限次测试
  -6, --ipv6               使用 IPv6，需搭配域名使用
  -p, --port int           端口，默认为80 (default 80)
  -s, --timeout duration   超时时间，默认为2秒 (default 2s)
```
