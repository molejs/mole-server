package mole

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

func DatabaseMiddleware(addr, db string) gin.HandlerFunc {
	session, err := mgo.Dial(addr)
	if err != nil {
		panic(err)
	}

	return func(c *gin.Context) {
		sessCopy := session.Copy()
		defer sessCopy.Close()

		c.Set("db", sessCopy.DB(db).C("logs"))
		c.Next()
	}
}
