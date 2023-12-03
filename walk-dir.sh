#!/bin/bash

all_files=false

while getopts a: flag
do
	case "${flag}" in
		a) all_files=true;;
	esac
done

list_files() {
	while IFS= read -r file; do
		if [[ -z "$file" || "$file" = "." || "$file" = ".." ]]; then # if line is empty or '.' or '..', skip
			continue
		fi
		FILE="$1/$file"
		if [ -d "$FILE" ]; then
			list_files "$FILE"
		else
			printf '%s\n' "$FILE"
		fi
	done <<< "$(ls $(if $all_files; then echo "-a"; fi) "$1")"
}

list_files "${@: -1}" # list files using the last cli argument
