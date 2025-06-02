package dtos

type PagerDto struct {
	Page           int `form:"page"`
	RecordsPerPage int `form:"records_per_page"`
}
