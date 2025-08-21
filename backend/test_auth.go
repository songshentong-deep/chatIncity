package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	fmt.Println("=== 测试注册登录功能 ===")
	
	// 1. 测试注册
	if err := testRegister(); err != nil {
		fmt.Printf("注册测试失败: %v\n", err)
	}
	
	// 2. 测试登录
	if err := testLogin(); err != nil {
		fmt.Printf("登录测试失败: %v\n", err)
	}
}

func testRegister() error {
	fmt.Println("\n--- 测试用户注册 ---")
	
	registerData := map[string]interface{}{
		"phone":    "13900139000", // 使用新的手机号
		"email":    "newuser@example.com",
		"password": "password123",
		"nickname": "新测试用户",
		"age":      25,
		"gender":   "male",
	}
	
	jsonData, _ := json.Marshal(registerData)
	
	resp, err := http.Post("http://localhost:8001/api/v1/auth/register", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("发送注册请求失败: %v", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取注册响应失败: %v", err)
	}
	
	fmt.Printf("注册响应状态: %s\n", resp.Status)
	fmt.Printf("注册响应内容: %s\n", string(body))
	
	if resp.StatusCode == 200 || resp.StatusCode == 409 {
		fmt.Println("✅ 注册功能正常")
		return nil
	}
	
	return fmt.Errorf("注册失败，状态码: %d", resp.StatusCode)
}

func testLogin() error {
	fmt.Println("\n--- 测试用户登录 ---")
	
	loginData := map[string]string{
		"phone":    "13800138000", // 使用已存在的手机号
		"password": "password123",
	}
	
	jsonData, _ := json.Marshal(loginData)
	
	resp, err := http.Post("http://localhost:8001/api/v1/auth/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("发送登录请求失败: %v", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取登录响应失败: %v", err)
	}
	
	fmt.Printf("登录响应状态: %s\n", resp.Status)
	fmt.Printf("登录响应内容: %s\n", string(body))
	
	if resp.StatusCode == 200 {
		fmt.Println("✅ 登录功能正常")
		
		// 解析token
		var loginResp map[string]interface{}
		if err := json.Unmarshal(body, &loginResp); err == nil {
			if data, ok := loginResp["data"].(map[string]interface{}); ok {
				if token, ok := data["token"].(string); ok {
					fmt.Printf("🔑 获取到JWT Token: %s...\n", token[:30])
					return testTokenValidation(token)
				}
			}
		}
		return nil
	}
	
	return fmt.Errorf("登录失败，状态码: %d", resp.StatusCode)
}

func testTokenValidation(token string) error {
	fmt.Println("\n--- 测试Token验证 ---")
	
	req, err := http.NewRequest("GET", "http://localhost:8001/api/v1/user/profile", nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}
	
	req.Header.Set("Authorization", "Bearer "+token)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %v", err)
	}
	
	fmt.Printf("Token验证响应状态: %s\n", resp.Status)
	fmt.Printf("Token验证响应内容: %s\n", string(body))
	
	if resp.StatusCode == 200 {
		fmt.Println("✅ Token验证正常")
	} else {
		fmt.Printf("⚠️  Token验证失败，状态码: %d\n", resp.StatusCode)
	}
	
	return nil
}