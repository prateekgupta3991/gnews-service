package repository

import (
	"fmt"

	"github.com/gocql/gocql"
	"github.com/prateekgupta3991/refresher/cassandra"
	"github.com/prateekgupta3991/refresher/entities"
)

type NewsRepo interface {
	GetTopNewsBySource(sid string, lim int) ([]entities.NewsBySource, error)
	InsertTopNews(m entities.NewsBySource) error
	InsertSources(m entities.NewsBySource) error
}

type DbSession struct {
	DbClient *gocql.Session
}

func GetNewDbSession() *DbSession {
	return &DbSession{
		DbClient: cassandra.Session,
	}
}

func (c *DbSession) GetTopNewsBySource(sid string, lim int) ([]entities.NewsBySource, error) {
	m := map[string]interface{}{}
	query := fmt.Sprintf("SELECT sid, created_at, title_hash, nauthor, ncontent, ndesc, npublished_at, ntitle, nurl, nurl_to_image, sname, scategory, sdesc, scountry, slang, surl from news_by_source where sid in (%s) limit ?", sid)
	iter := c.DbClient.Query(query, lim).Consistency(gocql.One).Iter()
	var newsEnt []entities.NewsBySource
	for iter.MapScan(m) {
		newsEnt = append(newsEnt, entities.NewsBySource{
			SourceId: fmt.Sprintf("%v", m["sid"]),
			// CreatedAt:       m["created_at"].(time.Time),
			NewsAuthor:        fmt.Sprintf("%v", m["nauthor"]),
			NewsContent:       fmt.Sprintf("%v", m["ncontent"]),
			NewsDescription:   fmt.Sprintf("%v", m["ndesc"]),
			NewsPublishedAt:   fmt.Sprintf("%v", m["npublished_at"]),
			NewsTitle:         fmt.Sprintf("%v", m["ntitle"]),
			NewsUrl:           fmt.Sprintf("%v", m["nurl"]),
			NewsUrlToImage:    fmt.Sprintf("%v", m["nurl_to_image"]),
			SourceName:        fmt.Sprintf("%v", m["sname"]),
			SourceCategory:    fmt.Sprintf("%v", m["ncategory"]),
			SourceDescription: fmt.Sprintf("%v", m["sdesc"]),
			SourceCountry:     fmt.Sprintf("%v", m["scountry"]),
			SourceLanguage:    fmt.Sprintf("%v", m["slang"]),
			SourceUrl:         fmt.Sprintf("%v", m["surl"]),
		})
		m = map[string]interface{}{}
	}
	return newsEnt, nil
}

func (c *DbSession) InsertTopNews(m entities.NewsBySource) error {
	query := "insert into news_by_source(sid, created_at, title_hash, nauthor, ncontent, ndesc, sname, npublished_at) values (?,?,?,?,?,?,?,?)"
	if err := c.DbClient.Query(query, m.SourceId, gocql.TimeUUID(), m.TitleHash, m.NewsAuthor, m.NewsContent, m.NewsDescription, m.SourceName, m.NewsPublishedAt).Consistency(gocql.One).Exec(); err != nil {
		fmt.Errorf("Error encountered : %s", err.Error())
		return err
	}
	return nil
}

func (c *DbSession) InsertSources(m entities.NewsBySource) error {
	query := "insert into news_sources(sid, created_at, sname, sdesc, surl, scategory, slang, scountry) values (?,?,?,?,?,?,?,?)"
	if err := c.DbClient.Query(query, m.SourceId, gocql.TimeUUID(), m.SourceName, m.SourceDescription, m.SourceUrl, m.SourceCategory, m.SourceLanguage, m.SourceCountry).Consistency(gocql.One).Exec(); err != nil {
		fmt.Errorf("Error encountered : %s", err.Error())
		return err
	}
	return nil
}
