package tag

import (
	//	"encoding/json"
	"strings"

	//	"github.com/firefirestyle/go.miniprop"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
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

const (
	TypeRootGroup = "RootGroup"
	TypeMainTag   = "MainTag"
	TypeSubTag    = "SubTag"
	TypeTargetId  = "TargetId"
	TypeInfo      = "Info"
	TypeCreated   = "Created"
	TypePriority  = "Priority"
	TypeType      = "Type"
)

func NewTagManager(kind string, rootGroup string) *TagManager {
	ret := new(TagManager)
	ret.kind = kind
	ret.rootGroup = rootGroup
	return ret
}

func (obj *TagManager) DeleteTagsFromTargetId(ctx context.Context, targetId string, cursor string) (string, error) {
	q := datastore.NewQuery(obj.kind)
	q = q.Filter("RootGroup =", obj.rootGroup)
	q = q.Filter("TargetId =", targetId)
	q = q.Order("-Created")
	r, _, eCursor := obj.FindTagKeyFromQuery(ctx, q, cursor)
	for _, v := range r {
		datastore.Delete(ctx, v)
	}
	return eCursor, nil
}

func (obj *TagManager) AddPairTags(ctx context.Context, tagList []string, info string, targetId string) error {
	vs := obj.MakeTags(ctx, tagList)
	for _, v := range vs {
		log.Infof(ctx, ">>"+v.MainTag+" : "+v.SubTag)
		tag := obj.NewTag(ctx, v.MainTag, v.SubTag, targetId, v.Type)
		err := tag.SaveOnDB(ctx)
		if err != nil {
			log.Infof(ctx, ">>"+err.Error())
		}
	}
	return nil
}

func (obj *TagManager) DeletePairTags(ctx context.Context, tagList []string, info string, targetId string) error {
	vs := obj.MakeTags(ctx, tagList)
	for _, v := range vs {
		key := obj.NewTagKey(ctx, v.MainTag, v.SubTag, targetId, v.Type)
		datastore.Delete(ctx, key)
	}
	return nil
}

func (obj *TagManager) AddMainTag(ctx context.Context, tag1 string, tag2 string, info string, targetId string) error {
	return obj.AddTag(ctx, tag1, tag2, info, targetId, "main")
}

func (obj *TagManager) AddSubTag(ctx context.Context, tag1 string, tag2 string, info string, targetId string) error {
	return obj.AddTag(ctx, tag1, tag2, info, targetId, "sub")
}

//
//
//
func (obj *TagManager) AddTag(ctx context.Context, tag1 string, tag2 string, info string, targetId string, Type string) error {
	mainTag := tag1
	subTag := tag2
	if subTag != "" && strings.Compare(tag1, tag2) <= 0 {
		mainTag = tag2
		subTag = tag1
	}
	tag := obj.NewTag(ctx, mainTag, subTag, targetId, Type)
	tag.gaeObject.Info = info
	return tag.SaveOnDB(ctx)
}

func (obj *TagManager) DeleteMainTag(ctx context.Context, MainTag string, SubTag string, TargetId string) error {
	return obj.DeleteTag(ctx, MainTag, SubTag, TargetId, "main")
}

func (obj *TagManager) DeleteSubTag(ctx context.Context, MainTag string, SubTag string, TargetId string) error {
	return obj.DeleteTag(ctx, MainTag, SubTag, TargetId, "sub")
}

func (obj *TagManager) DeleteTag(ctx context.Context, MainTag string, SubTag string, TargetId string, Type string) error {
	key := obj.NewTagKey(ctx, MainTag, SubTag, TargetId, Type)
	datastore.Delete(ctx, key)
	return nil
}

func (obj *TagManager) MakeTags(ctx context.Context, tagList []string) []TagSource {
	ret := make([]TagSource, 0)
	for _, x := range tagList {
		isSave := false
		for _, y := range tagList {
			if strings.Compare(x, y) > 0 {
				t := "sub"
				if isSave == false {
					t = "main"
				}
				ret = append(ret, TagSource{
					MainTag: x,
					SubTag:  y,
					Type:    t,
				})
				isSave = true
			}
		}
		if isSave == false {
			ret = append(ret, TagSource{
				MainTag: x,
				SubTag:  "",
				Type:    "main",
			})
		}
	}
	return ret
}
