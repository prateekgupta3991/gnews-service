package repository

import (
	// "time"

	"github.com/gocql/gocql"
	"github.com/prateekgupta3991/refresher/cassandra"
	"github.com/prateekgupta3991/refresher/entities"
)

type NewsRepo interface {
	GetTopNews(sid string, lim int) ([]entities.NewsBySource, error)
	InsertTopNews(m entities.NewsBySource) error
}

type DbSession struct {
	DbClient *gocql.Session
}

func GetNewDbSession() *DbSession {
	return &DbSession{
		DbClient: cassandra.Session,
	}
}

func (c *DbSession) GetTopNews(sid string, lim int) ([]entities.NewsBySource, error) {
	m := map[string]interface{}{}
	query := "SELECT sid, created_at, nauthor, ncontent, ndesc, sname from news_by_source where sid = ? limit ?"
	iter := c.DbClient.Query(query, sid, lim).Consistency(gocql.One).Iter()
	var newsEnt []entities.NewsBySource
	for iter.MapScan(m) {
		newsEnt = append(newsEnt, entities.NewsBySource{
			SourceId: sid,
			// CreatedAt:       m["created_at"].(time.Time),
			NewsAuthor:      m["nauthor"].(string),
			NewsContent:     m["ncontent"].(string),
			NewsDescription: m["ndesc"].(string),
			SourceName:      m["sname"].(string),
		})
		m = map[string]interface{}{}
	}
	return newsEnt, nil
}

func (c *DbSession) InsertTopNews(m entities.NewsBySource) error {
	query := "insert into news_by_source(sid, created_at, title_hash, nauthor, ncontent, ndesc, sname, npublished_at) values (?,?,?,?,?,?,?,?)"
	err := c.DbClient.Query(query, m.SourceId, gocql.TimeUUID(), m.TitleHash, m.NewsAuthor, m.NewsContent, m.NewsDescription, m.SourceName, m.NewsPublishedAt).Consistency(gocql.One).Exec()
	if err != nil {
		return err
	}
	return nil
}
