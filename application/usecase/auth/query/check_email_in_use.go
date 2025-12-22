package query

import (
	"context"

	"src/core/cqrs"
	"src/core/doc"
	"src/core/validator"
	"src/domain/exception"
	"src/domain/repository"
)

// #region errors

const (
	Err_CheckEmailInUseQuery_EmailInUse = "Email in use"
)

// #endregion
// #region CheckEmailInUseQuery

type CheckEmailInUseQuery struct {
	Email string `json:"email"`
}

var _ validator.IValidable = (*CheckEmailInUseQuery)(nil)

func (q *CheckEmailInUseQuery) Validate() error {
	return validator.Object(q,
		validator.String(&q.Email).Required().Email(),
	).Validate()
}

// #endregion
// #region CheckEmailInUseResult

type CheckEmailInUseResult struct {
	InUse bool `json:"in_use"`
}

// #endregion
// #region CheckEmailInUseQueryHandler

type CheckEmailInUseQueryHandler struct {
	accountRepository repository.IAccountRepository
}

var _ cqrs.IQueryHandler[
	*CheckEmailInUseQuery,
	*CheckEmailInUseResult,
] = (*CheckEmailInUseQueryHandler)(nil)

func NewCheckEmailInUseQueryHandler(
	accountRepository repository.IAccountRepository,
) *CheckEmailInUseQueryHandler {
	return &CheckEmailInUseQueryHandler{
		accountRepository: accountRepository,
	}
}

func (h *CheckEmailInUseQueryHandler) Handle(
	ctx context.Context,
	query *CheckEmailInUseQuery,
) (*CheckEmailInUseResult, error) {
	size, err := h.accountRepository.CountByEmail(ctx, query.Email)
	var result CheckEmailInUseResult
	if err != nil {
		return nil, err
	}
	if size > 0 {
		result.InUse = true
	}
	return &result, nil
}

// #endregion

func init() {
	query := CheckEmailInUseQuery{
		Email: "john.doe@email.com",
	}
	doc.Describe(&query,
		doc.Description("Check email in use query"),
		doc.Example(&query),
		doc.Throws[exception.ConflictException](Err_CheckEmailInUseQuery_EmailInUse))

	result := CheckEmailInUseResult{}
	doc.Describe(&result,
		doc.Description("Check email in use result"),
		doc.Example(&result),
		doc.Field(&result.InUse, doc.Description("Indicates if the email is in use")))

	cqrs.RegisterQueryHandler[
		*CheckEmailInUseQuery,
		*CheckEmailInUseResult,
		*CheckEmailInUseQueryHandler,
	](NewCheckEmailInUseQueryHandler)
}
