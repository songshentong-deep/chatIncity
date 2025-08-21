package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

func main() {
	// 登录获取token
	token, err := login()
	if err != nil {
		fmt.Printf("登录失败: %v\n", err)
		return
	}
	
	// 测试图片服务上传
	testPhotoUpload(token)
}

func login() (string, error) {
	loginData := map[string]string{
		"phone":    "13900139000",
		"password": "password123",
	}
	
	jsonData, _ := json.Marshal(loginData)
	resp, err := http.Post("http://localhost:8001/api/v1/auth/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("登录失败")
	}
	
	body, _ := io.ReadAll(resp.Body)
	var loginResp map[string]interface{}
	json.Unmarshal(body, &loginResp)
	data := loginResp["data"].(map[string]interface{})
	return data["token"].(string), nil
}

func testPhotoUpload(token string) {
	fmt.Println("=== 测试图片服务上传 ===")
	
	// 创建简单的文本文件作为测试（避免图片格式问题）
	testData := []byte("fake image data for testing")
	
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	
	writer.WriteField("type", "avatar")
	part, _ := writer.CreateFormFile("photo", "test.jpg")
	part.Write(testData)
	writer.Close()
	
	req, _ := http.NewRequest("POST", "http://localhost:8009/api/v1/photos", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("发送请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("图片服务状态: %s\n", resp.Status)
	fmt.Printf("图片服务响应: %s\n", string(body))
	
	if resp.StatusCode == 200 {
		fmt.Println("✅ 图片服务上传成功！")
	} else if resp.StatusCode == 401 {
		fmt.Println("❌ 认证失败 - 检查JWT token")
	} else {
		fmt.Printf("❌ 上传失败 - 状态码: %d\n", resp.StatusCode)
	}
}