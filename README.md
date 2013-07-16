go-daemon
=========

A library for writing system daemons in golang.

Installation
------------

	go get github.com/sevlyar/go-daemon

Documentation
-------------

After installation run godoc as http server:

	godoc -http:8080

And open address [http://127.0.0.1:8080/pkg/github.com/sevlyar/go-daemon/](http://127.0.0.1:8080/pkg/github.com/sevlyar/go-daemon/) in browser.

Usage
-----

	import "github.com/sevlyar/go-daemon"

	...
	daemon.Reborn(027, "/")
	...