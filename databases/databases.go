package database

import (
	"errors"

	"gopkg.in/mgo.v2"
)

var (
	DBName string
)

const (
	ERRConnectionstr = "connection string must be provided"
)

//GetConnection is to get database connection
func GetConnection(connections ...string) (interface{}, error) {
	if len(connections) < 1 {
		return nil, errors.New(ERRConnectionstr)
	}
	session, err := mgo.Dial(connections[0])
	if err != nil {
		return nil, err
	}
	return session, nil
}
