package main

import (
	"context"
	projectDomain "todo/project/domain"
	pbProject "todo/proto/project"
	"todo/shared/md"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type projectService struct {
	db *gorm.DB
	pbProject.UnimplementedProjectServiceServer
}

func (s *projectService) CreateProject(
	ctx context.Context,
	req *pbProject.CreateProjectRequest,
) (*pbProject.CreateProjectResponse, error) {
	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "project name is empty")
	}

	userID := md.GetUserIDFromContext(ctx)
	project := projectDomain.NewProject(req.GetName(), userID)
	projectDomain.NewProjectRepository(s.db).Save(project)

	return &pbProject.CreateProjectResponse{
		Project: buildPbProject(project),
	}, nil
}

func (s *projectService) FindProject(
	ctx context.Context,
	req *pbProject.FindProjectRequest,
) (*pbProject.FindProjectResponse, error) {
	userID := md.GetUserIDFromContext(ctx)
	project := projectDomain.NewProjectRepository(s.db).FindByID(userID, req.GetProjectId())
	if project == nil {
		return nil, status.Error(codes.NotFound, "project not found")
	}

	return &pbProject.FindProjectResponse{Project: buildPbProject(project)}, nil
}

func (s *projectService) FindProjects(ctx context.Context, _ *emptypb.Empty) (*pbProject.FindProjectsResponse, error) {
	userID := md.GetUserIDFromContext(ctx)

	projects := projectDomain.NewProjectRepository(s.db).FindAll(userID)
	var pbProjects []*pbProject.Project
	for _, p := range projects {
		pbProjects = append(pbProjects, buildPbProject(&p))
	}

	return &pbProject.FindProjectsResponse{Projects: pbProjects}, nil
}

func (s *projectService) UpdateProject(
	ctx context.Context,
	req *pbProject.UpdateProjectRequest,
) (*pbProject.UpdateProjectResponse, error) {
	if req.GetProjectName() == "" {
		return nil, status.Error(codes.InvalidArgument, "project name is empty")
	}
	userID := md.GetUserIDFromContext(ctx)
	repo := projectDomain.NewProjectRepository(s.db)
	project := repo.FindByID(userID, req.GetProjectId())
	if project == nil {
		return nil, status.Error(codes.NotFound, "project not found")
	}

	project.Name = req.GetProjectName()
	repo.Save(project)

	return &pbProject.UpdateProjectResponse{Project: buildPbProject(project)}, nil
}

func buildPbProject(p *projectDomain.Project) *pbProject.Project {
	var createdAt *timestamppb.Timestamp
	if p.CreatedAt != nil {
		createdAt, _ = ptypes.TimestampProto(*p.CreatedAt)
	}

	return &pbProject.Project{
		Id:        p.ID,
		Name:      p.Name,
		UserId:    p.UserID,
		CreatedAt: createdAt,
	}
}
