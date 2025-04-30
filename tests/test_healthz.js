import http from 'k6/http'
import { check } from 'k6'

export let options = {
    vus: 1000,
    duration: '60s',
}

let addr = __ENV.SERVER_ADDR || '127.0.0.1:8090';

export default function() {
    let res = http.get(`http://${addr}/api/v1/healthz`);
    check(res, {'status is 200': (r) => r.status == 200});
}
