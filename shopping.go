package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/PuerkitoBio/goquery"
	"github.com/avast/retry-go/v4"
	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
	"github.com/google/go-querystring/query"
	log "github.com/sirupsen/logrus"
)

const regular = "window\\.jsl\\.dh\\('%s',(.*?)\\);"

func doShopping(ctx context.Context, params *RequestParam) (*Response, error) {
	jar := tls_client.NewCookieJar()
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(profiles.Chrome_124),
		tls_client.WithCookieJar(jar),
	}
	client, _ := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	client.SetProxy(ProxyStr)

	resp, err := retry.DoWithData(func() (*http.Response, error) {
		urlQuery, _ := query.Values(params)
		urlQuery.Del("google_domain")
		urlQuery.Del("engine")
		urlQuery.Set("udm", "28")
		urlQuery.Set("sclient", "sclient=gws-wiz-modeless-shopping")
		urlStr := fmt.Sprintf("https://www.%s/search?", params.GoogleDomain) + urlQuery.Encode()
		req, reqError := http.NewRequest("GET", urlStr, nil)
		if reqError != nil {
			log.Errorf("request err:%v", reqError)
			return nil, errors.New("request err")
		}
		req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		req.Header.Set("accept-language", "en")
		req.Header.Set("cache-control", "no-cache")
		req.Header.Set("downlink", "10")
		req.Header.Set("pragma", "no-cache")
		req.Header.Set("priority", "u=0, i")
		req.Header.Set("referer", "https://www.google.com/")
		req.Header.Set("rtt", "200")
		req.Header.Set("sec-ch-prefers-color-scheme", "light")
		req.Header.Set("sec-ch-ua", `"Not(A:Brand";v="99", "Google Chrome";v="133", "Chromium";v="133"`)
		req.Header.Set("sec-ch-ua-arch", `"x86"`)
		req.Header.Set("sec-ch-ua-bitness", `"64"`)
		req.Header.Set("sec-ch-ua-form-factors", `"Desktop"`)
		req.Header.Set("sec-ch-ua-full-version", `"133.0.6943.98"`)
		req.Header.Set("sec-ch-ua-full-version-list", `"Not(A:Brand";v="99.0.0.0", "Google Chrome";v="133.0.6943.98", "Chromium";v="133.0.6943.98"`)
		req.Header.Set("sec-ch-ua-mobile", "?0")
		req.Header.Set("sec-ch-ua-model", `""`)
		req.Header.Set("sec-ch-ua-platform", `"Windows"`)
		req.Header.Set("sec-ch-ua-platform-version", `"10.0.0"`)
		req.Header.Set("sec-ch-ua-wow64", "?0")
		req.Header.Set("sec-fetch-dest", "document")
		req.Header.Set("sec-fetch-mode", "navigate")
		req.Header.Set("sec-fetch-site", "same-origin")
		req.Header.Set("sec-fetch-user", "?1")
		req.Header.Set("upgrade-insecure-requests", "1")
		req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36")
		req.Header.Set("x-browser-channel", "stable")
		req.Header.Set("x-browser-copyright", "Copyright 2025 Google LLC. All rights reserved.")
		req.Header.Set("x-browser-validation", "1nAW9Rb/M8Lkk97ILDg00FWYjns=")
		req.Header.Set("x-browser-year", "2025")
		resp, respError := client.Do(req)
		if respError != nil {
			return nil, errors.New("request err")
		}
		if resp.StatusCode != http.StatusOK {
			return nil, errors.New("request err")
		}
		return resp, nil
	},
		retry.Context(ctx),
		retry.Attempts(10),
		retry.Delay(100*time.Millisecond),
		retry.OnRetry(func(attempt uint, err error) {
			log.Warnf("fetch URI=%s count=%v error=%v", "google_shopping", attempt, err)
		}),
		retry.LastErrorOnly(true),
	)
	if err != nil {
		log.Errorf("request err:%v", err)
		return nil, errors.New("request err")
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("request err:%v", err)
		return nil, errors.New("request err")
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(bodyText))
	if err != nil {
		log.Errorf("resolve html error = %s\n", err)
		return nil, errors.New("request err")
	}

	var response Response

	if params.IsReturnRawHtml == "1" {
		searchMetaData := SearchMetadata{RawHtml: string(bodyText)}
		if !searchMetaData.IsEmpty() {
			response.SearchMetadata = &searchMetaData
		}
	}

	// 侧边栏
	var filterList []Filters
	doc.Find("div[id='appbar'] div[jsname='HJCfLb'] div[jsname='pYVSud'] ul").ChildrenFiltered("li").Each(func(i int, s *goquery.Selection) {
		filter := Filters{
			Type: s.Find("div[jsname='ARU61'] span[role='heading']").Text(), // eg: Refine results
		}
		var optionList []Options
		s.Find("ul[jsname='CbM3zb'] li").Each(func(i int, s1 *goquery.Selection) {
			option := Options{
				Text: s1.Find("a").Text(), // 侧边栏具体的子类 eg: In store
				Link: s1.Find("a").AttrOr("href", ""),
			}
			if option.Link != "" { // 可能为空？
				option.Link = fmt.Sprintf("https://www.google.com/%s", option.Link)
			}
			if !option.IsEmpty() {
				optionList = append(optionList, option)
			}
		})
		if len(optionList) > 0 {
			filter.Options = optionList
		}
		if !filter.IsEmpty() {
			filterList = append(filterList, filter)
		}
	})
	if len(filterList) > 0 {
		response.Filters = filterList
	}

	// 具体的商品解析方法定义 对应页面 More places下面的
	getShoppingResult := func(count int, s3 *goquery.Selection) ShoppingResult {
		s4 := s3.Find("div[jsname='luUKCc'] div[class='MUWJ8c']").Children()
		title := s4.Children().Eq(1).Text() // eg: Folgers Coffee Ground Classic Roast
		if title == "" {
			return ShoppingResult{}
		}
		shoppingResult := ShoppingResult{
			Position:   count,
			Title:      title,
			ProductID:  s3.Find("div[jsname='dQK82e']").AttrOr("data-cid", ""),
			Source:     s4.Children().Find("span[class='WJMUdc rw5ecc']").Text(),
			SourceIcon: s4.Children().Eq(3).Find("img").AttrOr("src", ""),
			Price:      s4.Children().Eq(2).Find("span").First().Text(),
			OldPrice:   s4.Children().Eq(2).Find("span").Eq(1).Text(),
			Thumbnail:  s4.Children().Eq(0).Find("img").AttrOr("src", ""),
			Tag:        s4.Children().Eq(0).Text(),
			Delivery:   s4.Children().Find("span[class='ybnj7e']").Text(),
		}
		if shoppingResult.ProductID != "" {
			shoppingResult.ProductLink = fmt.Sprintf("https://www.google.com/shopping/product/%s?gl=%s", shoppingResult.ProductID, params.Gl)
		}
		if shoppingResult.Price != "" {
			shoppingResult.ExtractedPrice, _ = strconv.ParseFloat(strings.ReplaceAll(strings.ReplaceAll(shoppingResult.Price, "$", ""), ",", ""), 64)
		}
		if shoppingResult.OldPrice != "" {
			shoppingResult.ExtractedOldPrice, _ = strconv.ParseFloat(strings.ReplaceAll(strings.ReplaceAll(shoppingResult.OldPrice, "$", ""), ",", ""), 64)
			if shoppingResult.ExtractedOldPrice == 0.0 {
				shoppingResult.SecondHandCondition = shoppingResult.OldPrice
				shoppingResult.OldPrice = ""
			}
		}
		shoppingResult.Rating, _ = strconv.ParseFloat(s4.Children().Find("div[class='LFROUd']").Find("span").First().Children().Eq(0).Text(), 64)
		reviewsText := s4.Children().Find("div[class='LFROUd']").Find("span").First().Children().Eq(2).Text()
		reviewsStr := strings.ReplaceAll(strings.ReplaceAll(reviewsText, "(", ""), ")", "")
		if strings.Contains(reviewsStr, "K") {
			reviews, _ := strconv.ParseFloat(strings.ReplaceAll(reviewsStr, "K", ""), 64)
			shoppingResult.Reviews = int(reviews * 1000)
		} else {
			shoppingResult.Reviews, _ = strconv.Atoi(strings.ReplaceAll(reviewsStr, "K", ""))
		}
		extensions := append(shoppingResult.Extensions, s4.Children().Eq(0).Text())
		if len(extensions) > 0 {
			shoppingResult.Extensions = extensions
		}

		id := s3.Find("div[jsname='uVFeEd']").AttrOr("id", "")
		html := ButtonsData(id, string(bodyText))
		doc2, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
		shoppingResult.Snippet = doc2.Find("div[jsname='luUKCc']").Eq(1).Text()
		matches := regexp.MustCompile(`google\.ldi\s*=\s*\{([\s\S]*?)};`).FindStringSubmatch(string(bodyText))
		var matchesStr string
		if len(matches) > 1 {
			matchesStr = "{" + matches[1] + "}"
		}
		var imageUrlMap map[string]string
		_ = json.Unmarshal([]byte(matchesStr), &imageUrlMap)
		if imageUrlMap != nil {
			var thumbnails []string
			doc2.Find("div[jsname='DzNtMd']").Each(func(i int, s5 *goquery.Selection) {
				id2 := s5.Find("img").AttrOr("id", "")
				thumbnails = append(thumbnails, imageUrlMap[id2])
			})
			if len(thumbnails) > 0 {
				shoppingResult.Thumbnails = thumbnails
			}
		}
		return shoppingResult
	}

	// More places下面的
	var count int
	var shoppingResultList, shoppingResultList2 []ShoppingResult
	var categorizedShoppingResultList []CategorizedShoppingResults
	doc.Find("div[jscontroller='wuEeed']").Each(func(i int, s2 *goquery.Selection) {

		s2.Find("g-card[jscontroller='XT8Clf'] ul li").Each(func(i int, s3 *goquery.Selection) {
			count++
			shoppingResult := getShoppingResult(count, s3)
			if !shoppingResult.IsEmpty() {
				shoppingResultList = append(shoppingResultList, shoppingResult)
			}
		})

		s2.Find("div[jsname='vyMcq']").Each(func(i int, s6 *goquery.Selection) {
			// todo 这个在页面上已经没有了
			categorizedShoppingResult := CategorizedShoppingResults{
				Title: s6.Find("div[jscontroller='YHhMSc'] div[jsname='tJHJj']").Find("div[role='heading']").First().Text(),
			}
			// 页面顶部的分类详情
			var count3 int
			s6.Find("div[jsname='s2gQvd'] div[role='listitem']").Each(func(i int, s7 *goquery.Selection) {
				count3++
				shoppingResult := getShoppingResult(count3, s7)
				if !shoppingResult.IsEmpty() {
					shoppingResultList2 = append(shoppingResultList2, shoppingResult)
				}
			})
			if len(shoppingResultList2) > 0 {
				categorizedShoppingResult.ShoppingResult = shoppingResultList2
			}
			if !categorizedShoppingResult.isEmpty() {
				categorizedShoppingResultList = append(categorizedShoppingResultList, categorizedShoppingResult)
			}
		})

	})
	if len(shoppingResultList) > 0 {
		response.ShoppingResults = shoppingResultList
	}
	if len(categorizedShoppingResultList) > 0 {
		response.CategorizedShoppingResults = categorizedShoppingResultList
	}

	var count2 int
	var inLineResultList []InlineShoppingResult
	doc.Find("div[id='rso'] g-scrolling-carousel[jscontroller='pgCXqb']").Find("div[jsname='s2gQvd'] div[jsname='U8yK8']").Each(func(i int, s5 *goquery.Selection) {
		title := s5.Find("div[class='orXoSd'] div[role='heading']").Find("a").Text()
		if title == "" {
			return
		}
		count2++
		inLineResult := InlineShoppingResult{
			Position:      count2,
			BlockPosition: "top",
			Title:         title,
			Price:         s5.Find("div[class='orXoSd']").Find("div[class='T4OwTb']").Text(),
			Link:          s5.Find("a").First().AttrOr("href", ""),
			Source:        "",
			Thumbnail:     s5.Find("img").First().AttrOr("src", ""),
		}
		if inLineResult.Price != "" {
			inLineResult.ExtractedPrice, _ = strconv.ParseFloat(extractNumbersUsingMap(inLineResult.Price), 64)
		}
		if inLineResult.Link != "" {
			inLineResult.Link = fmt.Sprintf("https://www.google.com/%s", inLineResult.Link)
		}
		var extensions []string
		extensions = append(extensions, s5.Find("div[class='orXoSd'] div[class='LbUacb']").Text())
		extensions = append(extensions, s5.Find("div[class='orXoSd']").Children().Eq(1).Text())
		if len(inLineResult.Extensions) > 0 {
			inLineResult.Extensions = extensions
		}
		if !inLineResult.IsEmpty() {
			inLineResultList = append(inLineResultList, inLineResult)
		}
	})
	if len(inLineResultList) > 0 {
		response.InlineShoppingResults = inLineResultList
	}

	return &response, nil
}

func ButtonsData(id string, html string) string {
	ariaControlsZZ := fmt.Sprintf(regular, id)
	regex, _ := regexp.Compile(ariaControlsZZ)
	matches := regex.FindAllStringSubmatch(html, -1)
	var jsonString string
	if len(matches) > 0 {
		for _, value := range matches {
			jsonString = value[1]
		}
	}
	newJson := strings.ReplaceAll(jsonString, "'", "")
	decodedString := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(string(newJson), `\x3c`, "<"), `\x3d`, "="), `\x22`, `"`), `\x3e`, ">")
	return decodedString
}

// 数字跟英文字符  只截取数字部分
func extractNumbersUsingMap(s string) string {
	filter := func(r rune) rune {
		if unicode.IsDigit(r) || r == '.' {
			return r
		}
		return -1
	}
	return strings.Map(filter, s)
}
