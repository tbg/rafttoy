#!/bin/bash
set -euo pipefail

STARTOPTS="--pipeline=basic --data-dir=/mnt/data1"
TESTS="BenchmarkRaft/conc=256/bytes=256"
COUNT="${COUNT-10}"

GO111MODULE=on go mod tidy -v
GO111MODULE=on go mod vendor
GO111MODULE=on GOOS=linux make build
roachprod stop $CLUSTER

PEERS=$(roachprod ip $CLUSTER | sed 's/$/:10000/' | paste -sd "," -)

for i in 2 3; do
	roachprod put $CLUSTER:$i rafttoy-follower
	roachprod run $CLUSTER:$i -- ./rafttoy-follower --id="${i}" --peers="${PEERS}" ${STARTOPTS} &> /dev/null &
done

CMD="./rafttoy-leader --id=1 --peers=${PEERS} ${STARTOPTS}"

roachprod put $CLUSTER:1 rafttoy-leader
if [ -z "${BENCHPROFILE-}" ]; then
	# Run benchmarks and collect results.
	roachprod run $CLUSTER:1 -- ${CMD} --test.bench=${TESTS} \
		--test.benchtime=2s --test.count=3
	exit 0
fi

roachprod run $CLUSTER:1 -- mkdir -p profiles
# Run only for the profiles.
roachprod run $CLUSTER:1 -- ${CMD} --test.bench=${TESTS} \
 	--test.memprofile "profiles/mem.pprof" \
 	--test.cpuprofile "profiles/cpu.pprof" \
 	--test.benchtime=2s --test.count=1

roachprod get $CLUSTER:1 -- profiles "${BENCHPROFILE}"
