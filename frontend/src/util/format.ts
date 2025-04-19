const UNITS = ["B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"];

// 字节数转换为人类可读的格式
export function prettyBytes(n: number) {
  if (n < 1000) {
    return n + " B";
  }
  const exponent = Math.min(Math.floor(Math.log10(n) / 3), UNITS.length - 1);
  n = Number((n / Math.pow(1000, exponent)).toPrecision(3));
  const unit = UNITS[exponent];
  return n + " " + unit;
}

// 拼接字符串
export function cJoin(arr: any, separator = ",") {
  let result = "";
  for (let i = 0; i < arr.length; i++) {
    result += arr[i];
    if (i < arr.length - 1) {
      result += separator;
    }
  }
  return result;
}

// 反向拼接字符串
export function rJoin(arr: any, separator = ",") {
  let result = "";
  for (let i = arr.length - 1; i >= 0; i--) {
    result += arr[i];
    if (i > 0) {
      result += separator;
    }
  }
  return result;
}

// 将 Date 对象格式化为 "YYYY-MM-DD HH:mm:ss.SSS" 的字符串
export function formatDate(date: Date): string {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, "0"); // 月份从 0 开始，需要加 1
  const day = String(date.getDate()).padStart(2, "0");
  const hours = String(date.getHours()).padStart(2, "0");
  const minutes = String(date.getMinutes()).padStart(2, "0");
  const seconds = String(date.getSeconds()).padStart(2, "0");
  const milliseconds = String(date.getMilliseconds()).padStart(3, "0"); // 毫秒部分

  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}.${milliseconds}`;
}

// 校验是否 url 格式
export function isHttpOrHttps(url: any) {
  const regex = /^(https?):\/\/[\w.-]+(?:\.[\w.-]+)+[\w\-._~:/?#[\]@!$&'()*+,;=.]+$/;
  return regex.test(url);
}
