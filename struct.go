package main

type RequestParam struct {
	Q            string `url:"q" json:"q"`
	Location     string `url:"location" json:"location"`
	ComeBack     string `url:"come_back" json:"come_back"`
	GoogleDomain string `url:"google_domain" json:"google_domain"`
	Gl           string `url:"gl" json:"gl"`
	Hl           string `url:"hl" json:"hl"`
	Tbs          string `url:"tbs" json:"tbs"`
	Shoprs       string `url:"shoprs" json:"shoprs"`
	DirectLink   string `url:"direct_link" json:"direct_link"`
	Start        string `url:"start" json:"start"`
	Engine       string `url:"engine" json:"engine"`

	IsReturnRawHtml string `json:"is_return_raw_html" json:"is_return_raw_html"` // "1" 为返回原始html
}

type Response struct {
	Filters                    []Filters                    `json:"filters,omitempty"`
	InlineShoppingResults      []InlineShoppingResult       `json:"inline_shopping_results,omitempty"`
	ShoppingResults            []ShoppingResult             `json:"shopping_results,omitempty"`
	RelatedShoppingResults     []RelatedShoppingResult      `json:"related_shopping_results,omitempty"`
	Categories                 []Category                   `json:"categories,omitempty"`
	FeaturedShoppingResults    []FeaturedShoppingResult     `json:"featured_shopping_results,omitempty"`
	Pagination                 *Pagination                  `json:"pagination,omitempty"`
	CategorizedShoppingResults []CategorizedShoppingResults `json:"categorized_shopping_results,omitempty"`
	SearchMetadata             *SearchMetadata              `json:"search_metadata,omitempty"`
}

type SearchMetadata struct {
	TaskId       string `json:"task_id,omitempty"`
	RawHtml      string `json:"raw_html,omitempty"`
	PrettifyHtml string `json:"prettify_html,omitempty"`
}

func (s *SearchMetadata) IsEmpty() bool {
	return s.TaskId == "" && s.RawHtml == "" && s.PrettifyHtml == ""
}

type CategorizedShoppingResults struct {
	Title          string           `json:"title,omitempty"`
	ShoppingResult []ShoppingResult `json:"shopping_result,omitempty"`
}

func (c *CategorizedShoppingResults) isEmpty() bool {
	return c.Title == "" && len(c.ShoppingResult) == 0
}

type Filters struct {
	Type      string    `json:"type,omitempty"`
	InputType string    `json:"input_type,omitempty"`
	Options   []Options `json:"option,omitempty"`
}

func (f *Filters) IsEmpty() bool {
	return f.Type == "" && f.InputType == "" && len(f.Options) == 0
}

type Options struct {
	Text string `json:"text,omitempty"`
	Tbs  string `json:"tbs,omitempty"`
	Link string `json:"link,omitempty"`
}

func (o *Options) IsEmpty() bool {
	return o.Text == "" && o.Tbs == "" && o.Link == ""
}

type InlineShoppingResult struct {
	Position       int      `json:"position,omitempty"`
	BlockPosition  string   `json:"block_position,omitempty"`
	Title          string   `json:"title,omitempty"`
	Price          string   `json:"price,omitempty"`
	ExtractedPrice float64  `json:"extracted_price,omitempty"`
	Link           string   `json:"link,omitempty"`
	Source         string   `json:"source,omitempty"`
	Thumbnail      string   `json:"thumbnail,omitempty"`
	Extensions     []string `json:"extensions,omitempty"`
}

func (i *InlineShoppingResult) IsEmpty() bool {
	return i.Title == "" && i.Price == "" && i.ExtractedPrice == 0 && i.Link == "" && i.Source == "" && i.Thumbnail == "" && len(i.Extensions) == 0
}

type ShoppingResult struct {
	Position            int      `json:"position,omitempty"`
	Title               string   `json:"title,omitempty"`
	Link                string   `json:"link,omitempty"`
	ProductLink         string   `json:"product_link,omitempty"`
	ProductID           string   `json:"product_id,omitempty"`
	Source              string   `json:"source,omitempty"`
	SourceIcon          string   `json:"source_icon,omitempty"`
	Price               string   `json:"price,omitempty"`
	ExtractedPrice      float64  `json:"extracted_price,omitempty"`
	OldPrice            string   `json:"old_price,omitempty"`
	ExtractedOldPrice   float64  `json:"extracted_old_price,omitempty"`
	SecondHandCondition string   `json:"second_hand_condition,omitempty"`
	Rating              float64  `json:"rating,omitempty"`
	Reviews             int      `json:"reviews,omitempty"`
	Snippet             string   `json:"snippet,omitempty"`
	Extensions          []string `json:"extensions,omitempty"`
	Badge               string   `json:"badge,omitempty"`
	Thumbnail           string   `json:"thumbnail,omitempty"`
	Thumbnails          []string `json:"thumbnails,omitempty"`
	Tag                 string   `json:"tag,omitempty"`
	Delivery            string   `json:"delivery,omitempty"`
	StoreRating         float64  `json:"store_rating,omitempty"`
	StoreReviews        int      `json:"store_reviews,omitempty"`
}

func (s ShoppingResult) IsEmpty() bool {
	return s.Position == 0 &&
		s.Title == "" &&
		s.Link == "" &&
		s.ProductLink == "" &&
		s.ProductID == "" &&
		s.Source == "" &&
		s.SourceIcon == "" &&
		s.Price == "" &&
		s.ExtractedPrice == 0 &&
		s.OldPrice == "" &&
		s.ExtractedOldPrice == 0 &&
		s.Rating == 0 &&
		s.Reviews == 0 &&
		len(s.Extensions) == 0 &&
		s.Badge == "" &&
		s.Thumbnail == "" &&
		len(s.Thumbnails) == 0 &&
		s.Tag == "" &&
		s.Delivery == "" &&
		s.StoreRating == 0 &&
		s.StoreReviews == 0
}

type RelatedShoppingResult struct {
	Position          int     `json:"position,omitempty"`
	Title             string  `json:"title,omitempty"`
	Link              string  `json:"link,omitempty"`
	ProductLink       string  `json:"product_link,omitempty"`
	ProductID         string  `json:"product_id,omitempty"`
	SerpapiProductAPI string  `json:"serpapi_product_api,omitempty"`
	Source            string  `json:"source,omitempty"`
	SourceIcon        string  `json:"source_icon,omitempty"`
	Price             string  `json:"price,omitempty"`
	ExtractedPrice    float64 `json:"extracted_price,omitempty"`
	AlternativePrice  struct {
		Price          string  `json:"price,omitempty"`
		ExtractedPrice float64 `json:"extracted_price,omitempty"`
		Currency       string  `json:"currency,omitempty"`
	} `json:"alternative_price,omitempty"`
	OldPrice          string   `json:"old_price,omitempty"`
	ExtractedOldPrice float64  `json:"extracted_old_price,omitempty"`
	Rating            float64  `json:"rating,omitempty"`
	Reviews           int      `json:"reviews,omitempty"`
	Thumbnail         string   `json:"thumbnail,omitempty"`
	Tag               string   `json:"tag,omitempty"`
	Extensions        []string `json:"extensions,omitempty"`
	Delivery          string   `json:"delivery,omitempty"`
}

type Category struct {
	Title   string `json:"title,omitempty"`
	Filters []struct {
		Title       string `json:"title,omitempty"`
		Thumbnail   string `json:"thumbnail,omitempty"`
		Link        string `json:"link,omitempty"`
		SerpapiLink string `json:"serpapi_link,omitempty"`
	} `json:"filters,omitempty"`
}

type FeaturedShoppingResult struct {
	Position            int      `json:"position,omitempty"`
	Title               string   `json:"title,omitempty"`
	Link                string   `json:"link,omitempty"`
	ProductLink         string   `json:"product_link,omitempty"`
	ProductID           string   `json:"product_id,omitempty"`
	SerpapiProductAPI   string   `json:"serpapi_product_api,omitempty"`
	Source              string   `json:"source,omitempty"`
	SourceIcon          string   `json:"source_icon,omitempty"`
	Price               string   `json:"price,omitempty"`
	ExtractedPrice      float64  `json:"extracted_price,omitempty"`
	OldPrice            string   `json:"old_price,omitempty"`
	ExtractedOldPrice   float64  `json:"extracted_old_price,omitempty"`
	SecondHandCondition string   `json:"second_hand_condition,omitempty"`
	Rating              float64  `json:"rating,omitempty"`
	Reviews             int      `json:"reviews,omitempty"`
	Extensions          []string `json:"extensions,omitempty"`
	Thumbnail           string   `json:"thumbnail,omitempty"`
	Tag                 string   `json:"tag,omitempty"`
	Delivery            string   `json:"delivery,omitempty"`
	StoreRating         float64  `json:"store_rating,omitempty"`
	StoreReviews        int      `json:"store_reviews,omitempty"`
	Remark              string   `json:"remark,omitempty"`
}

type Pagination struct {
	Current    int               `json:"current,omitempty"`
	Next       string            `json:"next,omitempty"`
	OtherPages map[string]string `json:"other_pages,omitempty"`
}
