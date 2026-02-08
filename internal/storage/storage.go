package storage

import (
	"sync"
	"practice2/internal/models"
)


type TaskStorage struct {
	tasks  map[int]models.Task
	nextID int
	mu     sync.RWMutex
}

func NewTaskStorage() *TaskStorage {
	return &TaskStorage{
		tasks:  make(map[int]models.Task),
		nextID: 1,
	}
}

func (s *TaskStorage) Create(title string) models.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	task := models.Task{
		ID:    s.nextID,
		Title: title,
		Done:  false,
	}
	s.tasks[s.nextID] = task
	s.nextID++

	return task
}

func (s *TaskStorage) GetByID(id int) (models.Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[id]
	return task, exists
}

func (s *TaskStorage) GetAll() []models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]models.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

func (s *TaskStorage) GetByStatus(done bool) []models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]models.Task, 0)
	for _, task := range s.tasks {
		if task.Done == done {
			tasks = append(tasks, task)
		}
	}
	return tasks
}

func (s *TaskStorage) Update(id int, done bool) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[id]
	if !exists {
		return false
	}

	task.Done = done
	s.tasks[id] = task
	return true
}


func (s *TaskStorage) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[id]; !exists {
		return false
	}

	delete(s.tasks, id)
	return true
}
