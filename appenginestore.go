package employees

import (
	"context"
	"net/http"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"
)

const collectionName = "Employee"

type AppEngineStore struct {
}

func (a *AppEngineStore) getContextFromRequest(req *http.Request) context.Context {
	return appengine.NewContext(req)
}

func (a *AppEngineStore) List(ctx context.Context) ([]*Employee, error) {
	q := datastore.NewQuery(collectionName)
	l := []*Employee{}
	for t := q.Run(ctx); ; {
		var e Employee
		key, err := t.Next(&e)
		if err == datastore.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		e.Id = key.Encode()
		l = append(l, &e)
	}
	return l, nil
}

func (a *AppEngineStore) Get(ctx context.Context, id string) (*Employee, error) {
	key, err := datastore.DecodeKey(id)
	if err != nil {
		return nil, err
	}
	var e Employee
	err = datastore.Get(ctx, key, &e)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (a *AppEngineStore) Put(ctx context.Context, e Employee) (string, error) {
	e.Account = user.Current(ctx).String()
	e.HireDate = time.Now()
	key, err := datastore.Put(ctx, datastore.NewIncompleteKey(ctx, collectionName, nil), &e)
	if err != nil {
		return "", err
	}
	return key.Encode(), nil
}

var _ EmployeeStore = (*AppEngineStore)(nil)
