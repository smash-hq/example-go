package play

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
)

var (
	// StoreCategories Parameter defines the apps store category
	gamesCategoryEnum = []string{"GAME", "GAME_ACTION", "GAME_ADVENTURE", "GAME_ARCADE", "GAME_BOARD", "GAME_CARD", "GAME_CASINO", "GAME_CASUAL",
		"GAME_EDUCATIONAL", "GAME_MUSIC", "GAME_PUZZLE", "GAME_RACING", "GAME_ROLE_PLAYING", "GAME_SIMULATION", "GAME_SPORTS", "GAME_STRATEGY", "GAME_TRIVIA", "GAME_WORD"}

	bookCategoryEnums = []string{"coll_1665", "subj_Art___Humor.AH_Art", "subj_Art___Humor.AH_Drama", "subj_Art___Humor.AH_Humor",
		"subj_Art___Humor.AH_Music", "subj_Art___Humor.AH_Performing_Arts", "coll_1204", "subj_Biography___Autobiography.General",
		"subj_Biography___Autobiography.Adventurers___Explorers", "subj_Biography___Autobiography.Business",
		"subj_Biography___Autobiography.Composers___Musicians", "subj_Biography___Autobiography.Criminals___Outlaws",
		"subj_Biography___Autobiography.Cultural_Heritage", "subj_Biography___Autobiography.Entertainment___Performing_Arts",
		"subj_Biography___Autobiography.Historical", "subj_Biography___Autobiography.Literary", "subj_Biography___Autobiography.Medical",
		"subj_Biography___Autobiography.Military", "subj_Biography___Autobiography.Personal_Memoirs", "subj_Biography___Autobiography.Political",
		"subj_Biography___Autobiography.Presidents___Heads_of_State", "subj_Biography___Autobiography.Religious",
		"subj_Biography___Autobiography.Rich___Famous", "subj_Biography___Autobiography.Royalty", "subj_Biography___Autobiography.Science___Technology",
		"subj_Biography___Autobiography.Sports", "subj_Biography___Autobiography.Women", "coll_1668", "coll_1205", "subj_Business___Economics.General",
		"subj_Business___Economics.Accounting", "subj_Business___Economics.Advertising___Promotion", "subj_Business___Economics.Business_Communication",
		"subj_Business___Economics.Careers", "subj_Business___Economics.Decision-Making___Problem_Solving", "subj_Business___Economics.Development",
		"subj_Business___Economics.E-Commerce", "subj_Business___Economics.Economic_History", "subj_Business___Economics.Economics",
		"subj_Business___Economics.Entrepreneurship", "subj_Business___Economics.Finance", "subj_Business___Economics.Human_Resources___Personnel_Management",
		"subj_Business___Economics.Industries", "subj_Business___Economics.International", "subj_Business___Economics.Investments___Securities",
		"subj_Business___Economics.Leadership", "subj_Business___Economics.Management", "subj_Business___Economics.Management_Science",
		"subj_Business___Economics.Marketing", "subj_Business___Economics.Motivational", "subj_Business___Economics.Organizational_Behavior",
		"subj_Business___Economics.Personal_Finance", "subj_Business___Economics.Real_Estate", "subj_Business___Economics.Sales___Selling",
		"subj_Business___Economics.Skills", "subj_Business___Economics.Small_Business", "subj_Business___Economics.Strategic_Planning",
		"subj_Business___Economics.Time_Management", "subj_Business___Economics.Workplace_Culture", "coll_1207", "subj_Computers.General",
		"subj_Computers.Desktop_Applications", "subj_Computers.Electronic_Publishing", "subj_Computers.Enterprise_Applications", "subj_Computers.Hardware",
		"subj_Computers.Internet", "subj_Computers.Networking", "subj_Computers.Operating_Systems", "subj_Computers.Programming", "subj_Computers.Programming_Languages",
		"subj_Computers.Security", "subj_Computers.Software_Development___Engineering", "subj_Computers.Systems_Architecture", "subj_Computers.Web",
		"coll_1208", "subj_Cooking.General", "subj_Cooking.Beverages", "subj_Cooking.Courses___Dishes", "subj_Cooking.Health___Healing", "subj_Cooking.Methods",
		"subj_Cooking.Regional___Ethnic", "subj_Cooking.Specific_Ingredients", "subj_Cooking.Tablesetting", "subj_Cooking.Vegetarian___Vegan", "coll_1670", "subj_Education.E_General",
		"subj_Education.English_as_a_Second_Language", "subj_Education.Higher", "subj_Education.Teaching_Methods___Materials", "coll_1261", "subj_Technology___Engineering.General",
		"subj_Technology___Engineering.Agriculture", "subj_Technology___Engineering.Electronics", "subj_Technology___Engineering.Industrial_Technology",
		"subj_Technology___Engineering.Military_Science", "subj_Technology___Engineering.Mobile___Wireless_Communications", "subj_Technology___Engineering.Nanotechnology___MEMS",
		"subj_Technology___Engineering.Remote_Sensing___Geographic_Information_Systems", "coll_1200", "subj_Fiction.General", "subj_Fiction.Action___Adventure",
		"subj_Fiction.African_American", "subj_Fiction.Alternative_History", "subj_Fiction.Anthologies__multiple_authors_", "subj_Fiction.Biographical", "subj_Fiction.Christian",
		"subj_Fiction.Classics", "subj_Fiction.Coming_of_Age", "subj_Fiction.Contemporary_Women", "subj_Fiction.Crime", "subj_Fiction.Cultural_Heritage", "subj_Fiction.Erotica",
		"subj_Fiction.Fairy_Tales__Folk_Tales__Legends___Mythology", "subj_Fiction.Family_Life", "subj_Fiction.Fantasy", "subj_Fiction.Gay", "subj_Fiction.Ghost",
		"subj_Fiction.Historical", "subj_Fiction.Horror", "subj_Fiction.Humorous", "subj_Fiction.Legal", "subj_Fiction.Lesbian", "subj_Fiction.Literary", "subj_Fiction.Media_Tie-In",
		"subj_Fiction.Medical", "subj_Fiction.Men_s_Adventure", "subj_Fiction.Mystery___Detective", "subj_Fiction.Occult___Supernatural", "subj_Fiction.Political",
		"subj_Fiction.Psychological", "subj_Fiction.Religious", "subj_Fiction.Romance", "subj_Fiction.Sagas", "subj_Fiction.Satire", "subj_Fiction.Science_Fiction",
		"subj_Fiction.Sea_Stories", "subj_Fiction.Short_Stories__single_author_", "subj_Fiction.Sports", "subj_Fiction.Thrillers", "subj_Fiction.War___Military",
		"subj_Fiction.Westerns", "coll_1605", "subj_Health__Mind___Body.HMB_General", "subj_Health__Mind___Body.Alternative_Therapies", "subj_Health__Mind___Body.Astrology",
		"subj_Health__Mind___Body.Beauty___Grooming", "subj_Health__Mind___Body.Creativity", "subj_Health__Mind___Body.Death__Grief__Bereavement",
		"subj_Health__Mind___Body.Exercise", "subj_Health__Mind___Body.Healing", "subj_Health__Mind___Body.Healthy_Living", "subj_Health__Mind___Body.Motivational___Inspirational",
		"subj_Health__Mind___Body.Personal_Growth", "subj_Health__Mind___Body.Spirituality", "subj_Health__Mind___Body.Yoga", "coll_1209", "subj_History.General", "subj_History.Africa",
		"subj_History.Americas__North__Central__South__West_Indies_", "subj_History.Ancient", "subj_History.Asia", "subj_History.Europe", "subj_History.Holocaust",
		"subj_History.Latin_America", "subj_History.Medieval", "subj_History.Middle_East", "subj_History.Military", "subj_History.Modern", "subj_History.Native_American",
		"subj_History.Revolutionary", "subj_History.Social_History", "subj_History.United_States", "subj_History.World", "coll_1211", "subj_House___Home.HH_General",
		"subj_House___Home.Crocheting", "subj_House___Home.Do-It-Yourself", "subj_House___Home.Fruit", "subj_House___Home.Furniture", "subj_House___Home.Knitting",
		"subj_House___Home.Models", "subj_House___Home.Quilts___Quilting", "subj_House___Home.Sustainable_Living", "subj_House___Home.Woodwork", "coll_1667", "subj_Law.General",
		"subj_Law.Criminal_Law", "subj_Law.Emigration___Immigration", "coll_1669", "subj_Medical.General", "subj_Medical.Cardiology", "subj_Medical.Caregiving",
		"subj_Medical.Family___General_Practice", "subj_Medical.Neuroscience", "subj_Medical.Nursing", "subj_Medical.Practice_Management___Reimbursement",
		"subj_Medical.Preventive_Medicine", "coll_1615", "subj_Fiction.Mystery___Detective.General", "subj_Fiction.Mystery___Detective.Hard-Boiled",
		"subj_Fiction.Mystery___Detective.Historical", "subj_Fiction.Mystery___Detective.Police_Procedural", "subj_Fiction.Mystery___Detective.Suspense",
		"subj_Fiction.Mystery___Detective.Traditional_British", "subj_Fiction.Mystery___Detective.True_Crime", "subj_Fiction.Mystery___Detective.Women_Sleuths",
		"coll_1214", "subj_Family___Relationships.General", "subj_Family___Relationships.Abuse", "subj_Family___Relationships.Children_with_Special_Needs",
		"subj_Family___Relationships.Conflict_Resolution", "subj_Family___Relationships.Divorce___Separation", "subj_Family___Relationships.Family_Relationships",
		"subj_Family___Relationships.Interpersonal_Relations", "subj_Family___Relationships.Life_Stages", "subj_Family___Relationships.Love___Romance",
		"subj_Family___Relationships.Marriage", "subj_Family___Relationships.Parenting", "coll_1215", "subj_Political_Science.General", "subj_Political_Science.Essays",
		"subj_Political_Science.History___Theory", "subj_Political_Science.International_Relations", "subj_Political_Science.Political_Ideologies",
		"subj_Political_Science.Political_Process", "subj_Political_Science.Public_Affairs___Administration", "subj_Political_Science.Public_Policy", "coll_1217",
		"subj_Religion.General", "subj_Religion.Biblical_Commentary", "subj_Religion.Biblical_Criticism___Interpretation", "subj_Religion.Biblical_Reference",
		"subj_Religion.Biblical_Studies", "subj_Religion.Buddhism", "subj_Religion.Christian_Church", "subj_Religion.Christian_Education", "subj_Religion.Christian_Life",
		"subj_Religion.Christian_Ministry", "subj_Religion.Christian_Rituals___Practice", "subj_Religion.Christian_Theology", "subj_Religion.Christianity",
		"subj_Religion.Comparative_Religion", "subj_Religion.Counseling", "subj_Religion.Eastern", "subj_Religion.Hinduism", "subj_Religion.History", "subj_Religion.Inspirational",
		"subj_Religion.Judaism", "subj_Religion.Philosophy", "subj_Religion.Prayer", "subj_Religion.Sikhism", "subj_Religion.Spirituality", "subj_Religion.Theology",
		"coll_1218", "subj_Fiction.Romance.General", "subj_Fiction.Romance.Contemporary", "subj_Fiction.Romance.Fantasy", "subj_Fiction.Romance.Historical",
		"subj_Fiction.Romance.Paranormal", "subj_Fiction.Romance.Suspense", "subj_Fiction.Romance.Time_Travel", "subj_Fiction.Romance.Western", "coll_1671", "coll_1604",
		"subj_Science_Fiction___Fantasy.SFF_General", "subj_Science_Fiction___Fantasy.Contemporary", "subj_Science_Fiction___Fantasy.Epic",
		"subj_Science_Fiction___Fantasy.SFF_Historical", "subj_Science_Fiction___Fantasy.Military", "subj_Science_Fiction___Fantasy.Paranormal",
		"subj_Science_Fiction___Fantasy.Space_Opera", "coll_1222", "subj_Sports___Recreation.General", "subj_Sports___Recreation.Baseball",
		"subj_Sports___Recreation.Basketball", "subj_Sports___Recreation.Bodybuilding___Weight_Training", "subj_Sports___Recreation.Coaching",
		"subj_Sports___Recreation.Essays", "subj_Sports___Recreation.Extreme_Sports", "subj_Sports___Recreation.Fishing", "subj_Sports___Recreation.Football",
		"subj_Sports___Recreation.Golf", "subj_Sports___Recreation.Hiking", "subj_Sports___Recreation.History", "coll_1673", "coll_1224", "subj_Travel.General",
		"subj_Travel.Essays___Travelogues", "subj_Travel.Europe", "subj_Travel.Special_Interest", "subj_Travel.United_States", "coll_1672", "audiobooks",
		"audiobooks_coll_1665", "audiobooks_subj_Arts___Entertainment.Celebrity_Bios", "audiobooks_subj_Arts___Entertainment.Humor",
		"audiobooks_subj_Arts___Entertainment.TV___Film", "audiobooks_coll_1204", "audiobooks_subj_Biographies___Memoirs.Artists__Writers____Musicians",
		"audiobooks_subj_Biographies___Memoirs.Business_Leaders", "audiobooks_subj_Biographies___Memoirs.Celebrities", "audiobooks_subj_Biographies___Memoirs.Criminal",
		"audiobooks_subj_Biographies___Memoirs.Personal_Memoirs", "audiobooks_subj_Biographies___Memoirs.Political_Figures", "audiobooks_subj_Biographies___Memoirs.Religious_Figures",
		"audiobooks_subj_Biographies___Memoirs.Science___Technology_Leaders", "audiobooks_coll_1205", "audiobooks_subj_Business___Investing.Career_Skills",
		"audiobooks_subj_Business___Investing.Commerce___Economy", "audiobooks_subj_Business___Investing.Leadership", "audiobooks_subj_Business___Investing.Management",
		"audiobooks_subj_Business___Investing.Marketing", "audiobooks_subj_Business___Investing.Personal_Finance___Investing", "audiobooks_subj_Business___Investing.Sales",
		"audiobooks_coll_1689", "audiobooks_coll_1200", "audiobooks_subj_Fiction___Literature.African_American", "audiobooks_subj_Fiction___Literature.Classics",
		"audiobooks_subj_Fiction___Literature.Historical", "audiobooks_subj_Fiction___Literature.Horror", "audiobooks_subj_Fiction___Literature.Humor",
		"audiobooks_subj_Fiction___Literature.Literary", "audiobooks_subj_Fiction___Literature.Media_Tie-In", "audiobooks_subj_Fiction___Literature.Religion___Spirituality",
		"audiobooks_subj_Fiction___Literature.Short_Stories___Anthologies", "audiobooks_subj_Fiction___Literature.Western", "audiobooks_coll_1605",
		"audiobooks_subj_Health__Mind___Body.Alternative_Therapies", "audiobooks_subj_Health__Mind___Body.Astrology", "audiobooks_subj_Health__Mind___Body.Beauty___Grooming",
		"audiobooks_subj_Health__Mind___Body.Creativity", "audiobooks_subj_Health__Mind___Body.Exercise", "audiobooks_subj_Health__Mind___Body.Healing",
		"audiobooks_subj_Health__Mind___Body.Healthy_Living", "audiobooks_subj_Health__Mind___Body.Spirituality", "audiobooks_coll_1209", "audiobooks_subj_History.20th_Century",
		"audiobooks_subj_History.21st_Century", "audiobooks_subj_History.American", "audiobooks_subj_History.Ancient", "audiobooks_subj_History.European",
		"audiobooks_subj_History.Kids___Young_Adults", "audiobooks_subj_History.Military", "audiobooks_subj_History.World", "audiobooks_coll_1676",
		"audiobooks_subj_Language_Instruction.French", "audiobooks_subj_Language_Instruction.German", "audiobooks_subj_Language_Instruction.Italian",
		"audiobooks_subj_Language_Instruction.Spanish", "audiobooks_coll_1615", "audiobooks_subj_Mystery___Thrillers.Cozy_Mysteries", "audiobooks_subj_Mystery___Thrillers.Espionage",
		"audiobooks_subj_Mystery___Thrillers.Historical", "audiobooks_subj_Mystery___Thrillers.International_Mystery___Crime", "audiobooks_subj_Mystery___Thrillers.Legal_Thrillers",
		"audiobooks_subj_Mystery___Thrillers.Medical_Thrillers", "audiobooks_subj_Mystery___Thrillers.Military", "audiobooks_subj_Mystery___Thrillers.Noir",
		"audiobooks_subj_Mystery___Thrillers.Police_Procedural", "audiobooks_subj_Mystery___Thrillers.Political", "audiobooks_subj_Mystery___Thrillers.Psychological",
		"audiobooks_subj_Mystery___Thrillers.Suspense", "audiobooks_subj_Mystery___Thrillers.True_Crime", "audiobooks_coll_1217",
		"audiobooks_subj_Religion___Spirituality.Buddhism___Eastern_Religions", "audiobooks_subj_Religion___Spirituality.Christianity", "audiobooks_subj_Religion___Spirituality.Islam",
		"audiobooks_subj_Religion___Spirituality.Judaism", "audiobooks_subj_Religion___Spirituality.Religious_Thought", "audiobooks_subj_Religion___Spirituality.Sermons___Ministries",
		"audiobooks_coll_1716", "audiobooks_subj_Romance_US.Christian_Romance", "audiobooks_subj_Romance_US.Contemporary", "audiobooks_subj_Romance_US.Erotica",
		"audiobooks_subj_Romance_US.Fantasy", "audiobooks_subj_Romance_US.Historical", "audiobooks_subj_Romance_US.LGBT", "audiobooks_subj_Romance_US.Paranormal",
		"audiobooks_subj_Romance_US.Regency", "audiobooks_subj_Romance_US.Romantic_Comedy", "audiobooks_subj_Romance_US.Science_Fiction___Fantasy", "audiobooks_subj_Romance_US.Sports",
		"audiobooks_subj_Romance_US.Suspense", "audiobooks_subj_Romance_US.Western", "audiobooks_subj_Romance_US.Young_Adult", "audiobooks_coll_1219",
		"audiobooks_subj_Science___Technology.Astronomy", "audiobooks_subj_Science___Technology.Biology", "audiobooks_subj_Science___Technology.Environment",
		"audiobooks_subj_Science___Technology.Physics", "audiobooks_coll_1604", "audiobooks_subj_Science_Fiction___Fantasy.Alternate_History",
		"audiobooks_subj_Science_Fiction___Fantasy.Contemporary", "audiobooks_subj_Science_Fiction___Fantasy.Dark_Fantasy", "audiobooks_subj_Science_Fiction___Fantasy.Fantasy__Epic",
		"audiobooks_subj_Science_Fiction___Fantasy.Military_Sci-Fi", "audiobooks_subj_Science_Fiction___Fantasy.Paranormal",
		"audiobooks_subj_Science_Fiction___Fantasy.Post_Apocalyptic", "audiobooks_subj_Science_Fiction___Fantasy.Short_Stories___Anthologies",
		"audiobooks_subj_Science_Fiction___Fantasy.Space_Opera", "audiobooks_subj_Science_Fiction___Fantasy.Steampunk", "audiobooks_subj_Science_Fiction___Fantasy.Superheroes",
		"audiobooks_subj_Science_Fiction___Fantasy.Time_Travel", "audiobooks_coll_1545", "audiobooks_subj_Self-Help.Communication_Skills", "audiobooks_subj_Self-Help.How-To",
		"audiobooks_subj_Self-Help.Hypnosis", "audiobooks_subj_Self-Help.Meditation", "audiobooks_subj_Self-Help.Motivation___Inspiration", "audiobooks_subj_Self-Help.Parenting",
		"audiobooks_subj_Self-Help.Personal_Finance", "audiobooks_subj_Self-Help.Relationships", "audiobooks_subj_Self-Help.Sexuality", "audiobooks_coll_1222",
		"audiobooks_subj_Sports.Baseball", "audiobooks_subj_Sports.Basketball", "audiobooks_subj_Sports.Football", "audiobooks_subj_Sports.Golf", "audiobooks_coll_1224",
		"audiobooks_subj_Travel.Adventure___Exploration", "audiobooks_subj_Travel.Essays___Travelogues", "audiobooks_subj_Travel.Guided_Tours", "audiobooks_coll_1672",
		"audiobooks_subj_Young_Adult.Family___Relationships", "audiobooks_subj_Young_Adult.Health___Sports", "audiobooks_subj_Young_Adult.History___Historical_Fiction",
		"audiobooks_subj_Young_Adult.Mystery___Thrillers", "audiobooks_subj_Young_Adult.Myths___Legends", "audiobooks_subj_Young_Adult.Religion___Spirituality",
		"audiobooks_subj_Young_Adult.Romance___Friendship", "audiobooks_subj_Young_Adult.Science_Fiction___Fantasy", "audiobooks_subj_Young_Adult.Social_Issues", "coll_1401",
		"subj_Comics___Graphic_Novels.General", "subj_Comics___Graphic_Novels.Crime___Mystery", "subj_Comics___Graphic_Novels.Fantasy", "subj_Comics___Graphic_Novels.Horror",
		"subj_Comics___Graphic_Novels.Literary", "subj_Comics___Graphic_Novels.Manga", "subj_Comics___Graphic_Novels.Media_Tie-In", "subj_Comics___Graphic_Novels.Science_Fiction",
		"subj_Comics___Graphic_Novels.Superheroes", "coll_1689", "coll_1690", "coll_1691", "coll_1693", "coll_1692", "coll_1694", "coll_1696", "coll_1698", "coll_1699", "coll_1702", "coll_1704", "coll_1705"}

	appsCategoryEnum = []string{"ART_AND_DESIGN", "AUTO_AND_VEHICLES", "BEAUTY", "BOOKS_AND_REFERENCE", "BUSINESS", "COMICS", "COMMUNICATION", "DATING", "EDUCATION",
		"ENTERTAINMENT", "EVENTS", "FINANCE", "FOOD_AND_DRINK", "HEALTH_AND_FITNESS", "HOUSE_AND_HOME", "LIBRARIES_AND_DEMO", "LIFESTYLE", "MAPS_AND_NAVIGATION",
		"MEDICAL", "MUSIC_AND_AUDIO", "NEWS_AND_MAGAZINES", "PARENTING", "PERSONALIZATION", "PHOTOGRAPHY", "PRODUCTIVITY", "SHOPPING", "SOCIAL", "SPORTS", "TOOLS",
		"TRAVEL_AND_LOCAL", "VIDEO_PLAYERS", "ANDROID_WEAR", "WATCH_FACE", "WEATHER", "FAMILY"}
	// ageRanges Parameter defines the device for sorting results.
	ageRanges = []string{"AGE_RANGE1", "AGE_RANGE2", "AGE_RANGE3"}
	// storeDevices Parameter defines age subcategory.
	storeDevices = []string{"phone", "tablet", "tv", "chromebook", "watch", "car", "windows", "movies", "family", "ebooks"}
)

// 使用在Query String Parameters中source-path参数常量，分别代表不同场景
const (
	AppsPath        = "/store/apps"
	StoreSearchPath = "/store/search"
	ProductPath     = "/store/apps/details"
	GamesPath       = "/store/games"
	MoviesPath      = "/store/movies"
	BooksPath       = "/store/books"
)

const (
	GooglePlay        = "scraper.google.play"
	GooglePlayProduct = "scraper.google.play.product"
	GooglePlayGames   = "scraper.google.play.games"
	GooglePlayMovies  = "scraper.google.play.movies"
	GooglePlayBooks   = "scraper.google.play.books"
)

// rpcids枚举
var (
	rpcidsEnum = map[string][]string{
		"default":       {"eIpeLd", "w3QCWb"},
		storeDevices[0]: {"eIpeLd", "w3QCWb"},
		storeDevices[1]: {"eIpeLd", "di6f4"},
		storeDevices[2]: {"eIpeLd", "di6f4"},
		storeDevices[3]: {"eIpeLd", "di6f4"},
		storeDevices[4]: {"eIpeLd", "di6f4"},
		storeDevices[5]: {"eIpeLd", "di6f4"},
		storeDevices[6]: {"eIpeLd", "di6f4", "w37aie"},
		storeDevices[7]: {"eIpeLd", "w3QCWb"},
		storeDevices[8]: {"eIpeLd", "w3QCWb"},
		storeDevices[9]: {"eIpeLd", "w3QCWb", "w37aie"},
		// 进行搜索时（即q不为空时）则使用
		"Query": {"AZO9Cb", "lGYRle"},
		// product场景搜索,								  badges、media				   similar game							    reviews
		//                                                about_this_app
		//    											  in_app_purchases
		//                                                released_on
		//                                                 updated_on
		//                                                 downloads                   content_rating
		//                                             interactive_elements
		//                                                 offered_by
		//                                                 permissions
		//												   categories
		//												   data_safety
		//												   what_s_new
		//"product": {"CLXjtf", "A6yuRe", "qjTkWb", "cBDeQe", "Ws7gDc", "ZittHe", "yowZ5", "ag2B9c", "aFXcAe", "e7uDs", "c5NYSc", "oCPfdb", "MTfLfb"},
		"product": {"CLXjtf", "A6yuRe", "Ws7gDc", "ZittHe", "yowZ5", "ag2B9c", "e7uDs", "Ws7gDc", "oCPfdb"},
	}
	engineEnum = map[string]int{
		GooglePlay:        2,
		GooglePlayProduct: 2,
		GooglePlayGames:   2,
		GooglePlayMovies:  1,
		GooglePlayBooks:   3,
	}
	engineSourceEnum = map[string]interface{}{
		GooglePlay:        "APPLICATION",
		GooglePlayProduct: "",
		GooglePlayGames:   "GAME",
		GooglePlayMovies:  "MOVIE",
		GooglePlayBooks:   nil,
	}
	productReviewPlatformEnum = map[string]int{
		storeDevices[0]: 2,
		storeDevices[1]: 3,
		storeDevices[2]: 6,
		storeDevices[3]: 5,
		storeDevices[4]: 4,
		storeDevices[5]: 7,
	}
)

// 生成对应的数据,mark2 eg: [1,5], 根据传入不同store device变化，car:[1,1], chromebook:[1,4], tv:[1,3], tablet:[1,5], phone:[1,1]
func generateMarkArray1(p *RequestParams, actor string) (mark1 []interface{}, mark2 []int) {
	var nextPageToken, chat, seeMoreToken, sectionPageToken string
	if p != nil {
		nextPageToken = p.NextPageToken
		chat = p.Chart
		seeMoreToken = p.SeeMoreToken
		sectionPageToken = p.SectionPageToken
	}

	q := p.Q
	arr := generalArr(nextPageToken, sectionPageToken, chat, seeMoreToken)
	if q != "" {
		var symbol = 4
		var priceArr []interface{}
		if actor == GooglePlayBooks {
			symbol = 2
			if p.Price != "" {
				atoi, err := strconv.Atoi(p.Price)
				if err != nil {
					atoi = 1
				}
				priceArr = append(priceArr, nil, atoi)
			}
		} else if actor == GooglePlayMovies {
			symbol = 1
		}
		frontArr := []interface{}{[]interface{}{}, arr, []interface{}{q}, symbol, []interface{}{nil, 1}, nil, nil, nil, []interface{}{1}}
		// 如果有价格相关参数则增加
		if len(priceArr) > 0 {
			frontArr = append(frontArr, priceArr)
		}
		return frontArr, []int{1}
	}
	storeDevice := p.StoreDevice
	switch actor {
	case GooglePlay:
		// 如果传入了category则直接返回
		if p.AppsCategory != "" {
			var category interface{}
			if p.AppsCategory != "" {
				category = p.AppsCategory
			} else {
				category = nil
			}
			return []interface{}{nil, 2, category, nil, arr, nil, 2}, []int{1, 1}
		}
		switch storeDevice {
		case storeDevices[0], "":
			return []interface{}{nil, 2, "APPLICATION", nil, arr, nil, 2}, []int{1, 1}
		case storeDevices[1]:
			return []interface{}{nil, arr, nil, 4, "APPLICATION"}, []int{1, 5}
		case storeDevices[2]:
			return []interface{}{nil, arr, nil, 1, "APPLICATION"}, []int{1, 3}
		case storeDevices[3]:
			return []interface{}{nil, arr, nil, 5, "APPLICATION"}, []int{1, 4}
		case storeDevices[4]:
			return []interface{}{nil, arr, nil, 2, "APPLICATION"}, []int{1, 1}
		case storeDevices[5]:
			return []interface{}{nil, arr, nil, 3, "APPLICATION"}, []int{1, 1}
		}
	case GooglePlayGames:
		// 如果传入了category则直接返回
		if p.GamesCategory != "" {
			var category interface{}
			if p.GamesCategory != "" {
				category = p.GamesCategory
			} else {
				category = "GAME"
			}
			return []interface{}{nil, 2, category, nil, arr, nil, 2}, []int{1, 2}
		}
		switch storeDevice {
		case storeDevices[0]:
			return []interface{}{nil, 2, "GAME", nil, arr, nil, 2}, []int{1, 2}
		case storeDevices[1]:
			return []interface{}{nil, arr, nil, 4, "GAME"}, []int{1, 6}
		case storeDevices[2]:
			return []interface{}{nil, arr, nil, 1, "GAME"}, []int{1, 3}
		case storeDevices[3]:
			return []interface{}{nil, arr, nil, 5, "GAME"}, []int{1, 4}
		case storeDevices[4]:
			return []interface{}{nil, arr, nil, 2, "GAME"}, []int{1, 7}
		case storeDevices[5], storeDevices[6], "":
			return []interface{}{nil, arr, 9}, []int{1, 2}
		}
	case GooglePlayBooks:
		var category interface{}
		if p.BooksCategory != "" {
			category = p.BooksCategory
		} else {
			category = nil
		}
		var ageSym = ageSymbol(p)
		return []interface{}{nil, 3, category, ageSym, arr, nil, 2}, []int{1}
	case GooglePlayMovies:
		category := p.MoviesCategory
		if category == "" {
			category = "MOVIE"
		}
		ageSym := ageSymbol(p)
		return []interface{}{nil, 1, category, ageSym, arr, nil, 2}, []int{1}
	}
	return nil, nil
}

func ageSymbol(p *RequestParams) interface{} {
	var ageSym interface{}
	switch p.Age {
	case ageRanges[0]:
		ageSym = 2
	case ageRanges[1]:
		ageSym = 3
	case ageRanges[2]:
		ageSym = 4
	default:
		ageSym = nil
	}
	return ageSym
}

// checks if a slice contains a specific string.
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func generalArr(nextPageToken, sectionPageToken, chat, seeMoreToken string) []interface{} {
	var arr []interface{}
	var pageArr = []int{20, 50}
	var outPageArr = []interface{}{8, pageArr}
	if nextPageToken != "" {
		outPageArr = append(outPageArr, nil, nextPageToken)
	}
	unknown1 := `[96,108,72,100,27,183,222,8,57,169,110,11,184,16,1,139,152,194,165,68,163,211,9,71,31,195,12,64,151,150,148,113,104,55,56,145,32,34,10,122]`
	var unknownArr1 []interface{}
	err := json.Unmarshal([]byte(unknown1), &unknownArr1)
	if err != nil {
		log.Error(err)
	}
	//unknown2 := `[null,null,[[[1,null,1],null,[[[]]],null,null,null,null,[null,2],null,null,null,null,null,null,null,null,null,null,null,null,null,null,[1]],[null,[[[]]],null,null,[1]],[null,[[[]]],null,[1]],[null,[[[]]]],null,null,null,null,[[[[]]]],[[[[]]]]],[[[[7,68],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,1],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,31],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,104],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,9],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,8],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,27],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,12],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,65],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,110],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,11],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,56],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,55],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,96],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,10],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,122],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,72],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,71],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,64],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,113],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,139],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,150],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,169],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,165],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,151],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,163],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,32],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,16],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,108],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,100],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,194],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,211],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,184],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,183],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[9,68],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,1],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,31],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,104],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,9],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,8],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,27],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,12],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,65],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,110],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,11],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,56],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,55],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,96],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,10],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,122],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,72],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,71],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,64],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,113],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,139],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,150],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,169],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,165],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,151],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,163],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,32],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,16],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,108],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,100],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,194],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,211],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,184],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,183],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[17,68],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,1],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,31],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,104],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,9],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,8],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,27],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,12],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,65],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,110],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,11],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,56],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,55],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,96],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,10],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,122],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,72],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,71],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,64],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,113],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,139],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,150],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,169],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,165],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,151],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,163],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,32],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,16],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,108],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,100],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,194],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,211],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,184],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,183],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[65,68],[[1,5,4,7,11,6]]],[[65,1],[[1,5,4,7,11,6]]],[[65,31],[[1,5,4,7,11,6]]],[[65,104],[[1,5,4,7,11,6]]],[[65,9],[[1,5,4,7,11,6]]],[[65,8],[[1,5,4,7,11,6]]],[[65,27],[[1,5,4,7,11,6]]],[[65,12],[[1,5,4,7,11,6]]],[[65,65],[[1,5,4,7,11,6]]],[[65,110],[[1,5,4,7,11,6]]],[[65,11],[[1,5,4,7,11,6]]],[[65,56],[[1,5,4,7,11,6]]],[[65,55],[[1,5,4,7,11,6]]],[[65,96],[[1,5,4,7,11,6]]],[[65,10],[[1,5,4,7,11,6]]],[[65,122],[[1,5,4,7,11,6]]],[[65,72],[[1,5,4,7,11,6]]],[[65,71],[[1,5,4,7,11,6]]],[[65,64],[[1,5,4,7,11,6]]],[[65,113],[[1,5,4,7,11,6]]],[[65,139],[[1,5,4,7,11,6]]],[[65,150],[[1,5,4,7,11,6]]],[[65,169],[[1,5,4,7,11,6]]],[[65,165],[[1,5,4,7,11,6]]],[[65,151],[[1,5,4,7,11,6]]],[[65,163],[[1,5,4,7,11,6]]],[[65,32],[[1,5,4,7,11,6]]],[[65,16],[[1,5,4,7,11,6]]],[[65,108],[[1,5,4,7,11,6]]],[[65,100],[[1,5,4,7,11,6]]],[[65,194],[[1,5,4,7,11,6]]],[[65,211],[[1,5,4,7,11,6]]],[[65,184],[[1,5,4,7,11,6]]],[[65,183],[[1,5,4,7,11,6]]],[[10,68],[[1,7,6,9,15,8]]],[[10,1],[[1,7,6,9,15,8]]],[[10,31],[[1,7,6,9,15,8]]],[[10,104],[[1,7,6,9,15,8]]],[[10,9],[[1,7,6,9,15,8]]],[[10,8],[[1,7,6,9,15,8]]],[[10,27],[[1,7,6,9,15,8]]],[[10,12],[[1,7,6,9,15,8]]],[[10,65],[[1,7,6,9,15,8]]],[[10,110],[[1,7,6,9,15,8]]],[[10,11],[[1,7,6,9,15,8]]],[[10,56],[[1,7,6,9,15,8]]],[[10,55],[[1,7,6,9,15,8]]],[[10,96],[[1,7,6,9,15,8]]],[[10,10],[[1,7,6,9,15,8]]],[[10,122],[[1,7,6,9,15,8]]],[[10,72],[[1,7,6,9,15,8]]],[[10,71],[[1,7,6,9,15,8]]],[[10,64],[[1,7,6,9,15,8]]],[[10,113],[[1,7,6,9,15,8]]],[[10,139],[[1,7,6,9,15,8]]],[[10,150],[[1,7,6,9,15,8]]],[[10,169],[[1,7,6,9,15,8]]],[[10,165],[[1,7,6,9,15,8]]],[[10,151],[[1,7,6,9,15,8]]],[[10,163],[[1,7,6,9,15,8]]],[[10,32],[[1,7,6,9,15,8]]],[[10,16],[[1,7,6,9,15,8]]],[[10,108],[[1,7,6,9,15,8]]],[[10,100],[[1,7,6,9,15,8]]],[[10,194],[[1,7,6,9,15,8]]],[[10,211],[[1,7,6,9,15,8]]],[[10,184],[[1,7,6,9,15,8]]],[[10,183],[[1,7,6,9,15,8]]],[[58,68],[[5,3,1,2,6,8]]],[[58,1],[[5,3,1,2,6,8]]],[[58,31],[[5,3,1,2,6,8]]],[[58,104],[[5,3,1,2,6,8]]],[[58,9],[[5,3,1,2,6,8]]],[[58,8],[[5,3,1,2,6,8]]],[[58,27],[[5,3,1,2,6,8]]],[[58,12],[[5,3,1,2,6,8]]],[[58,65],[[5,3,1,2,6,8]]],[[58,110],[[5,3,1,2,6,8]]],[[58,11],[[5,3,1,2,6,8]]],[[58,56],[[5,3,1,2,6,8]]],[[58,55],[[5,3,1,2,6,8]]],[[58,96],[[5,3,1,2,6,8]]],[[58,10],[[5,3,1,2,6,8]]],[[58,122],[[5,3,1,2,6,8]]],[[58,72],[[5,3,1,2,6,8]]],[[58,71],[[5,3,1,2,6,8]]],[[58,64],[[5,3,1,2,6,8]]],[[58,113],[[5,3,1,2,6,8]]],[[58,139],[[5,3,1,2,6,8]]],[[58,150],[[5,3,1,2,6,8]]],[[58,169],[[5,3,1,2,6,8]]],[[58,165],[[5,3,1,2,6,8]]],[[58,151],[[5,3,1,2,6,8]]],[[58,163],[[5,3,1,2,6,8]]],[[58,32],[[5,3,1,2,6,8]]],[[58,16],[[5,3,1,2,6,8]]],[[58,108],[[5,3,1,2,6,8]]],[[58,100],[[5,3,1,2,6,8]]],[[58,194],[[5,3,1,2,6,8]]],[[58,211],[[5,3,1,2,6,8]]],[[58,184],[[5,3,1,2,6,8]]],[[58,183],[[5,3,1,2,6,8]]],[[44,68],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,1],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,31],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,104],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,9],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,8],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,27],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,12],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,65],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,110],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,11],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,56],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,55],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,96],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,10],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,122],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,72],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,71],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,64],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,113],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,139],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,150],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,169],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,165],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,151],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,163],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,32],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,16],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,108],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,100],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,194],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,211],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,184],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,183],[[3,4,9,6,7,2,8,1,10,11,5]]],[[1,68],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,1],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,31],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,104],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,9],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,8],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,27],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,12],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,65],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,110],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,11],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,56],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,55],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,96],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,10],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,122],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,72],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,71],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,64],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,113],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,139],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,150],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,169],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,165],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,151],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,163],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,32],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,16],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,108],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,100],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,194],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,211],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,184],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,183],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[4,68],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,1],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,31],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,104],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,9],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,8],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,27],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,12],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,65],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,110],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,11],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,56],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,55],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,96],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,10],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,122],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,72],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,71],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,64],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,113],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,139],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,150],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,169],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,165],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,151],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,163],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,32],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,16],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,108],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,100],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,194],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,211],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,184],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,183],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[3,68],[[1,5,14,4,10,17]]],[[3,1],[[1,5,14,4,10,17]]],[[3,31],[[1,5,14,4,10,17]]],[[3,104],[[1,5,14,4,10,17]]],[[3,9],[[1,5,14,4,10,17]]],[[3,8],[[1,5,14,4,10,17]]],[[3,27],[[1,5,14,4,10,17]]],[[3,12],[[1,5,14,4,10,17]]],[[3,65],[[1,5,14,4,10,17]]],[[3,110],[[1,5,14,4,10,17]]],[[3,11],[[1,5,14,4,10,17]]],[[3,56],[[1,5,14,4,10,17]]],[[3,55],[[1,5,14,4,10,17]]],[[3,96],[[1,5,14,4,10,17]]],[[3,10],[[1,5,14,4,10,17]]],[[3,122],[[1,5,14,4,10,17]]],[[3,72],[[1,5,14,4,10,17]]],[[3,71],[[1,5,14,4,10,17]]],[[3,64],[[1,5,14,4,10,17]]],[[3,113],[[1,5,14,4,10,17]]],[[3,139],[[1,5,14,4,10,17]]],[[3,150],[[1,5,14,4,10,17]]],[[3,169],[[1,5,14,4,10,17]]],[[3,165],[[1,5,14,4,10,17]]],[[3,151],[[1,5,14,4,10,17]]],[[3,163],[[1,5,14,4,10,17]]],[[3,32],[[1,5,14,4,10,17]]],[[3,16],[[1,5,14,4,10,17]]],[[3,108],[[1,5,14,4,10,17]]],[[3,100],[[1,5,14,4,10,17]]],[[3,194],[[1,5,14,4,10,17]]],[[3,211],[[1,5,14,4,10,17]]],[[3,184],[[1,5,14,4,10,17]]],[[3,183],[[1,5,14,4,10,17]]],[[2,68],[[1,5,7,4,13,16,12,18]]],[[2,1],[[1,5,7,4,13,16,12,18]]],[[2,31],[[1,5,7,4,13,16,12,18]]],[[2,104],[[1,5,7,4,13,16,12,18]]],[[2,9],[[1,5,7,4,13,16,12,18]]],[[2,8],[[1,5,7,4,13,16,12,18]]],[[2,27],[[1,5,7,4,13,16,12,18]]],[[2,12],[[1,5,7,4,13,16,12,18]]],[[2,65],[[1,5,7,4,13,16,12,18]]],[[2,110],[[1,5,7,4,13,16,12,18]]],[[2,11],[[1,5,7,4,13,16,12,18]]],[[2,56],[[1,5,7,4,13,16,12,18]]],[[2,55],[[1,5,7,4,13,16,12,18]]],[[2,96],[[1,5,7,4,13,16,12,18]]],[[2,10],[[1,5,7,4,13,16,12,18]]],[[2,122],[[1,5,7,4,13,16,12,18]]],[[2,72],[[1,5,7,4,13,16,12,18]]],[[2,71],[[1,5,7,4,13,16,12,18]]],[[2,64],[[1,5,7,4,13,16,12,18]]],[[2,113],[[1,5,7,4,13,16,12,18]]],[[2,139],[[1,5,7,4,13,16,12,18]]],[[2,150],[[1,5,7,4,13,16,12,18]]],[[2,169],[[1,5,7,4,13,16,12,18]]],[[2,165],[[1,5,7,4,13,16,12,18]]],[[2,151],[[1,5,7,4,13,16,12,18]]],[[2,163],[[1,5,7,4,13,16,12,18]]],[[2,32],[[1,5,7,4,13,16,12,18]]],[[2,16],[[1,5,7,4,13,16,12,18]]],[[2,108],[[1,5,7,4,13,16,12,18]]],[[2,100],[[1,5,7,4,13,16,12,18]]],[[2,194],[[1,5,7,4,13,16,12,18]]],[[2,211],[[1,5,7,4,13,16,12,18]]],[[2,184],[[1,5,7,4,13,16,12,18]]],[[2,183],[[1,5,7,4,13,16,12,18]]]]]]`
	unknown2 := `[null,null,[[[1,null,1],null,[[null,[]]],null,null,null,null,[null,2],null,null,null,null,null,null,null,null,null,null,null,null,null,null,[1]],[null,[[null,[]]],null,null,[1]],[null,[[null,[]]],null,[1]],[null,[[null,[]]]],null,null,null,null,[[[null,[]]]],[[[null,[]]]]],[[[[7,68],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,1],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,31],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,104],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,9],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,8],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,27],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,12],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,65],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,110],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,11],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,56],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,55],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,96],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,10],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,122],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,72],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,71],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,64],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,113],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,139],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,150],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,169],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,165],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,151],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,163],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,32],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,16],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,108],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,100],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,194],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,211],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,184],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[7,183],[[1,73,96,103,97,58,50,92,52,112,69,19,31,101,123,74,49,80,20,10,14,79,43,42,139,63,169,95,155]]],[[9,68],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,1],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,31],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,104],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,9],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,8],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,27],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,12],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,65],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,110],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,11],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,56],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,55],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,96],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,10],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,122],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,72],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,71],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,64],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,113],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,139],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,150],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,169],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,165],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,151],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,163],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,32],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,16],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,108],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,100],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,194],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,211],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,184],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[9,183],[[1,7,9,24,12,31,5,15,27,8,13,10]]],[[17,68],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,1],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,31],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,104],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,9],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,8],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,27],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,12],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,65],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,110],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,11],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,56],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,55],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,96],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,10],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,122],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,72],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,71],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,64],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,113],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,139],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,150],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,169],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,165],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,151],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,163],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,32],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,16],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,108],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,100],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,194],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,211],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,184],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[17,183],[[1,7,9,25,13,31,5,41,27,8,14,10]]],[[65,68],[[1,5,4,7,11,6]]],[[65,1],[[1,5,4,7,11,6]]],[[65,31],[[1,5,4,7,11,6]]],[[65,104],[[1,5,4,7,11,6]]],[[65,9],[[1,5,4,7,11,6]]],[[65,8],[[1,5,4,7,11,6]]],[[65,27],[[1,5,4,7,11,6]]],[[65,12],[[1,5,4,7,11,6]]],[[65,65],[[1,5,4,7,11,6]]],[[65,110],[[1,5,4,7,11,6]]],[[65,11],[[1,5,4,7,11,6]]],[[65,56],[[1,5,4,7,11,6]]],[[65,55],[[1,5,4,7,11,6]]],[[65,96],[[1,5,4,7,11,6]]],[[65,10],[[1,5,4,7,11,6]]],[[65,122],[[1,5,4,7,11,6]]],[[65,72],[[1,5,4,7,11,6]]],[[65,71],[[1,5,4,7,11,6]]],[[65,64],[[1,5,4,7,11,6]]],[[65,113],[[1,5,4,7,11,6]]],[[65,139],[[1,5,4,7,11,6]]],[[65,150],[[1,5,4,7,11,6]]],[[65,169],[[1,5,4,7,11,6]]],[[65,165],[[1,5,4,7,11,6]]],[[65,151],[[1,5,4,7,11,6]]],[[65,163],[[1,5,4,7,11,6]]],[[65,32],[[1,5,4,7,11,6]]],[[65,16],[[1,5,4,7,11,6]]],[[65,108],[[1,5,4,7,11,6]]],[[65,100],[[1,5,4,7,11,6]]],[[65,194],[[1,5,4,7,11,6]]],[[65,211],[[1,5,4,7,11,6]]],[[65,184],[[1,5,4,7,11,6]]],[[65,183],[[1,5,4,7,11,6]]],[[10,68],[[1,7,6,9,15,8]]],[[10,1],[[1,7,6,9,15,8]]],[[10,31],[[1,7,6,9,15,8]]],[[10,104],[[1,7,6,9,15,8]]],[[10,9],[[1,7,6,9,15,8]]],[[10,8],[[1,7,6,9,15,8]]],[[10,27],[[1,7,6,9,15,8]]],[[10,12],[[1,7,6,9,15,8]]],[[10,65],[[1,7,6,9,15,8]]],[[10,110],[[1,7,6,9,15,8]]],[[10,11],[[1,7,6,9,15,8]]],[[10,56],[[1,7,6,9,15,8]]],[[10,55],[[1,7,6,9,15,8]]],[[10,96],[[1,7,6,9,15,8]]],[[10,10],[[1,7,6,9,15,8]]],[[10,122],[[1,7,6,9,15,8]]],[[10,72],[[1,7,6,9,15,8]]],[[10,71],[[1,7,6,9,15,8]]],[[10,64],[[1,7,6,9,15,8]]],[[10,113],[[1,7,6,9,15,8]]],[[10,139],[[1,7,6,9,15,8]]],[[10,150],[[1,7,6,9,15,8]]],[[10,169],[[1,7,6,9,15,8]]],[[10,165],[[1,7,6,9,15,8]]],[[10,151],[[1,7,6,9,15,8]]],[[10,163],[[1,7,6,9,15,8]]],[[10,32],[[1,7,6,9,15,8]]],[[10,16],[[1,7,6,9,15,8]]],[[10,108],[[1,7,6,9,15,8]]],[[10,100],[[1,7,6,9,15,8]]],[[10,194],[[1,7,6,9,15,8]]],[[10,211],[[1,7,6,9,15,8]]],[[10,184],[[1,7,6,9,15,8]]],[[10,183],[[1,7,6,9,15,8]]],[[58,68],[[5,3,1,2,6,8]]],[[58,1],[[5,3,1,2,6,8]]],[[58,31],[[5,3,1,2,6,8]]],[[58,104],[[5,3,1,2,6,8]]],[[58,9],[[5,3,1,2,6,8]]],[[58,8],[[5,3,1,2,6,8]]],[[58,27],[[5,3,1,2,6,8]]],[[58,12],[[5,3,1,2,6,8]]],[[58,65],[[5,3,1,2,6,8]]],[[58,110],[[5,3,1,2,6,8]]],[[58,11],[[5,3,1,2,6,8]]],[[58,56],[[5,3,1,2,6,8]]],[[58,55],[[5,3,1,2,6,8]]],[[58,96],[[5,3,1,2,6,8]]],[[58,10],[[5,3,1,2,6,8]]],[[58,122],[[5,3,1,2,6,8]]],[[58,72],[[5,3,1,2,6,8]]],[[58,71],[[5,3,1,2,6,8]]],[[58,64],[[5,3,1,2,6,8]]],[[58,113],[[5,3,1,2,6,8]]],[[58,139],[[5,3,1,2,6,8]]],[[58,150],[[5,3,1,2,6,8]]],[[58,169],[[5,3,1,2,6,8]]],[[58,165],[[5,3,1,2,6,8]]],[[58,151],[[5,3,1,2,6,8]]],[[58,163],[[5,3,1,2,6,8]]],[[58,32],[[5,3,1,2,6,8]]],[[58,16],[[5,3,1,2,6,8]]],[[58,108],[[5,3,1,2,6,8]]],[[58,100],[[5,3,1,2,6,8]]],[[58,194],[[5,3,1,2,6,8]]],[[58,211],[[5,3,1,2,6,8]]],[[58,184],[[5,3,1,2,6,8]]],[[58,183],[[5,3,1,2,6,8]]],[[44,68],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,1],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,31],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,104],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,9],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,8],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,27],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,12],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,65],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,110],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,11],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,56],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,55],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,96],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,10],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,122],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,72],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,71],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,64],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,113],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,139],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,150],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,169],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,165],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,151],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,163],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,32],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,16],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,108],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,100],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,194],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,211],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,184],[[3,4,9,6,7,2,8,1,10,11,5]]],[[44,183],[[3,4,9,6,7,2,8,1,10,11,5]]],[[1,68],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,1],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,31],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,104],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,9],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,8],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,27],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,12],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,65],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,110],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,11],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,56],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,55],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,96],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,10],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,122],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,72],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,71],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,64],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,113],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,139],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,150],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,169],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,165],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,151],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,163],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,32],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,16],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,108],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,100],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,194],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,211],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,184],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[1,183],[[1,5,14,38,19,29,34,4,12,11,6,30,43,40,42,16,10,7]]],[[4,68],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,1],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,31],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,104],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,9],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,8],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,27],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,12],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,65],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,110],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,11],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,56],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,55],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,96],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,10],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,122],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,72],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,71],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,64],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,113],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,139],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,150],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,169],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,165],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,151],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,163],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,32],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,16],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,108],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,100],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,194],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,211],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,184],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[4,183],[[1,3,5,4,7,6,11,19,21,17,15,12,16,20]]],[[3,68],[[1,5,14,4,10,17]]],[[3,1],[[1,5,14,4,10,17]]],[[3,31],[[1,5,14,4,10,17]]],[[3,104],[[1,5,14,4,10,17]]],[[3,9],[[1,5,14,4,10,17]]],[[3,8],[[1,5,14,4,10,17]]],[[3,27],[[1,5,14,4,10,17]]],[[3,12],[[1,5,14,4,10,17]]],[[3,65],[[1,5,14,4,10,17]]],[[3,110],[[1,5,14,4,10,17]]],[[3,11],[[1,5,14,4,10,17]]],[[3,56],[[1,5,14,4,10,17]]],[[3,55],[[1,5,14,4,10,17]]],[[3,96],[[1,5,14,4,10,17]]],[[3,10],[[1,5,14,4,10,17]]],[[3,122],[[1,5,14,4,10,17]]],[[3,72],[[1,5,14,4,10,17]]],[[3,71],[[1,5,14,4,10,17]]],[[3,64],[[1,5,14,4,10,17]]],[[3,113],[[1,5,14,4,10,17]]],[[3,139],[[1,5,14,4,10,17]]],[[3,150],[[1,5,14,4,10,17]]],[[3,169],[[1,5,14,4,10,17]]],[[3,165],[[1,5,14,4,10,17]]],[[3,151],[[1,5,14,4,10,17]]],[[3,163],[[1,5,14,4,10,17]]],[[3,32],[[1,5,14,4,10,17]]],[[3,16],[[1,5,14,4,10,17]]],[[3,108],[[1,5,14,4,10,17]]],[[3,100],[[1,5,14,4,10,17]]],[[3,194],[[1,5,14,4,10,17]]],[[3,211],[[1,5,14,4,10,17]]],[[3,184],[[1,5,14,4,10,17]]],[[3,183],[[1,5,14,4,10,17]]],[[2,68],[[1,5,7,4,13,16,12,18]]],[[2,1],[[1,5,7,4,13,16,12,18]]],[[2,31],[[1,5,7,4,13,16,12,18]]],[[2,104],[[1,5,7,4,13,16,12,18]]],[[2,9],[[1,5,7,4,13,16,12,18]]],[[2,8],[[1,5,7,4,13,16,12,18]]],[[2,27],[[1,5,7,4,13,16,12,18]]],[[2,12],[[1,5,7,4,13,16,12,18]]],[[2,65],[[1,5,7,4,13,16,12,18]]],[[2,110],[[1,5,7,4,13,16,12,18]]],[[2,11],[[1,5,7,4,13,16,12,18]]],[[2,56],[[1,5,7,4,13,16,12,18]]],[[2,55],[[1,5,7,4,13,16,12,18]]],[[2,96],[[1,5,7,4,13,16,12,18]]],[[2,10],[[1,5,7,4,13,16,12,18]]],[[2,122],[[1,5,7,4,13,16,12,18]]],[[2,72],[[1,5,7,4,13,16,12,18]]],[[2,71],[[1,5,7,4,13,16,12,18]]],[[2,64],[[1,5,7,4,13,16,12,18]]],[[2,113],[[1,5,7,4,13,16,12,18]]],[[2,139],[[1,5,7,4,13,16,12,18]]],[[2,150],[[1,5,7,4,13,16,12,18]]],[[2,169],[[1,5,7,4,13,16,12,18]]],[[2,165],[[1,5,7,4,13,16,12,18]]],[[2,151],[[1,5,7,4,13,16,12,18]]],[[2,163],[[1,5,7,4,13,16,12,18]]],[[2,32],[[1,5,7,4,13,16,12,18]]],[[2,16],[[1,5,7,4,13,16,12,18]]],[[2,108],[[1,5,7,4,13,16,12,18]]],[[2,100],[[1,5,7,4,13,16,12,18]]],[[2,194],[[1,5,7,4,13,16,12,18]]],[[2,211],[[1,5,7,4,13,16,12,18]]],[[2,184],[[1,5,7,4,13,16,12,18]]],[[2,183],[[1,5,7,4,13,16,12,18]]]]]]`
	var unknownArr2 []interface{}
	err = json.Unmarshal([]byte(unknown2), &unknownArr2)
	if err != nil {
		log.Error(err)
	}
	unknown3 := []interface{}{[]interface{}{[]interface{}{1, 2}, []interface{}{10, 8, 9}}}
	arr = append(arr, outPageArr, nil, nil, unknownArr1, unknownArr2, nil, nil, unknown3)
	return arr
}

// product第2个请求中的数组
func generalProductArr2() []interface{} {
	generalStr := `[null,null,[[1,9,10,11,13,14,19,20,38,43,47,49,52,58,59,63,69,70,73,74,75,78,79,80,91,92,95,96,97,100,101,103,106,112,119,129,137,139,141,145,146,149,151,155,169]],[[[1,null,1],null,[[[]]],null,null,null,null,[null,2],null,null,null,null,null,null,null,null,null,null,null,null,null,null,[1]],[null,[[[]]],null,null,[1]],[null,[[[]]],null,[1]],[null,[[[]]]],null,null,null,null,[[[[]]]],[[[[]]]]],null]`
	var generalArr []interface{}
	err := json.Unmarshal([]byte(generalStr), &generalArr)
	if err != nil {
		log.Error(err)
	}
	return generalArr
}

// product第5个请求中的数组
func generalProductArr34() []interface{} {
	generalStr := `[null,[[3,[10]],null,null,[184]]]`
	var generalArr []interface{}
	err := json.Unmarshal([]byte(generalStr), &generalArr)
	if err != nil {
		log.Error(err)
	}
	return generalArr
}

// product第12个请求中的数组
func generalProductArrForWs7gDc() []interface{} {
	generalStr := `[[null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,[2]]]`
	var generalArr []interface{}
	err := json.Unmarshal([]byte(generalStr), &generalArr)
	if err != nil {
		log.Error(err)
	}
	return generalArr
}

// 数组元素第二个，引号中的数组内容，最后一个元素
func productLst(productId string) (marshProductLstArr []interface{}) {
	return []interface{}{productId, 7}
}

func productMark(p *RequestParams, i int) []interface{} {
	rpcids := rpcidsEnum["product"]
	// 根据观察，标志为1开始的等差数列，2为差值
	markSym := fmt.Sprintf("%d", 2*(i+1)-1)
	lst := productLst(p.ProductID)
	reviewPlatform := productReviewPlatformEnum[p.Platform]
	var rpcid string
	rpcid = rpcids[i]
	quoMark := []interface{}{lst}
	switch i {
	case 2:
		quoMark = generalProductArr2()
		quoMark = append(quoMark, []interface{}{lst})
	case 3, 4:
		tempArr := append(generalProductArr34(), lst)
		quoMark = []interface{}{tempArr}
	case 5:
		quoMark = []interface{}{
			[]interface{}{nil, lst, nil, []interface{}{[]interface{}{3, []interface{}{6}}, nil, nil, []interface{}{1, 8}}}, []interface{}{1},
		}
	case 7:
		quoMark = []interface{}{nil, nil, [][]interface{}{{52}}, generalProductArrForWs7gDc(), nil, []interface{}{lst}}
	case 8:
		sortBy, err := strconv.Atoi(p.SortBy)
		if err != nil {
			sortBy = 1
		}
		var rating interface{}
		if p.Rating != "" {
			rating, err = strconv.Atoi(p.Rating)
			if err != nil {
				rating = nil
			}
		} else {
			rating = nil
		}
		num, err := strconv.Atoi(p.Num)
		page := []interface{}{num}
		if p.NextPageToken != "" {
			page = append(page, nil, p.NextPageToken)
		}
		quoMark = []interface{}{nil, []interface{}{2, sortBy, page, nil, []interface{}{nil, rating, nil, nil, nil, nil, nil, nil, reviewPlatform}}, lst}
	}
	quoMarkMarsh, err := json.Marshal(quoMark)
	if err != nil {
		log.Error(err)
	}
	return []interface{}{rpcid, string(quoMarkMarsh), nil, markSym}
}
