package handlers

import (
	"context"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/k-digital-project/mini-api/db/repository/models"
	"github.com/k-digital-project/mini-api/usecase"
	"github.com/k-digital-project/mini-api/usecase/requests"
)

// UserHandler ...
type UserHandler struct {
	Handler
}

func (h *UserHandler) EditUser(ctx *fiber.Ctx) error {
	c := ctx.Locals("ctx").(context.Context)
	id := ctx.Locals("user_id").(string)

	input := new(requests.UserRequest)
	if err := ctx.BodyParser(input); err != nil {
		return h.SendResponse(ctx, nil, nil, err, http.StatusBadRequest)
	}
	if err := h.Validator.Struct(input); err != nil {
		errMessage := h.ExtractErrorValidationMessages(err.(validator.ValidationErrors))
		return h.SendResponse(ctx, nil, nil, errMessage, http.StatusBadRequest)
	}
	userUC := usecase.UserUC{ContractUC: h.ContractUC}
	res, err := userUC.Edit(c, id, input)

	return h.SendResponse(ctx, res, nil, err, 0)
}

func (h *UserHandler) FindAllUser(ctx *fiber.Ctx) error {
	c := context.Background()
	param := models.UserParamater{
		Search: ctx.Query("search"),
		Status: ctx.Query("status"),
	}

	uc := usecase.UserUC{ContractUC: h.ContractUC}
	res, err := uc.FindAllUser(c, param)

	return h.SendResponse(ctx, res, nil, err, 0)
}
