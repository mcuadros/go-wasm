package wasm

import (
	"encoding/json"
	"errors"
	"fmt"
	"syscall/js"
)

type storageType int

const (
	localStorageType storageType = iota
	sessionStorageType
)

var (
	// LocalStorage is a Storage intance that represents `Window.localStorage`.
	LocalStorage = new(localStorageType)
	// SessionStorage is a Storage intance that represents `Window.sessionStorage`.
	SessionStorage = new(sessionStorageType)

	ErrKeyNotFound = errors.New("unable to find key")
)

// Storage provides access to a particular domain's session or local storage.
type Storage struct {
	js.Value
	Marshal   func(interface{}) ([]byte, error)
	Unmarshal func(data []byte, v interface{}) error
}

func new(t storageType) *Storage {
	var object string
	switch t {
	case localStorageType:
		object = "localStorage"
	case sessionStorageType:
		object = "sessionStorage"
	default:
		panic("unexpected storage type")
	}

	s := &Storage{}
	s.Value = js.Global().Get(object)
	if !s.Value.Truthy() {
		panic(fmt.Sprintf("unexpected error unable to find window.%s", object))
	}

	s.Marshal = json.Marshal
	s.Unmarshal = json.Unmarshal

	return s
}

// SetItem when passed a key name and value, will add that key to the sStorage,
// or update that key's value if it already exists.
func (s *Storage) SetItem(key string, value interface{}) error {
	data, err := s.Marshal(value)
	if err != nil {
		return err
	}

	s.Call("setItem", key, string(data))
	return nil
}

// GetItem return that key's value, or ErrKeyNotFound if the key does not exist.
func (s *Storage) GetItem(key string) (interface{}, error) {
	item := s.Call("getItem", key)
	if !item.Truthy() {
		return nil, fmt.Errorf("%w: %q", ErrKeyNotFound, key)
	}

	var value interface{}
	return value, s.Unmarshal([]byte(item.String()), &value)
}

// RemoveItem when passed a key name, will remove that key from the given
// Storage object if it exists.
func (s *Storage) RemoveItem(key string) {
	s.Call("removeItem", key)
}

// Clear clears all keys stored in a given Storage object.
func (s *Storage) Clear() {
	s.Call("clear")
}
