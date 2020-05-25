package token

import (
	"egnite.app/microservices/authentication/helpers"
	"egnite.app/microservices/authentication/services/user"
	"golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Server represents the gRPC server
type Server struct {
}

// RefreshToken generates response to a Ping request
func (s *Server) RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*RefreshTokenResponse, error) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		return &RefreshTokenResponse{Success: false, Err: "unable_to_connect_to_user_service"}, nil
	}
	defer conn.Close()

	userService := user.NewUserServiceClient(conn)
	response, err := userService.GetUserByID(context.Background(), &user.GetUserByIDRequest{Id: req.Id})
	if err != nil {
		return &RefreshTokenResponse{Success: false, Err: "internal_server_error"}, nil
	}
	if response.Success {
		token := make(map[string]string)
		userDataForToken := make(map[string]string)
		userDataForToken["username"] = response.User.Username
		userDataForToken["role"] = response.User.Role
		token = helpers.CreateToken(userDataForToken)
		return &RefreshTokenResponse{Success: true, AccessToken: token["accessToken"]}, nil
	}
	return &RefreshTokenResponse{Success: false, Err: response.Err}, nil
}
