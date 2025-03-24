package models

import (
	"database/sql"
	"errors"
	"time"
)

// Define a Message type to hold the data for an individual message. Notice how
// the fields of the struct correspond to the fields in our MySQL snippets
// table?
type Message struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// Define a MessageModel type which wraps a sql.DB connection pool.
type MessageModel struct {
	DB *sql.DB
}

// This will insert a new snippet into the database.
func (m *MessageModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO messages (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	result, err := m.DB.Exec(stmt, title, content, expires)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// This will return a specific snippet based on its id.
func (m *MessageModel) Get(id int) (*Message, error) {
	stmt := `SELECT id, title, content, created, expires FROM messages
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(stmt, id)

	s := &Message{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil

}

// This will return the 10 most recently created snippets.
func (m *MessageModel) Latest() ([]*Message, error) {
	stmt := `SELECT id, title, content, created, expires FROM messages
	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	messages := []*Message{}

	for rows.Next() {
		s := &Message{}

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		messages = append(messages, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
