package main

import (
	"context"
	"encoding/json"
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
	Actor    *actor.Actor
	ProxyStr string
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
	var param = &RequestParam{}

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
	ProxyStr = proxy
	log.Infof("proxy url:%s", ProxyStr)
	shopping, err := doShopping(context.TODO(), param)
	filters, _ := json.Marshal(shopping.Filters)
	results, _ := json.Marshal(shopping.ShoppingResults)
	inlineShoppingResults, _ := json.Marshal(shopping.InlineShoppingResults)
	serapMetadata, _ := json.Marshal(shopping.SearchMetadata)
	ok, err := Actor.AddItems(context.Background(), []map[string]any{
		{
			"filters":                 string(filters),
			"results":                 string(results),
			"inline_shopping_results": string(inlineShoppingResults),
			"serapMetadata":           string(serapMetadata),
		},
	})
	if err != nil {
		log.Errorf("add items error: %v", err)
	}
	log.Println(ok)
}
