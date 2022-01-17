#!/bin/bash

set -eu

my_dir=$(cd $(dirname $0); pwd)
target_spexec=${my_dir}/../dist/spexec_$(go env GOOS)_$(go env GOARCH)/spexec
tester_spexec="${target_spexec}"

have_error=no
for spec in $my_dir/spec/*.yaml; do
  echo $(basename ${spec})
  SPEXEC="${target_spexec}" "${target_spexec}" --strict "${spec}" || have_error=yes
  echo
done

test "${have_error}" = 'no'
