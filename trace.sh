#!/bin/bash

if [ $# -ne 2 ]; then
    echo "Usage: $0 [<trace_id>] <time_param>"
fi

# 如果提供了时间参数，则使用提供的参数；否则使用当前时间
if [ $# -eq 2 ]; then
    trace_id=$1
    time_param=$2
    current_date=$(date +%Y%m%d)
else
    trace_id=$1
    current_date=$(date +%Y%m%d)
fi

log_folder="logs"

# 根据输入的时间参数长度动态构建日期和小时
if [ ${#time_param} -eq 0 ]; then
    datetime=$(date +%Y%m%d%H)
elif [ ${#time_param} -eq 2 ]; then
    # 传入的是小时部分
    datetime="${current_date}${time_param}"
elif [ ${#time_param} -eq 4 ]; then
    # 传入的是日时部分
    datetime="${current_date:0:6}${time_param}"
elif [ ${#time_param} -eq 6 ]; then
    # 传入的是月日时部分
    datetime="${current_date:0:4}${time_param}"
else
    echo "Invalid time format. Please provide 2, 4, or 6 digits for hour, day-hour, or month-day-hour respectively."
    exit 1
fi

# 构建文件名
file_to_search="${log_folder}/${datetime}.txt"

# 检查文件是否存在
if [ -f "$file_to_search" ]; then
    # 在文件中查找包含指定 trace_id 的行
    matching_lines=$(grep "$trace_id" "$file_to_search")

    if [ -n "$matching_lines" ]; then
        # 输出匹配的行
        echo "File: $file_to_search"
        echo "$matching_lines"
        echo "------------------------"
    else
        echo "No matching lines in $file_to_search"
    fi
else
    echo "File $file_to_search not found"
fi