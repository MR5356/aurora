package datafilter

import (
	"bytes"
	"encoding/json"
	"github.com/MR5356/aurora/internal/domain/authentication"
	"github.com/MR5356/aurora/internal/domain/user"
	database2 "github.com/MR5356/aurora/internal/infrastructure/database"
	"github.com/MR5356/aurora/internal/response"
	"github.com/MR5356/aurora/pkg/util/ginutil"
	"github.com/MR5356/aurora/pkg/util/structutil"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"io"
	"reflect"
	"sync"
)

var filterMap = sync.Map{}

type Filter struct {
	Function  func(ctx *gin.Context)
	IsBefore  bool
	Action    []string
	Domain    string
	FiledName string `default:"ID"`
}

func RegisterFilter(filters []Filter) {
	for _, filter := range filters {
		key := reflect.ValueOf(filter.Function).Pointer()
		if len(filter.FiledName) == 0 {
			filter.FiledName = "id"
		}
		filterMap.Store(key, filter)
		logrus.Debugf("register filter: %d with %+v", key, filter)
	}
}

type authedWriter struct {
	gin.ResponseWriter
	body   *bytes.Buffer
	user   *user.User
	filter Filter
}

func (w *authedWriter) Write(body []byte) (int, error) {
	w.body.Write(body)

	errRes := &response.Response{
		Code:    response.CodeNoPermission,
		Message: response.MessageNoPermission,
		Data:    nil,
	}

	errResponse, _ := json.Marshal(errRes)

	var res response.Response
	err := json.Unmarshal(body, &res)

	if err != nil {
		return w.ResponseWriter.Write(body)
	}

	if w.filter.IsBefore || len(w.filter.Domain) == 0 || len(w.filter.Action) == 0 {
		return w.ResponseWriter.Write(body)
	}

	if res.Data == nil {
		return w.ResponseWriter.Write(body)
	}

	if isArrayOrSlice(res.Data) {
		data, ok := res.Data.([]any)
		if !ok {
			return w.ResponseWriter.Write(errResponse)
		}

		filteredData, err := authentication.GetPermission().FilterDataArray(data, w.filter.Action, w.filter.Domain, getRoles(w.user.ID), w.filter.FiledName)
		if err != nil {
			return w.ResponseWriter.Write(errResponse)
		}

		res.Data = filteredData
		body, _ = json.Marshal(res)
		return w.ResponseWriter.Write(body)
	}

	page, ok := isPage(res.Data.(map[string]any))
	if ok {
		//logrus.Debugf("page: %+v", page)
		filteredData, err := authentication.GetPermission().FilterDataArray(page.Data, w.filter.Action, w.filter.Domain, getRoles(w.user.ID), w.filter.FiledName)
		if err != nil {
			return w.ResponseWriter.Write(errResponse)
		}

		page.Data = filteredData
		res.Data = page
		body, _ = json.Marshal(res)
		return w.ResponseWriter.Write(body)
	}

	return w.ResponseWriter.Write(body)
}

func AutomationFilter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		u, err := user.GetJWTService().ParseToken(ginutil.GetToken(ctx))
		if err != nil {
			response.Error(ctx, response.CodeNotLogin)
			ctx.Abort()
			return
		}
		filterKey := reflect.ValueOf(ctx.Handler()).Pointer()
		f, ok := filterMap.Load(filterKey)
		if ok {
			filter := f.(Filter)
			logrus.Debugf("filter: %+v", filter)
			if !filter.IsBefore {
				writer := &authedWriter{
					ResponseWriter: ctx.Writer,
					body:           bytes.NewBuffer([]byte{}),
					user:           u,
					filter:         filter,
				}

				ctx.Writer = writer
			} else {
				var object string

				// path id or query id
				object = ctx.Param(filter.FiledName)
				if len(object) == 0 {
					object = ctx.Query(filter.FiledName)
				}

				// body ids
				ids := make([]uuid.UUID, 0)
				// read and set body
				body, err := io.ReadAll(ctx.Request.Body)
				if err == nil {
					err = json.Unmarshal(body, &ids)
					if err == nil {
						logrus.Debugf("not body ids")
					}
					ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
				}

				logrus.Debugf("filter: %d, object: %s, ids: %+v", filterKey, object, ids)
				if len(object) != 0 {
					logrus.Debugf("path id or query id: %s", object)
					ok, err = checkObjectPermission(u, object, filter)

					logrus.Debugf("checkObjectPermission: %v, %v", ok, err)
					if err != nil || !ok {
						response.Error(ctx, response.CodeNoPermission)
						ctx.Abort()
						return
					}
				} else if len(ids) > 0 {
					logrus.Debugf("body ids: %+v", ids)
					for _, id := range ids {
						ok, err = checkObjectPermission(u, id.String(), filter)
						logrus.Debugf("checkObjectPermission: %v, %v", ok, err)
						if err != nil || !ok {
							response.Error(ctx, response.CodeNoPermission)
							ctx.Abort()
							return
						}
					}
				} else {
					response.Error(ctx, response.CodeNoPermission)
					ctx.Abort()
					return
				}
			}
		}
		ctx.Next()
	}
}

func checkObjectPermission(u *user.User, object string, filter Filter) (ok bool, err error) {
	// check path id or query id
	ok = false
MainLoop:
	for _, role := range getRoles(u.ID) {
		for _, action := range filter.Action {
			ok, err = authentication.GetPermission().HasPermissionForRoleInDomain(filter.Domain, role, object, action)
			if err != nil {
				return false, err
			}
			if ok {
				ok = true
				break MainLoop
			}
		}
	}
	return ok, nil
}

func isArrayOrSlice(data any) bool {
	if data == nil {
		return false
	}
	switch reflect.TypeOf(data).Kind() {
	case reflect.Array, reflect.Slice:
		return true
	default:
		return false
	}
}

func isPage(data map[string]interface{}) (database2.Pager[any], bool) {
	if data == nil {
		return database2.Pager[any]{}, false
	}

	//logrus.Debugf("data: %+v", data)
	c, ok := structutil.GetMapFiledByName(data, "current")
	if !ok {
		return database2.Pager[any]{}, false
	}
	t, ok := structutil.GetMapFiledByName(data, "total")
	if !ok {
		return database2.Pager[any]{}, false
	}
	s, ok := structutil.GetMapFiledByName(data, "size")
	if !ok {
		return database2.Pager[any]{}, false
	}
	d, ok := structutil.GetMapFiledByName(data, "data")
	if !ok {
		return database2.Pager[any]{}, false
	}

	return database2.Pager[any]{
		CurrentPage: anyToInt64(c),
		Total:       anyToInt64(t),
		PageSize:    anyToInt64(s),
		Data:        d.([]any),
	}, true
}

func anyToInt64(num any) int64 {
	switch num.(type) {
	case int64:
		return num.(int64)
	case int:
		return int64(num.(int))
	case int32:
		return int64(num.(int32))
	case int16:
		return int64(num.(int16))
	case int8:
		return int64(num.(int8))
	case float64:
		return int64(num.(float64))
	case float32:
		return int64(num.(float32))
	case uint:
		return int64(num.(uint))
	case uint64:
		return int64(num.(uint64))
	case uint32:
		return int64(num.(uint32))
	case uint16:
		return int64(num.(uint16))
	case uint8:
		return int64(num.(uint8))
	default:
		return 0
	}
}

func getRoles(userID string) []string {
	res := []string{userID}
	userGroupRelationDB := database2.NewMapper(database2.GetDB(), &user.Relation{})

	userGroupRelations, err := userGroupRelationDB.List(&user.Relation{
		UserID: userID,
	})
	if err != nil {
		return res
	}

	for _, relation := range userGroupRelations {
		res = append(res, relation.GroupID.String())
	}
	return res
}
