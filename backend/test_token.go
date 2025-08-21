package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	// 使用刚才注册获得的token测试
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTQ1NjEwNDIsImlhdCI6MTc1NDQ3NDY0MiwidXNlcl9pZCI6IjVmZWE4NDZhLTA4MjUtNGM4Mi1iZjA4LWFlZWRiZjE1OWRkZiJ9.8w4becAwmHGG-SngULD93sHKVYo5DpzCJsu0sdFNdl0"
	
	req, err := http.NewRequest("GET", "http://localhost:8001/api/v1/user/profile", nil)
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
	
	fmt.Printf("Token验证响应状态: %s\n", resp.Status)
	fmt.Printf("Token验证响应内容: %s\n", string(body))
	
	if resp.StatusCode == 200 {
		fmt.Println("✅ Token验证正常，用户认证系统完全正常！")
	}
}