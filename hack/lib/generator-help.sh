#!/usr/bin/env bash

# Utility command based on 'find' command. The pipeline is as following:
#   1. find all the go files; (exclude specific path: vendor etc)
#   2. find all the files containing specific tags in contents;
#   3. extract related dirs;
#   4. remove duplicated paths;
#   5. merge all dirs in array with delimiter ,;
#
# Example:
#   find_dirs_containing_comment_tags("+k8s:")
# Return:
#   github.com/amit3512/descheduler_policy_master/a,github.com/amit3512/descheduler_policy_master/b,github.com/amit3512/descheduler_policy_master/c
function find_dirs_containing_comment_tags() {
   array=()
   while IFS='' read -r line; do array+=("$line"); done < <( \
     find . -type f -name \*.go -not -path "./vendor/*" -not -path "./_tmp/*" -print0  \
     | xargs -0 grep --color=never -l "$@" \
     | xargs -n1 dirname \
     | LC_ALL=C sort -u \
     )

   IFS=" ";
   printf '%s' "${array[*]}";
}
