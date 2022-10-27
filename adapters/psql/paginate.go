package psql

import (
	"fmt"

	"github.com/go-pg/pg/v10/orm"
)

type Paginate struct {
	PageIndex int
	PageSize  int
	OrderBy   string
	Order     string
}

func (p Paginate) calcOffset() int {
	return p.PageSize * (p.PageIndex - 1)
}

func (p Paginate) toQuery() func(q *orm.Query) (*orm.Query, error) {
	return func(q *orm.Query) (*orm.Query, error) {
		q = q.Limit(p.PageSize).
			Offset(p.calcOffset()).
			Order(fmt.Sprintf("%s %s", p.OrderBy, p.Order))

		return q, nil
	}
}
