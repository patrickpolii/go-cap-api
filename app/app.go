package app

import (
	"capi/domain"
	"capi/logger"
	"capi/service"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func sanityCheck() {
	envProps := []string{
		"SERVER_ADDRESS",
		"SERVER_PORT",
		"DB_USER",
		"DB_PASSWD",
		"DB_ADDR",
		"DB_PORT",
		"DB_NAME",
	}

	for _, envKey := range envProps {
		if os.Getenv(envKey) == "" {
			logger.Fatal(fmt.Sprintf("environment variable %s not defined. terminating application...", envKey))
		}
	}

	logger.Info("environment variables loaded...")

}

func Start() {

	err := godotenv.Load()
	if err != nil {
		logger.Fatal("error loading .env file")
	}
	logger.Info("load environment variables...")

	sanityCheck()

	dbClient := getClientDB()

	// * wiring
	// * setup repository
	customerRepositoryDB := domain.NewCustomerRepositoryDB(dbClient)
	accountRepositoryDB := domain.NewAccountRepositoryDB(dbClient)
	authRepositoryDB := domain.NewAuthRepositoryDB(dbClient)

	// * setup service
	customerService := service.NewCustomerService(customerRepositoryDB)
	accountService := service.NewAccountService(accountRepositoryDB)
	authService := service.NewAuthService(authRepositoryDB)

	// * setup handler
	ch := CustomerHandlers{customerService}
	ah := AccountHandler{accountService}
	authH := AuthHandler{authService}

	// * create ServeMux
	mux := mux.NewRouter()

	authR := mux.PathPrefix("/auth").Subrouter()
	authR.HandleFunc("/login", authH.Login).Methods(http.MethodPost)

	authR.Use(loggingMiddleware)
	// * defining routes
	// mux.HandleFunc("/auth/login", authH.Login).Methods(http.MethodPost)

	mux.HandleFunc("/customers", ch.getAllCustomers).Methods(http.MethodGet)
	mux.HandleFunc("/customers/{customer_id:[0-9]+}", ch.getCustomerByID).Methods(http.MethodGet)
	mux.HandleFunc("/customers/{customer_id:[0-9]+}/accounts", ah.NewAccount).Methods(http.MethodPost)
	mux.HandleFunc("/customers/{customer_id:[0-9]+}/accounts/{account_id:[0-9]+}", ah.MakeTransaction).Methods(http.MethodPost)

	mux.Use(authMiddleware)
	// * starting the server

	serverAddr := os.Getenv("SERVER_ADDRESS")
	serverPort := os.Getenv("SERVER_PORT")

	logger.Info(fmt.Sprintf("start server on %s:%s...", serverAddr, serverPort))
	http.ListenAndServe(fmt.Sprintf("%s:%s", serverAddr, serverPort), mux)
}

func getClientDB() *sqlx.DB {
	dbUser := os.Getenv("DB_USER")
	dbPasswd := os.Getenv("DB_PASSWD")
	dbAddr := os.Getenv("DB_ADDRESS")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPasswd, dbAddr, dbPort, dbName)
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		logger.Fatal(err.Error())
	}
	logger.Info("success connect to database...")

	return db
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timer := time.Now()
		next.ServeHTTP(w, r)
		logger.Info(fmt.Sprintf("%v %v %v", r.Method, r.URL, time.Since(timer)))
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		// check if token has bearer
		if !strings.Contains(token, "Bearer") {
			logger.Error("token invalid")
			return
		}
		// split token -> ambil tokennya buang "Bearer" nya
		getToken := ""
		tokenArray := strings.Split(token, " ")
		if len(tokenArray) == 2 {
			getToken = tokenArray[1]
		}
		fmt.Println(getToken)

		// parsing token, err := jwt.Parse(
		signedToken, err := jwt.Parse(getToken, func(token *jwt.Token) (interface{}, error) {
			return []byte("rahasia"), nil
		})
		if err != nil {
			logger.Error("Failed to parse token:" + err.Error())
		}

		fmt.Println(err)
		fmt.Println("test token", signedToken)
		// // check token validation
		if signedToken.Valid {
			fmt.Println("token is valid")
			writeResponse(w, http.StatusOK, signedToken)
		}

		logger.Info(token)

		next.ServeHTTP(w, r)
	})
}

// func Start() {

// 	// * wiring
// 	ch := CustomerHandler{service.NewCustomerService(domain.NewCustomerRepositoryDB())}

// 	// * create ServeMux
// 	mux := mux.NewRouter()

// 	// * defining routes
// 	// mux.HandleFunc("/greet", greet).Methods(http.MethodGet)
// 	mux.HandleFunc("/customers", ch.GetAllCustomers).Methods(http.MethodGet)
// 	// mux.HandleFunc("/customers", addCustomer).Methods(http.MethodPost)

// 	mux.HandleFunc("/customers/{customer_id:[0-9]+}", ch.GetCustomerByID).Methods(http.MethodGet)
// 	// mux.HandleFunc("/customers/{customer_id:[0-9]+}", updateCustomer).Methods(http.MethodPut)

// 	// * starting the server
// 	http.ListenAndServe(":8080", mux)
// }
