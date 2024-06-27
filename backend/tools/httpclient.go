package tools

import (
	"crypto/tls"
	"fmt"
	C "github.com/metacubex/mihomo/constant"
	"github.com/metacubex/mihomo/hub/executor"
	"github.com/metacubex/mihomo/log"
	"golang.org/x/net/context"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

var dialerBaidu = &net.Dialer{
	Resolver: &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Duration(5000) * time.Millisecond,
			}
			return d.DialContext(ctx, "udp", "180.76.76.76:53")
		},
	},
}

var dialBaiduContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
	return dialerBaidu.DialContext(ctx, network, addr)
}

// GetFileName 获取响应头中的文件名
func GetFileName(header http.Header) (fileName string) {
	disposition := header.Get("content-disposition")
	if disposition != "" {
		split := strings.Split(disposition, "=")
		fileName, _ = url.QueryUnescape(split[len(split)-1])
		fileName = strings.TrimLeft(fileName, "UTF-8''")
		fileName = strings.TrimLeft(fileName, "utf-8''")
	}

	return fileName
}

// HttpGetByProxy 使用代理访问指定的URL并返回响应数据
func HttpGetByProxy(requestUrl string) ([]byte, string, error) {
	// 拼接代理地址
	proxyUrl := fmt.Sprintf("http://127.0.0.1:%d", executor.GetGeneral().MixedPort)
	uri, _ := url.Parse(proxyUrl)

	// 创建一个带代理的HTTP客户端
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           http.ProxyURL(uri),
		},
		Timeout: 20 * time.Second,
	}
	priUrl := "https://github.com/snakem982/Pandora-Box/releases/download"
	if strings.HasPrefix(requestUrl, priUrl) {
		client.Timeout = time.Minute
	}

	// 创建一个GET请求
	req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
	if err != nil {
		log.Warnln("HttpGetByProxy http.NewRequest %s %v", requestUrl, err)
		return nil, "", err
	}
	req.Header.Set("Accept-Encoding", "utf-8")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", C.UA)

	// 发送请求并获取响应
	resp, err := client.Do(req)
	if err != nil {
		log.Warnln("HttpGetByProxy client.Do %s %v", requestUrl, err)
		return nil, "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			// 处理异常情况
		}
	}(resp.Body)
	// 读取响应数据
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Warnln("HttpGetByProxy io.ReadAll %s %v", requestUrl, err)
		return nil, "", err
	}

	if resp.StatusCode != http.StatusOK {
		log.Warnln("HttpGetByProxy StatusCode %s %d", requestUrl, resp.StatusCode)
		return nil, "", fmt.Errorf("StatusCode %d", resp.StatusCode)
	}

	return data, GetFileName(resp.Header), nil
}

// HttpGet 使用HTTP GET方法请求指定的URL，并返回响应的数据和可能的错误。
func HttpGet(requestUrl string) ([]byte, string, error) {
	timeOut := 30 * time.Second

	priUrl := "https://github.com/snakem982/Pandora-Box/releases/download"
	if strings.Contains(requestUrl, priUrl) {
		timeOut = time.Minute
	}

	return HttpGetWithTimeout(requestUrl, timeOut, true)
}

// HttpGetWithTimeout 使用HTTP GET方法请求指定的URL，并返回响应的数据和可能的错误。
func HttpGetWithTimeout(requestUrl string, outTime time.Duration, needDail bool) ([]byte, string, error) {
	client := http.Client{
		Timeout: outTime, // 请求超时时间
	}

	if needDail {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 安全证书验证关闭
			DialContext:     dialBaiduContext,
		}
	}

	req, err := http.NewRequest(http.MethodGet, requestUrl, nil) // 创建一个新的GET请求
	if err != nil {
		log.Warnln("HttpGetWithTimeout http.NewRequest %s %v", requestUrl, err)
		return nil, "", err
	}
	req.Header.Set("Accept-Encoding", "utf-8") // 设置响应内容编码为utf-8
	req.Header.Set("Accept", "*/*")            // 设置响应内容类型为全部
	req.Header.Set("User-Agent", C.UA)         // 设置用户代理为C.UA

	resp, err := client.Do(req) // 发送请求并获取响应
	if err != nil {
		log.Warnln("HttpGetWithTimeout client.Do %s %v", requestUrl, err)
		return nil, "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			// 处理关闭响应体的错误
		}
	}(resp.Body)
	data, err := io.ReadAll(resp.Body) // 读取响应体的数据
	if err != nil {
		log.Warnln("HttpGetWithTimeout io.ReadAll %s %v", requestUrl, err)
		return nil, "", err
	}

	if resp.StatusCode != http.StatusOK {
		log.Warnln("HttpGetWithTimeout StatusCode %s %d", requestUrl, resp.StatusCode)
		return nil, "", fmt.Errorf("StatusCode %d", resp.StatusCode)
	}

	return data, GetFileName(resp.Header), nil // 返回响应数据和无错误
}

// HttpGetWithDial 使用HTTP GET方法请求指定的URL，并返回响应的数据和可能的错误。
func HttpGetWithDial(requestUrl, fastlyHost, host string) ([]byte, error) {
	dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
		return net.Dial(network, fastlyHost)
	}

	client := http.Client{
		Transport: &http.Transport{
			DialContext:     dialContext,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 20 * time.Second,
	}

	req, _ := http.NewRequest("GET", requestUrl, nil)
	req.Header.Set("Accept-Encoding", "utf-8") // 设置响应内容编码为utf-8
	req.Header.Set("Accept", "*/*")            // 设置响应内容类型为全部
	req.Header.Set("User-Agent", C.UA)         // 设置用户代理为C.UA
	if host != "" {
		req.Header.Set("Host", host)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Warnln("HttpGetWithDial client.Do %s %v", requestUrl, err)
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Warnln("HttpGetWithDial io.ReadAll %s %v", requestUrl, err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Warnln("HttpGetWithDial StatusCode %s %d", requestUrl, resp.StatusCode)
		return nil, fmt.Errorf("StatusCode %d", resp.StatusCode)
	}

	return data, nil
}

// ConcurrentHttpGet 并发获取指定URL的HTTP内容
func ConcurrentHttpGet(url string) (all []byte, fileName string) {
	// 开启多线程请求
	cLock := sync.Mutex{}
	done := make(chan bool, 1)
	length := 128
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()

		content, name, err := HttpGetByProxy(url)
		if err == nil {
			cLock.Lock()
			if all == nil && len(content) > length {
				all = content
				fileName = name
			}
			cLock.Unlock()
			done <- true
		}
	}()
	go func() {
		defer wg.Done()

		content, name, err := HttpGet(url)
		if err == nil {
			cLock.Lock()
			if all == nil && len(content) > length {
				all = content
				fileName = name
			}
			cLock.Unlock()
			done <- true
		}
	}()

	go func() {
		time.Sleep(5 * time.Second)
		wg.Wait()
		done <- true
	}()

	select {
	case <-done:
	case <-time.After(30 * time.Second):
	}

	return
}