package timeuuid

import "github.com/gocql/gocql"

func NewString() string {
	return gocql.TimeUUID().String()
}
