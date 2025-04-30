#!/bin/bash

if [ -z "$1" ]; then
    echo "Usage: $0 <server addr>"
    exit 1
fi

SERVER_ADDR=$1
echo "Target Server Addr: $SERVER_ADDR"

TEST_DIR=./tests
GRAFANA_K6=swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/grafana/k6:0.53.0

for script in $TEST_DIR/test_*.js; do
    test_name=$(basename $script .js)
    echo "Running test: $test_name"
    docker run \
        --user 0 \
        --rm --network host \
        -e SERVER_ADDR=$SERVER_ADDR \
        -v $(pwd)/tests:/tests \
        $GRAFANA_K6 run /tests/$test_name.js --summary-export=/tests/${test_name}.json
done

echo "All tests finished"
