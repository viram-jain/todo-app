package repository

import (
	"encoding/json"
	"net/http"
	_ "net/http/pprof"
	"os"

	"todoapp/constant"
	"todoapp/logger"
	"todoapp/model"
	"todoapp/redis"

	kaleyra_bson "bitbucket.org/kaleyra/mongo-sdk/bson"
	"bitbucket.org/kaleyra/mongo-sdk/mongo"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	mongo_bson "gopkg.in/mgo.v2/bson"
)

var (
	collection  *mongo.Collection
	postCache   redis.PostCache = redis.NewRedisCache(os.Getenv("REDIS_URL"), 0, 10)
	sugarLogger                 = logger.InitLogger()
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   interface{} `json:"error,omitempty"`
}

type UpdateResult struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type DeleteResult struct {
	ID string `json:"id"`
}

// create connection with mongo db
func init() {
	err := godotenv.Load(constant.Envfile)
	if err != nil {
		sugarLogger.Errorf("Failed to load the env file : Error = %s", err.Error())
		return
	}

	db := mongo.URI{
		Username: "",
		Password: "",
		Host:     os.Getenv("HOST"),
		DB:       os.Getenv("DATABASE_NAME"),
		Port:     os.Getenv("MONGO_PORT"),
	}

	sugarLogger = logger.InitLogger()
	client, err := mongo.NewClient(db)
	if err != nil {
		sugarLogger.Errorf("Failed to connect to mongodb = %s", err.Error())
		return
	}

	zap.L().Info("Connected to MongoDB!")
	collection = client.Collection(os.Getenv("COLLECTION_NAME"))
	zap.L().Info("Collection instance created!")
}

// CreateTodo create a new Todo
func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var task model.TodoModel
	_ = json.NewDecoder(r.Body).Decode(&task)
	_, err := collection.InsertOne(task)
	if err != nil {
		responseData := &Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to insert todo",
			Error:   err,
		}
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(responseData)
		if err != nil {
			sugarLogger.Error(err)
			return
		}
	}

	responseData := &Response{
		Status:  http.StatusCreated,
		Message: "Todo inserted successfully",
		Data:    task,
	}
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(responseData)
	if err != nil {
		sugarLogger.Error(err)
		return
	}
}

// GetAllTodos get all the Todos
func GetAllTodos(w http.ResponseWriter, r *http.Request) {
	results, err := collection.FindAll()
	if err != nil {
		responseData := &Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to retrieve todos",
			Error:   err,
		}
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(responseData)
		if err != nil {
			sugarLogger.Error(err)
			return
		}
	}

	responseData := &Response{
		Status:  http.StatusOK,
		Message: "Retrieved todos",
		Data:    results,
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(responseData)
	if err != nil {
		sugarLogger.Error(err)
		return
	}
}

// GetTodoById get single todo by the given id
func GetTodoById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var result *model.TodoModel
	if !mongo_bson.IsObjectIdHex(params[constant.ID]) {
		responseData := &Response{
			Status:  http.StatusBadRequest,
			Message: "Failed to retrieve todo",
			Error:   "Invalid ID provided",
		}
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(responseData)
		if err != nil {
			sugarLogger.Error(err)
			return
		}
		return
	}

	id, _ := primitive.ObjectIDFromHex(params[constant.ID])
	post := postCache.Get(id.String())
	if post == nil {
		filter := kaleyra_bson.D{{Key: constant.Key, Value: id}}
		result, err := collection.FindOne(filter)
		if err != nil {
			responseData := &Response{
				Status:  http.StatusInternalServerError,
				Message: "Failed to retrieve todo",
				Error:   err,
			}
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(responseData)
			if err != nil {
				sugarLogger.Error(err)
				return
			}
		}
		postCache.Set(id.String(), result)
	}

	result = postCache.Get(id.String())
	responseData := &Response{
		Status:  http.StatusOK,
		Message: "Retrieved todo",
		Data:    result,
	}

	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(responseData)
	if err != nil {
		sugarLogger.Error(err)
		return
	}
}

// UpdateTodo update a Todo
func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if !mongo_bson.IsObjectIdHex(params[constant.ID]) {
		responseData := &Response{
			Status:  http.StatusBadRequest,
			Message: "Failed to update todo",
			Error:   "Invalid ID provided",
		}
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(responseData)
		if err != nil {
			sugarLogger.Error(err)
			return
		}
		return
	}

	id, _ := primitive.ObjectIDFromHex(params[constant.ID])
	var t model.TodoModel

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		responseData := &Response{
			Status:  http.StatusBadRequest,
			Message: "Failed to update todo",
			Error:   err,
		}
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(responseData)
		if err != nil {
			sugarLogger.Error(err)
			return
		}
		return
	}

	if t.Title == "" {
		responseData := &Response{
			Status:  http.StatusBadRequest,
			Message: "Failed to update todo",
			Error:   "Empty title provided",
		}
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(responseData)
		if err != nil {
			sugarLogger.Error(err)
			return
		}
		return
	}

	filter := kaleyra_bson.D{{Key: constant.Key, Value: id}}
	update := model.TodoModel{
		Title: t.Title,
	}

	_, err = collection.UpsertOne(filter, update)
	if err != nil {
		responseData := &Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to update todo",
			Error:   err,
		}
		w.WriteHeader(http.StatusInternalServerError)
		err := json.NewEncoder(w).Encode(responseData)
		if err != nil {
			sugarLogger.Error(err)
			return
		}
		return
	}

	insertResult := &UpdateResult{
		ID:    id.Hex(),
		Title: t.Title,
	}

	responseData := &Response{
		Status:  http.StatusOK,
		Message: "Todo updated successfully",
		Data:    insertResult,
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(responseData)
	if err != nil {
		sugarLogger.Error(err)
		return
	}
}

// DeleteTodo deletes a ToDo
func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if !mongo_bson.IsObjectIdHex(params[constant.ID]) {
		responseData := &Response{
			Status:  http.StatusBadRequest,
			Message: "Failed to delete todo",
			Error:   "Invalid ID provided",
		}
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(responseData)
		if err != nil {
			sugarLogger.Error(err)
			return
		}
		return
	}

	id, _ := primitive.ObjectIDFromHex(params[constant.ID])
	filter := kaleyra_bson.D{{Key: constant.Key, Value: id}}
	res, err := collection.DeleteOne(filter)
	if err != nil {
		responseData := &Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to delete todo",
			Error:   err,
		}
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(responseData)
		if err != nil {
			sugarLogger.Error(err)
			return
		}
		return
	}

	if res == 0 {
		responseData := &Response{
			Status:  http.StatusBadRequest,
			Message: "ID not found",
			Error:   err,
		}
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(responseData)
		if err != nil {
			sugarLogger.Error(err)
			return
		}
		return
	}

	deleteData := &DeleteResult{
		ID: id.Hex(),
	}

	responseData := &Response{
		Status:  http.StatusOK,
		Message: "Todo deleted successfully",
		Data:    deleteData,
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(responseData)
	if err != nil {
		sugarLogger.Error(err)
		return
	}
}
