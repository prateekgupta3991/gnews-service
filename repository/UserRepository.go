package repository

import (
	"fmt"

	"github.com/gocql/gocql"
	"github.com/prateekgupta3991/refresher/entities"
)

type UserRepo interface {
	GetUserByTgDetils(tgmId int64, tgmUn string) (entities.UserDetails, error)
	InsertUser(m entities.UserDetails) error
	GetAllUser() ([]entities.UserDetails, error)
	GetUserByTgUn(tgmUn string) (entities.UserDetails, error)
}

type UserDbSession struct {
	DbClient *gocql.Session
}

func NewUserDbSession(cassession *gocql.Session) *UserDbSession {
	return &UserDbSession{
		DbClient: cassession,
	}
}

func (c *UserDbSession) GetUserByTgDetils(tgmId int64, tgmUn string) (entities.UserDetails, error) {
	m := map[string]interface{}{}
	query := fmt.Sprintf("SELECT uid, name, t_un, chat_id from user where uid = %d and t_un = '%s'", tgmId, tgmUn)
	fmt.Println(query)
	iter := c.DbClient.Query(query).Consistency(gocql.One).Iter()
	var subscriber entities.UserDetails
	for iter.MapScan(m) {
		if id, ok := m["uid"].(int64); ok {
			if cid, ok := m["chat_id"].(int64); ok {
				subscriber = entities.UserDetails{
					ID:         int64(id),
					Name:       fmt.Sprintf("%v", m["name"]),
					TelegramId: fmt.Sprintf("%v", m["t_un"]),
					ChatId:     int64(cid),
				}
			}
		}
		m = map[string]interface{}{}
	}
	return subscriber, nil
}

func (c *UserDbSession) InsertUser(m entities.UserDetails) error {
	query := "insert into user(uid, name, t_un, chat_id) values (?,?,?,?)"
	if err := c.DbClient.Query(query, int64(m.ID), m.Name, m.TelegramId, m.ChatId).Consistency(gocql.One).Exec(); err != nil {
		fmt.Errorf("Error encountered : %s", err.Error())
		return err
	}
	return nil
}

func (c *UserDbSession) GetAllUser() ([]entities.UserDetails, error) {
	m := map[string]interface{}{}
	query := "SELECT uid, name, t_un, chat_id from user"
	iter := c.DbClient.Query(query).Consistency(gocql.One).Iter()
	fmt.Println("Query is done")
	defer iter.Close()
	var subscribers []entities.UserDetails
	for iter.MapScan(m) {
	    fmt.Println(m)
		if id, ok := m["uid"].(int64); ok {
			if cid, ok := m["chat_id"].(int64); ok {
				subscribers = append(subscribers, entities.UserDetails{
					ID:         int64(id),
					Name:       fmt.Sprintf("%v", m["name"]),
					TelegramId: fmt.Sprintf("%v", m["t_un"]),
					ChatId:     int64(cid),
				})
				m = map[string]interface{}{}
			} else {
			    fmt.Println("invalid chat id")
			}
		}
	}
	return subscribers, nil
}

func (c *UserDbSession) GetUserByTgUn(tgmUn string) (entities.UserDetails, error) {
	m := map[string]interface{}{}
	query := fmt.Sprintf("SELECT uid, name, t_un, chat_id from user where t_un = '%s' ALLOW FILTERING", tgmUn)
	fmt.Println(query)
	iter := c.DbClient.Query(query).Consistency(gocql.One).Iter()
	var subscriber entities.UserDetails
	for iter.MapScan(m) {
		if id, ok := m["uid"].(int64); ok {
			if cid, ok := m["chat_id"].(int64); ok {
				subscriber = entities.UserDetails{
					ID:         int64(id),
					Name:       fmt.Sprintf("%v", m["name"]),
					TelegramId: fmt.Sprintf("%v", m["t_un"]),
					ChatId:     int64(cid),
				}
			}
		}
		m = map[string]interface{}{}
	}
	return subscriber, nil
}
