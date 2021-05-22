package main

import (
	"log"
	"net/http"
	"os"

	"todo/front/handler"
	"todo/front/middleware"
	"todo/front/session"
	pbProject "todo/proto/project"
	pbTask "todo/proto/task"
	pbUser "todo/proto/user"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

const port = ":8080"

func main() {
	projectClient := pbProject.NewProjectServiceClient(getGRPCConn(os.Getenv("PROJECT_SERVICE_ADDR")))
	taskClient := pbTask.NewTaskServiceClient(getGRPCConn(os.Getenv("TASK_SERVICE_ADDR")))
	userClient := pbUser.NewUserServiceClient(getGRPCConn(os.Getenv("USER_SERVICE_ADDR")))
	sessionStore := session.NewStoreOnMemory()
	frontSrv := &handler.FrontServer{
		ProjectClient: projectClient,
		TaskClient:    taskClient,
		UserClient:    userClient,
		SessionStore:  sessionStore,
	}

	r := mux.NewRouter()
	// ミドルウェア
	r.Use(middleware.Tracing)
	r.Use(middleware.Logging)
	auth := middleware.NewAuthentication(userClient, sessionStore) // 認証チェック用のミドルウェア
	// ログイン必須のパスにはauthミドルウェアを追加する
	r.Path("/").Methods(http.MethodGet).
		HandlerFunc(auth(frontSrv.ViewHome))
	r.Path("/logout").Methods(http.MethodPost).
		HandlerFunc(auth(frontSrv.Logout))
	r.Path("/project").Methods(http.MethodPost).
		HandlerFunc(auth(frontSrv.CreateProject))
	r.Path("/project/{id}").Methods(http.MethodGet).
		HandlerFunc(auth(frontSrv.ViewProject))
	r.Path("/project/{id}").Methods(http.MethodPost).
		HandlerFunc(auth(frontSrv.UpdateProject))
	r.Path("/task").Methods(http.MethodPost).
		HandlerFunc(auth(frontSrv.CreateTask))
	r.Path("/task/{id}").Methods(http.MethodPost).
		HandlerFunc(auth(frontSrv.UpdateTask))
	r.Path("/signup").Methods(http.MethodGet).
		HandlerFunc(frontSrv.ViewSignup)
	r.Path("/signup").Methods(http.MethodPost).
		HandlerFunc(frontSrv.Signup)
	r.Path("/login").Methods(http.MethodGet).
		HandlerFunc(frontSrv.ViewLogin)
	r.Path("/login").Methods(http.MethodPost).
		HandlerFunc(frontSrv.Login)
	static := http.StripPrefix("/static",
		http.FileServer(http.Dir("static")))
	r.PathPrefix("/static/").Handler(static)
	log.Println("start server on port", port)
	err := http.ListenAndServe(port, r)
	if err != nil {
		log.Println("failed to exit serve: ", err)
	}
}

func getGRPCConn(target string) *grpc.ClientConn {
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial: %s", err)
	}

	return conn
}
