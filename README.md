# bitcaskDB-server

This repository implements a simple http server and supports bitcaskDB CRUD operations. And the grafana K6 tool is used to bench the server.


## Quick start

The server and all testcases are run in docker.

```shell
make build_docker
./run_tests.sh 127.0.0.1:8090
```

The results are as follows:

```
./run_tests.sh 127.0.0.1:8090
Target Server Addr: 127.0.0.1:8090
Running test: test_get4k
     execution: local
        script: /tests/test_get4k.js
        output: -

     scenarios: (100.00%) 1 scenario, 1000 max VUs, 1m30s max duration (incl. graceful stop):
              * default: 1000 looping VUs for 1m0s (gracefulStop: 30s)

time="2025-04-30T05:57:46Z" level=info msg="Preloading 10000 records" source=console
     ✓ status is 200

     checks.........................: 100.00% ✓ 708055       ✗ 0
     data_received..................: 3.0 GB  46 MB/s
     data_sent......................: 117 MB  1.8 MB/s
     http_req_blocked...............: avg=96.46µs  min=701ns    med=1.8µs   max=1.36s   p(90)=3.53µs   p(95)=4.53µs
     http_req_connecting............: avg=85.35µs  min=0s       med=0s      max=1.32s   p(90)=0s       p(95)=0s
     http_req_duration..............: avg=72.75ms  min=112.41µs med=55.94ms max=1.03s   p(90)=152.45ms p(95)=190.84ms
       { expected_response:true }...: avg=72.75ms  min=112.41µs med=55.94ms max=1.03s   p(90)=152.45ms p(95)=190.84ms
     http_req_failed................: 0.00%   ✓ 0            ✗ 718055
     http_req_receiving.............: avg=491.79µs min=8.37µs   med=25.03µs max=1.03s   p(90)=50.82µs  p(95)=119.26µs
     http_req_sending...............: avg=156.07µs min=4.24µs   med=8.36µs  max=1.02s   p(90)=18.17µs  p(95)=33.07µs
     http_req_tls_handshaking.......: avg=0s       min=0s       med=0s      max=0s      p(90)=0s       p(95)=0s
     http_req_waiting...............: avg=72.1ms   min=0s       med=55.67ms max=920.6ms p(90)=151.48ms p(95)=189.21ms
     http_reqs......................: 718055  10981.825883/s
     iteration_duration.............: avg=83.24ms  min=249.08µs med=64.65ms max=5.32s   p(90)=168.31ms p(95)=213.03ms
     iterations.....................: 708055  10828.887377/s
     vus............................: 1000    min=0          max=1000
     vus_max........................: 1000    min=1000       max=1000

Running test: test_healthz
     execution: local
        script: /tests/test_healthz.js
        output: -

     scenarios: (100.00%) 1 scenario, 1000 max VUs, 1m30s max duration (incl. graceful stop):
              * default: 1000 looping VUs for 1m0s (gracefulStop: 30s)

     ✓ status is 200

     checks.........................: 100.00% ✓ 1365739      ✗ 0
     data_received..................: 102 MB  1.7 MB/s
     data_sent......................: 128 MB  2.1 MB/s
     http_req_blocked...............: avg=41.67µs  min=672ns    med=1.5µs   max=195ms    p(90)=2.64µs  p(95)=3.31µs
     http_req_connecting............: avg=32.17µs  min=0s       med=0s      max=138.97ms p(90)=0s      p(95)=0s
     http_req_duration..............: avg=35.97ms  min=141.5µs  med=26.88ms max=411.24ms p(90)=69.59ms p(95)=93.95ms
       { expected_response:true }...: avg=35.97ms  min=141.5µs  med=26.88ms max=411.24ms p(90)=69.59ms p(95)=93.95ms
     http_req_failed................: 0.00%   ✓ 0            ✗ 1365739
     http_req_receiving.............: avg=1.93ms   min=5µs      med=12.54µs max=290.91ms p(90)=64.29µs p(95)=3.48ms
     http_req_sending...............: avg=446.03µs min=3.55µs   med=6.85µs  max=324.99ms p(90)=14.64µs p(95)=55.67µs
     http_req_tls_handshaking.......: avg=0s       min=0s       med=0s      max=0s       p(90)=0s      p(95)=0s
     http_req_waiting...............: avg=33.58ms  min=122.93µs med=26.47ms max=343.66ms p(90)=63.89ms p(95)=82.45ms
     http_reqs......................: 1365739 22732.334181/s
     iteration_duration.............: avg=43.49ms  min=191.87µs med=32.63ms max=426.56ms p(90)=84.41ms p(95)=112.58ms
     iterations.....................: 1365739 22732.334181/s
     vus............................: 1000    min=1000       max=1000
     vus_max........................: 1000    min=1000       max=1000

Running test: test_put4k
     execution: local
        script: /tests/test_put4k.js
        output: -

     scenarios: (100.00%) 1 scenario, 1000 max VUs, 1m30s max duration (incl. graceful stop):
              * default: 1000 looping VUs for 1m0s (gracefulStop: 30s)

     ✓ status is 200

     checks.........................: 100.00% ✓ 683182       ✗ 0
     data_received..................: 51 MB   854 kB/s
     data_sent......................: 2.9 GB  48 MB/s
     http_req_blocked...............: avg=132.73µs min=683ns    med=1.58µs  max=239.82ms p(90)=3.06µs   p(95)=3.87µs
     http_req_connecting............: avg=114.41µs min=0s       med=0s      max=239.73ms p(90)=0s       p(95)=0s
     http_req_duration..............: avg=81.78ms  min=413.12µs med=70.53ms max=575.18ms p(90)=140.12ms p(95)=171.41ms
       { expected_response:true }...: avg=81.78ms  min=413.12µs med=70.53ms max=575.18ms p(90)=140.12ms p(95)=171.41ms
     http_req_failed................: 0.00%   ✓ 0            ✗ 683182
     http_req_receiving.............: avg=1.71ms   min=7.31µs   med=14.94µs max=317.41ms p(90)=122.66µs p(95)=328.88µs
     http_req_sending...............: avg=701.95µs min=15.37µs  med=34.41µs max=408.5ms  p(90)=91µs     p(95)=183.26µs
     http_req_tls_handshaking.......: avg=0s       min=0s       med=0s      max=0s       p(90)=0s       p(95)=0s
     http_req_waiting...............: avg=79.37ms  min=0s       med=69.93ms max=441.43ms p(90)=135.02ms p(95)=161.67ms
     http_reqs......................: 683182  11380.654498/s
     iteration_duration.............: avg=87.48ms  min=1.72ms   med=73.73ms max=575.25ms p(90)=150.4ms  p(95)=190.2ms
     iterations.....................: 683182  11380.654498/s
     vus............................: 1000    min=1000       max=1000
     vus_max........................: 1000    min=1000       max=1000


running (1m00.0s), 0000/1000 VUs, 683182 complete and 0 interrupted iterations
default ✓ [ 100% ] 1000 VUs  1m0s
All tests finished
```
