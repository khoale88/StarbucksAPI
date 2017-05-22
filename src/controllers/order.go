package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"models"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type OrderController struct {
	session *mgo.Session
}

func NewOrderController(session *mgo.Session) *OrderController {
	return &OrderController{session}
}

type OrderError struct {
	Error   error
	Message string
	Code    int
}

type OrderErrorHandler func(http.ResponseWriter, *http.Request) *OrderError

func (oerh OrderErrorHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if err := oerh(w, req); err != nil {
		http.Error(w, err.Message, err.Code)
	}
}

func (oc OrderController) PlaceOrderHandler() OrderErrorHandler {
	return func(w http.ResponseWriter, req *http.Request) *OrderError {
		var newOrder models.Order

		json.NewDecoder(req.Body).Decode(&newOrder)

		newOrder.Id = bson.NewObjectId()
		newOrder.Links = make(map[string]string, 2)
		newOrder.Links["payment"] = "http://localhost:9090/starbucks/" + newOrder.Id.Hex() + "/pay"
		newOrder.Links["order"] = "http://localhost:9090/starbucks/order/" + newOrder.Id.Hex()
		newOrder.Status = "payment expected"
		newOrder.Message = "Order has been placed."

		if err := oc.session.DB("Restbucks").C("Order").Insert(&newOrder); err != nil {
			return &OrderError{err, "Try again later", 500}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(newOrder)
		return nil
	}
}

func (oc OrderController) GetOrderHandler() OrderErrorHandler {
	return func(w http.ResponseWriter, req *http.Request) *OrderError {
		params := mux.Vars(req)
		id := params["id"]

		if !bson.IsObjectIdHex(id) {
			return &OrderError{errors.New("Invalid Order Id"), "Please check your Order Id", 404}
		}

		oid := bson.ObjectIdHex(id)
		order := models.Order{}
		if err := oc.session.DB("Restbucks").C("Order").FindId(oid).One(&order); err != nil {
			return &OrderError{err, "database err", 404}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(order)
		return nil
	}
}

func (oc OrderController) GetAllOrdersHandler() OrderErrorHandler {
	return func(w http.ResponseWriter, req *http.Request) *OrderError {
		orders := []models.Order{}
		if err := oc.session.DB("Restbucks").C("Order").Find(nil).All(&orders); err != nil {
			return &OrderError{err, "database err", 404}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(orders)
		return nil
	}
}

func (oc OrderController) CancelOrderHandler() OrderErrorHandler {
	return func(w http.ResponseWriter, req *http.Request) *OrderError {
		params := mux.Vars(req)
		id := params["id"]
		if !bson.IsObjectIdHex(id) {
			return &OrderError{errors.New("invalid Order Id"), "Please check your Order Id", 404}
		}

		oid := bson.ObjectIdHex(id)
		order := models.Order{}
		if err := oc.session.DB("Restbucks").C("Order").FindId(oid).One(&order); err != nil {
			return &OrderError{err, "Order not found", 404}
		}

		if order.Status != "payment expected" {
			return &OrderError{errors.New(""), "order cannot be cancelled after paid", 403}
		}

		if err := oc.session.DB("Restbucks").C("Order").RemoveId(oid); err != nil {
			return &OrderError{err, "Try again later", 500}
		}

		w.WriteHeader(204)
		return nil
	}
}

func (oc OrderController) UpdateOrderHandler() OrderErrorHandler {
	return func(w http.ResponseWriter, req *http.Request) *OrderError {
		params := mux.Vars(req)
		id := params["id"]

		if !bson.IsObjectIdHex(id) {
			return &OrderError{errors.New("invalid Order Id"), "Check your Order Id", 403}
		}

		oid := bson.ObjectIdHex(id)
		order := models.Order{}
		if err := oc.session.DB("Restbucks").C("Order").FindId(oid).One(&order); err != nil {
			return &OrderError{err, "Order not found", 404}
		}

		if order.Status != "payment expected" {
			return &OrderError{errors.New("Action prohibited"), "Order can't be modified after paid", 403}
		}

		json.NewDecoder(req.Body).Decode(&order)
		order.Message = "Order has been updated."

		if err := oc.session.DB("Restbucks").C("Order").UpdateId(oid, order); err != nil {
			return &OrderError{err, "Sever is busy. Please try again later", 500}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(order)
		return nil
	}
}

func (oc OrderController) PayOrderHandler() OrderErrorHandler {
	return func(w http.ResponseWriter, req *http.Request) *OrderError {
		params := mux.Vars(req)
		id := params["id"]
		if !bson.IsObjectIdHex(id) {
			return &OrderError{errors.New("invalid Order Id"), "Check your Order Id", 403}
		}

		oid := bson.ObjectIdHex(params["id"])
		order := models.Order{}
		if err := oc.session.DB("Restbucks").C("Order").FindId(oid).One(&order); err != nil {
			return &OrderError{err, "Order not found", 404}
		}

		if order.Status != "payment expected" {
			return &OrderError{errors.New("Action prohibited"), "Order were paid. You can't play it twice", 403}
		}

		order.Status = "PAID"
		order.Message = "Order has been paid."
		delete(order.Links, "payment")

		if err := oc.session.DB("Restbucks").C("Order").UpdateId(oid, order); err != nil {
			return &OrderError{err, "Sever is busy. Please try again later", 500}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(order)

		oc.AssignBarista(oid)
		return nil
	}
}
