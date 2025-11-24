package repositories

import (
	"context"

	"it_rabotyagi/internal/data/database"
)

type MentorRepository struct {
	db *database.DB
}

func NewMentorRepository(db *database.DB) *MentorRepository {
	return &MentorRepository{db: db}
}

type MentorProfile struct {
	ID              int
	UserID          int
	Specialization  string
	Grade           *string
	ExperienceYears *int
	Description     *string
	Tags            []string
	Contacts        map[string]any
	Pricelist       map[string]any
}

func (r *MentorRepository) ListMentors(ctx context.Context, specialization *string, limit, offset int) ([]MentorProfile, int, error) {
	where := ""
	args := []any{}
	if specialization != nil && *specialization != "" {
		where = "WHERE specialization = $1"
		args = append(args, *specialization)
	}

	countQuery := "SELECT COUNT(*) FROM mentors " + where
	var total int
	if err := r.db.Pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := "SELECT id, user_id, specialization, grade, experience_years, description, tags, contacts, pricelist FROM mentors " + where + " ORDER BY id"
	if limit > 0 {
		query += " LIMIT $%d OFFSET $%d"
		if len(args) == 0 {
			query = "SELECT id, user_id, specialization, grade, experience_years, description, tags, contacts, pricelist FROM mentors ORDER BY id LIMIT $1 OFFSET $2"
			args = append(args, limit, offset)
		} else {
			args = append(args, limit, offset)
			// adjust placeholders if specialization used
			query = "SELECT id, user_id, specialization, grade, experience_years, description, tags, contacts, pricelist FROM mentors WHERE specialization = $1 ORDER BY id LIMIT $2 OFFSET $3"
		}
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var mentors []MentorProfile
	for rows.Next() {
		var m MentorProfile
		var tags []string
		var contacts map[string]any
		var pricelist map[string]any
		if err := rows.Scan(&m.ID, &m.UserID, &m.Specialization, &m.Grade, &m.ExperienceYears, &m.Description, &tags, &contacts, &pricelist); err != nil {
			return nil, 0, err
		}
		m.Tags = tags
		m.Contacts = contacts
		m.Pricelist = pricelist
		mentors = append(mentors, m)
	}

	return mentors, total, nil
}

func (r *MentorRepository) GetMentor(ctx context.Context, id int) (*MentorProfile, error) {
	query := `SELECT id, user_id, specialization, grade, experience_years, description, tags, contacts, pricelist FROM mentors WHERE id = $1`
	var m MentorProfile
	var tags []string
	var contacts map[string]any
	var pricelist map[string]any
	if err := r.db.Pool.QueryRow(ctx, query, id).Scan(&m.ID, &m.UserID, &m.Specialization, &m.Grade, &m.ExperienceYears, &m.Description, &tags, &contacts, &pricelist); err != nil {
		return nil, err
	}
	m.Tags = tags
	m.Contacts = contacts
	m.Pricelist = pricelist
	return &m, nil
}
