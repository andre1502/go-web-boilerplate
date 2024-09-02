package response

import (
	"boilerplate/utils"
	"boilerplate/utils/constant"
	cerror "boilerplate/utils/error"
	"boilerplate/utils/locale"
	"boilerplate/utils/logger"
	"net/http"
	"reflect"
	"strings"

	"math"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	locale     *locale.Locale
	Pagination *Pagination
	response   *ResponseResult
}

type ResponseResult struct {
	Result     Result      `json:"result"`
	Data       any         `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type Result struct {
	HttpStatus int         `json:"http_status"`
	Code       string      `json:"code"`
	Message    string      `json:"message"`
	Validator  []Validator `json:"validator,omitempty"`
}

type Validator struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Message string `json:"message"`
}

type Pagination struct {
	Page        int    `json:"page"`
	PageSize    int    `json:"page_size"`
	TotalPage   int    `json:"total_page"`
	TotalRecord uint64 `json:"total_record"`
}

func NewResponse(locale *locale.Locale) *Response {
	return &Response{
		locale:   locale,
		response: &ResponseResult{},
	}
}

func (r *Response) Json(c *gin.Context, httpStatus int, data any, err error) {
	c.Writer.Header().Add(constant.HEADER_CONTENT_TYPE, constant.APPLICATION_JSON)

	r.response.Data = nil
	r.response.Pagination = nil

	if err != nil {
		errC := cerror.ParseError(err)
		logger.Sugar.Error(errC)

		httpStatus = r.getHttpStatusCodeError(errC.Param.I18nMsg, httpStatus)
		r.response.Result.HttpStatus = httpStatus
		r.response.Result.Code = errC.Param.I18nMsg
		r.response.Result.Message = r.locale.Localize(r.response.Result.Code, errC.Param.Data)
		r.parseValidator(errC.Param.Data)

		c.AbortWithStatusJSON(httpStatus, r.response)
		return
	}

	r.response.Result.HttpStatus = httpStatus

	if httpStatus == http.StatusOK {
		r.response.Result.Code = "success"
	} else {
		r.response.Result.Code = "failed"
	}

	r.response.Result.Message = r.locale.Localize(r.response.Result.Code, nil)
	r.response.Data = data

	if (r.Pagination != nil) && (r.Pagination.Page > 0) {
		r.response.Pagination = r.calculatePagination()
	}

	c.IndentedJSON(httpStatus, r.response)
}

func (r *Response) getHttpStatusCodeError(errorCode string, httpStatus int) int {
	switch errorCode {
	case "failed_marshal":
		fallthrough
	case "failed_unmarshal":
		fallthrough
	case "failed_db_query":
		fallthrough
	case "failed_db_insert":
		fallthrough
	case "failed_db_update":
		fallthrough
	case "failed_db_transaction":
		fallthrough
	case "failed_db_commit":
		fallthrough
	case "failed_redis_set":
		fallthrough
	case "failed_redis_get":
		fallthrough
	case "failed_redis_remove":
		httpStatus = http.StatusInternalServerError
	}

	return httpStatus
}

func (r *Response) parseValidator(data map[string]any) {
	r.response.Result.Validator = nil

	if val, ok := data["validator"]; ok {
		logger.Sugar.Error(val)

		switch val := val.(type) {
		case validator.ValidationErrors:
			out := make([]Validator, len(val))
			for i, fe := range val {
				field := fe.Field()
				param := ""

				logger.Sugar.Debugf("field: %s", field)

				if len(fe.Param()) > 0 {
					param = fe.Param()
					logger.Sugar.Debugf("param: %s", param)

					if inputField, ok := data["input"]; ok {
						logger.Sugar.Debugf("inputField: %s", utils.DataString(inputField))

						if stField, ok := reflect.TypeOf(inputField).FieldByName(param); ok {
							logger.Sugar.Debugf("stField: %s", stField)

							param, _ = r.parseFieldName(stField)
						}
					}

					logger.Sugar.Debugf("param: %s", param)
				}

				out[i] = Validator{field, fe.Tag(), r.msgForTag(field, fe.Tag(), param)}
			}

			r.response.Result.Validator = out
		}
	}
}

func (r *Response) parseFieldName(f reflect.StructField) (name string, ignore bool) {
	tag := f.Tag.Get("json")

	if utils.IsEmptyString(tag) {
		return f.Name, false
	}

	if tag == "-" {
		return "", true
	}

	if i := strings.Index(tag, ","); i != -1 {
		if i == 0 {
			return f.Name, false
		} else {
			return tag[:i], false
		}
	}

	return tag, false
}

func (r *Response) msgForTag(field string, tag string, param string) string {
	switch tag {
	case "required":
		fallthrough
	case "empty_string":
		return r.locale.Localize(tag, map[string]any{"field": field})
	case "eqfield":
		fallthrough
	case "min":
		fallthrough
	case "max":
		return r.locale.Localize(tag, map[string]any{"field": field, "param": param})
	}

	return tag
}

func (r *Response) calculatePagination() *Pagination {
	if (r.Pagination != nil) && (r.Pagination.Page > 0) && (r.Pagination.PageSize > 0) {
		r.Pagination.TotalPage = int(math.Ceil(float64(r.Pagination.TotalRecord) / float64(r.Pagination.PageSize)))
		return r.Pagination
	}

	return nil
}
