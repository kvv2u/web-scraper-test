package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
)

type Restaurant struct {
	Name               string
	Location           string
	DeliveryTime       string
	URL                string // TODO: 店舗一覧から取得
	DeliveryFee        int64
	Category           string
	BusinessDay        []string
	BusinessHours      []string
	Menus              []Menu
	Rating             float64
	MinimumOrderAmount int64
}

type Menu struct {
	Restasrant       string // TODO: RestaurantName とかにしたほうが良さそう
	Name             string
	Price            int64
	Description      string
	ImageURL         string
	DeliveryProvider DeliveryProvider
}

type DeliveryProvider string

const (
	UberEats DeliveryProvider = "UberEats"
	Demaecan DeliveryProvider = "出前館"
	Wolt     DeliveryProvider = "Wolt"
)

func main() {
	fmt.Println("log - main")
	//CrawlUberEats()
	//goquerycrawl()
	//scrapeUberEatsMenu()

	// UberEats から全店舗URLを取得
	urls, err := GetRestaurantsURL()
	if err != nil {
		log.Println(err)
	}
	for _, v := range urls {
		log.Println(v.URL)
	}
	//fmt.Printf("%+v\n", urls)
}

func CrawlUberEats() {
	log.Println("crawluber")
	// 東京23区
	url := "https://www.ubereats.com/jp/feed?ps=1&pl=JTdCJTIyYWRkcmVzcyUyMiUzQSUyMiVFNiU5RCVCMSVFNCVCQSVBQzIzJUU1JThDJUJBJTIyJTJDJTIycmVmZXJlbmNlJTIyJTNBJTIyQ2hJSmRRdkplb1dMR0dBUk5ZNWktNkFEb0VZJTIyJTJDJTIycmVmZXJlbmNlVHlwZSUyMiUzQSUyMmdvb2dsZV9wbGFjZXMlMjIlMkMlMjJsYXRpdHVkZSUyMiUzQTM1LjcwOTAyNTklMkMlMjJsb25naXR1ZGUlMjIlM0ExMzkuNzMxOTkyNSU3RA%3D%3D"

	restauCollector := colly.NewCollector(
		colly.AllowedDomains("ubereats.com", "www.ubereats.com"),
	)

	//restauCollector := restoCol.Clone()

	/*restauCollector.Limit(&colly.LimitRule{
		Delay: 3 * time.Second,
	})*/

	restauCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL.String())
	})

	// さらに表示
	restauCollector.OnHTML("button.gn", func(e *colly.HTMLElement) {
		nextPage := e.Request.AbsoluteURL(e.Attr("button"))
		log.Println(nextPage)
		restauCollector.Visit(nextPage)
	})

	/*restauCollector.OnHTML("div.fu jg fw fx fy fz", func(e *colly.HTMLElement) {
		log.Println("success")

		e.ForEach("div.g0 jh ji", func(_ int, restau *colly.HTMLElement) {
			restaurantURL := restau.ChildAttr("div.af ie > a", "href")
			restaurantPage := restau.Request.AbsoluteURL(restaurantURL)
			restauCollector.Visit(restaurantPage)
		})
	})*/

	// テスト
	restauCollector.OnHTML("#main-content > div > div.ba > div.ba > div", func(e *colly.HTMLElement) {
		a := e.Attr("div.fu.ib.fw.fx.fy.fz > div.g0")
		log.Println(a)
		log.Println("ok")
		//e.ForEach("div.g0")
	})

	restauCollector.OnHTML("#main-content", func(e *colly.HTMLElement) {
		log.Println("oooooookkk")
		text := e.ChildText("div > div.ba.ag.du.dv.dw.dx > div.ba.c9.bn.ev.dz.e0.e1.e2 > div > div.fu.ji.fw.fx.fy.fz > div:nth-child(1) > div > figure > a > h3")
		log.Println(text)
	})

	restauCollector.OnHTML("#main-content > div > div.ba.ag.dd.de.df.dg > div.ba.bq.bn.ef.di.dj.dk.dl > div > div:nth-child(2)", func(e *colly.HTMLElement) {
		log.Println("yayyayayyyayaya")
	})

	if err := restauCollector.Visit(url); err != nil {
		log.Println("err: ", err)
	}
}

func scrapeUberEatsMenu() {
	//url := "https://www.ubereats.com/jp/tokyo/food-delivery/%E9%BA%BB%E8%BE%A3%E8%AA%98%E6%83%91%E5%8C%97%E5%8F%A3%E5%BA%97-mara-yuwaku-kitaguchiten/6jxjqvCJQxie48S8aU4zIw?pl=JTdCJTIyYWRkcmVzcyUyMiUzQSUyMiVFNiU5RCVCMSVFNCVCQSVBQzIzJUU1JThDJUJBJTIyJTJDJTIycmVmZXJlbmNlJTIyJTNBJTIyQ2hJSmRRdkplb1dMR0dBUk5ZNWktNkFEb0VZJTIyJTJDJTIycmVmZXJlbmNlVHlwZSUyMiUzQSUyMmdvb2dsZV9wbGFjZXMlMjIlMkMlMjJsYXRpdHVkZSUyMiUzQTM1LjcwOTAyNTklMkMlMjJsb25naXR1ZGUlMjIlM0ExMzkuNzMxOTkyNSU3RA%3D%3D"

	//url := "https://www.ubereats.com/tokyo/food-delivery/%E3%83%88%E3%83%9F%E3%83%8E%E3%83%BB%E3%83%92%E3%82%B5-%E6%B1%9F%E6%88%B8%E5%B7%9D%E6%A9%8B%E5%BA%97-dominos-pizza-edogawabashi-store/gTW7-VUNQ62MyZN4zpKpfA?pl=JTdCJTIyYWRkcmVzcyUyMiUzQSUyMiVFNiU5RCVCMSVFNCVCQSVBQzIzJUU1JThDJUJBJTIyJTJDJTIycmVmZXJlbmNlJTIyJTNBJTIyQ2hJSmRRdkplb1dMR0dBUk5ZNWktNkFEb0VZJTIyJTJDJTIycmVmZXJlbmNlVHlwZSUyMiUzQSUyMmdvb2dsZV9wbGFjZXMlMjIlMkMlMjJsYXRpdHVkZSUyMiUzQTM1LjcwOTAyNTklMkMlMjJsb25naXR1ZGUlMjIlM0ExMzkuNzMxOTkyNSU3RA%3D%3D"

	// 東京23区
	url := "https://www.ubereats.com/jp/feed?ps=1&pl=JTdCJTIyYWRkcmVzcyUyMiUzQSUyMiVFNiU5RCVCMSVFNCVCQSVBQzIzJUU1JThDJUJBJTIyJTJDJTIycmVmZXJlbmNlJTIyJTNBJTIyQ2hJSmRRdkplb1dMR0dBUk5ZNWktNkFEb0VZJTIyJTJDJTIycmVmZXJlbmNlVHlwZSUyMiUzQSUyMmdvb2dsZV9wbGFjZXMlMjIlMkMlMjJsYXRpdHVkZSUyMiUzQTM1LjcwOTAyNTklMkMlMjJsb25naXR1ZGUlMjIlM0ExMzkuNzMxOTkyNSU3RA%3D%3D"

	restoCol := colly.NewCollector(
		colly.AllowedDomains("ubereats.com", "www.ubereats.com"),
		colly.Debugger(&debug.LogDebugger{}),
	)

	restoCol.Limit(&colly.LimitRule{
		RandomDelay: 5 * time.Second,
	})

	moreInfoCol := restoCol.Clone()

	restaurant := Restaurant{}

	// TODO:「さらに表示」をクリックする処理を実装

	// 店舗一覧
	restoCol.OnHTML("#main-content > div > div.ba.ag.dd.de.df.dg > div:nth-child(2) > div > div:nth-child(2)", func(e *colly.HTMLElement) {
		log.Println("okだよ￥ーーーー")
		imgurls := []string{}
		restauranturls := []string{}
		e.ForEach("div", func(_ int, el *colly.HTMLElement) {
			imgURL := el.ChildAttr("div > figure > div > picture > img", "src")
			imgURL = el.Request.AbsoluteURL(imgURL)
			restoURL := el.ChildAttr("div > a", "href")
			restoURL = el.Request.AbsoluteURL(restoURL)
			imgurls = append(imgurls, imgURL)
			restauranturls = append(restauranturls, restoURL)
		})
		//log.Println("imgurls", imgurls)
		//log.Println("resturanturls", restauranturls)
	})

	/*restoCol.OnHTML("#main-content > div > div.ba.ag.dd.de.df.dg > div:nth-child(2) > div > div:nth-child(2) > div", func(e *colly.HTMLElement) {
		log.Println("ここまではいてる")
		log.Println(e.ChildText("div"))
	})*/

	restoCol.OnHTML("#main-content", func(e *colly.HTMLElement) {
		restaurant.Name = e.ChildText("h1.dy")
		restaurant.DeliveryTime = e.ChildText("div:nth-child(3) > div > div > div.eu > div:nth-child(1) > div.ag > div:nth-child(2) > div.cc > div:nth-child(7)")
		delivFeeStr := e.ChildText("div:nth-child(3) > div > div > div.eu > div:nth-child(1) > div.ag > div:nth-child(2) > div.cc > div:nth-child(4)")
		restaurant.DeliveryFee = numCheck(delivFeeStr)
		// TODO: 評価は店舗一覧からとっても良さそう
		restaurant.Rating, _ = strconv.ParseFloat(e.ChildText("div:nth-child(3) > div > div > div.eu > div:nth-child(1) > div.ag > div:nth-child(2) > div.cc > div:nth-child(9)"), 64)

		categoryStr := e.ChildText("div:nth-child(4) > div:nth-child(1) > div:nth-child(1) > div:nth-child(1)")
		// TODO: "&nbsp・&nbsp" の部分も後で取り除く
		if strings.Contains(categoryStr, "¥") {
			restaurant.Category = strings.TrimLeft(categoryStr, "¥")
		} else if strings.Contains(categoryStr, "¥¥") {
			restaurant.Category = strings.TrimLeft(categoryStr, "¥¥")
		} else if strings.Contains(categoryStr, "¥¥¥") {
			restaurant.Category = strings.TrimLeft(categoryStr, "¥¥¥")
		} else {
			restaurant.Category = strings.TrimLeft(categoryStr, "¥¥¥¥")
		}

		e.ForEach("div.b8.b9.ba.bb.bc > ul > li", func(_ int, el *colly.HTMLElement) {
			el.ForEach("ul > li", func(_ int, m *colly.HTMLElement) {
				tmpMenu := Menu{}
				tmpMenu.Restasrant = restaurant.Name
				tmpMenu.Name = m.ChildText("div > div > div > div:nth-child(1) > div:nth-child(1) > h4 > div")
				// TODO: 説明が無いメニューだとdivが1段無くてずれる
				tmpMenu.Description = m.ChildText("div > div > div > div:nth-child(1) > div:nth-child(2) > div")
				priceStr := strings.Trim(m.ChildText("div > div > div > div:nth-child(1) > div:nth-child(3) > div"), "¥")
				if strings.Contains(priceStr, ",") {
					priceStr = strings.Replace(priceStr, ",", "", -1)
				}
				tmpMenu.Price, _ = strconv.ParseInt(priceStr, 10, 64)
				// 何故か picture タグに到達できないから画像取れない
				tmpMenu.ImageURL = m.ChildAttr("div > div > div > div.hw > picture > img", "src")
				tmpMenu.DeliveryProvider = UberEats
				restaurant.Menus = append(restaurant.Menus, tmpMenu)
				log.Println("name", tmpMenu.Name)
				log.Println("price", tmpMenu.Price)
				log.Println("des", tmpMenu.Description)
				log.Println("image", tmpMenu.ImageURL)
			})
		})

		log.Println("Name: ", restaurant.Name)
		log.Println("DeliveryFee", restaurant.DeliveryFee)
		log.Println("DeliveryTime", restaurant.DeliveryTime)
		log.Println("Rating", restaurant.Rating)
		log.Println("Category", restaurant.Category)
		//log.Println("Menus", restaurant.Menus)
	})

	// Visit moreInfo
	restoCol.OnHTML("#main-content > div.b8 > div:nth-child(1) > div:nth-child(1) > p", func(e *colly.HTMLElement) {
		moreInfo := e.ChildAttr("a", "href")
		moreInfo = e.Request.AbsoluteURL(moreInfo)
		//log.Println("urlこっちはちゃんと取れてる", moreInfo)
		moreInfoCol.Visit(moreInfo)
	})

	// Location and hours
	moreInfoCol.OnHTML("#wrapper", func(e *colly.HTMLElement) {
		restaurant.Location = e.ChildText("div:nth-child(7) > div > div > div.ba > div:nth-child(4) > div:nth-child(3) > div.cc")

		e.ForEach("div:nth-child(7) > div > div > div.ba > div:nth-child(4) > div:nth-child(4) > table tbody > tr", func(_ int, el *colly.HTMLElement) {
			day := el.ChildText("td.crestoCol.df")
			hours := el.ChildText("td.crestoCol.c6")
			if day != "" {
				restaurant.BusinessDay = append(restaurant.BusinessDay, day)
			} else if hours != "" {
				restaurant.BusinessHours = append(restaurant.BusinessHours, hours)
			}
		})

		log.Println("Location", restaurant.Location)
		log.Println("day", restaurant.BusinessDay)
		log.Println("hours", restaurant.BusinessHours)
	})

	// メニュー
	/*restoCol.OnHTML("#main-content > div.b8.b9.ba.bb.bc > ul > li:nth-child(4) > ul", func(e *colly.HTMLElement) {
		log.Println("aaaaaaaaaaddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd")
	})*/

	if err := restoCol.Visit(url); err != nil {
		log.Println("err: ", err)
	}

	restoCol.Wait()
	moreInfoCol.Wait()
	//log.Printf("%#v/n", restaurant)
}

func numCheck(s string) int64 {
	n := 0
	for _, r := range s {
		if '0' <= r && r <= '9' {
			n = n*10 + int(r-'0')
		} else {
			continue
		}
	}
	return int64(n)
}
