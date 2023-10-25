package model

import (
	"log"

	"github.com/google/uuid"
)

var MetaInfo = map[string]interface{}{}

func init() {
	log.Println("init")
	MetaInfo["connection_status"] = false
	MetaInfo["last_commit_id"] = uuid.New()
}

func update() {
	MetaInfo["last_commit_id"] = uuid.New()
}

func Connected() {
	MetaInfo["connection_status"] = true
}

func DisConnected() {
	MetaInfo["connection_status"] = false
}
