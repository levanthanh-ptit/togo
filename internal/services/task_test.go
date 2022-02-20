package services

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"togo/internal/domain"
	"togo/internal/repository"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/mock"
)

func Test_taskService_Create(t *testing.T) {
	type fields struct {
		userRepo      repository.UserRepository
		taskRepo      repository.TaskRepository
		taskLimitRepo repository.TaskLimitRepository
	}
	type args struct {
		ctx  context.Context
		task *domain.Task
	}
	// Mocks
	notFoundUserID := uint(10)
	user := &domain.User{
		ID:          1,
		FullName:    faker.Name(),
		Username:    faker.Username(),
		Password:    faker.Password(),
		TasksPerDay: 2,
	}
	user2 := &domain.User{
		ID:          2,
		FullName:    faker.Name(),
		Username:    faker.Username(),
		Password:    faker.Password(),
		TasksPerDay: 2,
	}
	taskInput := &domain.Task{
		UserID:  user.ID,
		Content: "valid content",
	}
	task := &domain.Task{
		ID:      1,
		UserID:  user.ID,
		Content: "valid content",
	}
	userRepo := new(mockUserRepository)
	userRepo.On("FindOne", &domain.User{ID: notFoundUserID}).Return(nil, domain.ErrUserNotFound)
	userRepo.On("FindOne", &domain.User{ID: user.ID}).Return(user, nil)
	userRepo.On("FindOne", &domain.User{ID: user2.ID}).Return(user2, nil)
	taskLimitRepo := new(mockTaskLimitRepository)
	taskLimitRepo.On("Increase", user.ID, user.TasksPerDay).Return(1, nil)
	taskLimitRepo.On("Increase", user2.ID, user2.TasksPerDay).Return(0, domain.ErrTaskLimitExceed)
	taskRepo := new(mockTaskRepository)
	taskRepo.On("Create", taskInput).Return(task, nil)
	brokenTaskRepo := new(mockTaskRepository)
	brokenTaskRepo.On("Create", mock.Anything).Return(nil, errors.New("invalid"))
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Task
		wantErr bool
	}{
		{
			name:   "user not found",
			fields: fields{userRepo, taskRepo, taskLimitRepo},
			args: args{
				context.Background(),
				&domain.Task{
					UserID:  notFoundUserID,
					Content: "text",
				},
			},
			wantErr: true,
		},
		{
			name:   "tasks limit exceed",
			fields: fields{userRepo, taskRepo, taskLimitRepo},
			args: args{
				context.Background(),
				&domain.Task{
					UserID:  user2.ID,
					Content: "text",
				},
			},
			wantErr: true,
		},
		{
			name:   "task saving failed",
			fields: fields{userRepo, brokenTaskRepo, taskLimitRepo},
			args: args{
				context.Background(),
				taskInput,
			},
			wantErr: true,
		},
		{
			name:   "task create successfully",
			fields: fields{userRepo, taskRepo, taskLimitRepo},
			args: args{
				context.Background(),
				taskInput,
			},
			want: task,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := taskService{
				userRepo:      tt.fields.userRepo,
				taskRepo:      tt.fields.taskRepo,
				taskLimitRepo: tt.fields.taskLimitRepo,
			}
			got, err := s.Create(tt.args.ctx, tt.args.task)
			if (err != nil) != tt.wantErr {
				t.Errorf("taskService.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("taskService.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_taskService_UpdateByID(t *testing.T) {
	type fields struct {
		userRepo      repository.UserRepository
		taskRepo      repository.TaskRepository
		taskLimitRepo repository.TaskLimitRepository
	}
	type args struct {
		ctx    context.Context
		id     uint
		update *domain.Task
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Task
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := taskService{
				userRepo:      tt.fields.userRepo,
				taskRepo:      tt.fields.taskRepo,
				taskLimitRepo: tt.fields.taskLimitRepo,
			}
			got, err := s.UpdateByID(tt.args.ctx, tt.args.id, tt.args.update)
			if (err != nil) != tt.wantErr {
				t.Errorf("taskService.UpdateByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("taskService.UpdateByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_taskService_FindByUserID(t *testing.T) {
	type fields struct {
		userRepo      repository.UserRepository
		taskRepo      repository.TaskRepository
		taskLimitRepo repository.TaskLimitRepository
	}
	type args struct {
		ctx    context.Context
		userID uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*domain.Task
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := taskService{
				userRepo:      tt.fields.userRepo,
				taskRepo:      tt.fields.taskRepo,
				taskLimitRepo: tt.fields.taskLimitRepo,
			}
			got, err := s.FindByUserID(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("taskService.FindByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("taskService.FindByUserID() = %v, want %v", got, tt.want)
			}
		})
	}
}
