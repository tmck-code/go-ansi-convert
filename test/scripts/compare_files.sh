#!/bin/bash

set -euo pipefail

file1="$1"
file2="$2"

n_lines=$(wc -l < "$file1")

longest_file_name_length=$(( ${#file1} > ${#file2} ? ${#file1} : ${#file2} ))

file1Title=$(printf "%-${longest_file_name_length}s" "$file1")
file2Title=$(printf "%-${longest_file_name_length}s" "$file2")

for i in $(seq 1 "$n_lines"); do
    clear -x
    line1=$(sed -n "${i}p" "$file1")
    line2=$(sed -n "${i}p" "$file2")

    echo "line $i/$n_lines"
    echo -ne "$file1Title: $line1"
    echo -e "\x1b[0m\n"

    echo -ne "$file2Title: $line2"
    echo -e "\x1b[0m"

    read -p "press enter to continue"
done