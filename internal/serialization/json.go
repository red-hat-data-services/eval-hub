package serialization

import (
	"encoding/json"

	validator "github.com/go-playground/validator/v10"
	"github.com/julpayne/eval-hub/internal/executioncontext"
)

func Unmarshal(validate *validator.Validate, executionContext *executioncontext.ExecutionContext, jsonBytes []byte, v any) error {
	err := json.Unmarshal(jsonBytes, v)
	if err != nil {
		return err
	}
	// now validate the unmarshalled data
	err = validate.StructCtx(executionContext.Ctx, v)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, validationError := range validationErrors {
			// TODO: add the validation error to the response?
			executionContext.Logger.Info("Validation error", "field", validationError.Field(), "tag", validationError.Tag(), "value", validationError.Value())
		}
		return err
	}
	// if the validation is successful, return nil
	return nil
}
