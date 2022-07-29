package app

import (
	"capi/service"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type CustomerHandler struct {
	service service.CustomerService
}

func (ch *CustomerHandler) GetAllCustomers(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprint(w, "get customer endpoint\n")
	customerStatus := r.URL.Query().Get("status")
	fmt.Println("customer status", customerStatus)

	customers, err := ch.service.GetAllCustomer(customerStatus)

	if err != nil {
		writeResponse(w, err.Code, err.AsMessage())
		return
	} else {
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(customers)
	}
}

func (ch *CustomerHandler) getCustomerByID(w http.ResponseWriter, r *http.Request) {

	// * get route variable
	vars := mux.Vars(r)

	customerID := vars["customer_id"]

	customer, err := ch.service.GetCustomerByID(customerID)
	if err != nil {
		writeResponse(w, err.Code, err.AsMessage())
		return
	}

	// * return customer data
	writeResponse(w, http.StatusOK, customer)
}

func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}

// func greet(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprint(w, "Hello Celerates!")
// }

// func getCustomers(w http.ResponseWriter, r *http.Request) {
// 	//fmt.Fprint(w, "get customer endpoint")

// 	if r.Header.Get("Content-Type") == "application/xml" {
// 		w.Header().Add("Content-Type", "application/xml")
// 		xml.NewEncoder(w).Encode(customers)
// 	} else {
// 		w.Header().Add("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(customers)
// 	}

// }

// func getCustomerById(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	customerId := vars["customer_id"]
// 	//id, _ := strconv.Atoi(customerId)
// 	id, err := strconv.Atoi(customerId)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		fmt.Fprint(w, "invalid customer id")
// 		return
// 	}

// 	var cust Customer

// 	for _, data := range customers {
// 		if data.ID == id {
// 			cust = data
// 		}
// 	}
// 	if cust.ID == 0 {
// 		w.WriteHeader(http.StatusNotFound)
// 		fmt.Fprint(w, "customer data not found")
// 		return
// 	}

// 	// * return customer data
// 	w.Header().Add("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(cust)

// 	// if r.Header.Get("Content-Type") == "application/xml" {
// 	// 	w.Header().Add("Content-Type", "application/xml")
// 	// 	xml.NewEncoder(w).Encode(cust)
// 	// } else {
// 	// 	w.Header().Add("Content-Type", "application/json")
// 	// 	json.NewEncoder(w).Encode(cust)
// 	// }
// }

// func addCustomers(w http.ResponseWriter, r *http.Request) {
// 	// decode request body
// 	var cust Customer
// 	json.NewDecoder(r.Body).Decode(&cust)

// 	//generate new id
// 	nextID := getNextID()
// 	cust.ID = nextID

// 	// save data to array

// 	customers = append(customers, cust)

// 	w.WriteHeader(http.StatusCreated)
// }

// func getNextID() int {
// 	cust := customers[len(customers)-1]

// 	return cust.ID + 1
// }

// func updateCustomer(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	customerId := vars["customer_id"]
// 	//id, _ := strconv.Atoi(customerId)
// 	id, err := strconv.Atoi(customerId)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		fmt.Fprint(w, "invalid customer id")
// 		return
// 	}

// 	var cust Customer

// 	for customerIndex, data := range customers {
// 		if data.ID == id {
// 			cust = data

// 			var newCust Customer
// 			json.NewDecoder(r.Body).Decode(&newCust)

// 			customers[customerIndex].Name = newCust.Name
// 			customers[customerIndex].City = newCust.City
// 			customers[customerIndex].Zipcode = newCust.Zipcode

// 			w.WriteHeader(http.StatusOK)
// 			fmt.Fprintln(w, "customer data updated")
// 			return
// 		}
// 	}
// 	if cust.ID == 0 {
// 		w.WriteHeader(http.StatusNotFound)
// 		fmt.Fprint(w, "customer data not found")
// 		return
// 	}

// 	// w.Header().Set("Content-Type", "application/json")
// 	// params := mux.Vars(r)
// 	// for i, data := range customers {
// 	// 	if data.ID == params["id"] {
// 	// 		customers = append(customers[:i], customers[i+1]...)
// 	// 		var newCustomer Customer
// 	// 		json.NewDecoder(r.Body).Decode(&newCustomer)
// 	// 		newCustomer.ID = params["id"]
// 	// 		customers = append(customers, newCustomer)
// 	// 		json.NewEncoder(w).Encode(newCustomer)
// 	// 		return
// 	// 	}
// 	// }
// }

// func deleteCustomer(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	vars := mux.Vars(r)
// 	customerId := vars["customer_id"]
// 	id, err := strconv.Atoi(customerId)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		fmt.Fprint(w, "invalid customer id")
// 		return
// 	}

// for customerIndex, data := range customers {
// 	if data.ID == id {
// 		var newCust Customer
// 		json.NewDecoder(r.Body).Decode(&newCust)
// 		customers = append(customers[:customerIndex], customers[customerIndex+1:]...)
// 		w.WriteHeader(http.StatusOK)
// 		fmt.Fprintln(w, "customer data deleted")
// 		return
// 	}
// }

// }
