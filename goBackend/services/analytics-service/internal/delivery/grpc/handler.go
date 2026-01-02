package grpc

import (
	"context"
	"fmt"

	"github.com/portfolio/analytics-service/internal/usecase"
	pb "github.com/portfolio/proto/analytics"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AnalyticsServer implements the AnalyticsService gRPC server
type AnalyticsServer struct {
	pb.UnimplementedAnalyticsServiceServer
	analyticsUseCase *usecase.AnalyticsUseCase
}

// NewAnalyticsServer creates a new AnalyticsServer
func NewAnalyticsServer(
	analyticsUseCase *usecase.AnalyticsUseCase,
) *AnalyticsServer {
	return &AnalyticsServer{
		analyticsUseCase: analyticsUseCase,
	}
}


func (s *AnalyticsServer) RecordTaskActivity(ctx context.Context, req *pb.RecordTaskActivityRequest) (*pb.Empty, error) {


	err := s.analyticsUseCase.RecordTaskActivity(ctx, req.TaskId, req.UserId, req.Action)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.Empty{}, nil
}


// GetProjectStats returns project stats
func (s *AnalyticsServer) GetProjectStats(ctx context.Context, req *pb.GetProjectStatsRequest) (*pb.ProjectStatsResponse, error) {
	fmt.Println("GetProjectStats")
	fmt.Println( req.ProjectId)
	stats, err := s.analyticsUseCase.GetProjectStats(ctx, req.ProjectId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.ProjectStatsResponse{
		Stats: &pb.ProjectStats{
			ProjectId: stats.ProjectID,
			TotalTasks: int32(stats.TotalTasks),
			CompletedTasks: int32(stats.CompletedTasks),
			ProgressPercent: stats.ProgressPercent,
		},
	}, nil
}

func (s *AnalyticsServer) UpdateProjectStats(ctx context.Context, req *pb.UpdateProjectStatsRequest) (*pb.ProjectStatsResponse, error) {
	fmt.Println("UpdateProjectStats")
	fmt.Println( req.ProjectId)
	fmt.Println( req.TotalTasks)
	fmt.Println( req.CompletedTasks)
	_ , err := s.analyticsUseCase.UpdateProjectStats(ctx, req.ProjectId, int(req.TotalTasks), int(req.CompletedTasks))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.ProjectStatsResponse{}, nil
}
