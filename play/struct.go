package play

import (
	"errors"
	"strings"
)

type RequestParams struct {
	Type string `json:"type,omitempty"`
	/*
		google_play
		google_play_product
		google_play_games
		google_play_movies
		google_play_books
	*/
	//Actor string `json:"engine,omitempty"`
	Q  string `json:"q"`            // query查询内容
	QS string `json:"qs"`           // querys查询内容
	GL string `json:"gl,omitempty"` // country
	HL string `json:"hl,omitempty"` // Language
	/*
		apps_category分类如下
		ART_AND_DESIGN
		AUTO_AND_VEHICLES
		......
	*/
	AppsCategory   string `json:"apps_category,omitempty"`   // 应用程序类别
	MoviesCategory string `json:"movies_category,omitempty"` // 电影类别
	GamesCategory  string `json:"games_category,omitempty"`  // 电影类别
	BooksCategory  string `json:"books_category,omitempty"`  // 书籍类别
	// 1-->Free 2-->Paid
	Price string `json:"price,omitempty"` // 书籍价格
	/*
		store_device分类如下
		phone - Phone device (default)
		tablet - Tablet device
		tv - TV device
		chromebook - Chromebook device
		watch - Watch device
		car - Car device
		以上仅为application
	*/
	StoreDevice string `json:"store_device,omitempty"` // 使用设备类别
	/*
		// 年龄区间，age分类如下
		AGE_RANGE1 --> 0~5
		AGE_RANGE2 --> 6~8
		AGE_RANGE3 --> 9~12
	*/
	Age string `json:"age,omitempty"`
	//Pagination *Pagination `json:"pagination,omitempty"`
	// Parameter defines the next page token. It is used for retrieving the next page results. It shouldn't be used with the section_page_token, see_more_token, and chart parameters.
	// It should be used only when all_reviews parameter is set to true.
	NextPageToken string `json:"next_page_token,omitempty"`
	// Parameter defines the section page token used for retrieving the pagination results from individual sections.
	// This parameter is a safer version of see_more_token, and is found in every row you can paginate into. It shouldn't be used with the next_page_token, see_more_token, and chart parameters
	SectionPageToken string `json:"section_page_token,omitempty"`
	// Parameter is used for showing top charts. It can return up to 50 results. Each store contains different charts which require different values for retrieving results.
	// To get the value of a specific chart you can use our Google Play Apps Store API JSON output: chart_options[index].value (e.g. chart=topselling_free).
	// It shouldn't be used with the section_page_token, see_more_token, and next_page_token parameters
	Chart string `json:"chart,omitempty"`
	// Parameter defines the see more token used for retrieving the pagination results from individual sections It is usually found in next page results.
	// It shouldn't be used with the section_page_token, next_page_token, and chart, parameters
	SeeMoreToken    string `json:"see_more_token,omitempty"`
	IsReturnRawHtml string `json:"is_return_raw_html"`
	// product场景参数
	ProductID string `json:"product_id,omitempty"`
	Store     string `json:"store,omitempty"`
	// 仅当store参数设置为tv时，才应使用它。例如store=tv。
	SeasonID string `json:"season_id,omitempty"`
	// 以下为product场景参数中有关评论的参数
	// true or false，是否显示所有评论
	AllReviews string `json:"all_reviews,omitempty"`
	// 平台端，Phone、Watch、Chromebook、TV、Car
	Platform string `json:"platform,omitempty"`
	// 1、2、3、4、5star,接收1、2、3、4、5
	Rating string `json:"rating,omitempty"`
	// 1-->Most relevant、2-->Newest、3-->Rating
	SortBy string `json:"sort_by,omitempty"`
	// 要显示的评论数量
	Num string `json:"num,omitempty"`
}

func (params *RequestParams) FieldValidation(actor string) error {
	if actor == "" {
		actor = GooglePlay
		params.Type = GooglePlay
	} else {
		if actor == GooglePlayProduct {
			if params.ProductID == "" {
				return errors.New("`product_id` parameters can't be empty")
			}
			if params.SortBy == "" {
				params.SortBy = "1"
			} else {

			}
			if params.Platform == "" {
				params.Platform = "phone"
			}
			if params.NextPageToken != "" && params.AllReviews == "false" {
				return errors.New("`next_page_token`,It is used for retrieving the next page results. It should be used only when all_reviews parameter is set to true")
			}
			if params.SeasonID != "" && params.Store != "tv" {
				return errors.New("`season_id`,It should be used only when store parameter is set to tv. e.g. store=tv")
			}
			if params.Num == "" {
				params.Num = "20"
			}
		} else if actor == GooglePlayBooks {
			if params.BooksCategory != "" {
				b := contains(bookCategoryEnums, params.BooksCategory)
				if !b {
					return errors.New("`books_category` parameter is invalid")
				}
			}
			if params.Price != "" {
				if params.Price != "1" && params.Price != "2" {
					return errors.New("`price` parameter is invalid")
				}
				if params.Q != "" {
					return errors.New("`price` should be used only in combination with the q parameter")
				}
			}
			if params.Age != "" {
				b2 := contains(ageRanges, params.Age)
				if !b2 {
					return errors.New("`age` parameter is invalid")
				}
			}
		} else {
			if params.ProductID != "" {
				return errors.New("`product_id` parameters cannot have value")
			}
		}
	}
	if params.HL == "" {
		params.HL = "en-US"
	}
	if params.GL == "" {
		params.GL = "us"
	}

	if params.AppsCategory != "" && params.StoreDevice != "" {
		return errors.New("`apps_category` and `store_device` parameters can't be used together")
	} else {
		if params.StoreDevice != "" {
			b := contains(storeDevices, params.StoreDevice)
			if !b {
				return errors.New("`store_device` parameters is illegal, It should be included in [" + strings.Join(storeDevices, ",") + "]")
			}
		}
		if params.AppsCategory != "" {
			b := contains(appsCategoryEnum, params.AppsCategory)
			if !b {
				return errors.New("`apps_category` parameters is illegal, It should be included in [" + strings.Join(appsCategoryEnum, ",") + "]")
			}
		}
	}
	if params.GamesCategory != "" && params.StoreDevice != "" {
		return errors.New("`games_category` and `store_device` parameters can't be used together")
	} else {
		if params.StoreDevice != "" {
			b := contains(storeDevices, params.StoreDevice)
			if !b {
				return errors.New("`store_device` parameters is illegal, It should be included in [" + strings.Join(storeDevices, ",") + "]")
			}
		}
		if params.GamesCategory != "" {
			b := contains(gamesCategoryEnum, params.GamesCategory)
			if !b {
				return errors.New("`games_category` parameters is illegal, It should be included in [" + strings.Join(gamesCategoryEnum, ",") + "]")
			}
		}
	}
	var pagination = &Pagination{
		NextPageToken:    params.NextPageToken,
		SectionPageToken: params.SectionPageToken,
		Chart:            params.Chart,
		SeeMoreToken:     params.SeeMoreToken,
	}
	if pagination != nil {
		err := pagination.FieldValidation()
		if err != nil {
			return err
		}
	}
	if params.Age != "" {
		b := contains(ageRanges, params.Age)
		if !b {
			return errors.New("`age` parameters is illegal")
		}
	}
	return nil
}

// Pagination 四个参数只能选择其中一个使用，不可混用
type Pagination struct {
	// Parameter defines the next page token. It is used for retrieving the next page results. It shouldn't be used with the section_page_token, see_more_token, and chart parameters.
	NextPageToken string `json:"next_page_token,omitempty"`
	// Parameter defines the section page token used for retrieving the pagination results from individual sections.
	// This parameter is a safer version of see_more_token, and is found in every row you can paginate into. It shouldn't be used with the next_page_token, see_more_token, and chart parameters
	SectionPageToken string `json:"section_page_token,omitempty"`
	// Parameter is used for showing top charts. It can return up to 50 results. Each store contains different charts which require different values for retrieving results.
	// To get the value of a specific chart you can use our Google Play Apps Store API JSON output: chart_options[index].value (e.g. chart=topselling_free).
	// It shouldn't be used with the section_page_token, see_more_token, and next_page_token parameters
	Chart string `json:"chart,omitempty"`
	// Parameter defines the see more token used for retrieving the pagination results from individual sections It is usually found in next page results.
	// It shouldn't be used with the section_page_token, next_page_token, and chart, parameters
	SeeMoreToken string `json:"see_more_token,omitempty"`
}

func (p Pagination) FieldValidation() error {
	count := 0
	if p.NextPageToken != "" {
		count++
	}
	if p.SectionPageToken != "" {
		count++
	}
	//if p.Chart != "" {
	//	count++
	//}
	if p.SeeMoreToken != "" {
		count++
	}

	if count > 1 {
		return errors.New("only one of NextPageToken, SectionPageToken, or SeeMoreToken should be set")
		//return errors.New("only one of NextPageToken, SectionPageToken, Chart, or SeeMoreToken should be set")
	}
	return nil
}

type Response struct {
	SearchMetadata    *SearchMetadata           `json:"search_metadata,omitempty"`
	SearchParameters  *SearchParameters         `json:"search_parameters,omitempty"`
	ChartOption       []ChartOption             `json:"chart_option,omitempty"`
	HighlightItem     [][]HighlightItem         `json:"highlight_item,omitempty"`
	OrganicResults    []OrganicResults          `json:"organic_results,omitempty"`
	SerpapiPagination *SerpapiSectionPagination `json:"serpapi_pagination,omitempty"`
	// product场景信息
	ProductInfo  *ProductInfo  `json:"product_info,omitempty"`
	Media        *Media        `json:"media,omitempty"`
	AboutThisApp *AboutThisApp `json:"about_this_app,omitempty"`
	Badges       []Badges      `json:"badges,omitempty"`
	Categories   []Categories  `json:"categories,omitempty"`
	// 格式为 Mar 7, 2025
	UpdatedOn        string            `json:"updated_on,omitempty"`
	DataSafety       []DataSafety      `json:"data_safety,omitempty"`
	WhatIsNew        *WhatIsNew        `json:"what_s_new,omitempty"`
	Ratings          []Ratings         `json:"ratings,omitempty"`
	Reviews          []Review          `json:"reviews,omitempty"`
	DeveloperContact *DeveloperContact `json:"developer_contact,omitempty"`
	SimilarResults   []SimilarResults  `json:"similar_results,omitempty"`
}

type HighlightItem struct {
	Title       string `json:"title,omitempty"`
	Subtitle    string `json:"subtitle,omitempty"`
	Link        string `json:"link,omitempty"`
	ProductID   string `json:"product_id,omitempty"`
	SerpapiLink string `json:"serpapi_link,omitempty"`
	Thumbnail   string `json:"thumbnail,omitempty"`
}

type OrganicResults struct {
	Title                    string                    `json:"title,omitempty"`
	SerpapiSectionPagination *SerpapiSectionPagination `json:"serpapi_section_pagination,omitempty"`
	Item                     []Item                    `json:"item,omitempty"`
}

func (o OrganicResults) IsEmpty() bool {
	return o.Title == "" && o.SerpapiSectionPagination == nil && len(o.Item) == 0
}

type SearchMetadata struct {
	ID             string  `json:"id,omitempty"`
	Status         string  `json:"status,omitempty"`
	JSONEndpoint   string  `json:"json_endpoint,omitempty"`
	CreatedAt      string  `json:"created_at,omitempty"`
	ProcessedAt    string  `json:"processed_at,omitempty"`
	GooglePlayURL  string  `json:"google_play_url,omitempty"`
	HtmlFile       string  `json:"raw_html,omitempty"`
	TotalTimeTaken float64 `json:"total_time_taken,omitempty"`
}

func (m SearchMetadata) IsEmpty() bool {
	return m.GooglePlayURL == ""
}

type SearchParameters struct {
	Q                string `json:"q,omitempty"`
	ProductID        string `json:"product_id,omitempty"`
	Actor            string `json:"actor,omitempty"`
	Hl               string `json:"hl,omitempty"`
	Gl               string `json:"gl,omitempty"`
	StoreDevice      string `json:"store_device,omitempty"`
	Store            string `json:"store,omitempty"`
	AppCategory      string `json:"app_category,omitempty"`
	SectionPageToken string `json:"section_page_token,omitempty"`
	Platform         string `json:"platform,omitempty"`
	SortBy           string `json:"sort_by,omitempty"`
}

func (s SearchParameters) isEmpty() bool {
	return s.Q == "" && s.Actor == "" && s.Hl == ""
}

type ChartOption struct {
	Text  string `json:"text,omitempty"`
	Value string `json:"value,omitempty"`
}

type SerpapiSectionPagination struct {
	NextPageToken    string `json:"next_page_token,omitempty"`
	Next             string `json:"next,omitempty"`
	SectionPageToken string `json:"section_page_token,omitempty"`
}

func (p SerpapiSectionPagination) isEmpty() bool {
	return p.NextPageToken == "" && p.SectionPageToken == ""
}

type Item struct {
	Title       string   `json:"title,omitempty"`
	Link        string   `json:"link,omitempty"`
	ProductID   string   `json:"product_id,omitempty"`
	SerpapiLink string   `json:"serpapi_link,omitempty"`
	Rating      float64  `json:"rating,omitempty"`
	Author      string   `json:"author,omitempty"`
	Category    string   `json:"category,omitempty"`
	Downloads   string   `json:"downloads,omitempty"`
	Thumbnail   string   `json:"thumbnail,omitempty"`
	Description string   `json:"description,omitempty"`
	Video       string   `json:"video,omitempty"`
	Extension   []string `json:"extension,omitempty"`
}

func (i Item) isEmpty() bool {
	return i.Title == "" || i.Link == ""
}

type ProductInfo struct {
	Title         string        `json:"title,omitempty"`
	Authors       []Author      `json:"authors,omitempty"`
	Extensions    []string      `json:"extensions,omitempty"`
	Rating        float32       `json:"rating,omitempty"`
	Reviews       int64         `json:"reviews,omitempty"`
	ContentRating ContentRating `json:"content_rating,omitempty"`
	Downloads     string        `json:"downloads,omitempty"`
	Thumbnail     string        `json:"thumbnail,omitempty"`
	Offers        []Offer       `json:"offers,omitempty"`
}

type Author struct {
	Name string `json:"name,omitempty"`
	Link string `json:"link,omitempty"`
}

type ContentRating struct {
	Text      string `json:"text,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"`
}

type Offer struct {
	Text string `json:"text,omitempty"`
	Link string `json:"link,omitempty"`
}

type Media struct {
	Video  Video    `json:"video,omitempty"`
	Images []string `json:"images,omitempty"`
}

type Video struct {
	Thumbnail string `json:"thumbnail,omitempty"`
	Link      string `json:"link,omitempty"`
}

type AboutThisApp struct {
	Snippet             string       `json:"snippet,omitempty"`
	InAppPurchases      string       `json:"in_app_purchases,omitempty"`
	ReleasedOn          string       `json:"released_on,omitempty"`
	UpdatedOn           string       `json:"updated_on,omitempty"`
	Downloads           string       `json:"downloads,omitempty"`
	ContentRating       string       `json:"content_rating,omitempty"`
	InteractiveElements string       `json:"interactive_elements,omitempty"`
	OfferedBy           string       `json:"offered_by,omitempty"`
	Permissions         []Permission `json:"permissions,omitempty"`
}

type Permission struct {
	Type    string   `json:"type,omitempty"`
	Details []string `json:"details,omitempty"`
}

type Badges struct {
	Name string `json:"name,omitempty"`
}

type Categories struct {
	Name        string `json:"name,omitempty"`
	Link        string `json:"link,omitempty"`
	CategoryID  string `json:"category_id,omitempty"`
	SerpapiLink string `json:"serpapi_link,omitempty"`
}

type DataSafety struct {
	Text    string `json:"text,omitempty"`
	Subtext string `json:"subtext,omitempty"`
	Link    string `json:"link,omitempty"`
}

type WhatIsNew struct {
	Snippet string `json:"snippet,omitempty"`
}

type Ratings struct {
	Stars int   `json:"stars,omitempty"`
	Count int64 `json:"count,omitempty"`
}

type Review struct {
	ID      string  `json:"id,omitempty"`
	Title   string  `json:"title,omitempty"`
	Avatar  string  `json:"avatar,omitempty"`
	Rating  float32 `json:"rating,omitempty"`
	Snippet string  `json:"snippet,omitempty"`
	Likes   int     `json:"likes,omitempty"`
	Date    string  `json:"date,omitempty"`
	ISODate string  `json:"iso_date,omitempty"`
}

type SerpAPIPagination struct {
	Next          string `json:"next,omitempty"`
	NextPageToken string `json:"next_page_token,omitempty"`
}

type DeveloperContact struct {
	SupportEmail string `json:"support_email,omitempty"`
}

type SimilarResults struct {
	Title        string `json:"title,omitempty"`
	SeeMoreLink  string `json:"see_more_link,omitempty"`
	SeeMoreToken string `json:"see_more_token,omitempty"`
	SerpapiLink  string `json:"serpapi_link,omitempty"`
	Items        []Item `json:"items,omitempty"`
}

type BooksAndReference struct {
	Title                    string                    `json:"title,omitempty"`
	SerpapiSectionPagination *SerpapiSectionPagination `json:"serpapi_section_pagination,omitempty"`
	Items                    []Item                    `json:"items,omitempty"`
}

/*
  	>======================================<
	以下为通过请求接口返回数据时，请求的相应struct
*/

// QueryParams 查询参数
type QueryParams struct {
	RPCIDs      []string `json:"rpcids,omitempty"`
	SourcePath  string   `json:"source-path,omitempty"`
	Fsid        string   `json:"f.sid,omitempty"`
	BL          string   `json:"bl,omitempty"`
	HL          string   `json:"hl,omitempty"`
	GL          string   `json:"gl,omitempty"`
	Authuser    string   `json:"authuser,omitempty"`
	SocApp      string   `json:"soc-app,omitempty"`
	SocPlatform string   `json:"soc-platform,omitempty"`
	SocDevice   string   `json:"soc-device,omitempty"`
	ReqID       string   `json:"_reqid,omitempty"`
	RT          string   `json:"rt,omitempty"`
}

func (q *QueryParams) FieldValidation() error {
	if q.HL == "" {
		q.HL = "en"
	}
	if q.GL == "" {
		q.GL = "us"
	}
	if q.BL == "" {
		q.BL = "boq_playuiserver_20250615.18_p0"
	}
	if q.SourcePath == "" {
		return errors.New("`store_device` parameters can not be empty")
	}
	if q.Fsid == "" {
		q.Fsid = "3822332858809487112"
	}
	if q.SocApp == "" {
		q.SocApp = "121"
	}
	if q.SocPlatform == "" {
		q.SocPlatform = "1"
	}
	if q.SocDevice == "" {
		q.SocDevice = "1"
	}
	if q.RT == "" {
		q.RT = "c"
	}
	return nil
}
