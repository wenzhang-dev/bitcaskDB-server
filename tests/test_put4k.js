import { check } from 'k6'
import http from 'k6/http'
import { genRandomBytes, genKeyStr } from './utils.js'

export let options = {
    vus: 1000,
    duration: '60s',
}

let ns = ""
let key = genKeyStr(25)
let value = genRandomBytes(4096)

let addr = __ENV.SERVER_ADDR || '127.0.0.1:8090';
let url = `http://${addr}/api/v1/kv?ns=${ns}&key=${key}`

export default function() {
    let res = http.post(url, value);
    check(res, {'status is 200': (r) => r.status == 200});
}
