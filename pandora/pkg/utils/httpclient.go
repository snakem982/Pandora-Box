package utils

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// SendGet 发送 GET 请求
func SendGet(requestURL string, headers map[string]string, proxyURL string) (string, http.Header, error) {
	// 创建 HTTP 客户端
	client := &http.Client{}

	// 允许不安全链接
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
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
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 发送 HTTP 请求
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, fmt.Errorf("读取响应内容失败: %v", err)
	}

	return string(body), resp.Header, nil
}

type ResponseResult struct {
	Body    string
	Headers map[string][]string
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

	// 等待 Goroutine 完成
	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	// 优先返回第一个成功的结果
	select {
	case result := <-results:
		return result, nil
	case <-time.After(15 * time.Second): // 设置超时时间
		// 如果所有请求都失败，从错误通道中获取错误信息
		close(errors)
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
