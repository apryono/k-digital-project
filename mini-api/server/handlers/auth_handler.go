package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/k-digital-project/mini-api/usecase"
	"github.com/k-digital-project/mini-api/usecase/requests"
)

//AuthHandler ...
type AuthHandler struct {
	Handler
}

func (h *AuthHandler) RegisterByEmail(ctx *fiber.Ctx) error {
	c := ctx.Context()

	input := new(requests.RegisterByEmailRequest)
	if err := ctx.BodyParser(input); err != nil {
		return h.SendResponse(ctx, nil, nil, err, http.StatusBadRequest)
	}

	if err := h.Validator.Struct(input); err != nil {
		errMessage := h.ExtractErrorValidationMessages(err.(validator.ValidationErrors))
		return h.SendResponse(ctx, nil, nil, errMessage, http.StatusBadRequest)
	}

	uc := usecase.AuthUC{ContractUC: h.ContractUC}
	res, err := uc.RegisterByEmail(c, input)
	if err != nil {
		return h.SendResponse(ctx, nil, nil, err.Error(), http.StatusBadRequest)
	}

	return h.SendResponse(ctx, res, nil, err, 0)
}
