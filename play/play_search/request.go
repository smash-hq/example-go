package play_search

import (
	"context"
	"errors"
	"example/actor/play"
	log "github.com/sirupsen/logrus"
)

func Request(ctx context.Context, request *play.RequestParams) (*play.Response, error) {
	resp, err := play.DoPlay(ctx, request, request.Type)
	if err != nil {
		log.Errorf("success=false,  err=%v", err)
		return nil, errors.New("scraping failed")
	}

	return resp, nil
}
