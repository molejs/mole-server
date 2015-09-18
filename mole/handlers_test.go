package mole

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	. "gopkg.in/check.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

type HandlersSuite struct {
	db *mgo.Database
	r  *gin.Engine
}

func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&HandlersSuite{})

func (s *HandlersSuite) SetUpTest(c *C) {
	session, err := mgo.Dial("127.0.0.1:27017")
	c.Assert(err, IsNil)
	s.db = session.DB("unittest")
	gin.SetMode(gin.ReleaseMode)
	s.r = gin.Default()

	dbMiddleware := DatabaseMiddleware("127.0.0.1:27017", "unittest")
	s.r.POST("/logs", dbMiddleware, ReportHandler)
	s.r.GET("/logs", dbMiddleware, RetrieveHandler)

	logs := []Log{
		createLog(),
		createLog(),
		createLog(),
		createLog(),
		createLog(),
		createLog(),
	}

	for _, l := range logs {
		c.Assert(s.db.C("logs").Insert(l), IsNil)
	}

	go s.r.Run(":8769")
}

func (s *HandlersSuite) TestReportHandler(c *C) {
	log := createLog()
	log.Timestamp = "foo"
	log.Location = Location{
		"host", "href", "hash", "pathname", "port", "protocol", "search",
	}
	log.Error = Error{
		"foo msg",
		[]StracktraceLine{
			{"fn1", "file1", 1, 1},
			{"fn2", "file2", 2, 2},
		},
	}
	log.ActionStateHistory = []ActionStateHistory{
		{
			Action: bson.M{"action1": true},
			State:  bson.M{"state1": true},
		},
		{
			Action: bson.M{"action1": true},
			State:  bson.M{"state1": true},
		},
	}
	body, err := json.Marshal(log)
	c.Assert(err, IsNil)

	s.db.DropDatabase()

	response := request(c, string(body), "POST", 0, 0)
	c.Assert(response["error"].(bool), Equals, false)

	var resultLog Log
	c.Assert(s.db.C("logs").Find(bson.M{}).One(&resultLog), IsNil)

	c.Assert(resultLog.Timestamp, Equals, log.Timestamp)
	c.Assert(resultLog.Location, DeepEquals, log.Location)
	c.Assert(resultLog.Error, DeepEquals, log.Error)
	c.Assert(resultLog.ActionStateHistory, DeepEquals, log.ActionStateHistory)

	response = request(c, "", "POST", 0, 0)
	c.Assert(response["error"].(bool), Equals, true)
}

func (s *HandlersSuite) TestRetrieveHandler(c *C) {
	assertRetrieve(c, 3, 4)
	assertRetrieve(c, 5, 4)
	assertRetrieve(c, 20, 0)

	s.db.DropDatabase()
	response := request(c, "", "GET", -1, -1)
	c.Assert(response["error"].(bool), Equals, false)
}

func (s *HandlersSuite) TearDownTest(c *C) {
	s.db.DropDatabase()
}

func createLog() Log {
	return Log{
		Id:                 bson.NewObjectId(),
		Timestamp:          "t",
		CreatedAt:          time.Now(),
		Location:           Location{},
		Error:              Error{},
		ActionStateHistory: []ActionStateHistory{},
	}
}

func request(c *C, content, method string, limit, skip int) map[string]interface{} {
	var url = "http://0.0.0.0:8769/logs"
	if limit >= 0 && skip >= 0 {
		url = url + fmt.Sprintf("?limit=%d&skip=%d", limit, skip)
	}
	var buf = bytes.NewBuffer([]byte(content))
	req, err := http.NewRequest(method, url, buf)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	c.Assert(err, IsNil)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	var response map[string]interface{}
	c.Assert(json.Unmarshal(body, &response), IsNil)

	return response
}

func assertRetrieve(c *C, limit, skip int) {
	response := request(c, "", "GET", limit, skip)
	c.Assert(response["error"].(bool), Equals, false)
	c.Assert(len(response["logs"].([]interface{})), Equals, 6-skip)
	c.Assert(response["total"].(float64), Equals, float64(6))
	n := 6 - skip
	if n > limit {
		n = limit
	}
	c.Assert(response["count"].(float64), Equals, float64(6-skip))
}
