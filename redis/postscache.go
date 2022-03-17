package redis

import (
	"todoapp/model"
)

type PostCache interface {
	Set(key string, value map[string]interface{})
	Get(key string) *model.TodoModel
}
