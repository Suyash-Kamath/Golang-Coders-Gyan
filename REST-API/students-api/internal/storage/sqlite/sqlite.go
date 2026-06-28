package sqlite

// implementing the interface inside this

// inbuilt package
// database se kaam karnekeliye jo bhi interface hai wo provide karta hai
import (
	"database/sql"

	"github.com/Suyash-Kamath/Golang-Coders-Gyan/REST-API/students-api/internal/config"
	_ "github.com/mattn/go-sqlite3"
	// Dekho, we are not using driver ka variable , bass insitialisation keliye chahiye hota hai , we are using it indirectly behind the scene , so underscore lagaayaa 
)


type Sqlite struct {
	Db *sql.DB
}



func New(cfg *config.Config) (*Sqlite,error){
	// Open database , means connection
	// database ka driver pass and storage path
	db,err:=sql.Open("sqlite3",cfg.StoragePath)
	if err!=nil{
		return nil,err
	}

	_,err =db.Exec(`CREATE TABLE IF NOT EXISTS students(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	email TEXT,
	age INTEGER
	)`)
	// hey _,err:= me why did you remove the colon

	if err!=nil{
		return nil,err
	}

	return &Sqlite{
		Db: db,
	},nil

}


func (s *Sqlite) CreateStudent(name string,email string,age int64) (int64,error){

	stmt,err:=s.Db.Prepare("INSERT INTO students(name,email,age) VALUES (?,?,?)")
	// we dont pass the data directly , first we prepare the query and give the placeholder and then we bind the data , which prevenets the SQL injection attack

	if err!=nil{
		return 0,err
	}
	// stmt ko bhi close karna hota hai 
	defer stmt.Close()

	// as we have prepared the query so execute bhi karna hai 

	result,err:=stmt.Exec(name,email,age) // binds the stmt and replaces the placeholder

	if err!=nil{
		return 0,err
	}

	// result ke andhar , seeing the documentation , we get LastInsertID() , RowsAffected() function

	lastId,err:=result.LastInsertId()
	if err!=nil{
		return 0,err
		// int64 ka empty value is 0
	}
	return lastId,nil

}