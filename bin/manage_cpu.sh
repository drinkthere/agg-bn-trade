#!/bin/bash

cpu_index=44
start() {
    echo "Starting the process @cpu '$cpu_index'..."
    nohup taskset -c "$cpu_index"  ./aggbntrade ../config/config.json >> /data/dc/aggbntrade/nohup.log 2>&1 &
    echo "Process started."
}

stop() {
    echo "Stopping the process"
    pid=$(pgrep -f "aggbntrade ../config/config.json")
    if [ -n "$pid" ]; then
        kill -SIGINT $pid
        echo "Process stopping: "
        sleep 10
        echo "Process stopped."
    else
        echo "Process is not running."
    fi
}

restart() {
    stop
    sleep 10
    start
}

case "$1" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    restart)
        restart
        ;;
    *)
        echo "Usage: $0 {start|stop|restart}"
        exit 1
        ;;
esac

exit 0
