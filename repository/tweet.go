package repository

import (
	"errors"
	"fmt"

	"github.com/huandu/go-sqlbuilder"

	"github.com/svakode/svachan/dictionary"
)

// AddStreams is the repository for inserting twitter id and channel id to twitter_streams table
func (r *repository) AddStream(twitterID, twitterUsername, channelID string) (err error) {
	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ib.InsertInto("twitter_streams")
	ib.Cols("twitter_id", "twitter_username", "channel_id")
	ib.Values(twitterID, twitterUsername, channelID)

	query, args := ib.Build()
	res, err := r.db.Exec(query, args...)
	if err != nil {
		return errors.New(dictionary.DBError)
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected <= 0 {
		return fmt.Errorf(dictionary.AlreadyStreamingError, twitterID)
	}

	return
}

// RemoveStream is the repository for deleting twitter id and channel id from twitter_streams table
func (r *repository) RemoveStream(twitterID, channelID string) (err error) {
	delB := sqlbuilder.PostgreSQL.NewDeleteBuilder()
	delB.DeleteFrom("twitter_streams")
	delB.Where(
		delB.Equal("twitter_id", twitterID),
		delB.Equal("channel_id", channelID),
	)

	query, args := delB.Build()
	res, err := r.db.Exec(query, args...)
	if err != nil {
		return errors.New(dictionary.DBError)
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected <= 0 {
		return fmt.Errorf(dictionary.StreamingNotFoundError, twitterID)
	}

	return
}

// GetStreams is the repository for getting streams from DB
func (r *repository) GetStreams() (res map[string][]string, err error) {
	var twitterID, channelID string
	res = make(map[string][]string)

	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("twitter_id", "channel_id")
	sb.From("twitter_streams")
	query, args := sb.Build()

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, errors.New(dictionary.DBError)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&twitterID, &channelID)
		if err != nil {
			return nil, errors.New(dictionary.DBError)
		}
		res[twitterID] = append(res[twitterID], channelID)
	}

	return
}

func (r *repository) GetStreamsByChannel(channelID string) (res []string, err error) {
	var username string

	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("twitter_username")
	sb.From("twitter_streams")
	sb.Where(
		sb.Equal("channel_id", channelID),
	)
	query, args := sb.Build()

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, errors.New(dictionary.DBError)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&username)
		if err != nil {
			return nil, errors.New(dictionary.DBError)
		}
		res = append(res, username)
	}

	return
}
