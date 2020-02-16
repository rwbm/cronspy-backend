package db

import "cronspy/backend/pkg/util/model"

// GetJobs returns the list of jobs for a user
func (c *JobDB) GetJobs(idUser int, pageSize, page int) (jobs []model.Job, p model.Pagination, err error) {

	offset := 0
	if page > 1 {
		offset = ((page - 1) * pageSize)
	}

	// get total records
	totalRecords := 0
	if err = c.ds.Model(model.Job{}).Where("id_user = ?", idUser).Count(&totalRecords).Error; err != nil {
		return
	}

	// if page requested is invalid, return no results
	if page > totalRecords/pageSize {
		return
	}

	// get jobs
	q := c.ds.Model(model.Job{}).Where("id_user = ?", idUser).Order("date_created asc")
	q = q.Offset(offset).Limit(pageSize)
	err = q.Find(&jobs).Error

	if err == nil {
		p.Page = page
		p.PageSize = pageSize
		p.TotalRows = totalRecords
	}

	return
}
