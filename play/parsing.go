package play

import (
	"bytes"
	"encoding/json"
	jsoniter "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
	"regexp"
	"strings"
	"time"
)

// MakeConversionParsing 解析接口返回文本，进行分类组装
func MakeConversionParsing(respBytes []byte) (organics []OrganicResults, highlightItems [][]HighlightItem, chartOptions []ChartOption, pageToken string) {
	wrbs := GetWrbs(respBytes)
	var organicWrb []byte
	i := len(wrbs)
	if i > 0 {
		if i > 1 {
			organicWrb = wrbs[i-1]
		} else {
			organicWrb = wrbs[0]
		}
		var wrbFr [][]string
		err := json.Unmarshal(organicWrb, &wrbFr)
		if err != nil {
			log.Errorf("unmarshal organic wrb: %v", err)
		}
		playOrganic := wrbFr[0][2]
		// 获取到大类级别，eg:Popular Apps, Business tools
		clazz := jsoniter.Get([]byte(playOrganic), 0, 1).ToString()
		pageToken = jsoniter.Get([]byte(playOrganic), 0, 3, 1).ToString()
		var actual []interface{}
		err = json.Unmarshal([]byte(clazz), &actual)
		for actI, a := range actual {
			marshal, err := json.Marshal(a)
			if err != nil {
				log.Errorf("unmarshal actual[%d]: %v", actI, err)
				continue
			}
			cat := jsoniter.Get(marshal, 21, 1).ToString()
			if cat == "" {
				// 有的地方数组下标为22
				cat = jsoniter.Get(marshal, 22, 1).ToString()
			}
			// 以上两次判断如果有值则为organic_result，否则为items_highlight或者top_chart
			// 此处判断是否为highlight内容，如果为空则为highlight，一般highlight数组下标为34，top_chart数组下标为27
			if cat == "" {
				highlight := jsoniter.Get(marshal, 34).ToString()
				if highlight == "" {
					// 查看是否存在
					chartsJson := jsoniter.Get(marshal, 27, 1, 0).ToString()
					if chartsJson != "" {
						var charts []interface{}
						err := json.Unmarshal([]byte(chartsJson), &charts)
						if err != nil {
							log.Errorf("unmarshal charts option: %v", err)
							continue
						}
						// 如果序列化无数据不操作
						if len(charts) > 0 {
							for _, chart := range charts {
								chatJson, err := json.Marshal(chart)
								if err != nil {
									log.Errorf("unmarshal chatJson: %v", err)
									continue
								} else {
									text := jsoniter.Get(chatJson, 0).ToString()
									value := jsoniter.Get(chatJson, 1, 9, 0, 1).ToString()
									chartOption := ChartOption{
										Text:  text,
										Value: value,
									}
									chartOptions = append(chartOptions, chartOption)
								}
							}
						}
					}
				} else {
					// 解析每一个数组元素组装为highlight
					highs := jsoniter.Get([]byte(highlight), 0).ToString()
					var high []interface{}
					err := json.Unmarshal([]byte(highs), &high)
					if err != nil {
						log.Errorf("unmarshal organic wrb: %v", err)
						continue
					}
					var highItem []HighlightItem
					if len(high) > 0 {
						for _, v := range high {
							it, err := json.Marshal(v)
							if err != nil {
								log.Errorf("unmarshal highlight: %v", err)
								continue
							}
							productId := jsoniter.Get(it, 1, 0, 0, 0).ToString()
							title := jsoniter.Get(it, 1, 0, 3).ToString()
							var subTitle string
							var link string
							nail := jsoniter.Get(it, 1, 0, 1, 3, 2).ToString()
							if title == "" {
								title = jsoniter.Get(it, 0, 1, 1).ToString()
								link = jsoniter.Get(it, 0, 0, 4, 2).ToString()
								subTitle = jsoniter.Get(it, 0, 2, 1).ToString()
							} else {
								link = jsoniter.Get(it, 1, 0, 10, 3).ToString()
								subTitle = jsoniter.Get(it, 1, 2, 1).ToString()
							}
							highlightItem := HighlightItem{
								Title:       title,
								Subtitle:    subTitle,
								Link:        GooglePlayUrl + link,
								ProductID:   productId,
								SerpapiLink: "",
								Thumbnail:   nail,
							}
							highItem = append(highItem, highlightItem)
						}
					}
					highlightItems = append(highlightItems, highItem)
				}
			} else {
				// 组装为organic_result
				var organicResult OrganicResults
				orgaincs := jsoniter.Get(marshal, 21, 0).ToString()
				var organic []interface{}
				var items []Item
				// 有的数组长度为23，需要判断
				if orgaincs == "" {
					orgaincs = jsoniter.Get(marshal, 22, 0).ToString()
				}
				err := json.Unmarshal([]byte(orgaincs), &organic)
				if err != nil {
					log.Errorf("unmarshal organic wrb: %v", err)
					continue
				}
				title := jsoniter.Get([]byte(cat), 0).ToString()
				organicResult.Title = title
				items = assemOrganic(organic)
				if len(items) > 0 {
					organicResult.Item = items
				}
				//nextPageToken := jsoniter.Get([]byte(cat), 3, 1).ToString()
				//pagination := SerpapiSectionPagination{
				//	NextPageToken: nextPageToken,
				//SectionPageToken: sectionPageToken,
				//}
				//organicResult.SerpapiSectionPagination = &pagination
				// 循环组装organic中的item
				organics = append(organics, organicResult)
			}
		}
		return organics, highlightItems, chartOptions, pageToken
	} else {
		return organics, highlightItems, chartOptions, pageToken
	}

}

func GetWrbs(respBytes []byte) [][]byte {
	// 将 []byte 按行分割
	lines := bytes.Split(respBytes, []byte("\n"))
	var wrbs [][]byte
	// 遍历每一行
	for _, line := range lines {
		// 检查是否以 [["wrb.fr" 开头，以此为开头的为实际数据
		if bytes.HasPrefix(line, []byte(`[["wrb.fr"`)) {
			wrbs = append(wrbs, line)
		}
	}
	return wrbs
}

func assemOrganic(organic []interface{}) []Item {
	var items []Item
	if len(organic) > 0 {
		for _, org := range organic {
			orgMarsh, err := json.Marshal(org)
			if err != nil {
				log.Errorf("marshal organic: %v", err)
				continue
			}
			title := jsoniter.Get(orgMarsh, 3).ToString()
			var link, productId, author, category, download, video, thumbnail, desc string
			var rating float64
			if title == "" {
				// 有的再下一层才是数据，对应“有的数组长度为23，有的为22”，观察为23的在下一层
				author = jsoniter.Get(orgMarsh, 0, 14).ToString()
				link = jsoniter.Get(orgMarsh, 0, 10, 4, 2).ToString()
				productId = jsoniter.Get(orgMarsh, 0, 0, 0).ToString()
				rating = jsoniter.Get(orgMarsh, 0, 4, 0).ToFloat64()
				title = jsoniter.Get(orgMarsh, 0, 3).ToString()
				category = jsoniter.Get(orgMarsh, 0, 5).ToString()
				download = jsoniter.Get(orgMarsh, 0, 15).ToString()
				video = jsoniter.Get(orgMarsh, 0, 12, 0, 0, 3, 2).ToString()
				thumbnail = jsoniter.Get(orgMarsh, 0, 1, 3, 2).ToString()
				desc = jsoniter.Get(orgMarsh, 0, 13, 1).ToString()
			} else {
				link = jsoniter.Get(orgMarsh, 10, 4, 2).ToString()
				productId = jsoniter.Get(orgMarsh, 0, 0).ToString()
				rating = jsoniter.Get(orgMarsh, 4, 0).ToFloat64()
				author = jsoniter.Get(orgMarsh, 14).ToString()
				category = jsoniter.Get(orgMarsh, 5).ToString()
				download = jsoniter.Get(orgMarsh, 15).ToString()
				video = jsoniter.Get(orgMarsh, 12, 0, 0, 3, 2).ToString()
				thumbnail = jsoniter.Get(orgMarsh, 1, 3, 2).ToString()
				desc = jsoniter.Get(orgMarsh, 13, 1).ToString()
			}
			item := Item{
				Title:       title,
				Link:        GooglePlayUrl + link,
				ProductID:   productId,
				SerpapiLink: "",
				Rating:      rating,
				Author:      author,
				Category:    category,
				Downloads:   download,
				Thumbnail:   thumbnail,
				Description: desc,
				Video:       video,
			}
			if !item.isEmpty() {
				items = append(items, item)
			}
		}
	}
	return items
}

func makeSearchParam(params *RequestParams, actor string) SearchParameters {
	var store string
	switch actor {
	case GooglePlay:
		store = "apps"
	case GooglePlayGames:
		store = "games"
	case GooglePlayMovies:
		store = "movies"
	case GooglePlayBooks:
		store = "books"
	case GooglePlayProduct:
		if params.Store != "" {
			store = params.Store
		} else {
			store = "apps"
		}
	default:
		store = "apps"
	}
	searchInformation := SearchParameters{
		Q:           params.Q,
		ProductID:   params.ProductID,
		Actor:       actor,
		Hl:          params.HL,
		Gl:          params.GL,
		StoreDevice: params.StoreDevice,
		AppCategory: params.AppsCategory,
		Store:       store,
		SortBy:      params.SortBy,
		Platform:    params.Platform,
	}
	return searchInformation
}

// MakeProductInfo 获取product场景数据
func MakeProductInfo(wrbs [][]byte) (pctInfo ProductInfo, media Media, atApp AboutThisApp,
	badges []Badges, cgs []Categories, updateOn string, dataSafety []DataSafety, what WhatIsNew,
	ratings []Ratings, reviews []Review, contact DeveloperContact, apps []SimilarResults, pageToken string) {
	if len(wrbs) > 0 {
		// 此处做一个判断，获取ratings时一第一个Ws7gDc为准
		var ratingSym = 0
		for _, wrb := range wrbs {
			// 根据数据观测，Ws7gDc、ag2B9c、oCPfdb数据为有效数据
			rpcid := jsoniter.Get(wrb, 0, 1).ToString()
			switch rpcid {
			case "Ws7gDc":
				data := jsoniter.Get(wrb, 0, 2).ToString()
				dataBytes := []byte(data)
				pctInfo = getProductIndo(dataBytes)
				media = getMedia(dataBytes)
				atApp = getAboutThisApp(dataBytes)
				badges = getBadges(dataBytes)
				cgs = getCategories(dataBytes)
				updateOn = jsoniter.Get(dataBytes, 1, 2, 145, 0, 0).ToString()
				dataSafety = getDataSafety(dataBytes)
				what = getWhatIsNew(dataBytes)
				if ratingSym < 1 {
					ratings = getRatings(dataBytes)
					ratingSym = ratingSym + 1
				}
				contact = DeveloperContact{
					SupportEmail: jsoniter.Get(dataBytes, 1, 2, 69, 1, 0).ToString(),
				}
			case "ag2B9c":
				data := jsoniter.Get(wrb, 0, 2).ToString()
				dataBytes := []byte(data)
				apps = getApps(dataBytes)
			case "oCPfdb":
				data := jsoniter.Get(wrb, 0, 2).ToString()
				dataBytes := []byte(data)
				reviews, pageToken = getReviews(dataBytes)
			default:
				continue
			}
		}
	}
	return pctInfo, media, atApp, badges, cgs, updateOn, dataSafety, what, ratings, reviews, contact, apps, pageToken
}

func getApps(dataBytes []byte) (similarResults []SimilarResults) {
	moreAppStr := jsoniter.Get(dataBytes, 1, 1, 1, 21).ToString()
	SimilarAppStr := jsoniter.Get(dataBytes, 1, 1, 0, 21).ToString()
	link := jsoniter.Get([]byte(moreAppStr), 1, 2, 4, 2).ToString()

	moreAppItems := getProductItems(moreAppStr)
	SimilarAppItems := getProductItems(SimilarAppStr)
	moreApp := SimilarResults{
		Title:        jsoniter.Get([]byte(moreAppStr), 1, 0).ToString(),
		SeeMoreLink:  GooglePlayUrl + link,
		SeeMoreToken: jsoniter.Get([]byte(moreAppStr), 1, 3, 1).ToString(),
		SerpapiLink:  "",
		Items:        moreAppItems,
	}
	similarApp := SimilarResults{
		Title:        jsoniter.Get([]byte(SimilarAppStr), 1, 0).ToString(),
		SeeMoreLink:  GooglePlayUrl + link,
		SeeMoreToken: jsoniter.Get([]byte(SimilarAppStr), 1, 3, 1).ToString(),
		SerpapiLink:  "",
		Items:        SimilarAppItems,
	}
	similarResults = append(similarResults, moreApp, similarApp)
	return similarResults
}

func getProductItems(similarStr string) (similarItems []Item) {
	similarArrStr := jsoniter.Get([]byte(similarStr), 0).ToString()
	var similarArr []interface{}
	err := json.Unmarshal([]byte(similarArrStr), &similarArr)
	if err != nil {
		return similarItems
	}
	for _, it := range similarArr {
		marshal, err := json.Marshal(it)
		if err != nil {
			log.Errorf("json unmarsh failed,%v", err)
		}
		similarItems = append(similarItems, Item{
			Title:       jsoniter.Get(marshal, 3).ToString(),
			Link:        GooglePlayUrl + jsoniter.Get(marshal, 10, 4, 2).ToString(),
			ProductID:   jsoniter.Get(marshal, 0, 0).ToString(),
			SerpapiLink: "",
			Rating:      jsoniter.Get(marshal, 4, 0).ToFloat64(),
			Thumbnail:   jsoniter.Get(marshal, 1, 3, 2).ToString(),
			Extension:   []string{jsoniter.Get(marshal, 14).ToString()},
		})
	}
	return similarItems
}

func getReviews(dataBytes []byte) ([]Review, string) {
	pageToken := jsoniter.Get(dataBytes, 1, 1).ToString()
	var reviews []Review
	reviewStr := jsoniter.Get(dataBytes, 0).ToString()
	var reviewMarsh []interface{}
	err := json.Unmarshal([]byte(reviewStr), &reviewMarsh)
	if err != nil {
		// 如果序列化失败，则使用当前的rpcid，直接返回空内容
		return reviews, ""
	}
	for _, it := range reviewMarsh {
		marshal, err := json.Marshal(it)
		if err != nil {
			log.Error(err)
			continue
		}
		timestamp := jsoniter.Get(marshal, 5, 0).ToInt64()
		t := time.Unix(timestamp, 0)
		formattedDate := t.Format("January 2, 2006")
		formattedUTC := t.UTC().Format(time.RFC3339)
		reviews = append(reviews, Review{
			ID:      jsoniter.Get(marshal, 0).ToString(),
			Title:   jsoniter.Get(marshal, 1, 0).ToString(),
			Avatar:  jsoniter.Get(marshal, 1, 1, 3, 2).ToString(),
			Rating:  jsoniter.Get(marshal, 2).ToFloat32(),
			Snippet: jsoniter.Get(marshal, 4).ToString(),
			Likes:   jsoniter.Get(marshal, 6).ToInt(),
			Date:    formattedDate,
			ISODate: formattedUTC,
		})
	}
	return reviews, pageToken
}

func getRatings(dataBytes []byte) []Ratings {
	var ratings []Ratings
	rateStr := jsoniter.Get(dataBytes, 1, 2, 51, 1).ToString()
	var rateStrMarsh []interface{}
	err := json.Unmarshal([]byte(rateStr), &rateStrMarsh)
	if err != nil {
		log.Errorf("unmarshal json: %v", err)
		return ratings
	}
	for i, marsh := range rateStrMarsh {
		if marsh == nil {
			continue
		}
		marshal, err := json.Marshal(marsh)
		if err != nil {
			log.Errorf("marshal json: %v", err)
			continue
		}
		count := jsoniter.Get(marshal, 1).ToInt64()
		ratings = append(ratings, Ratings{
			Stars: i,
			Count: count,
		})
	}
	return ratings
}

func getWhatIsNew(dataBytes []byte) WhatIsNew {
	return WhatIsNew{
		Snippet: jsoniter.Get(dataBytes, 1, 2, 144, 1, 1).ToString(),
	}
}

func getDataSafety(dataBytes []byte) (safeties []DataSafety) {
	dataStr := jsoniter.Get(dataBytes, 1, 2, 136, 1).ToString()
	var dataStrMarsh []interface{}
	err := json.Unmarshal([]byte(dataStr), &dataStrMarsh)
	if err != nil {
		return safeties
	}
	for _, item := range dataStrMarsh {
		marshal, err := json.Marshal(item)
		if err != nil {
			log.Errorf("marshal json: %v", err)
			continue
		}
		subtext := jsoniter.Get(marshal, 2, 1).ToString()
		unescapeString := html.UnescapeString(subtext)
		var link = ""
		if strings.Contains(unescapeString, "href") {
			// 定义正则表达式模式
			hrefPattern := `href=["']([^"']+)["']`
			// 这里我们匹配 <a> 标签内的内容和紧随其后的文本
			pattern := `<a[^>]*>(.*?)<\/a>\s*(.*)`

			hrefRegex := regexp.MustCompile(hrefPattern)
			aRegex := regexp.MustCompile(pattern)

			// 查找 href 和 target
			hrefMatch := hrefRegex.FindStringSubmatch(unescapeString)
			match := aRegex.FindStringSubmatch(unescapeString)
			if len(hrefMatch) < 2 || len(match) < 2 {
				log.Errorf("无法匹配 href、target 或链接文本")
				continue
			}
			if len(match) > 2 {
				subtext = match[1] + match[2]
			} else {
				subtext = match[1]
			}
			link = hrefMatch[1]
		}
		safeties = append(safeties, DataSafety{
			Text:    jsoniter.Get(marshal, 1).ToString(),
			Subtext: subtext,
			Link:    link,
		})
	}
	return safeties
}

func getCategories(dataBytes []byte) (categories []Categories) {
	cgsStr := jsoniter.Get(dataBytes, 1, 2, 79, 0).ToString()
	var cgsStrMarsh []interface{}
	err := json.Unmarshal([]byte(cgsStr), &cgsStrMarsh)
	if err != nil {
		return categories
	}
	for _, item := range cgsStrMarsh {
		marshal, err := json.Marshal(item)
		if err != nil {
			log.Errorf("marshal json: %v", err)
			continue
		}
		categories = append(categories, Categories{
			Name:        jsoniter.Get(marshal, 0).ToString(),
			Link:        GooglePlayUrl + jsoniter.Get(marshal, 1, 4, 2).ToString(),
			CategoryID:  jsoniter.Get(marshal, 2).ToString(),
			SerpapiLink: "",
		})
	}
	return categories
}

func getBadges(dataBytes []byte) []Badges {
	var badges []Badges
	badge2 := jsoniter.Get(dataBytes, 1, 2, 58, 0).ToString()
	badge1 := jsoniter.Get(dataBytes, 1, 2, 58, 2).ToString()
	return append(badges, Badges{Name: badge1 + " " + badge2})
}

func getAboutThisApp(dataBytes []byte) (about AboutThisApp) {
	var permissions []Permission
	detailsStr := jsoniter.Get(dataBytes, 1, 2, 74, 2, 0).ToString()
	var detailStrMarsh []interface{}
	err := json.Unmarshal([]byte(detailsStr), &detailStrMarsh)
	if err != nil {
		return about
	}
	for _, item := range detailStrMarsh {
		marshal, err := json.Marshal(item)
		if err != nil {
			log.Errorf("marshal json: %v", err)
			continue
		}
		detStr := jsoniter.Get(marshal, 2).ToString()
		var detStrArr []interface{}
		err = json.Unmarshal([]byte(detStr), &detStrArr)
		if err != nil {
			log.Errorf("unmarshal json: %v", err)
			continue
		}
		var details []string
		for _, item := range detStrArr {
			it, err := json.Marshal(item)
			if err != nil {
				log.Errorf("marshal json: %v", err)
			}
			det := jsoniter.Get(it, 1).ToString()
			details = append(details, det)
		}
		var permission = Permission{
			Type:    jsoniter.Get(marshal, 0).ToString(),
			Details: details,
		}
		permissions = append(permissions, permission)
	}
	return AboutThisApp{
		Snippet:             jsoniter.Get(dataBytes, 1, 2, 72, 0, 1).ToString(),
		InAppPurchases:      jsoniter.Get(dataBytes, 1, 2, 19, 0).ToString(),
		ReleasedOn:          jsoniter.Get(dataBytes, 1, 2, 10, 0).ToString(),
		UpdatedOn:           jsoniter.Get(dataBytes, 1, 2, 145, 0, 0).ToString(),
		Downloads:           jsoniter.Get(dataBytes, 1, 2, 13, 0).ToString(),
		ContentRating:       jsoniter.Get(dataBytes, 1, 2, 9, 0).ToString(),
		InteractiveElements: jsoniter.Get(dataBytes, 1, 2, 9, 3, 1).ToString(),
		OfferedBy:           jsoniter.Get(dataBytes, 1, 2, 37, 0).ToString(),
		Permissions:         permissions,
	}

}

func getMedia(dataBytes []byte) (media Media) {
	var video = Video{
		Thumbnail: jsoniter.Get(dataBytes, 1, 2, 100, 0, 1, 3, 2).ToString(),
		Link:      jsoniter.Get(dataBytes, 1, 2, 100, 0, 0, 3, 2).ToString(),
	}
	var images []string
	imgStr := jsoniter.Get(dataBytes, 1, 2, 78, 0).ToString()
	var imgStrMarsh []interface{}
	err := json.Unmarshal([]byte(imgStr), &imgStrMarsh)
	if err != nil {
		return media
	}
	for _, marsh := range imgStrMarsh {
		marshal, err := json.Marshal(marsh)
		if err != nil {
		}
		imgUrl := jsoniter.Get(marshal, 3, 2).ToString()
		images = append(images, imgUrl)
	}
	media = Media{
		Video:  video,
		Images: images,
	}
	return media
}

func getProductIndo(dataByte []byte) (pctInfo ProductInfo) {
	rating := jsoniter.Get(dataByte, 1, 2, 51, 0, 0).ToFloat32()
	title := jsoniter.Get(dataByte, 1, 2, 0, 0).ToString()
	author := Author{
		Name: jsoniter.Get(dataByte, 1, 2, 68, 0).ToString(),
		Link: GooglePlayUrl + jsoniter.Get(dataByte, 1, 2, 68, 1, 4, 2).ToString(),
	}
	extension := jsoniter.Get(dataByte, 1, 2, 48, 0).ToString()
	reviews := jsoniter.Get(dataByte, 1, 2, 51, 21).ToInt64()
	contentRating := ContentRating{
		Text:      jsoniter.Get(dataByte, 1, 2, 9, 0).ToString(),
		Thumbnail: jsoniter.Get(dataByte, 1, 2, 9, 1, 3, 2).ToString(),
	}
	downloads := jsoniter.Get(dataByte, 1, 2, 13, 2).ToString()
	thumbnail := jsoniter.Get(dataByte, 1, 2, 95, 0, 3, 2).ToString()
	offersStr := jsoniter.Get(dataByte, 1, 2, 57, 0).ToString()
	var offersUnMarsh []interface{}
	err := json.Unmarshal([]byte(offersStr), &offersUnMarsh)
	if err != nil {
		return pctInfo
	}
	var offers []Offer
	for _, marsh := range offersUnMarsh {
		marshal, err := json.Marshal(marsh)
		if err != nil {
			log.Errorf("marsh marshal: %v", err)
		}
		text := jsoniter.Get(marshal, 1, 0).ToString()
		link := jsoniter.Get(marshal, 0, 0, 6, 5, 2).ToString()
		offers = append(offers, Offer{text, link})
	}
	pctInfo = ProductInfo{
		Title:         title,
		Authors:       []Author{author},
		Extensions:    []string{extension},
		Rating:        rating,
		Reviews:       reviews,
		ContentRating: contentRating,
		Downloads:     downloads,
		Thumbnail:     thumbnail,
		Offers:        offers,
	}
	return pctInfo
}
