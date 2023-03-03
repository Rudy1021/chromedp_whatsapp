package main

import (
	chrome "chromedp_whatsapp/CHROME"
	"fmt"
	"time"
)

func main() {
	ctx := chrome.InitChromedp()
	fmt.Println("Chrome is running now...")
	chrome.SendMessage(ctx, "Ruei", "headless模式成功")
	fmt.Println("finish")
	time.Sleep(10 * time.Second)
}
