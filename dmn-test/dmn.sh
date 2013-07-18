#!/bin/bash


WORK_DIR=../../../bin
PID_FILE=dmn.pid
LOG_FILE=dmn.log


cd $WORK_DIR


PID=
getpid() {
	if ./dmn-test --status --silent; then
		echo "daemon is not running"
		exit
	else
		PID=`cat $PID_FILE`
	fi
}

case "$1" in
	start)
		if ./dmn-test; then
			echo "starting daemon: OK"
		else
			echo "daemon return error code: $?"
		fi
		;;

	stop)
		getpid
		kill -TERM $PID
		echo "stopping daemon: OK"
		cat $LOG_FILE
		;;

	status)
		getpid
		echo "daemon pid: $PID"
		;;

	reload)
		getpid
		kill -HUP $PID
		echo "reloading daemon config: OK"
		;;

	crash)
		getpid
		kill -USR1 $PID
		echo "crashing daemon: OK"
		cat $LOG_FILE
		;;

	*)
		echo "Usage: dmn.sh {start|stop|status|reload|crash}"
esac