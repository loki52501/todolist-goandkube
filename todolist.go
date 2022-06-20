package main
import("fmt"
"io"
"os"
"net/http"
"github.com/gorilla/mux"
log "github.com/sirupsen/logrus" 
//"github.com/lib/pq"
"gorm.io/driver/postgres"
"gorm.io/gorm"
"encoding/json"
"strconv"
"github.com/rs/cors"
"github.com/joho/godotenv"
)

type todo struct{
	ID uint `gorm:"primaryKey"`
	Todos string 
	Completed bool
}



func Healthz(w http.ResponseWriter, r *http.Request){
log.Info("API Health is ok");
w.Header().Set("Content-Type","application/json")
io.WriteString(w,`{"alive":true}`)
}
func readinessHandler(w http.ResponseWriter, r *http.Request){
w.WriteHeader(http.StatusOK)
}

func init(){
log.SetFormatter(&log.TextFormatter{})
log.SetReportCaller(true)
}
func CreateItem(w http.ResponseWriter, r *http.Request){
	todos:=r.FormValue("description")
	log.WithFields(log.Fields{"Todos":todos}).Info("add new todoitem. Save it")
	todo1:=&todo{Todos:todos,Completed:false}
	db.Create(&todo1)
	result := map[string]interface{}{}
db.Model(&todo1).First(&result)
	log.Info(result["Todos"])
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(result)
}
func UpdateItem(w http.ResponseWriter, r *http.Request){

  vars := mux.Vars(r)
        id, _ := strconv.Atoi(vars["id"])

        // Test if the TodoItem exist in DB
        err := GetItemByID(id)
        if err == false {
               w.Header().Set("Content-Type", "application/json")
                io.WriteString(w, `{"updated": false, "error": "Record Not Found"}`)
      } else {
               completed, _ := strconv.ParseBool(r.FormValue("completed"))
               log.WithFields(log.Fields{"Id": id, "Completed": completed}).Info("Updating TodoItem")
               todo1 := &todo{}
               db.First(&todo1, id)
               todo1.Completed = completed
               db.Save(&todo1)
               w.Header().Set("Content-Type", "application/json")
                io.WriteString(w, `{"updated": true}`)
}
}
func DeleteItem(w http.ResponseWriter, r *http.Request) {
       // Get URL parameter from mux
       vars := mux.Vars(r)
       id, _ := strconv.Atoi(vars["id"])

       // Test if the TodoItem exist in DB
       err := GetItemByID(id)
       if err == false {
               w.Header().Set("Content-Type", "application/json")
               io.WriteString(w, `{"deleted": false, "error": "Record Not Found"}`)
       } else {
               log.WithFields(log.Fields{"Id": id}).Info("Deleting TodoItem")
               todo1 := &todo{}
               db.First(&todo1, id)
               db.Delete(&todo1)
               w.Header().Set("Content-Type", "application/json")
                io.WriteString(w, `{"deleted": true}`)
       }
}
func GetItemByID(Id int) bool {
       todo1 := &todo{}
       result := db.First(&todo1, Id)
       if result.Error != nil{
               log.Warn("TodoItem not found in database")
               return false
       }
       return true
}

func GetCompletedItems(w http.ResponseWriter, r *http.Request) {
       log.Info("Get completed TodoItems")
       completedTodoItems := GetTodoItems(true)
       w.Header().Set("Content-Type", "application/json")
       json.NewEncoder(w).Encode(completedTodoItems)
}

func GetIncompleteItems(w http.ResponseWriter, r *http.Request) {
       log.Info("Get Incomplete TodoItems")
       IncompleteTodoItems := GetTodoItems(false)
       w.Header().Set("Content-Type", "application/json")
       json.NewEncoder(w).Encode(IncompleteTodoItems)
}
 
func GetTodoItems(completed bool) interface{} {
       var todos1 []todo
       TodoItems := db.Where("completed = ?", completed).Find(&todos1)
       result := map[string]interface{}{}
db.Table("todos").Where("completed=?",completed).Find(&result)
log.Info(TodoItems,"hi",result,"hsi",todos1)
       return todos1
}
 
var err =godotenv.Load()

var (
host=os.Getenv("DB_HOST")	
driver=os.Getenv("DB_DRIVER")
port = os.Getenv("DB_PORT")
user=os.Getenv("DB_USER")
password=os.Getenv("DB_PASSWORD")
dbname=os.Getenv("DB_NAME")
)
var dsn =fmt.Sprintf( "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Kolkata",host,user,password,dbname,port)
var db,_ = gorm.Open(postgres.Open(dsn), &gorm.Config{})
func main(){
log.Info(host,driver,"chheese",port)
	sqlDB,err:=db.DB()
	if err!=nil {
panic(err)
	}
	defer sqlDB.Close()
        
	//db..DropTableIfExists(&todo{})
	db.AutoMigrate(&todo{})

	log.Info("successfully connected to postgres")

log.Info("starting todolist api server")
router:=mux.NewRouter()
router.HandleFunc("/healthz",Healthz).Methods("GET")
router.HandleFunc("/readiness",readinessHandler).Methods("GET")
router.HandleFunc("/todo-completed",GetCompletedItems).Methods("GET")
router.HandleFunc("/todo-incomplete",GetIncompleteItems).Methods("GET")
router.HandleFunc("/todo/{id}",UpdateItem).Methods("POST")
router.HandleFunc("/todo/{id}",DeleteItem).Methods("DELETE")
router.HandleFunc("/todo",CreateItem).Methods("POST")
handler := cors.New(cors.Options{AllowedMethods: []string{"GET", "POST", "DELETE", "PATCH", "OPTIONS"},       }).Handler(router)
http.ListenAndServe(":8080",handler)
}
