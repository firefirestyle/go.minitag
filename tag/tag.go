package tag

import (
	"encoding/json"
	//	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
)

type GaeObjectTag struct {
	ProjectId string
	MainTag   string
	SubTag    string
	TargetId  string
	Info      string `datastore:",noindex"`
	Created   time.Time
	Priority  int
	Type      string
}

type Tag struct {
	gaeObject    *GaeObjectTag
	gaeObjectKey *datastore.Key
	kind         string
}

func (obj *Tag) GetProjectId() string {
	return obj.gaeObject.ProjectId
}

func (obj *Tag) GetMainTag() string {
	return obj.gaeObject.MainTag
}

func (obj *Tag) GetSubTag() string {
	return obj.gaeObject.SubTag
}

func (obj *Tag) GetTargetId() string {
	return obj.gaeObject.TargetId
}

func (obj *Tag) GetCreated() time.Time {
	return obj.gaeObject.Created
}

func (obj *Tag) GetPriority() int {
	return obj.gaeObject.Priority
}

func (obj *Tag) GetGaeObjectKey() *datastore.Key {
	return obj.gaeObjectKey
}

func (obj *Tag) SaveOnDB(ctx context.Context) error {
	_, err := datastore.Put(ctx, obj.gaeObjectKey, obj.gaeObject)
	memSrc, errMemSrc := obj.toJson()
	if err == nil && errMemSrc == nil {
		objMem := &memcache.Item{
			Key:   obj.gaeObjectKey.StringID(),
			Value: []byte(memSrc), //
		}
		memcache.Set(ctx, objMem)
	}
	return err
}

func (obj *Tag) toJson() (string, error) {
	v := map[string]interface{}{
		TypeRootGroup: obj.gaeObject.ProjectId,
		TypeMainTag:   obj.gaeObject.MainTag,
		TypeSubTag:    obj.gaeObject.SubTag,
		TypeTargetId:  obj.gaeObject.TargetId,
		TypeInfo:      obj.gaeObject.Info,
		TypeCreated:   obj.GetCreated().UnixNano(),
		TypePriority:  obj.gaeObject.Priority,
		TypeType:      obj.gaeObject.Type,
	}
	vv, e := json.Marshal(v)
	return string(vv), e
}

//
func (obj *Tag) SetParamFromsJson(ctx context.Context, source []byte) error {
	v := make(map[string]interface{})
	e := json.Unmarshal(source, &v)
	if e != nil {
		return e
	}
	//
	obj.gaeObject.ProjectId = v[TypeRootGroup].(string)
	obj.gaeObject.MainTag = v[TypeMainTag].(string)
	obj.gaeObject.SubTag = v[TypeSubTag].(string)
	obj.gaeObject.TargetId = v[TypeTargetId].(string)
	obj.gaeObject.Info = v[TypeInfo].(string)
	obj.gaeObject.Created = time.Unix(0, int64(v[TypeCreated].(float64))) //srcCreated
	obj.gaeObject.Priority = int(v[TypePriority].(float64))
	obj.gaeObject.Type = v[TypeType].(string)

	return nil
}
