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

var UA = "clash-verge/v2.0.2"

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

func SetUA(ua string) {
	if ua != "clash.meta/"+C.Version {
		UA = ua
	} else {
		UA = "clash-verge/v2.0.2"
	}
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
func HttpGetByProxy(requestUrl string, headers map[string]string) ([]byte, string, error) {
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
	req.Header.Set("User-Agent", UA)
	if headers != nil && len(headers) > 0 {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

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
func HttpGet(requestUrl string, headers map[string]string) ([]byte, string, error) {
	timeOut := 30 * time.Second

	priUrl := "https://github.com/snakem982/Pandora-Box/releases/download"
	if strings.Contains(requestUrl, priUrl) {
		timeOut = time.Minute
	}

	return HttpGetWithTimeout(requestUrl, timeOut, true, headers)
}

// HttpGetWithTimeout 使用HTTP GET方法请求指定的URL，并返回响应的数据和可能的错误。
func HttpGetWithTimeout(requestUrl string, outTime time.Duration, needDail bool, headers map[string]string) ([]byte, string, error) {
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
	req.Header.Set("Accept-Encoding", "utf-8")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", UA)
	if headers != nil && len(headers) > 0 {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

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

// ConcurrentHttpGet 并发获取指定URL的HTTP内容
func ConcurrentHttpGet(url string, headers map[string]string) (all []byte, fileName string) {
	// 开启多线程请求
	cLock := sync.Mutex{}
	done := make(chan bool, 1)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()

		content, name, err := HttpGetByProxy(url, headers)
		if err == nil {
			cLock.Lock()
			if all == nil && len(content) > 0 {
				all = content
				fileName = name
			}
			cLock.Unlock()
			done <- true
		}
	}()

	go func() {
		defer wg.Done()

		content, name, err := HttpGet(url, headers)
		if err == nil {
			cLock.Lock()
			if all == nil && len(content) > 0 {
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
	case <-time.After(20 * time.Second):
	}

	return
}
