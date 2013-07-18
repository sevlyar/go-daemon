#!/bin/bash


WORK_DIR=../../../bin
PID_FILE=dmn.pid
LOG_FILE=dmn.log


cd $WORK_DIR


if [ -f $PID_FILE ]; then
	PID=`cat $PID_FILE`
else
	PID="not found"
fi


case "$1" in
	start)
		if ./dmn-test; then
			echo "starting daemon: OK"
		else
			echo "daemon return error code: $?"
		fi
		;;

	stop)
		if [ -f $PID_FILE ]; then
			PID=`cat $PID_FILE`
			kill -TERM $PID
			cat $LOG_FILE
		fi
		;;

	status)
		echo "daemon pid: $PID"
		;;

	*)
		echo "Usage: dmn.sh {start|stop|status}"
esac