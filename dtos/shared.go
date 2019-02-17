package dtos

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v8"
	"math"
	"net/http"
)

type BaseDto struct {
	Success      bool     `json:"success"`
	FullMessages []string `json:"full_messages"`
}

type ErrorDto struct {
	BaseDto
	Errors map[string]interface{} `json:"errors"`
}

func CreatePageMeta(request *http.Request, loadedItemsCount, page, page_size, totalItemsCount int) map[string]interface{} {
	page_meta := map[string]interface{}{}
	page_meta["offset"] = (page - 1) * page_size
	page_meta["requested_page_size"] = page_size
	page_meta["current_page_number"] = page
	page_meta["current_items_count"] = loadedItemsCount

	page_meta["prev_page_number"] = 1
	total_pages_count := int(math.Ceil(float64(totalItemsCount) / float64(page_size)))
	page_meta["total_pages_count"] = total_pages_count

	if page < total_pages_count {
		page_meta["has_next_page"] = true
		page_meta["next_page_number"] = page + 1
	} else {
		page_meta["has_next_page"] = false
		page_meta["next_page_number"] = 1
	}
	if page > 1 {
		page_meta["prev_page_number"] = page - 1
	} else {
		page_meta["has_prev_page"] = false
		page_meta["prev_page_number"] = 1
	}

	page_meta["next_page_url"] = fmt.Sprintf("%v?page=%d&page_size=%d", request.URL.Path, page_meta["next_page_number"], page_meta["requested_page_size"])
	page_meta["prev_page_url"] = fmt.Sprintf("%s?page=%d&page_size=%d", request.URL.Path, page_meta["prev_page_number"], page_meta["requested_page_size"])

	response := gin.H{
		"success":   true,
		"page_meta": page_meta,
	}

	return response
}

func CreatePagedResponse(request *http.Request, resources []interface{}, resource_name string, page, page_size, totalItemsCount int) map[string]interface{} {

	response := CreatePageMeta(request, len(resources), page, page_size, totalItemsCount)
	response[resource_name] = resources
	return response
}

func CreateDetailedErrorDto(key string, err error) map[string]interface{} {
	return map[string]interface{}{
		"success":       false,
		"full_messages": []string{fmt.Sprintf("s -> %v", key, err.Error())},
		"errors":        err,
	}
}

func CreateErrorDtoWithMessage(message string) map[string]interface{} {
	return map[string]interface{}{
		"success":       false,
		"full_messages": []string{message},
	}
}

// This should only be called when we have an Error that is returned from a ShouldBind which contains a lot of information
// other kind of errors should use other functions such as CreateDetailedErrorDto
func CreateBadRequestErrorDto(err error) ErrorDto {
	res := ErrorDto{}
	res.Errors = make(map[string]interface{})
	errs := err.(validator.ValidationErrors)
	res.FullMessages = make([]string, len(errs))
	count := 0
	for _, v := range errs {
		if v.ActualTag == "required" {
			var message = fmt.Sprintf("%v is required", v.Field)
			res.Errors[v.Field] = message
			res.FullMessages[count] = message
		} else {
			var message = fmt.Sprintf("%v has to be %v", v.Field, v.ActualTag)
			res.Errors[v.Field] = message
			res.FullMessages = append(res.FullMessages, message)
		}
		count++
	}
	return res
}

func CreateSuccessDto(result map[string]interface{}) map[string]interface{} {
	result["success"] = true
	return result
}

func CreateSuccessWithMessageDto(message string) interface{} {
	return CreateSuccessWithMessagesDto([]string{message})
}

func CreateSuccessWithMessagesDto(messages []string) interface{} {
	return gin.H{
		"success":       true,
		"full_messages": messages,
	}
}

func CreateSuccessWithDtoAndMessagesDto(data map[string]interface{}, messages []string) map[string]interface{} {
	data["success"] = true
	data["full_messages"] = messages
	return data
}
func CreateSuccessWithDtoAndMessageDto(data map[string]interface{}, message string) map[string]interface{} {
	return CreateSuccessWithDtoAndMessagesDto(data, []string{message})
}
