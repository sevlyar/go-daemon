go-daemon
=========

A library for writing system daemons in golang.

Installation
------------

	go get github.com/sevlyar/go-daemon

Documentation
-------------

After installation run godoc ad http server:

	godoc -http:8080

And open in browser address 127.0.0.1:8080.

Usage
-----

	import "github.com/sevlyar/go-daemon"

	...
	daemon.Reborn(027, "/")
	...