package main

import (
	"log"
	"time"

	"github.com/mxschmitt/playwright-go"
)

func main() {
	URL := "https://www.bloomberg.co.jp/markets/stocks"
	// URL := "https://finviz.com/map.ashx"

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
		// WaitUntil: playwright.WaitUntilStateNetworkidle,
	}); err != nil {
		log.Fatalf("could not goto: %v", err)
	}

	time.Sleep(3 * time.Second)
	if _, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path:     playwright.String("foo.png"),
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
