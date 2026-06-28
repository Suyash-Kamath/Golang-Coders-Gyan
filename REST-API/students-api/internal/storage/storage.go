package storage

// using interface concept , so that we can make it plug-in type
// also during testing , we can switch it to fake database

type Storage interface{
	CreateStudent(name string,email string, age int) (int64,error)

}


