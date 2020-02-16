package db

import "cronspy/backend/pkg/util/model"

// GetJobs returns the list of jobs for a user
func (c *JobDB) GetJobs(idUser int, count, offset int) (jobs []model.Job, err error) {
	q := c.ds.Model(model.Job{}).Where("id_user = ?", idUser).Order("date_created asc")

	if count > 0 {
		q = q.Limit(count)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}

	err = q.Find(&jobs).Error
	return
}
