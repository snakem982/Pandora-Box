package utils

import (
	"crypto/tls"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"sync"
	"time"
)

// ConnTimeOut 请求时间
var ConnTimeOut = time.Second * 15

// DialTimeOut 拨号时间
var DialTimeOut = time.Second * 5

// FastTimeOut 并发请求时间
var FastTimeOut = time.Second * 16

// SendGet 发送 GET 请求
func SendGet(requestURL string, headers map[string]string, proxyURL string) (string, http.Header, error) {
	// 创建 HTTP 客户端
	client := &http.Client{
		Timeout: ConnTimeOut,
	}

	// 创建 Transport 并允许不安全链接
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout: DialTimeOut, // 仅拨号阶段超时设置
		}).DialContext,
	}

	// 如果提供了代理路径，则设置代理
	if proxyURL != "" {
		proxy, err := url.Parse(proxyURL)
		if err != nil {
			return "", nil, fmt.Errorf("解析代理路径失败: %v", err)
		}
		transport.Proxy = http.ProxyURL(proxy)
	}

	// 设置 Transport
	client.Transport = transport

	// 创建 HTTP 请求
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return "", nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	if _, ok := headers["User-Agent"]; !ok {
		req.Header.Set("User-Agent", "clash-verge/v2.2.3")
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 发送 HTTP 请求
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer func() {
		if resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, fmt.Errorf("读取响应内容失败: %v", err)
	}

	return html.UnescapeString(string(body)), resp.Header, nil
}

type ResponseResult struct {
	Body    string
	Headers http.Header
}

// FastGet 并发 GET 请求
func FastGet(requestURL string, headers map[string]string, proxyURL string) (*ResponseResult, error) {
	var wg sync.WaitGroup
	results := make(chan *ResponseResult, 2)
	errors := make(chan error, 2)

	// 并发发送请求
	wg.Add(2)
	go func() {
		defer wg.Done()
		body, responseHeaders, err := SendGet(requestURL, headers, proxyURL)
		if err != nil {
			errors <- fmt.Errorf("使用代理请求失败: %v", err)
			return
		}
		if len(body) > 0 {
			results <- &ResponseResult{Body: body, Headers: responseHeaders}
		}
	}()
	go func() {
		defer wg.Done()
		body, responseHeaders, err := SendGet(requestURL, headers, "")
		if err != nil {
			errors <- fmt.Errorf("不使用代理请求失败: %v", err)
			return
		}
		if len(body) > 0 {
			results <- &ResponseResult{Body: body, Headers: responseHeaders}
		}
	}()

	// 等待 Goroutines 完成
	go func() {
		wg.Wait()
		close(results) // 确保只在所有 Goroutines 完成后关闭通道
		close(errors)
	}()

	// 优先返回第一个成功的结果
	select {
	case result := <-results:
		return result, nil
	case <-time.After(FastTimeOut): // 设置超时时间
		// 如果所有请求都失败，从错误通道中获取错误信息
		var proxyErr, directErr error
		for err := range errors {
			if err.Error()[:4] == "使用代理" {
				proxyErr = err
			} else {
				directErr = err
			}
		}
		// 优先返回使用代理的错误
		if proxyErr != nil {
			return nil, proxyErr
		}
		return nil, directErr
	}
}

// 创建正则表达式
var headPattern = `204|blank|generate|gstatic`

// SendHead 发送 Head 请求
func SendHead(requestURL string, proxyURL string) (int, error) {

	// 设置代理地址
	proxy, err := url.Parse(proxyURL) // 替换为你的代理地址
	if err != nil {
		return 500, fmt.Errorf("解析代理路径失败: %v", err)
	}

	// 创建Transport并设置代理
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxy),
	}

	// 创建HTTP客户端并设置Transport
	client := &http.Client{
		Transport: transport,
		Timeout:   8 * time.Second,
	}

	// 创建请求
	re, _ := regexp.Compile(headPattern)
	method := "GET"
	if re.MatchString(requestURL) {
		method = "HEAD"
	}
	req, err := http.NewRequest(method, requestURL, nil)
	if err != nil {
		return 500, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 13_3) AppleWebKit/537.36 (KHTML, like Gecko) Version/16.4 Safari/537.36")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return 500, fmt.Errorf("发送请求失败: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// 打印响应状态码
	return resp.StatusCode, nil
}
