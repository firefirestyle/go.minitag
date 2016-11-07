package tag

import (
	//	"encoding/json"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
)

func (obj *MiniTagManager) NewTag(ctx context.Context, mainTag string, //
	subTag string, targetId string, t string) *MiniTag {
	ret := new(MiniTag)
	ret.gaeObject = new(GaeObjectTag)
	ret.gaeObject.ProjectId = obj.projectId
	ret.gaeObject.MainTag = mainTag
	ret.gaeObject.SubTag = subTag
	ret.gaeObject.TargetId = targetId
	ret.gaeObjectKey = obj.NewTagKey(ctx, mainTag, subTag, targetId, t)
	ret.gaeObject.Created = time.Now()
	ret.gaeObject.Type = t
	return ret
}

func (obj *MiniTagManager) NewTagKey(ctx context.Context, mainTag string, //
	subTag string, targetId string, ttype string) *datastore.Key {
	ret := datastore.NewKey(ctx, obj.kind, ""+obj.projectId+"::"+targetId+"::"+mainTag+"::"+subTag+"::"+ttype, 0, nil)
	return ret
}

func (obj *MiniTagManager) NewTagFromKey(ctx context.Context, gaeKey *datastore.Key) (*MiniTag, error) {

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
