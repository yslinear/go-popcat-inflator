package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func getIP() string {
	url := "https://api.ipify.org?format=text"
	fmt.Printf("Getting IP address from ipify\n")
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("My IP is:%s\n", ip)
	return string(ip)
}

// get geolocation with ip address
func getGeolocation(ip_address string) {
	fmt.Print("Getting geolocation\n")
	url := "https://ipwhois.app/json/" + ip_address
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Geolocation:%s\n", body)
}

func openNewPage() {

	dir, err := ioutil.TempDir("", "chromedp-example")
	if err != nil {
		panic(err)
	}

	defer os.RemoveAll(dir)
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("headless", false),
		chromedp.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3830.0 Safari/537.36"),
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.Flag("window-size", "50,400"),
		chromedp.UserDataDir(dir),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	// also set up a custom logger
	taskCtx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()
	// create a timeout
	taskCtx, cancel = context.WithTimeout(taskCtx, 300*time.Second)
	defer cancel()
	// ensure that the browser process is started
	if err := chromedp.Run(taskCtx); err != nil {
		panic(err)
	}
	// listen network event
	listenForNetworkEvent(taskCtx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)

	var res interface{}

	chromedp.Run(taskCtx,
		network.Enable(),
		chromedp.Navigate(`https://popcat.click/`),
		chromedp.Evaluate(`
		var event = new KeyboardEvent('keydown', {
			key: 'g',
			ctrlKey: true
		});
		setInterval(function () {
			document.dispatchEvent(event);
			document.getElementById('app').__vue__.accumulator = 800
		}, 1000);`,
			&res),
		chromedp.WaitNotVisible("body", chromedp.BySearch),
	)
	fmt.Println("Press Ctrl+C to exit")
}

//监听
func listenForNetworkEvent(ctx context.Context) {
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *network.EventResponseReceived:
			resp := ev.Response
			if len(resp.Headers) != 0 {
				if strings.Contains(resp.URL, "pop_count") {
					log.Printf("received status code: %d", resp.Status)
				}
			}
		}
		// other needed network Event
	})
}

func main() {
	my_ip := getIP()
	getGeolocation(my_ip)
	openNewPage()
}
