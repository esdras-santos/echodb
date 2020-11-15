package echodb

import (
	"io/ioutil"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)


// Load/reload the database
func Load(location string, isAutoDump bool) *echodb {
	db := echodb{location, isAutoDump, nil}
	db.Load(location,isAutoDump)
	return &db
}

type echodb struct {
	location string
	isAutoDump bool
	db       map[string]interface{}
}

// Loads, reloadds or changes the path to the db file
func (database *echodb) Load(location string, isAutoDump bool) bool {
	loca, err := filepath.Abs(location)
	handler(err)
	database.location = loca
	database.isAutoDump = isAutoDump
	_, err = os.Stat(database.location)
	if os.IsNotExist(err) {
		database.db = map[string]interface{}{}
	} else {
		database.loaddb()
	}

	return true
}

// Dump memory of db to the file
func (database *echodb) Dump(){
	data,err := json.Marshal(database.db)
	handler(err)
	err = ioutil.WriteFile(database.location,data, 0644)
	handler(err)
}

// Load/reload the json file
func (database *echodb) loaddb() {
	file,err := ioutil.ReadFile(database.location)
	handler(err)
	err = json.Unmarshal([]byte(file),&database.db)	
	handler(err)
}

// Write/Save the json automatically when "isAutoDump" is enabled
func (database *echodb) autoDumpDb() {
	if database.isAutoDump {
		database.Dump()
	}
}

// Set/add a key value to the database
func (database *echodb) Set(key, value string) bool{
	database.db[key] = value
	database.autoDumpDb()
	return true
}

// Get the value of a key
func (database *echodb) Get(key string) interface{} {
	value, ok := database.db[key]
	if ok {
		return value
	}
	return false
}
// Get all the keys in a slice of strings
func (database *echodb) GetAll() []string {
	
	keys := []string{}
	for k := range database.db {
		keys = append(keys, k)
	}
	return keys
}

// Return true is key exists in db, retrn false if not
func (database *echodb) Exists(key string) bool {
	
	_, ok := database.db[key]
	return ok
}

// Delete key
func (database *echodb) Remove(key string) bool {
	_, ok := database.db[key]
	if !ok {
		return false
	}
	delete(database.db, key)
	database.autoDumpDb()
	return true
}

// Get a total number of keys inside the db
func (database *echodb) TotalKeys() int {
	total := len(database.db)
	return total
}

// Add more information to a key value
func (database *echodb) Append(key, more string) bool {
	database.db[key] = database.db[key].(string) + more
	database.autoDumpDb()
	return true
}

// Create a list
func (database *echodb) CreateList(name string) bool {
	database.db[name] = []string{}
	database.autoDumpDb()
	return true
}

// Add a value to a list
func (database *echodb) ListAdd(name string, value string) bool {
	database.db[name] = append(database.db[name].([]string), value)
	database.autoDumpDb()
	return true
}

// Extend a list with a sequence
func (database *echodb) ListExtend(name string, seq []string) bool {
	database.db[name] = append(database.db[name].([]string), seq...)
	database.autoDumpDb()
	return true
}

// Return all values in a list
func (database *echodb) ListGetAll(name string) interface{} {
	return database.db[name]
}

// Return one value in a list
func (database *echodb) ListGet(name string, pos int) interface{} {
	return database.db[name].(string)[pos]
}

// Return range of values in a list
func (database *echodb) ListRange(name string, start, end int) []string {
	return database.db[name].([]string)[start:end]
}

// Remove a list
func (database *echodb) RemList(listName string) {
	delete(database.db,listName)
	database.autoDumpDb()
}

// Remove an element from a list by the value
func (database *echodb) RemElemByValue(listName, elem string){
	for i, elem := range database.db[listName].([]string){
		if database.db[listName].([]string)[i] == elem {
			database.db[listName] = append(database.db[listName].([]string)[:i-1], database.db[listName].([]string)[i:] ...)
		}
	}
	database.autoDumpDb()
}

// Remove an alement from a list by the position
func (database *echodb) RemElemByPos(listName string, pos int){
	database.db[listName] = append(database.db[listName].([]string)[:pos-1], database.db[listName].([]string)[pos:] ...)
	database.autoDumpDb()
}

// Return the length of the list
func (database *echodb) ListLen(listName string) int{
	total := len(database.db[listName].([]string))
	return total	
}

// Add more information to a value
func (database *echodb) ListAppend(listName string, pos int, more string){
	database.db[listName].([]string)[pos] = database.db[listName].([]string)[pos] + more
	database.autoDumpDb()
}

// Check if a value exists in a list
func (database *echodb) ListValueExists(listName, value string) bool{
	for i, item := range database.db[listName].([]string){
		if item == database.db[listName].([]string)[i]{
			return true
		}
	}
	return false
}

// Create a map 
func (database *echodb) MapCreate(mapName string){
	database.db[mapName] = map[string]string{}
	database.autoDumpDb()
}

// Add a key-value to the map 
func (database *echodb) MapAdd(mapName string, newPair map[string]string){
	database.db[mapName] = newPair
	database.autoDumpDb()
}

// Return value of a key in a map
func (database *echodb) MapGet(mapName, mapKey string) string{
	return database.db[mapName].(map[string]string)[mapKey]
}

// Return the entire map
func (database *echodb) MapGetAll(mapName string) map[string]string{
	return database.db[mapName].(map[string]string)
}

// Remove a map from the database
func (database *echodb) MapRemove(mapName string) bool{
	delete(database.db, mapName)
	database.autoDumpDb()
	return true
}

// Remove one key-value from the map
func (database *echodb) MapRemoveOne(mapName, mapKey string) bool{
	delete(database.db[mapName].(map[string]string),database.db[mapName].(map[string]string)[mapKey])
	database.autoDumpDb()
	return true
}

// Return all the keys from a map
func (database *echodb) MapKeys(mapName string) []string{
	keys := []string{}
	for k := range database.db[mapName].(map[string]string) {
		keys = append(keys, k)
	}
	return keys
}

//	Return all the values from a map
func (database *echodb) MapValues(mapName string) []string{
	values := []string{}
	for value := range database.db[mapName].(map[string]string){
		values = append(values, value )
	}
	return values
}

// Checks if a key exist in the map
func (database *echodb) MapExists(mapName, mapKey string) bool{
	_,ok := database.db[mapName].(map[string]interface{})[mapKey].(map[string]string)
	return ok 
}

// Merge two maps together into mapName1
func (database *echodb) MapMerge(mapName1, mapName2 string) bool{
	for i,value := range database.db[mapName2].(map[string]string){
		database.db[mapName2].(map[string]string)[i] = value
	}
	return true
}

// Delete all the database
func (database *echodb) DeleteDatabase() bool{
	database.db = map[string]interface{}{}
	database.autoDumpDb()
	return true
}

func handler(err error) {
	if err != nil {
		log.Panic(err)
	}
}
