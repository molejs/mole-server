package mole

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
	"time"
)

func getDB(c *gin.Context) *mgo.Collection {
	return c.MustGet("db").(*mgo.Collection)
}

func intQuery(c *gin.Context, key string, defaultValue int) int {
	v := c.DefaultQuery(key, "")
	if intVal, err := strconv.Atoi(v); err == nil {
		return intVal
	}
	return defaultValue
}

func ReportHandler(c *gin.Context) {
	var (
		db       = getDB(c)
		log      Log
		status   = http.StatusOK
		response gin.H
	)

	if c.BindJSON(&log) == nil {
		log.CreatedAt = time.Now()
		log.Id = bson.NewObjectId()
		if err := db.Insert(log); err != nil {
			status = http.StatusInternalServerError
			response = gin.H{
				"error": true,
				"msg":   err.Error(),
			}
		} else {
			response = gin.H{"error": false}
		}
	} else {
		status = http.StatusBadRequest
		response = gin.H{
			"error": true,
			"msg":   "missing required fields",
		}
	}

	c.JSON(status, response)
}

func RetrieveHandler(c *gin.Context) {
	var (
		db    = getDB(c)
		skip  = intQuery(c, "skip", 0)
		limit = intQuery(c, "limit", 25)
		logs  = []Log{}
	)

	count, err := db.Find(bson.M{}).Count()
	err2 := db.Find(bson.M{}).Sort("-created_at").Skip(skip).Limit(limit).All(&logs)

	if err != nil || err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   "internal server error",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"error": false,
			"logs":  logs,
			"total": count,
			"count": len(logs),
		})
	}
}
