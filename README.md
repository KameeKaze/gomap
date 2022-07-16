# gomap

## Simple port scanner written in Golang

### Usage
`go run main.go -i 127.0.0.1 -p 22,80 -t 250` 

```
  -i string
        IP address or domain
  -p string
        Ports separated with comma 
        example: -p 22,80,443
  -t int
        Set the timeout in milliseconds (default 500)
```