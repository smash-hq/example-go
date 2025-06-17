package play

import (
	"context"
	"encoding/json"
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/avast/retry-go/v4"
	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
	log "github.com/sirupsen/logrus"
)

const (
	GooglePlayUrl        = "https://play.google.com"
	GooglePlayStoreUI    = "https://play.google.com/_/PlayStoreUi/data/batchexecute"
	GoogleAppUrl         = GooglePlayUrl + "/store/apps"
	GooglePlayProductUrl = GooglePlayUrl + "/store/apps/details"
	GooglePlayGamesUrl   = GooglePlayUrl + "/store/games"
	GooglePlayMoviesUrl  = GooglePlayUrl + "/store/movies"
	GooglePlayBooksUrl   = GooglePlayUrl + "/store/books"
	GooglePlaySearchUrl  = GooglePlayUrl + "/store/search"
)

func DoPlay(ctx context.Context, params *RequestParams, actor, proxy string) (*Response, error) {
	var beginTime = time.Now().UnixMilli()
	queryUrl := urlByQuery(params, actor)
	fetchUrl, err := urlByFormData(GooglePlayStoreUI, params, actor)
	if err != nil {
		log.Errorf("field validation error: %v", err)
		return nil, err
	}
	fReq := constructFReq(params, actor)
	respBytes, err := batchExecutor(ctx, fetchUrl, fReq, proxy)
	if err != nil {
		log.Errorf("request err:%v", err)
		return nil, err
	}
	var response Response
	searchMetaData := searchMetadata(queryUrl, beginTime)
	if !searchMetaData.IsEmpty() {
		response.SearchMetadata = &searchMetaData
	}

	searchInformation := makeSearchParam(params, actor)
	if !searchInformation.isEmpty() {
		response.SearchParameters = &searchInformation
	}
	// 组装organic_results、highlights、chartOptions、 pageToken
	organics, highlights, chartOptions, pageToken := MakeConversionParsing(respBytes)
	if len(organics) > 0 {
		response.OrganicResults = organics
	}
	if len(highlights) > 0 {
		response.HighlightItem = highlights
	}
	if len(chartOptions) > 0 {
		response.ChartOption = chartOptions
	}
	serpapiPagination := serpapiSectionPag(params, pageToken)
	if !serpapiPagination.isEmpty() {
		response.SerpapiPagination = &serpapiPagination
	}
	return &response, nil
}

func DoProductPlay(ctx context.Context, params *RequestParams, actor, proxy string) (*Response, error) {
	var beginTime = time.Now().UnixMilli()
	fetchUrl, err := urlByFormData(GooglePlayStoreUI, params, actor)
	queryUrl := urlByQuery(params, actor)
	if err != nil {
		log.Errorf("field validation error: %v", err)
		return nil, err
	}
	fReq := constructFReqToProduct(params)
	respBytes, err := batchExecutor(ctx, fetchUrl, fReq, proxy)
	//lens.SaveFile("resp_product.json", string(respBytes))
	wrbs := GetWrbs(respBytes)
	info, media, app, badges, cgs, on, safety, what, ratings, reviews, contact, apps, pageToken := MakeProductInfo(wrbs)
	var response = Response{
		ProductInfo:      &info,
		Media:            &media,
		AboutThisApp:     &app,
		Badges:           badges,
		Categories:       cgs,
		UpdatedOn:        on,
		DataSafety:       safety,
		WhatIsNew:        &what,
		Ratings:          ratings,
		Reviews:          reviews,
		DeveloperContact: &contact,
		SimilarResults:   apps,
	}
	response.SerpapiPagination = &SerpapiSectionPagination{
		NextPageToken: pageToken,
	}
	searchMetaData := searchMetadata(queryUrl, beginTime)
	if !searchMetaData.IsEmpty() {
		response.SearchMetadata = &searchMetaData
	}

	searchInformation := makeSearchParam(params, actor)
	if !searchInformation.isEmpty() {
		response.SearchParameters = &searchInformation
	}
	return &response, nil
}

func serpapiSectionPag(params *RequestParams, pageToken string) SerpapiSectionPagination {
	//var sectionPageToken string
	//if params.NextPageToken != "" {
	//	sectionPageToken = params.NextPageToken
	//}
	var serpapiPagination = SerpapiSectionPagination{
		NextPageToken: pageToken,
		//Next:             "",
		//SectionPageToken: sectionPageToken,
	}
	return serpapiPagination
}

func searchMetadata(queryUrl string, beginTime int64) SearchMetadata {
	searchMetaData := SearchMetadata{
		//ID:             "todo 生成id",
		Status: "success",
		//JSONEndpoint:   "No object storage",
		CreatedAt:      time.Now().Format("2006-01-02 15:04:05 UTC"),
		ProcessedAt:    time.Now().Format("2006-01-02 15:04:05 UTC"),
		GooglePlayURL:  queryUrl,
		TotalTimeTaken: float64(time.Now().UnixMilli()-beginTime) / 1000,
	}
	return searchMetaData
}

func urlByQuery(params *RequestParams, actor string) string {
	var path string
	values := url.Values{}
	values.Set("hl", params.HL)
	values.Set("gl", params.GL)
	switch actor {
	case GooglePlayBooks:
		if params.BooksCategory == "" {
			if params.Q != "" {
				values.Set("q", params.Q)
				path = GooglePlaySearchUrl
			} else {
				path = GooglePlayBooksUrl
			}
		} else {
			path = GooglePlayBooksUrl + "/category/" + params.BooksCategory
		}
		if params.Age != "" {
			values.Set("age", params.Age)
		}
	case GooglePlayGames:
		if params.GamesCategory != "" {
			path = GoogleAppUrl + "/category/" + params.GamesCategory
		} else {
			if params.Q != "" {
				values.Set("q", params.Q)
				path = GooglePlaySearchUrl
			} else {
				if params.StoreDevice != "" {
					values.Set("device", params.StoreDevice)
				}
				path = GooglePlayGamesUrl
			}
		}
	case GooglePlayMovies:
		if params.Age != "" {
			values.Set("age", params.Age)
		}
		if params.Q != "" {
			values.Set("q", params.Q)
			path = GooglePlaySearchUrl
		} else {
			path = GooglePlayMoviesUrl + "/category/" + params.MoviesCategory
		}
	case GooglePlayProduct:
		values.Set("id", params.ProductID)
		path = GooglePlayProductUrl
	case GooglePlay:
		if params.AppsCategory != "" {
			path = GoogleAppUrl + "/category/" + params.AppsCategory
		} else {
			if params.Q != "" {
				values.Set("q", params.Q)
				path = GooglePlaySearchUrl
			} else {
				if params.StoreDevice != "" {
					values.Set("device", params.StoreDevice)
				}
				path = GoogleAppUrl
			}
		}
	}

	return path + "?" + values.Encode()
}

// 根据参数组装req
func constructFReq(params *RequestParams, actor string) (fReq string) {
	// 一共有三个数组
	var fReqArray [][]interface{}
	// 默认使用的rpcids
	var rpcids = rpcidFur(params, actor)
	firstArr := firstArrFur(params, actor, rpcids)
	secondArr := secondArrFur(params, actor, rpcids)
	thirdArr := thirdArrFur(params, actor, rpcids)
	fReqArray = append(fReqArray, firstArr, secondArr)
	if len(thirdArr) > 0 {
		fReqArray = append(fReqArray, thirdArr)
	}
	var require = []interface{}{
		fReqArray,
	}
	fReqMarsh, _ := json.Marshal(require)
	//lens.SaveFile("freq.json", string(fReqMarsh))
	escape := url.QueryEscape(string(fReqMarsh))
	return "f.req=" + escape
}

func rpcidFur(params *RequestParams, actor string) []string {
	var rpcids []string
	if params.Q != "" {
		rpcids = []string{"AZO9Cb", "lGYRle"}
	} else {
		switch {
		case actor == GooglePlayMovies:
			rpcids = []string{"eIpeLd", "w3QCWb", "w37aie"}
		case actor == GooglePlayBooks:
			rpcids = []string{"eIpeLd", "w3QCWb", "w37aie"}
		case actor == GooglePlayProduct:
			rpcids = []string{"CLXjtf", "A6yuRe", "Ws7gDc", "ZittHe", "yowZ5", "ag2B9c", "e7uDs", "Ws7gDc", "oCPfdb"}
		default:
			if params.StoreDevice != "" {
				if params.StoreDevice == "phone" {
					rpcids = []string{"eIpeLd", "w3QCWb"}
				} else {
					rpcids = []string{"eIpeLd", "di6f4"}
				}
			} else {
				if actor == GooglePlayGames {
					if params.GamesCategory != "" {
						rpcids = []string{"eIpeLd", "w3QCWb"}
					} else {
						rpcids = []string{"eIpeLd", "di6f4"}
					}
				} else {
					rpcids = []string{"eIpeLd", "w3QCWb"}
				}
			}
		}
	}
	return rpcids
}

func thirdArrFur(params *RequestParams, actor string, rpcids []string) []interface{} {
	var third []interface{}
	if actor == GooglePlayBooks || params.Q == "" {
		// eg: ["w37aie","[null,2,\"subj_Comics___Graphic_Novels.Crime___Mystery\"]",null,"3"]
		var category = params.BooksCategory
		if category == "" {
			return third
		}
		markArr3 := []interface{}{nil, 2, category}
		marshal, err := json.Marshal(markArr3)
		if err != nil {
			log.Errorln("marshal err:", err)
		} else {
			third = append(third, rpcids[2], string(marshal), nil, "3")
		}
	}
	return third
}

func secondArrFur(params *RequestParams, actor string, rpcids []string) []interface{} {
	var markArr1, markArr2 = generateMarkArray1(params, actor)
	var marksArray = []interface{}{
		markArr1, markArr2,
	}
	arrMarsh, _ := json.Marshal(marksArray)
	// >====引号内的内容===<
	var second = []interface{}{
		rpcids[1], string(arrMarsh), nil, "2",
	}
	return second
}

func firstArrFur(params *RequestParams, actor string, rpcids []string) (firstArr []interface{}) {
	var sceneCategoryArr []interface{}
	var scene, category string
	if params.Q != "" {
		sceneCategoryArr = []interface{}{params.Q}
		var querySym = 4
		if actor == GooglePlayBooks {
			querySym = 2
		} else if actor == GooglePlayMovies {
			querySym = 1
		}
		sceneCategoryArr = append(sceneCategoryArr, querySym)
		// 对于书籍有价格、参数，进行特殊处理
		if params.Price != "" {
			atoi, err := strconv.Atoi(params.Price)
			if err != nil {
				// 默认为免费
				atoi = 1
			}
			sceneCategoryArr = append(sceneCategoryArr, nil, []interface{}{nil, atoi})
		}
	} else {
		switch {
		case actor == GooglePlayMovies:
			if params.MoviesCategory != "" {
				scene = params.MoviesCategory
				category = params.MoviesCategory
			} else {
				scene = "MOVIE"
				category = "MOVIE"
			}
		case actor == GooglePlayBooks:
			if params.BooksCategory != "" {
				scene = params.BooksCategory
				category = params.BooksCategory
			} else {
				scene = "ebooks"
				category = "ebooks"
			}
		default:
			if params.StoreDevice != "" {
				if actor == GooglePlay {
					scene = "APPLICATION"
				} else if actor == GooglePlayGames {
					scene = "GAME"
				}
				category = params.StoreDevice
			} else {
				if params.AppsCategory != "" {
					scene = params.AppsCategory
					category = params.AppsCategory
				} else if params.GamesCategory != "" {
					scene = params.GamesCategory
					category = params.GamesCategory
				} else {
					if actor == GooglePlay {
						scene = "APPLICATION"
						category = "phone"
					} else if actor == GooglePlayGames {
						scene = "GAME"
						category = "windows"
					}
				}
			}
		}
		sceneIndex := engineEnum[actor]
		sceneCategoryArr = []interface{}{
			sceneIndex, scene, category,
		}
	}
	sceneCategory, _ := json.Marshal(sceneCategoryArr)

	first := []interface{}{
		rpcids[0], string(sceneCategory), nil, "1",
	}
	return first
}

// 根据参数组装req，product场景
func constructFReqToProduct(params *RequestParams) (fReq string) {
	var fReqArray [][]interface{}
	for i := 0; i < 9; i++ {
		mark := productMark(params, i)
		fReqArray = append(fReqArray, mark)
	}
	var require = []interface{}{
		fReqArray,
	}
	fReqMarsh, _ := json.Marshal(require)
	//lens.SaveFile("freq.json", string(fReqMarsh))
	escape := url.QueryEscape(string(fReqMarsh))
	return "f.req=" + escape
}

// 使用http获取数组方式组装url进行数据爬取
func urlByFormData(baseUrl string, params *RequestParams, actor string) (fetchUrl string, err error) {
	// 构造查询字符串
	queryValues := url.Values{}
	var q = &QueryParams{}
	if params.HL != "" {
		q.HL = params.HL
		queryValues.Set("hl", params.HL)
	}
	if params.GL != "" {
		q.GL = params.GL
		queryValues.Set("gl", params.GL)
	}
	device := params.StoreDevice
	var rpcids []string
	if device == "" {
		rpcids = rpcidsEnum[storeDevices[0]]
	} else {
		rpcids = rpcidsEnum[device]
	}
	var path string
	switch actor {
	case GooglePlay:
		if params.Q != "" {
			path = StoreSearchPath
		} else {
			path = AppsPath
		}
	case GooglePlayBooks:
		if params.Q != "" {
			path = StoreSearchPath
		}
		// 如果不存在books_category参数，则认为默认为/store/category，否则为/books/category/%s
		if params.BooksCategory != "" {
			path = BooksPath + "/category/" + params.BooksCategory
		} else {
			path = BooksPath
		}
		if params.Price != "" {
			queryValues.Set("price", params.Price)
		}
		if params.Age != "" {
			queryValues.Set("age", params.Age)
		}
	case GooglePlayGames:
		if params.Q != "" {
			path = StoreSearchPath
		}
		// 如果不存在games_category参数，则认为默认为/store/category，否则为/games/category/%s
		if params.GamesCategory != "" {
			rpcids = rpcidsEnum["default"]
			path = AppsPath + "/category/" + params.GamesCategory
		} else {
			path = GamesPath
		}
	case GooglePlayMovies:
		if params.Q != "" {
			path = StoreSearchPath
		} else {
			path = MoviesPath
		}
	case GooglePlayProduct:
		path = ProductPath
	default:
		path = AppsPath
	}

	if params.ProductID != "" {
		rpcids = rpcidsEnum["product"]
	}
	q.SourcePath = path
	q.RPCIDs = rpcids
	// 此处判断engine，如果有则拼接下来的参数
	err = q.FieldValidation()
	if err != nil {
		log.Errorf("field validation error: %v", err)
		return "", err
	}
	if len(rpcids) > 0 {
		queryValues.Set("rpcids", strings.Join(rpcids, ","))
	}
	if q.BL != "" {
		queryValues.Set("bl", q.BL)
	}
	if q.SourcePath != "" {
		queryValues.Set("source-path", q.SourcePath)
	}
	if q.SocApp != "" {
		queryValues.Set("soc-app", q.SocApp)
	}
	if q.SocPlatform != "" {
		queryValues.Set("soc-platform", q.SocPlatform)
	}
	if q.SocDevice != "" {
		queryValues.Set("soc-device", q.SocDevice)
	}
	if q.RT != "" {
		queryValues.Set("rt", q.RT)
	}
	queryValues.Set("f.sid", "-1194886866173724340")
	queryValues.Set("_reqid", "62273")
	queryValues.Set("authuser", "")
	return baseUrl + "?" + queryValues.Encode(), nil
}

// 调用https://play.google.com/_/PlayStoreUi/data/batchexecute接口获取数组
func batchExecutor(ctx context.Context, url, fReq, proxy string) (respBytes []byte, err error) {
	jar := tls_client.NewCookieJar()
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(profiles.Chrome_124),
		tls_client.WithCookieJar(jar),
	}
	client, _ := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	client.SetProxy(proxy)
	//client.SetProxy("http://127.0.0.1:7890")

	resp, err := retry.DoWithData(func() (*http.Response, error) {
		req, reqError := http.NewRequest("POST", url, strings.NewReader(fReq))
		if reqError != nil {
			log.Errorf("request err:%v", reqError)
			return nil, err
		}
		req.Header.Set("accept", "*/*")
		req.Header.Set("accept-language", "en")
		req.Header.Set("cache-control", "no-cache")
		req.Header.Set("content-type", "application/x-www-form-urlencoded;charset=UTF-8")
		req.Header.Set("origin", "https://play.google.com")
		req.Header.Set("pragma", "no-cache")
		req.Header.Set("priority", "u=1, i")
		req.Header.Set("referer", "https://play.google.com/")
		req.Header.Set("sec-ch-ua", `"Chromium";v="134", "Not:A-Brand";v="24", "Google Chrome";v="134"`)
		req.Header.Set("sec-ch-ua-arch", `"x86"`)
		req.Header.Set("sec-ch-ua-bitness", `"64"`)
		req.Header.Set("sec-ch-ua-form-factors", `"Desktop"`)
		req.Header.Set("sec-ch-ua-full-version", `"134.0.6998.35"`)
		req.Header.Set("sec-ch-ua-full-version-list", `"Chromium";v="134.0.6998.35", "Not:A-Brand";v="24.0.0.0", "Google Chrome";v="134.0.6998.35"`)
		req.Header.Set("sec-ch-ua-mobile", "?0")
		req.Header.Set("sec-ch-ua-model", `""`)
		req.Header.Set("sec-ch-ua-platform", `"Windows"`)
		req.Header.Set("sec-ch-ua-platform-version", `"10.0.0"`)
		req.Header.Set("sec-ch-ua-wow64", "?0")
		req.Header.Set("sec-fetch-dest", "empty")
		req.Header.Set("sec-fetch-mode", "cors")
		req.Header.Set("sec-fetch-site", "same-origin")
		req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36")
		req.Header.Set("x-same-domain", "1")
		resp, respErr := client.Do(req)
		log.Infof("%s-->request url: %s", GooglePlay, req.URL.String())
		if respErr != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			return nil, err
		}
		return resp, nil
	},
		retry.Context(ctx),
		retry.Attempts(10),
		retry.Delay(100*time.Millisecond),
		retry.OnRetry(func(attempt uint, err error) {
			log.Warnf("fetch URI=%s count=%v error=%v", GooglePlay, attempt, err)
		}),
		retry.LastErrorOnly(true),
	)
	if err != nil {
		log.Errorf("request err:%v", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Errorf("body close err:%v", err)
		}
	}(resp.Body)
	respBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("request err:%v", err)
		return nil, err
	}
	return respBytes, nil
}
