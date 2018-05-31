// Package container defined some variables and interfaces related to container library.
package container

import "errors"

// Lesser implements how data to be compared.
type Lesser interface {
	// if kv is less than receiver by comparing corresponding field values, it returns false, otherwise true.
	Less(kv interface{}) bool
}

// Finder implements how to exactly match the data to be found.
type Finder interface {

	// Find tests whether the key matches the value to be found, if the values is exactly
	// what we want to find, returns true, otherwise false.
	Find(key interface{}) bool
}

// Setter is used for updating values.
type Setter interface {
	// Set sets fields values pointed by receiver to v.
	Set(v interface{})
}

// Interface implements how to insert, find, and modify Data in a container. Note that the key should be used as a basis for all Data operation related
// with tree, the key used to find Data should also be used to determine the direction of insertion, as well in modifing Data, also, because of needing
// to modify Data in the tree, all methods must be implemented with pointer receiver.
type Interface interface {
	// Less is used to determine the direction of walking in the tree, Note that the kv must be unique, When using structure field as key,
	// as long as the key is identical to corresponding field of receiver, regardless of whether the structure values are identical,
	// the Data to be inserted onto tree will be treated as already exsit in the tree. If Less returns true, it means that we along the path
	//from that node to its right child when performing operation of inserting, finding or modifying Data, and vice versa.
	Lesser

	Finder

	Setter
}

var (
	// ErrNotExist means that The data to found by key doesn't exist.
	ErrNotExist = errors.New("data doesn't exist")

	// ErrDataExists is returned by Insert when data to be inserted already exists.
	ErrDataExists = errors.New("Inserted data already exists")

	// ErrEmptyList means that the list is empty.
	ErrEmptyList = errors.New("List is empty")

	// ErrEmptyTree means that the tree is empty.
	ErrEmptyTree = errors.New("Tree is empty")
)
