package hashtable

import (
	responsestr "main/Struct/responseStr"
)

type HashTable struct {
	elem []responsestr.WeatherToFront
	size int
}

func (ht *HashTable) Init(size int) {
	ht.elem = make([]responsestr.WeatherToFront, size)
	ht.size = size
}

// For INT tables
func newHashINT(key int, size int) int {
	hash := 0
	for key != 0 {
		digit := key % 10
		hash += digit
		key /= 10
	}
	return hash % size
}

// For INT tables
func (ht *HashTable) InsertINT(key int, value responsestr.WeatherToFront) {
	index := newHashINT(key, ht.size)
	ht.elem[index] = value
}

// For INT tables
func (ht *HashTable) FindINT(key int) (responsestr.WeatherToFront, bool) {
	index := newHashINT(key, ht.size)

	if ht.elem[index].Id == key {

		return ht.elem[index], true
	}
	return responsestr.WeatherToFront{}, false
}

// For INT tables
func newHashSTRING(key string, size int) int {
	hash := 0
	for _, char := range key {
		hash += int(char)
	}
	return hash % size
}

//====================================================================================\\

// For STRING tables
func (ht *HashTable) InsertSTRING(key string, value responsestr.WeatherToFront) {
	index := newHashSTRING(key, ht.size)
	ht.elem[index] = value
}

// For STRING tables
func (ht *HashTable) FindSTRING(key string) (responsestr.WeatherToFront, bool) {
	index := newHashSTRING(key, ht.size)

	if ht.elem[index].Name == key {

		return ht.elem[index], true
	}
	return responsestr.WeatherToFront{}, false
}
