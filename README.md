# go-daemon 

A library for writing system daemons in golang.

## Installation

	go get github.com/sevlyar/go-daemon

## Documentation

[http://godoc.org/github.com/sevlyar/go-daemon](http://godoc.org/github.com/sevlyar/go-daemon)

[http://gowalker.org/github.com/sevlyar/go-daemon](http://gowalker.org/github.com/sevlyar/go-daemon)

## Idea

```go
func main() {
	Pre()

	context := new(Context)
	child, _ := context.Reborn()

	if child != nil {
		PostParent()
	} else {
		defer context.Release()
		PostChild()
	}
}
```

![](https://github.com/sevlyar/go-daemon/raw/master/img/idea.png)

## Build status

[![Build Status](https://travis-ci.org/sevlyar/go-daemon.png?branch=master)](https://travis-ci.org/sevlyar/go-daemon)

## History

### 14.01.12
* released new major version, old version moved to github.com/sevlyar/go-daemon/oldapi

