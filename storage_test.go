package wasm

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorageGetItem(t *testing.T) {
	defer SessionStorage.Clear()

	err := SessionStorage.SetItem("foo", 42)
	assert.Nil(t, err)

	v, err := SessionStorage.GetItem("foo")
	assert.Nil(t, err)
	assert.Equal(t, 42.0, v)
}

func TestStorageGetItem_NotFound(t *testing.T) {
	defer SessionStorage.Clear()

	v, err := SessionStorage.GetItem("bar")
	assert.True(t, errors.Is(err, ErrKeyNotFound))
	assert.Equal(t, nil, v)
}

func TestStorageRemoveItem(t *testing.T) {
	defer SessionStorage.Clear()

	err := SessionStorage.SetItem("foo", 42)
	assert.Nil(t, err)

	SessionStorage.RemoveItem("foo")

	_, err = SessionStorage.GetItem("foo")
	assert.True(t, errors.Is(err, ErrKeyNotFound))
}

func TestStorageClear(t *testing.T) {
	err := SessionStorage.SetItem("foo", 42)
	assert.Nil(t, err)

	err = SessionStorage.SetItem("bar", 42)
	assert.Nil(t, err)

	SessionStorage.Clear()

	_, err = SessionStorage.GetItem("foo")
	assert.True(t, errors.Is(err, ErrKeyNotFound))

	_, err = SessionStorage.GetItem("bar")
	assert.True(t, errors.Is(err, ErrKeyNotFound))
}
