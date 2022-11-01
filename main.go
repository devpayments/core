package main

import (
	"context"
	"fmt"
	"github.com/devpayments/core/config"
	"github.com/devpayments/core/datastore"
	"github.com/devpayments/core/datastore/db"
	"github.com/devpayments/core/payments"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func main() {
	initCtx := context.TODO()

	dbConfig := config.DatabaseConfig{
		Driver:   "postgres",
		Host:     "localhost",
		User:     "remi",
		Password: "root1234",
		SSLMode:  "disable",
		Name:     "payments",
		Port:     "5432",
	}
	d := db.New(dbConfig)
	dbCon, err := d.GetDbConnection()
	if err != nil {
		panic(err)
	}
	defer dbCon.Close()
	store := datastore.NewStore(*dbCon)

	paymentService := payments.NewPaymentService(*store)
	payment, err := paymentService.Initiate(initCtx)

	//err = paymentService.Authorize(initCtx, uuid.MustParse("ebb0fa19-b9e2-4718-b7d3-358e55b519f8"))
	err = paymentService.Complete(initCtx, uuid.MustParse(payment.ID))

	//res, err := store.Transactions.FindByID(initCtx, uuid.MustParse("90c36afe-f425-42a5-a111-c99a81ef6fd3"))

	//fmt.Printf("%+v\n", res)
	fmt.Printf("%+v\n", err)

	//fmt.Println(res)
}
