package db

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

/*
	Library boltdb/bolt is is chosen as persistent lightweight key-value database
	It is based on a principle of "Buckets" which are collections of key/value pairs within the database, where keys are unique
*/

//VARIABLE DEFINITION

//Structure definition of database
//Application structure will work as a key/value map, using NAME field as key and all the structure as value
type Application struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Author  string `json:"author"`
}

var (
	db *bolt.DB
)

//database file name
var dbFile = "app.db"

//our bucket name
var bucketName = []byte("MyBucket")

//init function, starting database
//it will create the bucket if it doesn't exist
func init() {
	var err error
	//open database file r/w
	db, err = bolt.Open(dbFile, 0600, nil)
	if err != nil {
		fmt.Errorf("could not open db file, %v", err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return fmt.Errorf("could not create bucket: %v", err)
		}
		return nil
	})
	if err != nil {
		fmt.Errorf("could not setup buckets, %v", err)
	}
	fmt.Println("DB Setup Done")
}

//Function to return all information of a single app based on key (app name)
func GetApp(appName string) (app []byte, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		app = tx.Bucket(bucketName).Get([]byte(appName))
		fmt.Printf("corresponding values for key  %s : %s\n", appName, app)
		return err
	})
	if err != nil {
		log.Fatal(err)
	}
	return app, err
}

//Function to return all information in database
func GetApps() (apps []byte, err error) {
	var a []Application
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		b.ForEach(func(k, v []byte) error {
			var app Application
			err = json.Unmarshal(v, &app)
			if err != nil {
				return fmt.Errorf("could not unmarshal")
			}
			a = append(a, app)
			return nil
		})
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	result, err := json.Marshal(a)
	if err != nil {
		fmt.Errorf("could not marshal")
	}
	apps = []byte(result)
	return apps, err
}

//Delete whole information in database
func DeleteDB() error {
	err := db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(bucketName)
		if err != nil {
			return fmt.Errorf("could not delete db: %v", err)
		}
		fmt.Printf("BD deleted ")
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

//Add a new application in database
func AddApp(newApp Application) (err error) {
	encoded, err := json.Marshal(newApp)
	if err != nil {
		return err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket(bucketName).Put([]byte(newApp.Name), encoded)
		if err != nil {
			return fmt.Errorf("could not add app: %v", err)
		}
		return nil
	})
	fmt.Println("App added", newApp.Name)
	return err
}

//Delete the application
func DeleteApp(appName string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket(bucketName).Delete([]byte(appName))

		if err != nil {
			return fmt.Errorf("could not delete app: %v", err)
		}
		return nil
	})
	fmt.Println("App deleted", appName)
	return err
}
