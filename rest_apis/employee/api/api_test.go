package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"utils"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type fakeEmployee struct {
	err error
}

func (f *fakeEmployee) jsonMarshal(v interface{}) ([]byte, error) {
	return nil, f.err
}

func TestRestApiCrud(t *testing.T) {
	t.Run("testCreateEmployee", func(t *testing.T) {
		testCreateEmployee(t)
	})

	t.Run("testDeleteEmployee", func(t *testing.T) {
		testDeleteEmployee(t)
	})

	t.Run("testGetEmployee", func(t *testing.T) {
		testGetEmployee(t)
	})

	t.Run("testUpdateEmployee", func(t *testing.T) {
		testUpdateEmployee(t)
	})

	t.Run("testListAllEmployees", func(t *testing.T) {
		testListAllEmployees(t)
	})
}

func mockDb() {
	session, err := mgo.Dial(utils.TEST_DB_URL)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(utils.TEST_DB_NAME)
}

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/employees", ListAllEmployees).Methods("GET")
	router.HandleFunc("/employees", CreateEmployee).Methods("POST")
	router.HandleFunc("/employee/{id}", GetEmployee).Methods("GET")
	router.HandleFunc("/employee/{id}", DeleteEmployee).Methods("DELETE")
	router.HandleFunc("/employee/{id}", UpdateEmployee).Methods("PUT")

	return router
}

func testUpdateEmployee(t *testing.T) {
	mockDb()
	t.Run("it updates employee record with given id", func(t *testing.T) {
		var emp1 []byte
		emp1 = []byte(`{"Name":"Jonny","Practice":"New Test"}`)
		empID := addEmployee()
		request, err := http.NewRequest("PUT", "/employee/"+empID, bytes.NewBuffer(emp1))

		if err != nil {
			t.Fatal(err)
		}

		response := httptest.NewRecorder()

		Router().ServeHTTP(response, request)
		if status := response.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	})

	t.Run("it fails to update employee when json marshal failed", func(t *testing.T) {
		var emp1 []byte
		emp1 = []byte(`{"Name":"Jonny","Practice":"New Test"}`)
		empID := addEmployee()
		request, err := http.NewRequest("PUT", "/employee/"+empID, bytes.NewBuffer(emp1))

		if err != nil {
			t.Fatal(err)
		}

		response := httptest.NewRecorder()

		f := &fakeEmployee{err: errors.New("Failed")}
		marshal = f.jsonMarshal

		Router().ServeHTTP(response, request)
		if status := response.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}
		defer func() {
			marshal = json.Marshal
		}()

	})

	t.Run("it fails to update employee when json decode failed", func(t *testing.T) {
		var emp1 []byte
		emp1 = []byte(`{"Name":"Jonny","Practice":"New Test}`)
		empID := addEmployee()
		request, err := http.NewRequest("PUT", "/employee/"+empID, bytes.NewBuffer(emp1))

		if err != nil {
			t.Fatal(err)
		}

		response := httptest.NewRecorder()

		Router().ServeHTTP(response, request)
		if status := response.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}
	})

	t.Run("it fails to update employee if employee not present", func(t *testing.T) {
		var emp1 []byte
		emp1 = []byte(`{"Name":"Jonny","Practice":"New Test"}`)
		empID := addEmployee()
		err := removeEmployee(empID)
		request, err := http.NewRequest("PUT", "/employee/"+empID, bytes.NewBuffer(emp1))

		if err != nil {
			t.Fatal(err)
		}

		response := httptest.NewRecorder()

		Router().ServeHTTP(response, request)
		if status := response.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}
	})
}

func testCreateEmployee(t *testing.T) {
	mockDb()
	t.Run("it creates new entry for employee into the database", func(t *testing.T) {
		var employee []byte
		employee = []byte(`{"Name":"John","Practice":"Test", "Designstion": "SE", "EmpId": "EP-19", "Mobile": 432232223}`)
		request, err := http.NewRequest("POST", "/employees", bytes.NewBuffer(employee))
		if err != nil {
			t.Fatal("Error")
		}
		response := httptest.NewRecorder()
		Router().ServeHTTP(response, request)
		assert.Equal(t, 200, response.Code, "OK response is expected")
	})

	t.Run("it fails to create new entry for invalid json", func(t *testing.T) {
		var employee []byte
		employee = []byte(`{"Name":"John","Practice":"Test, "Designation": "SE", }`)
		request, err := http.NewRequest("POST", "/employees", bytes.NewBuffer(employee))
		if err != nil {
			t.Fatal("Error")
		}
		response := httptest.NewRecorder()
		Router().ServeHTTP(response, request)
		assert.Equal(t, 500, response.Code, "Failed to decode employee data")
	})

	t.Run("it fails to create new entry for same empID", func(t *testing.T) {
		addEmployee()
		var employee []byte
		employee = []byte(`{"Name":"John","Practice":"Test", "Designstion": "SE", "EmpId": "EP-19"}`)
		request, err := http.NewRequest("POST", "/employees", bytes.NewBuffer(employee))
		if err != nil {
			t.Fatal("Error")
		}
		response := httptest.NewRecorder()
		Router().ServeHTTP(response, request)
		assert.Equal(t, 500, response.Code, "Failed to decode employee data")
	})

	t.Run("it fails to create employee when json marshal failed", func(t *testing.T) {
		var employee []byte

		employee = []byte(`{"Name":"John","Practice":"Test", "Designation": "SE"}`)
		request, err := http.NewRequest("POST", "/employees", bytes.NewBuffer(employee))
		if err != nil {
			t.Fatal("Error")
		}
		response := httptest.NewRecorder()

		f := &fakeEmployee{err: errors.New("Failed")}
		marshal = f.jsonMarshal

		Router().ServeHTTP(response, request)
		if status := response.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}

		defer func() {
			marshal = json.Marshal
		}()
	})

}

func testDeleteEmployee(t *testing.T) {
	mockDb()
	t.Run("it deletes employee record with given id", func(t *testing.T) {
		empID := addEmployee()
		request, err := http.NewRequest("DELETE", "/employee/"+empID, nil)
		if err != nil {
			t.Fatal(err)
		}
		response := httptest.NewRecorder()

		Router().ServeHTTP(response, request)
		if status := response.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	})

	t.Run("it fails to delete employee if id is invalid", func(t *testing.T) {
		empID := addEmployee()
		err := removeEmployee(empID)
		request, err := http.NewRequest("DELETE", "/employee/"+empID, nil)
		if err != nil {
			t.Fatal(err)
		}
		response := httptest.NewRecorder()

		Router().ServeHTTP(response, request)
		if status := response.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}
	})
}

func testGetEmployee(t *testing.T) {
	mockDb()

	t.Run("it fails to give employee when json marshal failed", func(t *testing.T) {
		empID := addEmployee()
		request, err := http.NewRequest("GET", "/employee/"+empID, nil)
		if err != nil {
			t.Fatal(err)
		}
		response := httptest.NewRecorder()

		f := &fakeEmployee{err: errors.New("Failed")}
		marshal = f.jsonMarshal

		Router().ServeHTTP(response, request)
		if status := response.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}
		defer func() {
			marshal = json.Marshal
		}()

	})

	t.Run("it gives employee record with given id", func(t *testing.T) {
		empID := addEmployee()
		// check response for employee where id is not blank
		request, err := http.NewRequest("GET", "/employee/"+empID, nil)
		if err != nil {
			t.Fatal(err)
		}

		response := httptest.NewRecorder()

		Router().ServeHTTP(response, request)

		expected := `{"_id":"` + empID + `","name":"John","empid":"EP-19","mobile":432232223,"designation":"SE","practice":"Test"}`
		if response.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				response.Body.String(), expected)
		}
	})

	t.Run("it fails to give employee if id is invalid", func(t *testing.T) {
		empID := addEmployee()
		err := removeEmployee(empID)
		request, err := http.NewRequest("GET", "/employee/"+empID, nil)
		if err != nil {
			t.Fatal(err)
		}
		response := httptest.NewRecorder()

		Router().ServeHTTP(response, request)
		if status := response.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusNotFound)
		}
	})

}

func testListAllEmployees(t *testing.T) {
	mockDb()
	t.Run("it fails if page number is invalid", func(t *testing.T) {

		request, err := http.NewRequest("GET", "/employees?page=invalid", nil)
		if err != nil {
			t.Fatal(err)
		}
		response := httptest.NewRecorder()

		Router().ServeHTTP(response, request)
		if status := response.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})

	t.Run("it fails to load data if page number is negative integer", func(t *testing.T) {

		request, err := http.NewRequest("GET", "/employees?page=-1", nil)
		if err != nil {
			t.Fatal(err)
		}
		response := httptest.NewRecorder()

		Router().ServeHTTP(response, request)
		if status := response.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})

	t.Run("it fails to list employees when json marshal failed", func(t *testing.T) {
		request, err := http.NewRequest("GET", "/employees", nil)
		if err != nil {
			t.Fatal(err)
		}
		response := httptest.NewRecorder()

		f := &fakeEmployee{err: errors.New("Failed")}
		marshal = f.jsonMarshal

		Router().ServeHTTP(response, request)
		if status := response.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}

		defer func() {
			marshal = json.Marshal
		}()
	})

	t.Run("it gives all employees list", func(t *testing.T) {
		request, err := http.NewRequest("GET", "/employees", nil)
		if err != nil {
			t.Fatal(err)
		}
		response := httptest.NewRecorder()
		Router().ServeHTTP(response, request)

		if status := response.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	})

	// t.Run("it fails to list employees when db is closed", func(t *testing.T) {
	// 	// session.Close()
	// 	request, err := http.NewRequest("GET", "/employees", nil)
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	response := httptest.NewRecorder()

	// 	Router().ServeHTTP(response, request)
	// 	if status := response.Code; status != http.StatusInternalServerError {
	// 		t.Errorf("handler returned wrong status code: got %v want %v",
	// 			status, http.StatusInternalServerError)
	// 	}

	// })

}

func addEmployee() string {
	var employee []byte
	var emp Employee

	employee = []byte(`{"Name":"John","Practice":"Test", "Designation": "SE", "EmpId": "EP-19", "Mobile": 432232223}`)
	err := json.Unmarshal(employee, &emp)
	emp.ID = bson.NewObjectId()

	collection := db.C("employees")
	_ = collection.DropCollection()
	_ = collection.DropIndex("empid")
	index := mgo.Index{
		Key:      []string{"empid"},
		Unique:   true,
		DropDups: true,
	}
	err = collection.EnsureIndex(index)
	if err != nil {
		log.Fatal(err)
	}
	err = collection.Insert(emp)
	if err != nil {
		log.Fatal(err)
	}
	return emp.ID.Hex()
}

func removeEmployee(id string) error {
	collection := db.C("employees")
	empID := bson.ObjectIdHex(id)
	err := collection.Remove(bson.M{"_id": empID})
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
