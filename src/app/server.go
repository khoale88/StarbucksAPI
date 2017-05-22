package main

import (
	"fmt"
	"log"
	"net/http"

	"os"

	"controllers"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
)

var mongo0 = "10.0.0.27:27017"
var mongo1 = "10.0.2.34.224:27017"
var mongo2 = "10.0.0.8:27017"

func main() {

	router := mux.NewRouter()
	oc := controllers.NewOrderController(getSession())

	router.Handle("/starbucks/order", oc.PlaceOrderHandler()).Methods("POST")
	router.Handle("/starbucks/orders", oc.GetAllOrdersHandler()).Methods("GET")
	router.Handle("/starbucks/order/{id}", oc.GetOrderHandler()).Methods("GET")
	router.Handle("/starbucks/order/{id}", oc.CancelOrderHandler()).Methods("DELETE")
	router.Handle("/starbucks/order/{id}", oc.UpdateOrderHandler()).Methods("PUT")
	router.Handle("/starbucks/order/{id}/pay", oc.PayOrderHandler()).Methods("POST")

	fmt.Println("serving on port 9090")

	err := http.ListenAndServe(":9090", router)

	log.Fatal(err)
}

func getSession() *mgo.Session {
	// Connect to our local mongo
	mongohost := os.Getenv("MONGOHOST")
	mongoDial := "mongodb://" + mongohost + ":27017"
	s, err := mgo.Dial(mongoDial)
	//s, err := mgo.Dial("mongodb://node0:30001,node1:30002,node2:30003")
	// s, err := mgo.Dial("mongodb://" + mongo0 + "," + mongo1 + "," + mongo2)
	//s, err := mgo.Dial("mongodb://127.23.234.32:27017")
	if err != nil {
		log.Fatal(err)
	}
	return s
}
