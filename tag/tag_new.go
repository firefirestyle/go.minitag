package tag

import (
	//	"encoding/json"
	"time"

	"github.com/firefirestyle/go.miniprop"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
)

func (obj *TagManager) NewTag(ctx context.Context, mainTag string, //
	subTag string, target string, tagType string) *MiniTag {
	ret := new(MiniTag)
	ret.gaeObject = new(GaeObjectTag)
	ret.gaeObject.ProjectId = obj.rootGroup
	ret.gaeObject.MainTag = mainTag
	ret.gaeObject.SubTag = subTag
	ret.gaeObject.TargetId = target
	ret.gaeObjectKey = obj.NewTagKey(ctx, mainTag, subTag, target, tagType)
	ret.gaeObject.Created = time.Now()
	ret.gaeObject.Type = tagType
	return ret
}

func (obj *TagManager) NewTagKey(ctx context.Context, mainTag string, //
	subTag string, targetId string, ttype string) *datastore.Key {
	ret := datastore.NewKey(ctx, obj.kind, obj.MakeStringId(mainTag, subTag, targetId, ttype), 0, nil)
	return ret
}

func (obj *TagManager) MakeStringId(mainTag string, //
	subTag string, targetId string, ttype string) string {
	propObj := miniprop.NewMiniProp()
	propObj.SetString("p", obj.rootGroup)
	propObj.SetString("v", targetId)
	propObj.SetString("m", mainTag)
	propObj.SetString("s", subTag)
	propObj.SetString("t", ttype)
	return string(propObj.ToJson())
}

type TagKeyInfo struct {
	RootGroup string
	Value     string
	MainTag   string
	SubTag    string
	TagType   string
}

func (obj *TagManager) GetKeyInfoFromStringId(stringId string) TagKeyInfo {
	propObj := miniprop.NewMiniPropFromJson([]byte(stringId))
	return TagKeyInfo{
		RootGroup: propObj.GetString("p", ""),
		Value:     propObj.GetString("v", ""),
		MainTag:   propObj.GetString("m", ""),
		SubTag:    propObj.GetString("s", ""),
		TagType:   propObj.GetString("t", ""),
	}
}

func (obj *TagManager) NewTagFromKey(ctx context.Context, gaeKey *datastore.Key) (*MiniTag, error) {

	ret := new(MiniTag)
	ret.kind = obj.kind
	ret.gaeObject = new(GaeObjectTag)
	ret.gaeObjectKey = gaeKey
	//
	//
	memObjSrc, memObjErr := memcache.Get(ctx, gaeKey.StringID())
	if memObjErr == nil {
		err := ret.SetParamFromsJson(ctx, memObjSrc.Value)
		if err == nil {
			return ret, nil
		}
	}
	//
	//
	err := datastore.Get(ctx, gaeKey, ret.gaeObject)
	if err != nil {
		return nil, err
	}
	//
	//	ret.SetParamFromsJson(ctx)

	return ret, nil
}
