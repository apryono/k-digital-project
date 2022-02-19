package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/k-digital-project/mini-api/pkg/responsedto"
	"github.com/k-digital-project/mini-api/pkg/str"
	"github.com/k-digital-project/mini-api/usecase"
)

type Handler struct {
	FiberApp   *fiber.App
	ContractUC *usecase.ContractUC
	DB         *sql.DB
	Validator  *validator.Validate
	Translator ut.Translator
}

//SendResponse base send response
func (h Handler) SendResponse(ctx *fiber.Ctx, data interface{}, meta interface{}, err interface{}, code int) error {
	if code == 0 && err != nil {
		code = http.StatusUnprocessableEntity
		err = err.(error).Error()
	}

	if code != http.StatusOK && err != nil {
		return h.SendErrorResponse(ctx, err, code)
	}

	return h.SendSuccessResponse(ctx, data, meta)
}

//SendErrorResponse send response if status code != 200
func (h Handler) SendErrorResponse(ctx *fiber.Ctx, err interface{}, code int) error {
	response := responsedto.ErrorResponse(err)

	return ctx.Status(code).JSON(response)
}

//SendSuccessResponse send response if status code 200
func (h Handler) SendSuccessResponse(ctx *fiber.Ctx, data interface{}, meta interface{}) error {
	response := responsedto.SuccessResponse(data, meta)

	return ctx.Status(http.StatusOK).JSON(response)
}

//extract error message from validator
func (h Handler) ExtractErrorValidationMessages(error validator.ValidationErrors) map[string][]string {
	errorMessage := map[string][]string{}
	errorTranslation := error.Translate(h.Translator)

	for _, err := range error {
		errKey := str.Underscore(err.StructField())
		errorMessage[errKey] = append(
			errorMessage[errKey],
			strings.Replace(errorTranslation[err.Namespace()], err.StructField(), err.StructField(), -1),
		)
	}

	return errorMessage
}
