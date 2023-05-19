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

func getTencentCloudStackArg(byte []byte, arg string) string {
	unsafe := blackfriday.MarkdownCommon(byte)
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	if html == nil || strings.Contains(string(html), "404: Not Found") {
		return ""
	}
	re := regexp.MustCompile(`(?s)<h2>Argument Reference</h2>.*?</ul>`)
	result := re.FindStringSubmatch(string(html))[0]

	codeRe := regexp.MustCompile(`<code>(.*?)</code>`)
	rt := codeRe.FindAllStringSubmatch(result, -1)
	for _, r := range rt {
		if r[1] == arg {
			return arg
		}
	}
	return ""
}

func getTencentCloudStackMarkdown(resourceName string, te string) []byte {
	url := "https://raw.githubusercontent.com/tencentcloudstack/terraform-provider-tencentcloud/master/website/docs/" + te + "/" + resourceName + ".html.markdown"
	// 根据URL获取资源
	res, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
	}
	// 读取资源数据 body: []byte
	body, err := ioutil.ReadAll(res.Body)
	// 关闭资源流
	res.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
	}
	//println(url)
	return body
}
