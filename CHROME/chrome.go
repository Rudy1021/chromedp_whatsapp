package chrome

/*
目前問題：印出網頁元素的話就無法waitVisible

可能的解法：先創建一個瀏覽器 如果有ＱＲＣＯＤＥ就關掉重開 然後實施登入 沒有的話就關掉重開 然後輸入訊息
*/

/*
第二種解法：
*/
import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/chromedp"
	"github.com/mdp/qrterminal"
)

func InitChromedp() context.Context {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-sync", false),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("blink-settings", "imagesEnabled=true"),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.5481.77 Safari/537.36"),
		// chromedp.UserDataDir("/Users/rudysmacbook/Library/Application Support/Google/Chrome/Default"),
		// chromedp.UserDataDir("/Users/rueiyang/projects/whatsapp_push/data"),
		chromedp.UserDataDir("./data"),
		// chromedp.UserDataDir("./dataInDocker/"),
		// chromedp.UserDataDir("/app/data/Default"),
	)

	alloCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)

	ctx, _ := chromedp.NewContext(
		alloCtx,
		chromedp.WithLogf(log.Printf),
	)

	chromedp.Run(ctx,
		chromedp.Navigate("https://web.whatsapp.com/"),
	)

	// --------------------------看現在正在qrcode還是chatroom畫面--------------------------
	var page string = func() string {
		var hasFound bool
		for {
			// 找是否在qrcode頁面
			chromedp.Run(ctx,
				chromedp.Evaluate(`document.getElementsByClassName('_19vUU').length > 0`, &hasFound),
			)
			if hasFound {
				return "qrcode"
			}

			// 找是否在chatroom頁面
			chromedp.Run(ctx,
				chromedp.Evaluate(`document.getElementsByClassName('ggj6brxn').length > 0`, &hasFound),
			)
			if hasFound {
				return "chatroom"
			}
			time.Sleep(1 * time.Second) // 等兩秒再找
		}
	}()

	// --------------------------確認畫面--------------------------
	if page == "qrcode" {
		// 正在登入畫面等待手機刷qrcode
		// waitVisible等qrcode出現，一定等得到
		fmt.Println("In qrcode page.")
		fmt.Println("waitng for qrcode cavas...")
		chromedp.Run(ctx,
			chromedp.WaitVisible(`#app > div > div > div.landing-window > div.landing-main > div > div > div._2I5ox > div > canvas`))
		fmt.Println("qrcode cavas has shown!!!")

		// 等到qrcode出現時
		// AttributeValue獲取qrcode真正的值，但有可能會爛掉
		// 但不管是否爛掉，只要把最新的qrcode值print出來就可以
		var prevQrcode string
		var ok, isInLoginPage bool
		for {
			// 每次for都先判斷是否在登入畫面
			// 如果登入成功，不再login畫面，AttributeValue找qrcode時，找不到會卡死
			chromedp.Run(ctx,
				chromedp.Evaluate(`document.getElementsByClassName('_19vUU').length > 0`, &isInLoginPage),
			)
			if isInLoginPage {
				var newQrcode string
				err := chromedp.Run(ctx,
					chromedp.AttributeValue(
						`#app > div > div > div.landing-window > div.landing-main > div > div > div._2I5ox > div`,
						"data-ref",
						&newQrcode, &ok))

				if prevQrcode != newQrcode {
					prevQrcode = newQrcode
					fmt.Println("get qrcode value: ", newQrcode)
					qrterminal.Generate(newQrcode, qrterminal.L, os.Stdout)
				}

				if err != nil {
					fmt.Println("qrcode error!!!", err.Error())
					log.Fatal(err)
				}
			} else {
				fmt.Println("Login successed, go to chatroom.")
				break
			}
			time.Sleep(2 * time.Second) // 讓for休息
		}
	} else { // page == "chatroom" 已在聊天室畫面
		fmt.Println("In chatroom page.")
	}
	return ctx
}

// if strings.Contains(res, "_21S-L") == true {
// 	//表示直接登入
// 	status = "logined"
// } else {
// 	//表示沒登入過
// 	status = "not login"
// }

/*
WhatsApp Web
*/
// func ScanQrCode(ctx context.Context) string {
// 	var test string
// 	var ok bool

// 	err := chromedp.Run(ctx,
// 		chromedp.WaitVisible(`#app > div > div > div.landing-window > div.landing-main > div > div > div._2I5ox > div`),
// 	)
// 	fmt.Println("@@@@@test1")
// 	if err != nil && err != context.DeadlineExceeded {
// 		log.Fatal(err)
// 	}

// 	if err == context.DeadlineExceeded {
// 		return "logined"
// 	}

// 	time.Sleep(10 * time.Second)

// 	err = chromedp.Run(ctx,
// 		chromedp.AttributeValue(`#app > div > div > div.landing-window > div.landing-main > div > div > div._2I5ox > div`, "data-ref", &test, &ok))
// 	if err != nil && err != context.DeadlineExceeded {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("@@@@@test2")
// 	qrterminal.Generate(test, qrterminal.L, os.Stdout)

// 	err = chromedp.Run(ctx,
// 		chromedp.WaitVisible(`#pane-side > div:nth-child(1) > div > div > div:nth-child(1) > div > div > div > div._8nE1Y > div.y_sn4 > div._21S-L > span`))

// 	if err == context.DeadlineExceeded {
// 		return "retry"
// 	}
// 	return "success"
// }

func SendMessage(ctx context.Context, title string, insert string) {
	group := fmt.Sprintf(`span.ggj6brxn[title='%s']`, title)

	err := chromedp.Run(ctx,
		chromedp.WaitVisible(`#pane-side > div:nth-child(1) > div > div > div:nth-child(1) > div > div > div > div._8nE1Y > div.y_sn4 > div._21S-L > span`),

		chromedp.Click(group),
		chromedp.Click(`p.selectable-text`),

		input.InsertText(insert),
		chromedp.Click(`button.tvf2evcx`),
	)

	if err != nil && err != context.DeadlineExceeded {
		log.Fatal(err)
	}
}

// func getWholeHtml(ctx context.Context, res string) {
// 	err := chromedp.Run(ctx,
// 		chromedp.ActionFunc(func(ctx context.Context) error {
// 			node, err := dom.GetDocument().Do(ctx)
// 			if err != nil {
// 				return err
// 			}
// 			res, err = dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
// 			return err
// 		}),
// 	)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// }

// var res string
// err = chromedp.Run(ctx,
// 	chromedp.ActionFunc(func(ctx context.Context) error {
// 		node, err := dom.GetDocument().Do(ctx)
// 		if err != nil {
// 			return err
// 		}
// 		res, err = dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
// 		return err
// 	}),
// )
// fmt.Println(res)
