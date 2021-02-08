package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/sclevine/agouti"
)

func GetRestaurantsURL() ([]Restaurant, error) {
	driver := agouti.ChromeDriver(
		agouti.ChromeOptions("args", []string{
			"--headless",
			"--disable-gpu",
			//"--start-fullscreen",
		}),
	)
	defer driver.Stop()

	if err := driver.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return nil, err
	}

	page, err := driver.NewPage()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return nil, err
	}

	url := "https://www.ubereats.com/jp/feed?ps=1&pl=JTdCJTIyYWRkcmVzcyUyMiUzQSUyMiVFNiU5RCVCMSVFNCVCQSVBQzIzJUU1JThDJUJBJTIyJTJDJTIycmVmZXJlbmNlJTIyJTNBJTIyQ2hJSmRRdkplb1dMR0dBUk5ZNWktNkFEb0VZJTIyJTJDJTIycmVmZXJlbmNlVHlwZSUyMiUzQSUyMmdvb2dsZV9wbGFjZXMlMjIlMkMlMjJsYXRpdHVkZSUyMiUzQTM1LjcwOTAyNTklMkMlMjJsb25naXR1ZGUlMjIlM0ExMzkuNzMxOTkyNSU3RA%3D%3D"
	page.Navigate(url)

	if err := page.Session().SetImplicitWait(10); err != nil {
		log.Println(err)
	}

	// 「さらに表示」ボタンをなくなるまでクリックする
	var pageCount int
	for {
		pageCount++
		log.Printf("現在%dページ目\n", pageCount)

		if err := page.All("#main-content > div > div:nth-child(3) > div:nth-child(2) > div > button").Click(); err != nil {
			log.Println("End: ", err)
			break
		}
		if err := page.Session().SetImplicitWait(5); err != nil {
			log.Println(err)
		}
		time.Sleep(time.Second * 5)
	}

	html, err := page.HTML()
	if err != nil {
		log.Println(err)
	}
	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)

	// log
	listDiv := doc.Find("#main-content > div > div:nth-child(3) > div:nth-child(2) > div > div:nth-child(2)").Children()
	listLen := listDiv.Length()
	log.Println("listLen:", listLen)

	allRestaurants := doc.Find("#main-content > div > div:nth-child(3) > div:nth-child(2) > div > div:nth-child(2)").Children()
	restaurants := []Restaurant{}
	allRestaurants.Each(func(_ int, s *goquery.Selection) {
		restaurant := Restaurant{}
		if restaurantURL, exists := s.Find("div > a").Attr("href"); exists {
			restaurant.URL = "https://www.ubereats.com" + restaurantURL
			restaurants = append(restaurants, restaurant)
			//log.Println("restoURL", restaurant.URL)
		} else {
			log.Println("URL does not exist.")
		}
	})

	//fmt.Printf("%+v\n", restaurants)

	//restoList := doc.Find("#main-content > div > div:nth-child(3) > div:nth-child(2) > div > div:nth-child(2) > div:nth-child(2) > div")
	//restoList.Children().Each(func(_ int, s *goquery.Selection) {})
	//resto := restoList.Find("a")
	//restoURL, exists := resto.Attr("href")
	/*if exists == false {
		log.Println("見つからないよう...")
	}
	log.Println("url:", restoURL)*/

	//page.Screenshot("uber5.png")
	fmt.Println("finish")

	return restaurants, nil
}

//1180, 8760

//TODO: agoutiのhtml取得して、goquery で取得する (やった)
//あとスクロールしてボタンを押す自動化できる方法を探す (これはやらなくていいかも)
//goquety でpicture figure 取れるか確認する (取れた https://www.ubereats.com/jp/ これつければおｋ) できた
//大体レストラン1839位
