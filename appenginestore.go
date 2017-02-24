// +build appengine

package employees

import (
	"net/http"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
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
	e.Id = id
	return &e, nil
}

func (a *AppEngineStore) Put(ctx context.Context, e Employee) (string, error) {
	u := user.Current(ctx)
	if u != nil {
		e.Account = u.String()
	}
	e.HireDate = time.Now()
	key, err := datastore.Put(ctx, datastore.NewIncompleteKey(ctx, collectionName, nil), &e)
	if err != nil {
		log.Errorf(ctx, "could not put into datastore: %v", err)
		return "", err
	}
	return key.Encode(), nil
}

var _ EmployeeStore = (*AppEngineStore)(nil)
