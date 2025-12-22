#!/bin/bash

set -euo pipefail

file1="$1"
file2="$2"

n_lines=$(wc -l < "$file1")

file1Title=$(printf "%-30s" "File 1: $file1")
file2Title=$(printf "%-30s" "File 2: $file2")

for i in $(seq 1 "$n_lines"); do
    clear -x
    line1=$(sed -n "${i}p" "$file1")
    line2=$(sed -n "${i}p" "$file2")

    echo -e "\
line $i/$n_lines

$file1Title: $line1\x1b[0m

$file2Title: $line2\x1b[0m

"
    read -p "press enter to continue"
done