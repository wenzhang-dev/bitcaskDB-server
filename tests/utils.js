import crypto from 'k6/crypto'

export function genRandomBytes(length) {
    const bytes = crypto.randomBytes(length)
    return new Uint8Array(bytes)
}

export function genKeyStr(length) {
    const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
    let result = '';
    for (let i = 0; i < length; i++) {
        result += chars.charAt(Math.floor(Math.random() * chars.length));
    }
    return result;
}
