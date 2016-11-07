package tag

import (
	//	"encoding/json"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

type TagManager struct {
	kind      string
	rootGroup string
}

type TagSource struct {
	MainTag string
	SubTag  string
	Type    string
}

func NewTagManager(kind string, rootGroup string) *TagManager {
	ret := new(TagManager)
	ret.kind = kind
	ret.rootGroup = rootGroup
	return ret
}

func Debug(ctx context.Context, message string) {
	log.Infof(ctx, message)
}
