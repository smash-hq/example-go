package play_product

import (
	"context"
	"errors"
	"example/actor/play"
	log "github.com/sirupsen/logrus"
)

func Request(ctx context.Context, request *play.RequestParams, proxy string) (*play.Response, error) {
	resp, err := play.DoProductPlay(ctx, request, request.Type, proxy)
	if err != nil {
		log.Errorf("success=false,  err=%v", err)
		return nil, errors.New("scraping failed")
	}

	return resp, nil
}
