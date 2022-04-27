<div align="center" >
    <img src="https://user-images.githubusercontent.com/26270009/165508782-4783f9c2-bb55-405f-ab83-e9c30f82072e.png" alt="Hyuga" />
</div>

## xssfinder 是什么？

基于被动代理和 ChromeBrowser 的 XSS 漏洞发现工具。

- 动态地语义分析网页中的`JavaScript`源码，Hook关键点，利用污点分析检出 Dom-Based XSS。

## 立即尝试

## 安装

## 用法

```bash
Usage of xssfinder:
  -bexecpath string
        Set browser exec path
  -bproxy string
        Set browser proxy addr
  -debug
        Set debug mode
  -incognito
        Set browser incognito mode
  -maddr string
        Set mitm-server listen address (default ":8222")
  -mhosts string
        Set mitm-server target hosts .e.g. foo.com,bar.io
  -mporxy string
        Set mitm-server parent proxy
  -mverbose
        Set mitm-server verbose mode
  -nocolor
        Set logger no-color mode
  -noheadless
        Set browser no-headless mode
  -outjson
        Set logger output json format
  -vv
        Set very-verbose mode
```

## TODO

- [ ] 支持检测反射XSS
- [ ] 优化Runner & Worker
- [ ] cmd parse