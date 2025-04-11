const UNITS = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];

export function prettyBytes(n: number) {
    if (n < 1000) {
        return n + ' B';
    }
    const exponent = Math.min(Math.floor(Math.log10(n) / 3), UNITS.length - 1);
    n = Number((n / Math.pow(1000, exponent)).toPrecision(3));
    const unit = UNITS[exponent];
    return n + ' ' + unit;
}

export function cJoin(arr: any, separator = ',') {
    let result = '';
    for (let i = 0; i < arr.length; i++) {
        result += arr[i];
        if (i < arr.length - 1) {
            result += separator;
        }
    }
    return result;
}

export function rJoin(arr: any, separator = ',') {
    let result = '';
    for (let i = arr.length - 1; i >= 0; i--) {
        result += arr[i];
        if (i > 0) {
            result += separator;
        }
    }
    return result;
}