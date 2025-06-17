package main

import (
	"context"
	"encoding/json"
	"github.com/scrapeless-ai/scrapeless-actor-sdk-go/scrapeless"
	"github.com/scrapeless-ai/scrapeless-actor-sdk-go/scrapeless/actor"
	"github.com/scrapeless-ai/scrapeless-actor-sdk-go/scrapeless/services/proxies"
	"log"
)

var (
	Actor       *actor.Actor
	ActorClient *scrapeless.Client
	ProxyStr    string
)

func main() {
	// new Actor
	Actor = actor.New()
	defer Actor.Close()
	var param = &RequestParam{}

	if err := Actor.Input(param); err != nil {
		log.Fatal(err)
	}
	// get proxy url
	proxy, err := Actor.Proxy.Proxy(context.TODO(), proxies.ProxyActor{
		Country:         "US",
		SessionDuration: 10,
	})
	if err != nil {
		log.Fatal(err)
	}
	ProxyStr = proxy
	shopping, err := doShopping(context.TODO(), param)
	filters, _ := json.Marshal(shopping.Filters)
	results, _ := json.Marshal(shopping.ShoppingResults)
	inlineShoppingResults, _ := json.Marshal(shopping.InlineShoppingResults)
	ok, err := Actor.Storage.GetDataset().AddItems(context.Background(), []map[string]any{
		{
			"filters":                 filters,
			"results":                 results,
			"inline_shopping_results": inlineShoppingResults,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(ok)
}
