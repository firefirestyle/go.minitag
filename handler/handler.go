package handler

import (
	"net/http"

	"github.com/firefirestyle/go.miniprop"
	"github.com/firefirestyle/go.minitag/tag"
	"google.golang.org/appengine"
)

type TagHandlerConfig struct {
	RootGroup string
	Kind      string
}

type TagHandler struct {
	config     TagHandlerConfig
	tagManager *tag.TagManager
}

func NewTagHandler(config TagHandlerConfig) *TagHandler {
	if config.Kind == "" {
		config.Kind = "fftag"
	}
	if config.Kind == "" {
		config.Kind = "ffstyle"
	}

	return &TagHandler{
		config:     config,
		tagManager: tag.NewTagManager(config.Kind, config.RootGroup),
	}
}

func (obj *TagHandler) HandleAdd(w http.ResponseWriter, r *http.Request) {
	inputProp := miniprop.NewMiniPropFromJsonReader(r.Body)
	tags := inputProp.GetPropStringList("", "tags", nil)
	value := inputProp.GetString("value", "")
	//	token := inputProp.GetString("token", "")
	ctx := appengine.NewContext(r)
	for _, v := range tags {
		obj.tagManager.AddBasicTag(ctx, v, "", value)
	}
}
