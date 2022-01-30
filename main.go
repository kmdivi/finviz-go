package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mxschmitt/playwright-go"
	"gopkg.in/ini.v1"
)

type slack struct {
	TOKEN   string
	CHANNEL string
}

func main() {
	// URL := "https://www.bloomberg.com/markets/stocks"
	// URL := "https://finviz.com/map.ashx"
	URL := ""
	if len(os.Args) < 1 {
		URL = os.Args[1]
	} else {
		fmt.Println("usage:\n\t finviz [URL]")
		os.Exit(0)
	}

	cfg, err := ini.Load("ini.config")
	if err != nil {
		log.Fatal(err)
	}
	TOKEN := cfg.Section("slack").Key("token").String()
	CHANNEL := cfg.Section("slack").Key("channel").String()

	fn := take_screenshot(URL)
	post_slack(TOKEN, CHANNEL, fn)
}

func take_screenshot(URL string) string {
	fn := get_filename(URL)

	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not launch playwright: %v", err)
	}
	browser, err := pw.WebKit.Launch()
	if err != nil {
		log.Fatalf("could not launch Chromium: %v", err)
	}
	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	if _, err = page.Goto(URL, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateLoad,
	}); err != nil {
		log.Fatalf("could not goto: %v", err)
	}

	time.Sleep(3 * time.Second)
	if _, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path:     playwright.String(fn),
		FullPage: playwright.Bool(true),
	}); err != nil {
		log.Fatalf("could not create screenshot: %v", err)
	}
	if err = browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}
	if err = pw.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}

	return fn
}

func get_filename(s string) string {
	fn := ""
	if strings.Contains(s, "http://") {
		fn = s[strings.Index(s, "http://")+7:]
	} else if strings.Contains(s, "https://") {
		fn = s[strings.Index(s, "https://")+8:]
	}
	fn = fn[:strings.Index(fn, "/")]
	fn += ".png"

	return fn
}

func post_slack(TOKEN string, CHANNEL string, fn string) {
	f, err := os.Open(fn)
	if err != nil {
		panic("error")
	}
	defer f.Close()

	bodyBuf := &bytes.Buffer{}
	writer := multipart.NewWriter(bodyBuf)
	part, err := writer.CreateFormFile("file", fn)
	if err != nil {
		panic("error")
	}
	if _, err := io.Copy(part, f); err != nil {
		panic("error")
	}

	err = writer.WriteField("token", TOKEN)
	if err != nil {
		panic("error")
	}

	err = writer.WriteField("channels", CHANNEL)
	if err != nil {
		panic("error")
	}

	err = writer.Close()
	if err != nil {
		panic("error")
	}

	requestSlack, err := http.NewRequest(
		"POST",
		"https://slack.com/api/files.upload",
		bodyBuf)
	if err != nil {
		panic("error")
	}

	requestSlack.Header.Set("Content-Type", writer.FormDataContentType())

	clientSlack := new(http.Client)
	responseSlack, err := clientSlack.Do(requestSlack)
	if err != nil {
		panic("error")
	}
	defer responseSlack.Body.Close()
}
