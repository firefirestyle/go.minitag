package handler

import (
	"net/http"

	"github.com/firefirestyle/go.minitag/tag"
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
	return &TagHandler{
		config:     config,
		tagManager: tag.NewTagManager(config.Kind, config.RootGroup),
	}
}

func (obj *TagHandler) HandleAdd(w http.ResponseWriter, r *http.Request) {

}
