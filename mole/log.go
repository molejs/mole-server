package mole

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Log struct {
	Id                 bson.ObjectId      `bson:"_id" json:"id"`
	Timestamp          string             `json:"timestamp" binding:"required"`
	CreatedAt          time.Time          `bson:"created_at" json:"-"`
	Location           Location           `json:"location" binding:"required"`
	Error              Error              `json:"error" binding:"required"`
	ActionStateHistory []ActionStateHistory `json:"action_state_history"`
}

type Location struct {
	Host     string `json:"host"`
	Href     string `json:"href"`
	Hash     string `json:"hash"`
	Pathname string `json:"pathname"`
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
	Search   string `json:"search"`
}

type Error struct {
	Message    string            `json:"message"`
	Stacktrace []StracktraceLine `json:"stacktrace"`
}

type StracktraceLine struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int `json:"line"`
	Column   int `json:"column"`
}

type ActionStateHistory struct {
	Action interface{} `json:"action"`
	State  interface{} `json:"state"`
}
