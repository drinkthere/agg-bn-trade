#!/bin/bash

# 指定日志文件所在的目录
log_directory="/data/data/aggbntrade"

# 保留最近的两个okxmm.log文件
ls -t $log_directory/aggbntrade* | tail -n +7 | xargs rm -f


