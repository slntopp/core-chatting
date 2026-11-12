package graph

import (
	"context"
	"fmt"
)

// AccountSpending is what NoCloud billing stores in account.spending.
type AccountSpending struct {
	Amount     float64 `json:"amount"`
	PeriodDays int32   `json:"period_days"`
	Currency   struct {
		Code  string `json:"code"`
		Title string `json:"title"`
	} `json:"currency"`
}

// AccountGradeInfo is the account data a grade is derived from.
type AccountGradeInfo struct {
	Uuid          string           `json:"uuid"`
	Spending      *AccountSpending `json:"spending"`
	GradeOverride string           `json:"grade_override"`
}

const accountGradeInfoQuery = `
FOR a IN @@accounts
  FILTER a._key IN @keys
  FILTER a.spending != null || (a.grade_override != null && a.grade_override != "")
  RETURN { uuid: a._key, spending: a.spending, grade_override: a.grade_override }
`

// NoCloud keeps currencies in the same database as accounts.
const defaultCurrencyQuery = `
FOR c IN Currencies
  FILTER c.default == true
  LIMIT 1
  RETURN c.code
`

// GradeInfo reads spending and the manual grade of the given accounts, and the
// default currency code ("" when it can't be read). Accounts with neither are
// left out.
func (c *UsersController) GradeInfo(ctx context.Context, uuids []string) (map[string]AccountGradeInfo, string, error) {
	cur, err := c.db.Query(ctx, accountGradeInfoQuery, map[string]interface{}{
		"@accounts": c.colname,
		"keys":      uuids,
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to read account grades: %w", err)
	}
	defer cur.Close()

	infos := map[string]AccountGradeInfo{}
	for cur.HasMore() {
		var info AccountGradeInfo
		if _, err := cur.ReadDocument(ctx, &info); err != nil {
			return nil, "", fmt.Errorf("failed to read account grade: %w", err)
		}
		infos[info.Uuid] = info
	}
	if len(infos) == 0 {
		return infos, "", nil
	}

	code := ""
	if ccur, err := c.db.Query(ctx, defaultCurrencyQuery, nil); err == nil {
		if ccur.HasMore() {
			_, _ = ccur.ReadDocument(ctx, &code)
		}
		_ = ccur.Close()
	}
	return infos, code, nil
}
