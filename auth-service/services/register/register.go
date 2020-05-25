package register

import (
	"egnite.app/microservices/authentication/helpers"
	"egnite.app/microservices/authentication/services/user"
	"golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Server represents the gRPC server
type Server struct {
}

// Register generates response to a Ping request
func (s *Server) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		return &RegisterResponse{Success: false, Err: "unable_to_connect_to_user_service"}, nil
	}
	defer conn.Close()

	userService := user.NewUserServiceClient(conn)
	response, err := userService.GetUserByUsername(context.Background(), &user.GetUserByUsernameRequest{Username: req.Username})
	if err != nil {
		return &RegisterResponse{Success: false, Err: "internal_server_error"}, nil
	}
	if response.Success {
		return &RegisterResponse{Success: false, Err: "user_already_exist"}, nil

	}
	var userData user.User
	userData.Name = req.Name
	userData.Username = req.Username
	userData.Email = req.Email
	userData.Phone = req.Phone
	userData.Role = req.Role

	createResponse, err := userService.CreateUser(context.Background(), &user.CreateUserRequest{User: &userData, Password: req.Password})

	if err != nil {
		return &RegisterResponse{Success: false, Err: "unable_to_create_user"}, nil
	}

	if createResponse.Success {
		token := make(map[string]string)
		userDataForTOken := make(map[string]string)
		userDataForTOken["username"] = req.Username
		userDataForTOken["role"] = req.Role
		token = helpers.CreateToken(userDataForTOken)
		return &RegisterResponse{AccessToken: token["accessToken"], Success: true, Id: createResponse.Id}, nil
	}
	return &RegisterResponse{Success: false, Err: createResponse.Err}, nil

}
