package repository

import (
	"boilerplate/model"
	"boilerplate/server/response"
	"boilerplate/utils/config"
	"boilerplate/utils/database"
	cerror "boilerplate/utils/error"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type Repository struct {
	config     *config.Config
	db         *database.Database
	pagination *response.Pagination
}

func NewRepository(cfg *config.Config, db *database.Database, pagination *response.Pagination) *Repository {
	return &Repository{
		config:     cfg,
		db:         db,
		pagination: pagination,
	}
}

func (repo *Repository) Paginate(db *gorm.DB) *gorm.DB {
	page := repo.pagination.Page

	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * repo.pagination.PageSize

	return db.Offset(offset).Limit(repo.pagination.PageSize)
}

func (repo *Repository) GetTotalRecord(db *gorm.DB, countFieldName string) (uint64, error) {
	var dummyRes []map[string]interface{}
	var sql strings.Builder
	additionalVars := 0

	stmt := db.Session(&gorm.Session{DryRun: true}).Find(&dummyRes).Statement
	stmtSql := stmt.SQL.String()
	stmtVars := stmt.Vars
	totalVars := len(stmt.Vars)

	indexFrom := strings.Index(stmtSql, "FROM")
	indexLimit := strings.Index(stmtSql, "LIMIT")
	indexOffset := strings.Index(stmtSql, "OFFSET")
	if indexFrom > -1 {
		if indexLimit > -1 {
			additionalVars++
			stmtSql = stmtSql[indexFrom:indexLimit]
		} else {
			stmtSql = stmtSql[indexFrom:]
		}

		if indexOffset > -1 {
			additionalVars++
		}

		totalVars -= additionalVars

		if totalVars > 0 {
			stmtVars = stmtVars[0:(totalVars)]
		}

		sql.WriteString(fmt.Sprintf("SELECT COUNT(%s) AS total_rows ", countFieldName))
		sql.WriteString(stmtSql)
	}

	var rows *model.Base
	output := db.Session(&gorm.Session{}).Raw(sql.String(), stmtVars...).Scan(&rows)

	if output.Error != nil {
		return 0, cerror.Fail(cerror.FuncName(), "failed_db_query", nil, output.Error)
	}

	return rows.TotalRows, nil
}
