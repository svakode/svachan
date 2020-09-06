package repository

import (
	"database/sql"
	"errors"

	"github.com/huandu/go-sqlbuilder"

	"github.com/svakode/svachan/dictionary"
)

// SetEmail is the repository for inserting username and email to members table
func (r *repository) SetMemberEmail(username, email string) (err error) {
	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ib.InsertInto("members")
	ib.Cols("username", "email")
	ib.Values(username, email)

	query, args := ib.Build()
	res, err := r.db.Exec(query, args...)
	if err != nil {
		return errors.New(dictionary.DBError)
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected > 0 {
		return
	}

	ub := sqlbuilder.PostgreSQL.NewUpdateBuilder()
	ub.Update("members")
	ub.Set("email", email)
	ub.Where(
		ub.Equal("username", username),
	)

	query, args = ub.Build()
	res, err = r.db.Exec(query, args...)
	if err != nil {
		return errors.New(dictionary.DBError)
	}

	return
}

// GetMemberEmail is the repository for getting email from members table for a specific user
func (r *repository) GetMemberEmail(username string) (res string, err error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("email")
	sb.From("members")
	sb.Where(sb.Equal("username", username))

	query, args := sb.Build()
	row := r.db.QueryRow(query, args...)
	err = row.Scan(&res)
	if err == sql.ErrNoRows {
		return "", errors.New(dictionary.MemberNotFoundMessage)
	} else if err != nil {
		return "", errors.New(dictionary.DBError)
	}

	return
}

// GetMembersEmail is the repository for getting email from members table for more than one user
func (r *repository) GetMembersEmail(usernames []interface{}) (res map[string]string, err error) {
	var username, email string
	res = make(map[string]string)

	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("username", "email")
	sb.From("members")
	sb.Where(sb.In("username", usernames...))

	query, args := sb.Build()
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return res, errors.New(dictionary.DBError)
	}

	for rows.Next() {
		err := rows.Scan(&username, &email)
		if err != nil {
			return res, errors.New(dictionary.DBError)
		}

		res[username] = email
	}

	if len(res) == 0 {
		return res, errors.New(dictionary.MemberNotFoundMessage)
	}

	return
}
