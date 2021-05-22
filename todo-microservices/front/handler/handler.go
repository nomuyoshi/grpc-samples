package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"todo/front/session"
	"todo/front/template"
	pbProject "todo/proto/project"
	pbTask "todo/proto/task"
	pbUser "todo/proto/user"
	"todo/support"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/gorilla/mux"
)

type FrontServer struct {
	TaskClient    pbTask.TaskServiceClient
	ProjectClient pbProject.ProjectServiceClient
	UserClient    pbUser.UserServiceClient
	SessionStore  session.Store
}

// signup
func (s *FrontServer) ViewSignup(w http.ResponseWriter, r *http.Request) {
	template.Render(w, "signup.html", nil)
}

func (s *FrontServer) SignUp(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	resp, err := s.UserClient.CreateUser(r.Context(), &pbUser.CreateUserRequest{
		Email:    r.Form.Get("email"),
		Password: []byte(r.Form.Get("password")),
	})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	sessionID := session.ID()
	s.SessionStore.Set(sessionID, resp.GetUser().GetId())
	session.SetSessionIDToResponse(w, sessionID)
	http.Redirect(w, r, "/", http.StatusFound)
}

// login
func (s *FrontServer) ViewLogin(w http.ResponseWriter, r *http.Request) {
	template.Render(w, "login.html", nil)
}

func (s *FrontServer) Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	resp, err := s.UserClient.VerifyUser(r.Context(), &pbUser.VerifyUserRequest{
		Email:    r.Form.Get("email"),
		Password: []byte(r.Form.Get("password")),
	})

	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	sessionID := session.ID()
	s.SessionStore.Set(sessionID, resp.GetUser().GetId())
	session.SetSessionIDToResponse(w, sessionID)
	http.Redirect(w, r, "/", http.StatusFound)
}

// logout
func (s *FrontServer) Logout(w http.ResponseWriter, r *http.Request) {
	sessionID := session.GetSessionIDFromRequest(r)
	s.SessionStore.Delete(sessionID)
	session.DeleteSessionIDFromResponse(w)
	http.Redirect(w, r, "/login", http.StatusFound)
}

func (s *FrontServer) ViewProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	projectResp, err := s.ProjectClient.FindProject(ctx, &pbProject.FindProjectRequest{
		ProjectId: projectID,
	})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	taskResp, err := s.TaskClient.FindProjectTasks(ctx, &pbTask.FindProjectTasksRequest{
		ProjectId: projectID,
	})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	var taskRows []*TaskRow
	for _, task := range taskResp.Tasks {
		taskRows = append(taskRows, &TaskRow{task, projectResp.Project})
	}
	user := support.GetUserFromContext(ctx)
	template.Render(w, "project.html", &ProjectContent{
		PageName:     "Project",
		IsLoggedIn:   true,
		UserEmail:    user.Email,
		TaskStatuses: taskStatuses,
		Project:      projectResp.Project,
		TaskRows:     taskRows,
	})
}

type ProjectContent struct {
	PageName     string
	IsLoggedIn   bool
	TaskStatuses []TaskStatus
	UserEmail    string
	Project      *pbProject.Project
	TaskRows     []*TaskRow
}

type TaskRow struct {
	task    *pbTask.Task
	project *pbProject.Project
}

func (r *TaskRow) ID() uint64 {
	return r.task.Id
}

func (r *TaskRow) Name() string {
	return r.task.Name
}

func (r *TaskRow) ProjectName() string {
	return r.project.Name
}

func (r *TaskRow) Status() int32 {
	return int32(r.task.Status)
}

func (r *TaskRow) StatusName() string {
	return r.task.Status.String()
}

var taskStatuses = []TaskStatus{
	TaskStatus(pbTask.Status_WAITING),
	TaskStatus(pbTask.Status_WORKING),
	TaskStatus(pbTask.Status_COMPLETED),
}

type TaskStatus pbTask.Status

func (s TaskStatus) Status() int32 {
	return int32(s)
}

func (s *TaskStatus) StatusName() string {
	return pbTask.Status_name[s.Status()]
}

func (s *FrontServer) CreateProject(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if _, err := s.ProjectClient.CreateProject(r.Context(), &pbProject.CreateProjectRequest{
		Name: r.Form.Get("name"),
	}); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *FrontServer) UpdateProject(w http.ResponseWriter, r *http.Request) {
	projectIDStr := mux.Vars(r)["id"]
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	r.ParseForm()
	if _, err := s.ProjectClient.UpdateProject(r.Context(), &pbProject.UpdateProjectRequest{
		ProjectId:   projectID,
		ProjectName: r.Form.Get("name"),
	}); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/project/"+projectIDStr, http.StatusFound)
}

func (s *FrontServer) CreateTask(
	w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	projectIDStr := r.Form.Get("project_id")
	projectID, err := strconv.
		ParseUint(projectIDStr, 10, 64)
	if err != nil {
		http.Error(w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest)
		return
	}
	if _, err := s.TaskClient.CreateTask(
		r.Context(),
		&pbTask.CreateTaskRequest{
			Name:      r.Form.Get("name"),
			ProjectId: projectID,
		}); err != nil {
		http.Error(w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest)
		return
	}
	redirectURL := "/"
	if strings.Contains(r.Referer(), "/project/") {
		redirectURL += "project/" + projectIDStr
	}
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (s *FrontServer) UpdateTask(w http.ResponseWriter, r *http.Request) {
	taskID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	r.ParseForm()
	status, err := strconv.ParseUint(r.Form.Get("status_id"), 10, 32)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	resp, err := s.TaskClient.UpdateTask(r.Context(), &pbTask.UpdateTaskRequest{
		TaskId: taskID,
		Name:   r.Form.Get("name"),
		Status: pbTask.Status(status),
	})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	redirectURL := "/"
	if strings.Contains(r.Referer(), "/project/") {
		redirectURL += fmt.Sprintf("project/%d", resp.Task.ProjectId)
	}
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (s *FrontServer) ViewHome(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var in empty.Empty
	projectResp, err := s.ProjectClient.FindProjects(ctx, &in)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	taskResp, err := s.TaskClient.FindTasks(ctx, &in)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	idToPj := make(map[uint64]*pbProject.Project)
	for _, project := range projectResp.GetProjects() {
		idToPj[project.GetId()] = project
	}
	var taskRows []*TaskRow
	for _, task := range taskResp.GetTasks() {
		project := idToPj[task.GetProjectId()]
		taskRows = append(taskRows, &TaskRow{
			task, project})
	}
	user := support.GetUserFromContext(ctx)
	template.Render(w, "home.html",
		&HomeContent{
			PageName:     "Home",
			IsLoggedIn:   true,
			UserEmail:    user.Email,
			TaskStatuses: taskStatuses,
			Projects:     projectResp.Projects,
			TaskRows:     taskRows,
		})
}

type HomeContent struct {
	PageName     string
	IsLoggedIn   bool
	TaskStatuses []TaskStatus
	UserEmail    string
	Projects     []*pbProject.Project
	TaskRows     []*TaskRow
}
