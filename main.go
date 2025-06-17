package main

import (
	"context"
	"encoding/json"
	"example/actor/play"
	"example/actor/play/play_books"
	"example/actor/play/play_games"
	"example/actor/play/play_movies"
	"example/actor/play/play_product"
	"fmt"
	"github.com/scrapeless-ai/sdk-go/scrapeless/actor"
	"github.com/scrapeless-ai/sdk-go/scrapeless/services/proxies"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"runtime"
	"time"
)

var (
	Actor *actor.Actor
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		DisableColors: false,
		ForceColors:   true,
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			filename := path.Base(f.File)
			fc := path.Base(f.Function)
			return fmt.Sprintf("%s()", fc), fmt.Sprintf(" - %s:%d", filename, f.Line)
		},
		TimestampFormat: time.DateTime,
	})
	log.SetReportCaller(true)
	log.SetLevel(log.InfoLevel)
}

func main() {
	// new Actor
	Actor = actor.New()
	defer Actor.Close()
	var param = &play.RequestParams{}

	if err := Actor.Input(param); err != nil {
		log.Errorf("input error: %v", err)
	}
	marshal, _ := json.Marshal(param)
	log.Infof("params info:%v", string(marshal))
	// get proxy url
	proxy, err := Actor.Proxy.Proxy(context.TODO(), proxies.ProxyActor{
		Country:         "US",
		SessionDuration: 10,
	})
	if err != nil {
		log.Errorf("get proxy error: %v", err)
	}
	//proxy = "http://group_scraper_google_trneds:c8d2279d492a@pm-gw-us.scrapeless.io:24125"
	log.Infof("proxy url:%s", proxy)

	paramErr := param.FieldValidation(param.Type)
	if paramErr != nil {
		log.Warnf("param error: %v", paramErr)
	}
	var res *play.Response
	var resErr error
	switch param.Type {
	case play.GooglePlayGames:
		res, resErr = play_games.Request(context.TODO(), param, proxy)
	case play.GooglePlayProduct:
		res, resErr = play_product.Request(context.TODO(), param, proxy)
	case play.GooglePlayMovies:
		res, resErr = play_movies.Request(context.TODO(), param, proxy)
	case play.GooglePlayBooks:
		res, resErr = play_books.Request(context.TODO(), param, proxy)
	case play.GooglePlay:
		res, resErr = play_books.Request(context.TODO(), param, proxy)
	default:
		res, resErr = play_games.Request(context.TODO(), param, proxy)
	}
	if resErr != nil {
		log.Errorf("success=false,  err=%v", err)
		return
	}
	bytes, err := json.Marshal(res)
	if err != nil {
		log.Errorf("success=false,  err=%v", err)
		return
	}
	log.Infof("success=true, res=%s", bytes)

	if res == nil {
		log.Warnf("res is nil")
		return
	}
	items, err := Actor.AddItems(context.TODO(), []map[string]any{
		{
			"title":             "Play Store",
			"search_parameters": toString(res.SearchParameters),
			"search_metadata":   toString(res.SearchMetadata),
			"chart_option":      toString(res.ChartOption),
			"highlight_item":    toString(res.HighlightItem),
			"organic_results":   toString(res.OrganicResults),
			"product_info":      toString(res.ProductInfo),
			"Media":             toString(res.Media),
			"about_this_app":    toString(res.AboutThisApp),
			"categories":        toString(res.Categories),
			"what_s_new":        toString(res.WhatIsNew),
			"ratings":           toString(res.Ratings),
			"reviews":           toString(res.Reviews),
			"developer_contact": toString(res.DeveloperContact),
			"similar_results":   toString(res.SimilarResults),
		},
	})
	if err != nil {
		log.Warnf("add items error: %v", err)
	}
	log.Infof("add items success: %v", items)
}

func toString(obj any) string {
	marshal, _ := json.Marshal(obj)
	return string(marshal)
}
