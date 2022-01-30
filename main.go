package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/mxschmitt/playwright-go"
)

func main() {
	// URL := "https://www.bloomberg.com/markets/stocks"
	// URL := "https://finviz.com/map.ashx"
	URL := os.Args[1]
	take_screenshot(URL)
}

func take_screenshot(URL string) {
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
