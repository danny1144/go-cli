# cobra 命令行

- redis connect test
- mysql connect test
- kafka connect test
- websocket connect test

- 实时天气查询根据城市名称或者城市代码

````
根据城市代码  go run main.go weather -c 330100
根据城市名称  go run main.go weather -n 杭州市
````

## 打包

- linux 环境
```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build 
```
- window环境

```
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build 
```