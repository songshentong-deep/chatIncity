package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type LoginResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Token string `json:"token"`
		User  struct {
			ID       string `json:"id"`
			Username string `json:"username"`
		} `json:"user"`
	} `json:"data"`
}

func main() {
	// 1. 先登录获取token
	token, err := login()
	if err != nil {
		fmt.Printf("登录失败: %v\n", err)
		return
	}
	
	fmt.Printf("登录成功，获取到token: %s...\n", token[:20])
	
	// 2. 测试图片上传
	testPhotoUpload(token)
}

func login() (string, error) {
	loginData := map[string]string{
		"username": "testuser",
		"password": "password123",
	}
	
	jsonData, _ := json.Marshal(loginData)
	
	resp, err := http.Post("http://localhost:8002/api/v1/auth/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("发送登录请求失败: %v", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取登录响应失败: %v", err)
	}
	
	fmt.Printf("登录响应状态: %s\n", resp.Status)
	fmt.Printf("登录响应内容: %s\n", string(body))
	
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("登录失败，状态码: %d", resp.StatusCode)
	}
	
	var loginResp LoginResponse
	if err := json.Unmarshal(body, &loginResp); err != nil {
		return "", fmt.Errorf("解析登录响应失败: %v", err)
	}
	
	return loginResp.Data.Token, nil
}

func testPhotoUpload(token string) {
	fmt.Println("\n=== 测试图片上传 ===")
	
	// 创建测试图片数据
	testImageData := createTestImage()
	
	// 创建multipart表单
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	
	// 添加照片类型字段
	writer.WriteField("type", "avatar")
	
	// 添加文件字段
	part, err := writer.CreateFormFile("photo", "test_avatar.jpg")
	if err != nil {
		fmt.Printf("创建表单文件失败: %v\n", err)
		return
	}
	
	part.Write(testImageData)
	writer.Close()
	
	// 发送POST请求
	req, err := http.NewRequest("POST", "http://localhost:8009/api/v1/photos", &buf)
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
		return
	}
	
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	
	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("发送请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取响应失败: %v\n", err)
		return
	}
	
	fmt.Printf("上传响应状态: %s\n", resp.Status)
	fmt.Printf("上传响应内容: %s\n", string(body))
	
	// 如果上传成功，测试获取用户照片
	if resp.StatusCode == 200 {
		testGetUserPhotos(token)
	}
}

func testGetUserPhotos(token string) {
	fmt.Println("\n=== 测试获取用户照片 ===")
	
	req, err := http.NewRequest("GET", "http://localhost:8009/api/v1/photos?type=avatar", nil)
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
		return
	}
	
	req.Header.Set("Authorization", "Bearer "+token)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("发送请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取响应失败: %v\n", err)
		return
	}
	
	fmt.Printf("获取照片响应状态: %s\n", resp.Status)
	fmt.Printf("获取照片响应内容: %s\n", string(body))
}

// 创建一个简单的测试图片数据（JPEG格式）
func createTestImage() []byte {
	return []byte{
		0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46, 0x00, 0x01,
		0x01, 0x01, 0x00, 0x48, 0x00, 0x48, 0x00, 0x00, 0xFF, 0xDB, 0x00, 0x43,
		0x00, 0x08, 0x06, 0x06, 0x07, 0x06, 0x05, 0x08, 0x07, 0x07, 0x07, 0x09,
		0x09, 0x08, 0x0A, 0x0C, 0x14, 0x0D, 0x0C, 0x0B, 0x0B, 0x0C, 0x19, 0x12,
		0x13, 0x0F, 0x14, 0x1D, 0x1A, 0x1F, 0x1E, 0x1D, 0x1A, 0x1C, 0x1C, 0x20,
		0x24, 0x2E, 0x27, 0x20, 0x22, 0x2C, 0x23, 0x1C, 0x1C, 0x28, 0x37, 0x29,
		0x2C, 0x30, 0x31, 0x34, 0x34, 0x34, 0x1F, 0x27, 0x39, 0x3D, 0x38, 0x32,
		0x3C, 0x2E, 0x33, 0x34, 0x32, 0xFF, 0xC0, 0x00, 0x11, 0x08, 0x00, 0x01,
		0x00, 0x01, 0x01, 0x01, 0x11, 0x00, 0x02, 0x11, 0x01, 0x03, 0x11, 0x01,
		0xFF, 0xC4, 0x00, 0x14, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08, 0xFF, 0xC4,
		0x00, 0x14, 0x10, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xDA, 0x00, 0x0C,
		0x03, 0x01, 0x00, 0x02, 0x11, 0x03, 0x11, 0x00, 0x3F, 0x00, 0x8A, 0x00,
		0xFF, 0xD9,
	}
}