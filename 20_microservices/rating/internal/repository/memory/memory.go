package memory

import (
	"context"
	"sync"

	"movie.com/rating/internal/repository"
	"movie.com/rating/pkg/model"
)

type Repository struct {
	sync.RWMutex
	data map[model.RecordType]map[model.RecordID][]model.Rating
}

func New() *Repository {
	return &Repository{data: map[model.RecordType]map[model.RecordID][]model.Rating{}}
}

func (r *Repository) Get(_ context.Context, recordID model.RecordID, recordType model.RecordType) ([]model.Rating, error) {
	r.RLock()
	defer r.RUnlock()
	if _, ok := r.data[recordType]; !ok {
		return nil, repository.ErrNotFound
	}
	if all, ok := r.data[recordType][recordID]; !ok || len(all) == 0 {
		return nil, repository.ErrNotFound
	}
	return r.data[recordType][recordID], nil
}

func (r *Repository) Put(_ context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.data[recordType]; !ok {
		r.data[recordType] = map[model.RecordID][]model.Rating{}
	}
	r.data[recordType][recordID] =
		append(r.data[recordType][recordID], *rating)
	return nil
}
