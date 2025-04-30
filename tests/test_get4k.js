import { check } from 'k6'
import http from 'k6/http'
import { genRandomBytes } from './utils.js'

export let options = {
    vus: 1000,
    duration: '60s',
}

let ns = ""
let value = genRandomBytes(4096)

let addr = __ENV.SERVER_ADDR || '127.0.0.1:8090';

const total_keys = 10000;

export function setup() {
    console.log(`Preloading ${total_keys} records`)

    let keys = [];
    if (__ENV.SKIP_PRELOAD) {
        for (let i = 0; i < total_keys; i++) {
            keys.push(`key_${i}`);
        }
        return keys;
    }

    for (let i = 0; i < total_keys; i++) {
        let key = `key_${i}`;
        let url = `http://${addr}/api/v1/kv?ns=${ns}&key=${key}`
        let res = http.post(url, value);

        if (res.status == 200) {
            keys.push(key);
        }
    }

    console.log(`Loaded ${keys.length} records`);

    return keys;
}

export default function(keys) {
    let key = keys[Math.floor(Math.random() * keys.length)];
    let url = `http://${addr}/api/v1/kv?ns=${ns}&key=${key}`;

    let res = http.get(url);
    check(res, {'status is 200': (r) => r.status == 200});
}
