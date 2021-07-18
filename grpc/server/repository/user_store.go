package repository

import (
	"fmt"
	model2 "github.com/care0717/deepthought-api/grpc/server/model"
	"sync"
)

type UserStore interface {
	Save(user *model2.User) error
	Find(username string) (*model2.User, error)
}

type inMemoryUserStore struct {
	mutex sync.RWMutex
	users map[string]*model2.User
}

func NewInMemoryUserStore() UserStore {
	return &inMemoryUserStore{
		users: make(map[string]*model2.User),
	}
}

func (store *inMemoryUserStore) Save(user *model2.User) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if store.users[user.Username] != nil {
		return fmt.Errorf("already exists %s", user.Username)
	}

	store.users[user.Username] = user.Clone()
	return nil
}

func (store *inMemoryUserStore) Find(username string) (*model2.User, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	user := store.users[username]
	if user == nil {
		return nil, nil
	}

	return user.Clone(), nil
}
