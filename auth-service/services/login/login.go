package login

import (
	"egnite.app/microservices/authentication/helpers"
	"egnite.app/microservices/authentication/services/user"
	"golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Server represents the gRPC server
type Server struct {
}

// Login generates response to a Ping request
func (s *Server) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		return &LoginResponse{Success: false, Err: "unable_to_connect_to_user_service"}, nil
	}
	defer conn.Close()

	userService := user.NewUserServiceClient(conn)
	response, err := userService.GetUserByUsernameWithPassword(context.Background(), &user.GetUserByUsernameWithPasswordRequest{Username: req.Username})
	if err != nil {
		return &LoginResponse{Success: false, Err: "internal_server_error"}, nil
	}
	if response.Success {
		token := make(map[string]string)
		if helpers.CheckPasswordHash(req.Password, response.Password) {
			user := make(map[string]string)
			user["username"] = response.User.Username
			user["role"] = response.User.Role
			token = helpers.CreateToken(user)
			return &LoginResponse{AccessToken: token["accessToken"], Id: response.User.Id, Success: true}, nil
		}
		return &LoginResponse{Success: false, Err: "password_incorrect"}, nil

	}
	return &LoginResponse{Success: false, Err: response.Err}, nil

}
