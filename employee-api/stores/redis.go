package stores

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

type RedisStore struct {
	DBUrl string
	DB    *redis.Client
}

func NewRedisStore(dbURL string) *RedisStore {
	r := RedisStore{dbURL, nil}
	return &r
}

func (r *RedisStore) Connect() error {
	r.DB = redis.NewClient(&redis.Options{
		Addr:     r.DBUrl,
		Password: "",
		DB:       0,
	})

	_, err := r.DB.Ping().Result()
	if err != nil {
		return err
	}

	return nil
}

func (r RedisStore) Create(listID string, t *Employee) error {
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	t.ID = id.String()

	tb, err := json.Marshal(t)
	if err != nil {
		return err
	}

	err = r.DB.Set(fmt.Sprintf("%s-%s", listID, t.ID), tb, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r RedisStore) Delete(listID string, id string) error {
	err := r.DB.Del(fmt.Sprintf("%s-%s", listID, id)).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r RedisStore) Update(listID string, id string, newT *Employee) (*Employee, error) {
	tb, err := r.DB.Get(fmt.Sprintf("%s-%s", listID, id)).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var t Employee
	json.Unmarshal(tb, &t)

	if newT.First_Name != "" {
		t.First_Name = newT.First_Name
	}
	if newT.Last_Name != "" {
		t.Last_Name = newT.Last_Name
	}
	if newT.Department != "" {
		t.Department = newT.Department
	}
	if newT.Salary != 0 {
		t.Salary = newT.Salary
	}
	if newT.Age != 0 {
		t.Age = newT.Age
	}

	tn, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	err = r.DB.Set(fmt.Sprintf("%s-%s", listID, t.ID), tn, 0).Err()
	return &t, nil

}

func (r RedisStore) Get(listID string, id string) (*Employee, error) {
	tb, err := r.DB.Get(fmt.Sprintf("%s-%s", listID, id)).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var t Employee
	json.Unmarshal(tb, &t)

	return &t, nil
}

func (r RedisStore) Clear(listID string) error {
	keys, err := r.DB.Keys(fmt.Sprintf("%s-*", listID)).Result()
	if err != nil {
		return err
	}

	for _, k := range keys {
		r.DB.Del(k).Err()
	}
	return nil
}

func (r RedisStore) List(listID string) ([]Employee, error) {
	keys, err := r.DB.Keys(fmt.Sprintf("%s-*", listID)).Result()
	if err != nil {
		return nil, err
	}
	result := []Employee{}
	for _, k := range keys {
		tb, err := r.DB.Get(k).Bytes()
		if err == redis.Nil {
			continue
		}
		if err != nil {
			return nil, err
		}
		var t Employee
		json.Unmarshal(tb, &t)
		result = append(result, t)
	}
	return result, nil
}
