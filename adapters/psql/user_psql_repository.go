package psql

import (
	"context"
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type User struct {
	//nolint:structcheck,gocritic,unused
	tableName struct{} `sql:"users"`
	ID        string   `sql:"id,notnull"`
	UserName  string   `sql:"user_name,notnull"`
	LastName  string   `sql:"last_name,notnull"`
	FirstName string   `sql:"first_name,notnull"`
	Password  string   `sql:"password,notnull"`
	Email     string   `sql:"email,notnull"`
	Mobile    string   `sql:"mobile,notnull"`
	ASZF      bool     `sql:"aszf,notnull" pg:",use_zero"`
}

type UserListFilter struct {
	Email *string
}

func (f UserListFilter) toQuery() func(q *orm.Query) (*orm.Query, error) {
	return func(q *orm.Query) (*orm.Query, error) {
		if f.Email != nil {
			q = q.Where("email ILIKE ?", fmt.Sprintf("%%%s%%", *f.Email))
		}

		return q, nil
	}
}

type UserPSQLRepository struct {
	db *pg.DB
}

func NewUserPSQLRepository(db *pg.DB) *UserPSQLRepository {
	return &UserPSQLRepository{db: db}
}

func (r *UserPSQLRepository) GetUserCount(
	ctx context.Context) (int, error) {
	user := new(User)
	count, err := r.db.WithContext(ctx).
		Model(user).
		Count()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *UserPSQLRepository) InsertUser(ctx context.Context, user *User) error {
	_, err := r.db.WithContext(ctx).
		Model(user).
		Insert()
	if err != nil {
		return err
	}

	return nil
}

func (r *UserPSQLRepository) UpdateUser(ctx context.Context, user *User) error {
	_, err := r.db.WithContext(ctx).
		Model(user).
		Column("user_name", "last_name", "first_name", "password", "mobile").
		WherePK().
		Returning("*").
		Update()
	if err != nil {
		return err
	}

	return nil
}

func (r *UserPSQLRepository) ListUser(ctx context.Context, paginate *Paginate, filter *UserListFilter) ([]User, int, error) {
	var users []User

	count, err := r.db.WithContext(ctx).
		Model(&users).
		Apply(filter.toQuery()).
		Apply(paginate.toQuery()).
		SelectAndCount()
	if err != nil {
		return users, count, nil
	}

	return users, count, nil
}
