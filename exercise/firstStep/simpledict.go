// 任务：发送词典，获取词典解释
// 涉及到调用API，发送http请求（用marshal请求序列化），解析json（响应反序列化）
// 基本思路都是用结构体序列化或者解析json，复杂的api用代码生成工具
// 解析api的post请求curl：https://curlconverter.com/go/
// 解析response body：https://transform.tools/json-to-go

package firststep

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// 为了序列化json，需要定义结构体
type DictRequest struct {
	TranslatType string `json:"trans_type"`
	Source       string `json:"source"`
	UserID       string `json:"user_id"`
}

// 解析responsebody，需要定义结构体
type DictResponse struct {
	Rc   int `json:"rc"`
	Wiki struct {
	} `json:"wiki"`
	Dictionary struct {
		Prons struct {
			EnUs string `json:"en-us"`
			En   string `json:"en"`
		} `json:"prons"`
		Explanations []string      `json:"explanations"`
		Synonym      []string      `json:"synonym"`
		Antonym      []string      `json:"antonym"`
		WqxExample   [][]string    `json:"wqx_example"`
		Entry        string        `json:"entry"`
		Type         string        `json:"type"`
		Related      []interface{} `json:"related"`
		Source       string        `json:"source"`
	} `json:"dictionary"`
}

// 找现成检查里面的network请求，找线程调用的api，copy curl(bash) -> go
func query(word string) {
	// 创建请求，Client内置参数，其中常用timeout表示最大超时
	client := &http.Client{}
	// json序列化，构建结构体对应json字段后调用json.Marshal
	// var data = strings.NewReader(`{"trans_type":"en2zh","source":"good"}`) // 创建请求
	// 使用序列化的json数据创建请求
	request := DictRequest{TranslatType: "en2zh", Source: word}
	buf, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
	}
	// Marshal后的数据是byte数组，需要转换为io.Reader
	var data = bytes.NewBuffer(buf)
	req, err := http.NewRequest("POST", "https://api.interpreter.caiyunai.com/v1/dict", data)
	if err != nil {
		log.Fatal(err)
	}
	// 设置请求头
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("DNT", "1")
	req.Header.Set("os-version", "")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36")
	req.Header.Set("app-name", "xy")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("device-id", "")
	req.Header.Set("os-type", "web")
	req.Header.Set("X-Authorization", "token:qgemv4jr1y38jyq6vhvi")
	req.Header.Set("Origin", "https://fanyi.caiyunapp.com")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://fanyi.caiyunapp.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cookie", "_ym_uid=16456948721020430059; _ym_d=1645694872")
	// 发起请求，拿到响应
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	// 拿到请求后的第一步，为了避免资源泄露，需要手动关闭资源
	defer resp.Body.Close()
	// 读取响应体，把流转换为字符串
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	// 防止请求失败，打印状态码
	if resp.StatusCode != 200 {
		log.Fatal("bad StatusCode", resp.StatusCode, "body:", string(bodyText))
	}
	// fmt.Printf("%s\n", bodyText)
	// 解析bodyText，反序列化json
	var dictResponse DictResponse
	// 反序列化，把json字符串转换为结构体（&dictResponse是指针）
	err = json.Unmarshal(bodyText, &dictResponse)
	if err != nil {

	}
	// 筛选我们需要的信息
	fmt.Println(word, "UK:", dictResponse.Dictionary.Prons.En, "US:", dictResponse.Dictionary.Prons.EnUs)
	for _, item := range dictResponse.Dictionary.Explanations {
		fmt.Println(item)
	}
}

func simpleDict() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, `usage: simpleDict WORD
example: simpleDict hello
		`)
		os.Exit(1)
	}
	word := os.Args[1]
	query(word)
}
