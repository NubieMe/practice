package pagination

import (
	"gorm.io/gorm"
)

type PaginationParams struct {
	Limit int
	Page  int
	Sort  string
}

type Paginator interface {
	GetLimit() int
	GetPage() int
	GetSort() string
	GetOffset() int
	Apply(db *gorm.DB) *gorm.DB
}

type Pagination struct {
	Limit      int    `json:"limit,omitempty"`
	Page       int    `json:"page,omitempty"`
	Sort       string `json:"sort,omitempty"`
	TotalRows  int64  `json:"total_rows,omitempty"`
	TotalPages int    `json:"total_pages,omitempty"`
}

func NewPagination(params *PaginationParams) *Pagination {
	p := &Pagination{
		Limit: 10,
		Page:  1,
		Sort:  "created_at desc",
	}

	if params != nil {
		if params.Limit > 0 {
			p.Limit = params.Limit
		}
		if params.Page > 0 {
			p.Page = params.Page
		}
		if params.Sort != "" {
			p.Sort = params.Sort
		}
	}

	return p
}

func (p *Pagination) GetLimit() int {
	return p.Limit
}

func (p *Pagination) GetPage() int {
	return p.Page
}

func (p *Pagination) GetSort() string {
	return p.Sort
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) Apply(db *gorm.DB) *gorm.DB {
	return db.Offset(p.GetOffset()).Limit(p.GetLimit()).Order(p.GetSort())
}

func Paginate(model interface{}, p *Pagination, db *gorm.DB) (*gorm.DB, error) {
	var totalRows int64
	query := db.Model(model)

	// count total records
	if err := query.Count(&totalRows).Error; err != nil {
		return db, err
	}

	p.TotalRows = totalRows
	p.TotalPages = int((totalRows + int64(p.Limit) - 1) / int64(p.Limit))

	return p.Apply(query), nil
}
