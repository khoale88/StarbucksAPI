package controllers

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type BaristaController struct {
	session *mgo.Session
	oid     bson.ObjectId
}

func NewBaristaController(session *mgo.Session, oid bson.ObjectId) *BaristaController {
	return &BaristaController{session, oid}
}

func (bc BaristaController) ProcessOrder() {
	go func() {
		time.Sleep(time.Second * 5)
		err := bc.session.DB("Restbucks").C("Order").UpdateId(bc.oid, bson.M{"$set": bson.M{"status": "PREPARING"}})
		if err != nil {
			fmt.Println(err)
			return
		}
		time.Sleep(time.Second * 20)
		err = bc.session.DB("Restbucks").C("Order").UpdateId(bc.oid, bson.M{"$set": bson.M{"status": "SERVED"}})
		if err != nil {
			fmt.Println(err)
			return
		}
		time.Sleep(time.Second * 10)
		err = bc.session.DB("Restbucks").C("Order").UpdateId(bc.oid, bson.M{"$set": bson.M{"status": "COLLECTED"}})
		if err != nil {
			fmt.Println(err)
			return
		}
	}()
}

func (oc OrderController) AssignBarista(oid bson.ObjectId) {
	bc := NewBaristaController(oc.session, oid)
	bc.ProcessOrder()
}
