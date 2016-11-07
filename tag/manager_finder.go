package tag

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type FoundTags struct {
	Keys       []string
	CursorOne  string
	CursorNext string
}

/*
https://cloud.google.com/appengine/docs/go/config/indexconfig#updating_indexes
*/
func (obj *TagManager) FindTags(ctx context.Context, mainTag string, subTag string, cursorSrc string) FoundTags {
	q := datastore.NewQuery(obj.kind)
	q = q.Filter("ProjectId =", obj.rootGroup)
	q = q.Filter("MainTag =", mainTag)

	if subTag != "" {
		q = q.Filter("SubTag =", subTag)
	} else {
		q = q.Filter("Type =", "main")
	}

	q = q.Order("-Created").Limit(10)
	return obj.FindTagFromQuery(ctx, q, cursorSrc)
}

/*
- kind: ArticleTag
  properties:
  - name: ProjectId
  - name: TargetId
  - name: Created
    direction: desc
https://cloud.google.com/appengine/docs/go/config/indexconfig#updating_indexes
*/
func (obj *TagManager) FindTagFromTargetId(ctx context.Context, targetTag string, cursorSrc string) FoundTags {
	q := datastore.NewQuery(obj.kind)
	q = q.Filter("ProjectId =", obj.rootGroup)
	q = q.Filter("TargetId =", targetTag)
	q = q.Order("-Created").Limit(10)
	return obj.FindTagFromQuery(ctx, q, cursorSrc)
}

func (obj *TagManager) FindTagFromOwner(ctx context.Context, owner string, cursorSrc string) FoundTags {
	q := datastore.NewQuery(obj.kind)
	q = q.Filter("ProjectId =", obj.rootGroup)
	q = q.Filter("Owner =", owner)
	q = q.Order("-Created").Limit(10)
	return obj.FindTagFromQuery(ctx, q, cursorSrc)
}

func (obj *TagManager) FindTagFromQuery(ctx context.Context, q *datastore.Query, cursorSrc string) FoundTags {
	cursor := obj.newCursorFromSrc(cursorSrc)
	if cursor != nil {
		q = q.Start(*cursor)
	}
	q = q.KeysOnly()
	founds := q.Run(ctx)

	var retUser []string

	var cursorNext string = ""
	var cursorOne string = ""

	for i := 0; ; i++ {
		var d GaeObjectTag
		key, err := founds.Next(&d)
		if err != nil || err == datastore.Done {
			break
		} else {
			retUser = append(retUser, key.StringID())
		}
		if i == 0 {
			cursorOne = obj.makeCursorSrc(founds)
		}
	}
	cursorNext = obj.makeCursorSrc(founds)
	return FoundTags{
		Keys:       retUser,
		CursorOne:  cursorOne,
		CursorNext: cursorNext,
	}
}

func (obj *TagManager) newCursorFromSrc(cursorSrc string) *datastore.Cursor {
	c1, e := datastore.DecodeCursor(cursorSrc)
	if e != nil {
		return nil
	} else {
		return &c1
	}
}

func (obj *TagManager) makeCursorSrc(founds *datastore.Iterator) string {
	c, e := founds.Cursor()
	if e == nil {
		return c.String()
	} else {
		return ""
	}
}
