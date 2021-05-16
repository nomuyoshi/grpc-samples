package main

import (
	"context"
	pbProject "todo/proto/project"
	pbTask "todo/proto/task"
	"todo/shared/md"
	taskDomain "todo/task/domain"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type TaskService struct {
	db            *gorm.DB
	projectClient pbProject.ProjectServiceClient
	pbTask.UnimplementedTaskServiceServer
}

func (s *TaskService) CreateTask(ctx context.Context, req *pbTask.CreateTaskRequest) (*pbTask.CreateTaskResponse, error) {
	name := req.GetName()
	if name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is empty")
	}

	prResp, err := s.projectClient.FindProject(ctx, &pbProject.FindProjectRequest{ProjectId: req.GetProjectId()})
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	userID := md.GetUserIDFromContext(ctx)
	t := taskDomain.NewTask(name, userID, prResp.Project.Id)
	taskDomain.NewTaskRepository(s.db).Save(t)
	return &pbTask.CreateTaskResponse{Task: buildPbTask(t)}, nil
}

func (s *TaskService) FindTasks(ctx context.Context, _ *emptypb.Empty) (*pbTask.FindTasksResponse, error) {
	userID := md.GetUserIDFromContext(ctx)
	tasks := taskDomain.NewTaskRepository(s.db).FindByUserID(userID)
	var pbTasks []*pbTask.Task
	for _, t := range tasks {
		pbTasks = append(pbTasks, buildPbTask(&t))
	}
	return &pbTask.FindTasksResponse{Tasks: pbTasks}, nil
}

func (s *TaskService) FindProjectTasks(
	ctx context.Context,
	req *pbTask.FindProjectTasksRequest,
) (*pbTask.FindProjectTasksResponse, error) {
	projectID := req.GetProjectId()

	if projectID == 0 {
		return nil, status.Error(codes.InvalidArgument, "projectID is empty")
	}

	userID := md.GetUserIDFromContext(ctx)
	tasks := taskDomain.NewTaskRepository(s.db).FindByUserIDAndProjectID(userID, projectID)

	var pbTasks []*pbTask.Task
	for _, t := range tasks {
		pbTasks = append(pbTasks, buildPbTask(&t))
	}
	return &pbTask.FindProjectTasksResponse{Tasks: pbTasks}, nil
}

func (s *TaskService) UpdateTask(
	ctx context.Context,
	req *pbTask.UpdateTaskRequest,
) (*pbTask.UpdateTaskResponse, error) {
	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "name is empty")
	}

	if req.GetStatus() == pbTask.Status_UNKNOWN {
		return nil, status.Error(codes.InvalidArgument, "status is unknown")
	}

	repo := taskDomain.NewTaskRepository(s.db)
	userID := md.GetUserIDFromContext(ctx)
	task := repo.FindByID(req.TaskId, userID)
	if task == nil {
		return nil, status.Error(codes.NotFound, "task not found")
	}

	task.Name = req.GetName()
	task.Status = req.GetStatus()
	repo.Save(task)

	return &pbTask.UpdateTaskResponse{Task: buildPbTask(task)}, nil
}

func buildPbTask(t *taskDomain.Task) *pbTask.Task {
	var createdAt, updatedAt *timestamppb.Timestamp
	if t.CreatedAt != nil {
		createdAt, _ = ptypes.TimestampProto(*t.CreatedAt)
	}
	if t.UpdatedAt != nil {
		updatedAt, _ = ptypes.TimestampProto(*t.UpdatedAt)
	}

	return &pbTask.Task{
		Id:        t.ID,
		Name:      t.Name,
		Status:    t.Status,
		ProjectId: t.ProjectID,
		UserId:    t.UserID,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
