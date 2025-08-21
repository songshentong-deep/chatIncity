package storage

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"
)

type OBSClient struct {
	endpoint  string
	bucket    string
	accessKey string
	secretKey string
}

func NewOBSClient(endpoint, bucket string) *OBSClient {
	return &OBSClient{
		endpoint: endpoint,
		bucket:   bucket,
	}
}

func NewOBSClientWithAuth(endpoint, bucket, accessKey, secretKey string) *OBSClient {
	return &OBSClient{
		endpoint:  endpoint,
		bucket:    bucket,
		accessKey: accessKey,
		secretKey: secretKey,
	}
}

func (c *OBSClient) UploadFile(key string, data []byte, contentType string) (string, error) {
	// 如果没有配置认证信息，返回模拟URL
	if c.accessKey == "" || c.secretKey == "" {
		url := fmt.Sprintf("https://%s.obs.cn-south-4.myhuaweicloud.com/%s", c.bucket, key)
		fmt.Printf("模拟上传到OBS: %s, 大小: %d bytes\n", url, len(data))
		return url, nil
	}
	
	// 使用虚拟主机域名访问桶
	url := fmt.Sprintf("https://%s.obs.cn-south-4.myhuaweicloud.com/%s", c.bucket, key)
	
	// 创建 HTTP PUT 请求
	req, err := http.NewRequest("PUT", url, bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}
	
	// 设置必要的头信息
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(data)))
	req.Header.Set("x-obs-acl", "public-read")
	
	// 生成签名
	date := time.Now().UTC().Format(time.RFC1123)
	req.Header.Set("Date", date)
	
	auth := c.generateSignature("PUT", fmt.Sprintf("/%s/%s", c.bucket, key), contentType, date, req.Header)
	req.Header.Set("Authorization", fmt.Sprintf("OBS %s:%s", c.accessKey, auth))
	
	fmt.Printf("上传到OBS: %s, 大小: %d bytes\n", url, len(data))
	
	// 发送请求
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()
	
	fmt.Printf("OBS响应状态码: %d\n", resp.StatusCode)
	
	// 检查响应状态
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Printf("上传成功: %s\n", url)
		return url, nil
	}
	
	// 读取错误响应
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("OBS错误响应: %s\n", string(body))
	
	return "", fmt.Errorf("上传失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
}

// 生成OBS签名
func (c *OBSClient) generateSignature(method, resource, contentType, date string, headers http.Header) string {
	// 构建OBS头部
	var obsHeaders []string
	for key, values := range headers {
		key = strings.ToLower(key)
		if strings.HasPrefix(key, "x-obs-") {
			for _, value := range values {
				obsHeaders = append(obsHeaders, fmt.Sprintf("%s:%s", key, value))
			}
		}
	}
	sort.Strings(obsHeaders)
	
	// 构建StringToSign - 注意格式
	var stringToSign string
	if len(obsHeaders) > 0 {
		stringToSign = fmt.Sprintf("%s\n\n%s\n%s\n%s\n%s",
			method,
			contentType,
			date,
			strings.Join(obsHeaders, "\n"),
			resource)
	} else {
		stringToSign = fmt.Sprintf("%s\n\n%s\n%s\n%s",
			method,
			contentType,
			date,
			resource)
	}
	
	fmt.Printf("StringToSign: %s\n", stringToSign)
	
	// HMAC-SHA1签名
	h := hmac.New(sha1.New, []byte(c.secretKey))
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}