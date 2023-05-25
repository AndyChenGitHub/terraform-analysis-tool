package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

// getTencentCloudStackArgDesc 获取腾讯云的资源描述
func getTencentCloudStackArgDesc(tencentByte []byte, arg string) (string, string) {
	tencentByteMark := blackfriday.MarkdownCommon(tencentByte)
	markHtmlTencentClioud := bluemonday.UGCPolicy().SanitizeBytes(tencentByteMark)
	if markHtmlTencentClioud == nil || strings.Contains(string(markHtmlTencentClioud), "404: Not Found") {
		return "", ""
	}

	re := regexp.MustCompile(`(?s)<h2>Argument Reference</h2>.*?</ul>`)
	argumentTencent := re.FindStringSubmatch(string(markHtmlTencentClioud))[0]
	codeRe := regexp.MustCompile("<li><code>" + arg + "</code> - (.*?)</li>")
	rtTencent := codeRe.FindAllStringSubmatch(argumentTencent, -1)

	if len(rtTencent) > 0 {
		return arg, "tencentcloud: " + rtTencent[0][1]
	}
	return "", ""
}

// getAliyunArgDesc 获取阿里云的参数描述
func getAliyunArgDesc(aliyunByte []byte, arg string) string {
	aliyunByteMark := blackfriday.MarkdownCommon(aliyunByte)
	markHtmlAliyun := bluemonday.UGCPolicy().SanitizeBytes(aliyunByteMark)
	if markHtmlAliyun == nil || strings.Contains(string(markHtmlAliyun), "404: Not Found") {
		return ""
	}

	re := regexp.MustCompile(`(?s)<h2>Argument Reference</h2>.*?<h2>Attributes Reference</h2>`)
	argumentAliyun := re.FindStringSubmatch(string(markHtmlAliyun))[0]
	codeRe := regexp.MustCompile("<code>" + arg + "</code> - (.+?)(</li>|</p>|\\.\\n)")
	rtAliyun := codeRe.FindAllStringSubmatch(argumentAliyun, -1)
	if len(rtAliyun) > 0 {
		return "aliyun: " + rtAliyun[0][1] + "\n\n"
	}

	return ""
}

// getTencentCloudStackMarkdown 获取腾讯云的文档资源
func getTencentCloudStackMarkdown(resourceName string, te string) []byte {
	url := "https://raw.githubusercontent.com/tencentcloudstack/terraform-provider-tencentcloud/master/website/docs/" + te + "/" + resourceName + ".html.markdown"
	return getCloudMarkdown(url)
}

// getAliyunMarkdown 获取阿里云的文档资源
func getAliyunMarkdown(resourceName string, te string) []byte {
	url := "https://raw.githubusercontent.com/aliyun/terraform-provider-alicloud/master/website/docs/" + te + "/" + resourceName + ".html.markdown"
	return getCloudMarkdown(url)
}

// getCloudMarkdown 获取云的文档资源
func getCloudMarkdown(url string) []byte {
	// 根据URL获取资源
	res, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
	}
	// 读取资源数据 body: []byte
	body, err := ioutil.ReadAll(res.Body)
	// 关闭资源流
	defer res.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
	}
	return body
}
