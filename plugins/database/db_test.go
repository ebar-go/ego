package database

import (
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

type DBPluginSuite struct {
	suite.Suite
	db *gorm.DB
}

func (suite *DBPluginSuite) SetupTest() {
	// initialize..
	db, err := gorm.Open(sqlite.Open("/tmp/test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	suite.db = db
}

func (suite *DBPluginSuite) TestInstance() {
	Register(suite.db)
	suite.Equal(suite.db, Instance())
}

func TestDBPluginSuite(t *testing.T) {
	suite.Run(t, new(DBPluginSuite))
}
