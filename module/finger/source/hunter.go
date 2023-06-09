package source

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gookit/color"
)

type HunterResults struct {
	IsRisk         string `json:"is_risk"`
	Url            string `json:"url"`
	Ip             string `json:"ip"`
	Port           int    `json:"port"`
	WebTitle       string `json:"web_title"`
	Domain         string `json:"domain"`
	IsRiskProtocol string `json:"is_risk_protocol"`
	Protocol       string `json:"protocol"`
	BaseProtocol   string `json:"base_protocol"`
	StatusCode     int    `json:"status_code"`
	Component      []struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"component"`
}
type HunterData struct {
	AccountType string          `json:"account_type"`
	Total       int             `json:"total"`
	Time        int             `json:"time"`
	Arr         []HunterResults `json:"arr"`
}

type HunterAutoGenerated struct {
	Code int        `json:"code"`
	Data HunterData `json:"data"`
}

func GetHunterConfig() Config {
	//创建一个空的结构体,将本地文件读取的信息放入
	c := &Config{}
	//创建一个结构体变量的反射
	cr := reflect.ValueOf(c).Elem()
	//打开文件io流
	f, err := os.Open(GetCurrentAbPathByExecutable() + "/config.ini")
	if err != nil {
		log.Fatal(err)
		color.RGBStyleFromString("237,64,35").Println("[Error] Hunter configuration file error!!!")
		os.Exit(1)
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	//我们要逐行读取文件内容
	s := bufio.NewScanner(f)

	for s.Scan() {
		//以=分割,前面为key,后面为value
		var str = s.Text()
		var index = strings.Index(str, "=")
		var key = strings.TrimSpace(str[0:index])
		var value = strings.TrimSpace(str[index+1:])
		//通过反射将字段设置进去
		cr.FieldByName(key).Set(reflect.ValueOf(value))
	}
	err = s.Err()
	if err != nil {
		log.Fatal(err)
	}
	//返回Config结构体变量
	return *c
}
func hunter_api(keyword string, apiKey string, page int) string {
	// 获取当前时间
	now := time.Now()
	// 计算六个月前的时间
	sixMonthsAgo := now.AddDate(0, -1, 0)
	// 将时间按照指定格式转换成字符串
	layout := "2006-01-02"
	nowStr := now.Format(layout)
	sixMonthsAgoStr := sixMonthsAgo.Format(layout)

	search := base64.URLEncoding.EncodeToString([]byte(keyword))
	api_request := fmt.Sprintf("https://hunter.qianxin.com/openApi/search?api-key=%v&search=%v&page=%v&page_size=100&is_web=1&start_time=%v&end_time=%v", apiKey, search, page, sixMonthsAgoStr, nowStr)

	return api_request
}

// 请求api
func hunterHttp(url string, timeout string) *HunterAutoGenerated {
	var itime, err = strconv.Atoi(timeout)
	if err != nil {
		log.Println("hunter超时参数错误: ", err)
	}
	transport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{
		Timeout:   time.Duration(itime) * time.Second,
		Transport: transport,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Accept", "*/*;q=0.8")
	req.Header.Set("Connection", "close")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2227.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)
	res := &HunterAutoGenerated{}
	json.Unmarshal(result, &res)
	return res
}

func Hunterip(ips string) (urls []string) {
	hunter := GetHunterConfig()

	keyword := `ip="` + ips + `"`

	for i := 1; i <= 10; i++ {
		//获取请求的api
		url := hunter_api(keyword, hunter.Hunter_key, i)
		//发起请求
		res := hunterHttp(url, hunter.Fofa_timeout)
		results := res.Data.Arr
		num := len(results)

		for _, result := range results {
			//fmt.Printf("%v\n", result.Url)
			urls = append(urls, result.Url)
		}
		if num < 100 {
			break
		}
	}

	return urls
}

func Hunterall(keyword string) (urls []string) {
	color.RGBStyleFromString("244,211,49").Println("请耐心等待hunter搜索......")
	hunter := GetHunterConfig()
	for i := 1; i <= 10; i++ {
		//获取请求的api
		url := hunter_api(keyword, hunter.Hunter_key, i)
		//发起请求
		res := hunterHttp(url, hunter.Fofa_timeout)
		results := res.Data.Arr
		num := len(results)

		for _, result := range results {
			//fmt.Printf("%v\n", result.Url)
			urls = append(urls, result.Url)
		}
		if num < 100 {
			break
		}
	}

	return urls
}
