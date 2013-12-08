go-daemon 
=========

A library for writing system daemons in golang.

Installation
------------

	go get github.com/sevlyar/go-daemon

Documentation
-------------

[http://godoc.org/github.com/sevlyar/go-daemon](http://godoc.org/github.com/sevlyar/go-daemon)

Usage
-----

	import "github.com/sevlyar/go-daemon"

	...
	daemon.Reborn(027, "/")
	...

Build status
------------
[![Build Status](https://travis-ci.org/sevlyar/go-daemon.png?branch=master)](https://travis-ci.org/sevlyar/go-daemon)
