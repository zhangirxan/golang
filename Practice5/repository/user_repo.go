package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"practice5/models"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

var allowedColumns = map[string]bool{
	"id":        true,
	"name":      true,
	"email":     true,
	"gender":    true,
	"birthdate": true,
}

func (r *Repository) GetPaginatedUsers(p models.FilterParams) (models.PaginatedResponse, error) {
	args := []interface{}{}
	argIdx := 1
	whereClauses := []string{}

	if p.ID != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("id = $%d", argIdx))
		args = append(args, *p.ID)
		argIdx++
	}
	if p.Name != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("name ILIKE $%d", argIdx))
		args = append(args, "%"+*p.Name+"%")
		argIdx++
	}
	if p.Email != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("email ILIKE $%d", argIdx))
		args = append(args, "%"+*p.Email+"%")
		argIdx++
	}
	if p.Gender != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("gender = $%d", argIdx))
		args = append(args, *p.Gender)
		argIdx++
	}
	if p.Birthdate != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("birthdate = $%d", argIdx))
		args = append(args, *p.Birthdate)
		argIdx++
	}

	where := ""
	if len(whereClauses) > 0 {
		where = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	orderCol := "id" // default
	if p.OrderBy != "" && allowedColumns[p.OrderBy] {
		orderCol = p.OrderBy
	}
	orderDir := "ASC" // default
	if strings.ToUpper(p.OrderDir) == "DESC" {
		orderDir = "DESC"
	}

	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM users %s`, where)
	var total int
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return models.PaginatedResponse{}, err
	}

	offset := (p.Page - 1) * p.PageSize
	dataQuery := fmt.Sprintf(
		`SELECT id, name, email, gender, birthdate FROM users %s ORDER BY %s %s LIMIT $%d OFFSET $%d`,
		where, orderCol, orderDir, argIdx, argIdx+1,
	)
	dataArgs := append(args, p.PageSize, offset)

	rows, err := r.db.Query(dataQuery, dataArgs...)
	if err != nil {
		return models.PaginatedResponse{}, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.Birthdate); err != nil {
			return models.PaginatedResponse{}, err
		}
		users = append(users, u)
	}

	return models.PaginatedResponse{
		Data:       users,
		TotalCount: total,
		Page:       p.Page,
		PageSize:   p.PageSize,
	}, nil
}

func (r *Repository) GetCommonFriends(userID1, userID2 int) ([]models.User, error) {
	query := `
		SELECT u.id, u.name, u.email, u.gender, u.birthdate
		FROM user_friends uf1
		JOIN user_friends uf2 ON uf1.friend_id = uf2.friend_id
		JOIN users u          ON u.id = uf1.friend_id
		WHERE uf1.user_id = $1
		  AND uf2.user_id = $2
	`
	rows, err := r.db.Query(query, userID1, userID2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.Birthdate); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
