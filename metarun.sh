#!/bin/bash
set -euo pipefail

trap "kill -- -$$" EXIT

DIR="$1-$2"
mkdir -p "${DIR}"

roachprod list --details $CLUSTER &> "${DIR}/roachprod-list.txt"

(cd $(go env GOPATH)/src/go.etcd.io/etcd && git checkout $1)
./run.sh | tee "${DIR}/${1}.txt"
BENCHPROFILE="${DIR}/${1}" ./run.sh
(cd $(go env GOPATH)/src/go.etcd.io/etcd && git checkout $2)
./run.sh | tee "${DIR}/${2}.txt"
BENCHPROFILE="${DIR}/${2}" ./run.sh

benchstat "${DIR}/${1}.txt" "${DIR}/${2}.txt" > "${DIR}/result.txt"

git add .
git ci -m "$0 $1 $2"

cat "${DIR}/result.txt"
