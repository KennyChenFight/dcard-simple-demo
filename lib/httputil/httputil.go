package httputil

import (
	"encoding/json"
	"github.com/KennyChenFight/dcard-simple-demo/lib/validate"
	"github.com/gin-gonic/gin"
	"xorm.io/core"
)

var columnNameMapper core.IMapper

func Init(mapper core.IMapper) {
	columnNameMapper = mapper
}

// bind the http request to a struct and return struct validation result
func BindForUpdate(c *gin.Context, obj interface{}) (map[string]bool, error) {
	var input map[string]interface{}
	c.ShouldBind(&input)
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(inputBytes, obj); err != nil {
		return nil, err
	}

	dbFieldNames := make(map[string]bool)
	structFieldNames := make(map[string]bool)
	for fieldName := range input {
		structFieldNames[fieldName] = true
		fieldName = columnNameMapper.Obj2Table(fieldName)
		dbFieldNames[fieldName] = true
	}

	return dbFieldNames, validate.StructForUpdate(obj, structFieldNames)
}
