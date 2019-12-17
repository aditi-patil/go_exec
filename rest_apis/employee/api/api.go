// Testing go-swagger
//
// The purpose of this application is to create basic CRUD operation with mongo database, mux and go-swagger.
//
//     Schemes: http, https
//     Host: localhost:3000
//
//     Header:
//      - Access-Control-Allow-Methods: GET, POST, PUT
//     Produces:
//     - application/json
//
// swagger:meta
package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"utils"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var db *mgo.Database

var decode = json.NewDecoder
var marshal = json.Marshal

// var update = collection.Update

// swagger:model Employee
type Employee struct {
	// id of this employee
	ID bson.ObjectId `bson:"_id,omitempty" json:"_id,omitempty"`
	// name of this employee
	Name string `bson:"name,omitempty" json:"name,omitempty"`
	// empid of this employee
	EmpID string `bson:"empid,omitempty" json:"empid,omitempty"`
	// mobile of this employee
	Mobile int `bson:"mobile,omitempty" json:"mobile,omitempty"`
	// destignation of this employee
	Designation string `bson:"designation,omitempty" json:"designation,omitempty"`
	// practice of this employee
	Practice string `bson:"practice,omitempty" json:"practice,omitempty"`
}

// EmployeeCollection is a collection of all employess and its total count.
// swagger:model EmployeeCollection
type EmployeeCollection struct {
	// total of this employeeCollection
	Total int `json:"total,omitempty"`
	//employees of this employeeCollection
	Employees []Employee `json:"employees,omitempty"`
}

func init() {
	session, err := mgo.Dial(utils.DATABASE_URL)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(utils.DATABASE_NAME)
}

// ListAllEmployees gives list of all employees.
func ListAllEmployees(response http.ResponseWriter, request *http.Request) {
	//   swagger:operation GET /employees ListAllEmployees
	//
	//   Returns all employees from the collection db
	//   ---
	//   produces:
	//     - application/json
	//   parameters:
	//     - name: page
	//       in: query
	//       type: int
	//   responses:
	//     '200':
	//       description: employee response
	//       schema:
	//         type: Object
	//         $ref: "#/definitions/EmployeeCollection"
	response.Header().Set("content-type", "application/json")
	var employees []Employee
	var skips, page int
	var err error
	page = 1
	if request.FormValue("page") != "" {
		page, err = strconv.Atoi(request.FormValue("page"))
		if err != nil {
			handleError(err, "Failed to convert string to int", response, http.StatusBadRequest)
			return
		}
	}
	skips = (page - 1) * 20

	collection := db.C("employees")
	err = collection.Find(nil).Skip(skips).Limit(20).All(&employees)
	if err != nil {
		handleError(err, "Failed to load employees from database: %v", response, http.StatusBadRequest)
		return
	}
	empCollection := EmployeeCollection{
		Employees: employees,
		Total:     len(employees),
	}

	empResp, err := marshal(&empCollection)
	if err != nil {
		handleError(err, "Failed to marshal employees data: %v", response, http.StatusInternalServerError)
		return
	}
	response.Write(empResp)
}

// CreateEmployee creates new employee.
func CreateEmployee(response http.ResponseWriter, request *http.Request) {
	//   swagger:operation POST /employees CreateEmployee
	//
	//   Create new employee record into coolection db
	//   ---
	//   consumes:
	//     - application/json
	//   produces:
	//     - application/json
	//   parameters:
	//     - name: employee
	//       in: body
	//       required: true
	//       schema:
	//         type: object
	//         properties:
	//           id:
	//             type: bson.ObjectId
	//           name:
	//             type: string
	//           designation:
	//             type: string
	//           practice:
	//             type: string
	//           mobile:
	//             type: int
	//           empid:
	//             type: string
	//   responses:
	//     '200':
	//       description: employee response
	//       schema:
	//         $ref: "#/definitions/Employee"
	response.Header().Set("content-type", "application/json")
	var emp Employee
	err := json.NewDecoder(request.Body).Decode(&emp)
	if err != nil {
		handleError(err, "Failed to decode employee data: %v", response, http.StatusInternalServerError)
		return
	}
	emp.ID = bson.NewObjectId()
	collection := db.C("employees")
	err = collection.Insert(emp)
	if err != nil {
		handleError(err, "Failed to create employee: %v", response, http.StatusInternalServerError)
		return
	}
	result, er := marshal(&emp)
	if er != nil {
		handleError(err, "Failed to marshal employees data: %v", response, http.StatusInternalServerError)
		return
	}
	response.Write(result)

}

// GetEmployee gives employee with given id.
func GetEmployee(response http.ResponseWriter, request *http.Request) {
	//   swagger:operation GET /employee/{id} GetEmployee
	//
	//   Returns Employee with specific id from the collection db
	//   ---
	//   produces:
	//     - application/json
	//   parameters:
	//     - name: id
	//       in: path
	//       description: ID of employee to return
	//       required: true
	//       type: string
	//   responses:
	//     '200':
	//       description: Employee response
	//       schema:
	//         $ref: "#/definitions/Employee"
	response.Header().Set("content-type", "application/json")
	var employee Employee
	collection := db.C("employees")
	params := mux.Vars(request)
	fmt.Println(params)

	id := bson.ObjectIdHex(params["id"])
	err := collection.Find(bson.M{"_id": id}).One(&employee)
	if err != nil {
		handleError(err, "Failed to find employee from database: %v", response, http.StatusNotFound)
		return
	}
	result, er := marshal(&employee)
	if er != nil {
		handleError(err, "Failed to marshal employee data: %v", response, http.StatusInternalServerError)
		return
	}
	response.Write(result)

}

// DeleteEmployee deletes employee with given id from database.
func DeleteEmployee(response http.ResponseWriter, request *http.Request) {
	//   swagger:operation DELETE /employee/{id} DeleteEmployee
	//
	//   Deletes employee with specific id from the collection db
	//   ---
	//   produces:
	//     - application/json
	//   parameters:
	//     - name: id
	//       in: path
	//       description: ID of employee to return
	//       required: true
	//       type: string
	//   responses:
	//     '200':
	//       description: employee response
	response.Header().Set("content-type", "application/json")
	collection := db.C("employees")
	params := mux.Vars(request)

	id := bson.ObjectIdHex(params["id"])
	err := collection.Remove(bson.M{"_id": id})
	if err != nil {
		handleError(err, "Failed to remove employee from database: %v", response, http.StatusInternalServerError)
		return
	}
	response.Write([]byte("Employee deleted successfully."))
}

// UpdateEmployee updates employee record.
func UpdateEmployee(response http.ResponseWriter, request *http.Request) {
	//   swagger:operation PUT /employee/{id} UpdateEmployee
	//
	//   Update book parameters with specific id from the collection db
	//   ---
	//   produces:
	//     - application/json
	//   parameters:
	//     - name: id
	//       in: path
	//       description: ID of employee to return
	//       required: true
	//       type: string
	//     - name: employee
	//       in: body
	//       required: true
	//       schema:
	//         type: object
	//         properties:
	//           name:
	//             type: string
	//           mobile:
	//             type: int
	//           designation:
	//             type: string
	//           practice:
	//             type: string
	//   responses:
	//     '200':
	//       description: employee updated
	response.Header().Set("content-type", "application/json")
	var employee Employee
	collection := db.C("employees")
	params := mux.Vars(request)

	err := json.NewDecoder(request.Body).Decode(&employee)
	if err != nil {
		handleError(err, "Failed to decode employee data: %v", response, http.StatusInternalServerError)
		return
	}

	id := bson.ObjectIdHex(params["id"])
	err = collection.Update(bson.M{"_id": id}, bson.M{"$set": employee})
	if err != nil {
		handleError(err, "Failed to update employee: %v", response, http.StatusInternalServerError)
		return
	}
	result, err := marshal(&employee)
	if err != nil {
		handleError(err, "Failed to marshal employee data: %v", response, http.StatusInternalServerError)
		return
	}
	response.Write(result)
}

func handleError(err error, message string, w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(fmt.Sprintf(message, err)))
}
