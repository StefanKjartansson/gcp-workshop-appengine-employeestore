// +build !appengine

package employees

import (
	"errors"
	"math/rand"
	"net/http"
	"time"

	"golang.org/x/net/context"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func randString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

type MemoryStore struct {
	store map[string]*Employee
}

func NewMemoryStore() *MemoryStore {
	ms := MemoryStore{
		store: make(map[string]*Employee),
	}
	return &ms
}

func (ms *MemoryStore) getContextFromRequest(req *http.Request) context.Context {
	return context.TODO()
}

func (ms *MemoryStore) List(ctx context.Context) ([]*Employee, error) {
	l := []*Employee{}
	for key, e := range ms.store {
		e.Id = key
		l = append(l, e)
	}
	return l, nil
}

func (ms *MemoryStore) Get(ctx context.Context, id string) (*Employee, error) {
	e, ok := ms.store[id]
	if !ok {
		return nil, errors.New("No employee found")
	}
	return e, nil
}

func (ms *MemoryStore) Put(ctx context.Context, e Employee) (string, error) {
	key := randString(16)
	e.HireDate = time.Now()
	ms.store[key] = &e
	return key, nil
}

var _ EmployeeStore = (*MemoryStore)(nil)
