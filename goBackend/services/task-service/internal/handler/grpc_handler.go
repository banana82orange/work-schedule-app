package handler

import (
	"context"
	"time"

	pb "github.com/portfolio/proto/task"
	"github.com/portfolio/task-service/internal/domain/entity"
	"github.com/portfolio/task-service/internal/usecase"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TaskHandler handles gRPC requests for task service
type TaskHandler struct {
	pb.UnimplementedTaskServiceServer
	taskUC       *usecase.TaskUseCase
	subtaskUC    *usecase.SubtaskUseCase
	commentUC    *usecase.CommentUseCase
	attachmentUC *usecase.AttachmentUseCase
	tagUC        *usecase.TagUseCase
}

// NewTaskHandler creates a new TaskHandler
func NewTaskHandler(
	taskUC *usecase.TaskUseCase,
	subtaskUC *usecase.SubtaskUseCase,
	commentUC *usecase.CommentUseCase,
	attachmentUC *usecase.AttachmentUseCase,
	tagUC *usecase.TagUseCase,
) *TaskHandler {
	return &TaskHandler{
		taskUC:       taskUC,
		subtaskUC:    subtaskUC,
		commentUC:    commentUC,
		attachmentUC: attachmentUC,
		tagUC:        tagUC,
	}
}

// --- Task CRUD ---

func (h *TaskHandler) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.TaskResponse, error) {
	var dueDate *time.Time
	if req.DueDate != nil {
		t := req.DueDate.AsTime()
		dueDate = &t
	}

	task, err := h.taskUC.CreateTask(ctx, req.ProjectId, req.Title, req.Description, req.Status, int(req.Priority), req.AssignedTo, dueDate)
	if err != nil {
		return nil, err
	}

	return &pb.TaskResponse{Task: mapTaskToProto(task)}, nil
}

func (h *TaskHandler) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.TaskResponse, error) {
	task, err := h.taskUC.GetTask(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.TaskResponse{Task: mapTaskToProto(task)}, nil
}

func (h *TaskHandler) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.TaskResponse, error) {
	var dueDate *time.Time
	if req.DueDate != nil {
		t := req.DueDate.AsTime()
		dueDate = &t
	}

	task, err := h.taskUC.UpdateTask(ctx, req.Id, req.Title, req.Description, req.Status, int(req.Priority), req.AssignedTo, dueDate)
	if err != nil {
		return nil, err
	}

	return &pb.TaskResponse{Task: mapTaskToProto(task)}, nil
}

func (h *TaskHandler) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*pb.Empty, error) {
	err := h.taskUC.DeleteTask(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (h *TaskHandler) ListTasks(ctx context.Context, req *pb.ListTasksRequest) (*pb.ListTasksResponse, error) {
	tasks, total, err := h.taskUC.ListTasks(ctx, req.ProjectId, int(req.Page), int(req.Limit), req.Status, req.AssignedTo)
	if err != nil {
		return nil, err
	}

	var protoTasks []*pb.Task
	for _, t := range tasks {
		protoTasks = append(protoTasks, mapTaskToProto(t))
	}

	return &pb.ListTasksResponse{
		Tasks: protoTasks,
		Total: int32(total),
	}, nil
}

// --- Subtasks ---

func (h *TaskHandler) CreateSubtask(ctx context.Context, req *pb.CreateSubtaskRequest) (*pb.SubtaskResponse, error) {
	var dueDate *time.Time
	if req.DueDate != nil {
		t := req.DueDate.AsTime()
		dueDate = &t
	}

	subtask, err := h.subtaskUC.CreateSubtask(ctx, req.TaskId, req.Title, req.AssignedTo, dueDate)
	if err != nil {
		return nil, err
	}
	return &pb.SubtaskResponse{Subtask: mapSubtaskToProto(subtask)}, nil
}

func (h *TaskHandler) UpdateSubtask(ctx context.Context, req *pb.UpdateSubtaskRequest) (*pb.SubtaskResponse, error) {
	var dueDate *time.Time
	if req.DueDate != nil {
		t := req.DueDate.AsTime()
		dueDate = &t
	}

	subtask, err := h.subtaskUC.UpdateSubtask(ctx, req.Id, req.Title, req.Status, req.AssignedTo, dueDate)
	if err != nil {
		return nil, err
	}
	return &pb.SubtaskResponse{Subtask: mapSubtaskToProto(subtask)}, nil
}

func (h *TaskHandler) DeleteSubtask(ctx context.Context, req *pb.DeleteSubtaskRequest) (*pb.Empty, error) {
	err := h.subtaskUC.DeleteSubtask(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (h *TaskHandler) ListSubtasks(ctx context.Context, req *pb.ListSubtasksRequest) (*pb.ListSubtasksResponse, error) {
	subtasks, err := h.subtaskUC.GetSubtasks(ctx, req.TaskId)
	if err != nil {
		return nil, err
	}

	var protoSubtasks []*pb.Subtask
	for _, s := range subtasks {
		protoSubtasks = append(protoSubtasks, mapSubtaskToProto(s))
	}

	return &pb.ListSubtasksResponse{Subtasks: protoSubtasks}, nil
}

// --- Comments ---

func (h *TaskHandler) AddComment(ctx context.Context, req *pb.AddCommentRequest) (*pb.CommentResponse, error) {
	comment, err := h.commentUC.AddComment(ctx, req.TaskId, req.UserId, req.Comment)
	if err != nil {
		return nil, err
	}
	return &pb.CommentResponse{
		Comment: &pb.Comment{
			Id:        comment.ID,
			TaskId:    comment.TaskID,
			UserId:    comment.UserID,
			Comment:   comment.Comment,
			CreatedAt: timestamppb.New(comment.CreatedAt),
		},
	}, nil
}

func (h *TaskHandler) DeleteComment(ctx context.Context, req *pb.DeleteCommentRequest) (*pb.Empty, error) {
	err := h.commentUC.DeleteComment(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (h *TaskHandler) ListComments(ctx context.Context, req *pb.ListCommentsRequest) (*pb.ListCommentsResponse, error) {
	comments, err := h.commentUC.GetComments(ctx, req.TaskId)
	if err != nil {
		return nil, err
	}

	var protoComments []*pb.Comment
	for _, c := range comments {
		protoComments = append(protoComments, &pb.Comment{
			Id:        c.ID,
			TaskId:    c.TaskID,
			UserId:    c.UserID,
			Comment:   c.Comment,
			CreatedAt: timestamppb.New(c.CreatedAt),
		})
	}

	return &pb.ListCommentsResponse{Comments: protoComments}, nil
}

// --- Attachments ---

func (h *TaskHandler) AddAttachment(ctx context.Context, req *pb.AddAttachmentRequest) (*pb.AttachmentResponse, error) {
	attachment, err := h.attachmentUC.AddAttachment(ctx, req.TaskId, req.FileUrl)
	if err != nil {
		return nil, err
	}
	return &pb.AttachmentResponse{
		Attachment: &pb.Attachment{
			Id:         attachment.ID,
			TaskId:     attachment.TaskID,
			FileUrl:    attachment.FileURL,
			UploadedAt: timestamppb.New(attachment.UploadedAt),
		},
	}, nil
}

func (h *TaskHandler) DeleteAttachment(ctx context.Context, req *pb.DeleteAttachmentRequest) (*pb.Empty, error) {
	err := h.attachmentUC.DeleteAttachment(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (h *TaskHandler) ListAttachments(ctx context.Context, req *pb.ListAttachmentsRequest) (*pb.ListAttachmentsResponse, error) {
	attachments, err := h.attachmentUC.GetAttachments(ctx, req.TaskId)
	if err != nil {
		return nil, err
	}

	var protoAttachments []*pb.Attachment
	for _, a := range attachments {
		protoAttachments = append(protoAttachments, &pb.Attachment{
			Id:         a.ID,
			TaskId:     a.TaskID,
			FileUrl:    a.FileURL,
			UploadedAt: timestamppb.New(a.UploadedAt),
		})
	}

	return &pb.ListAttachmentsResponse{Attachments: protoAttachments}, nil
}

// --- Tags ---

func (h *TaskHandler) CreateTag(ctx context.Context, req *pb.CreateTagRequest) (*pb.TagResponse, error) {
	tag, err := h.tagUC.CreateTag(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	return &pb.TagResponse{Tag: &pb.Tag{Id: tag.ID, Name: tag.Name}}, nil
}

func (h *TaskHandler) ListTags(ctx context.Context, req *pb.Empty) (*pb.ListTagsResponse, error) {
	tags, err := h.tagUC.ListTags(ctx)
	if err != nil {
		return nil, err
	}

	var protoTags []*pb.Tag
	for _, t := range tags {
		protoTags = append(protoTags, &pb.Tag{Id: t.ID, Name: t.Name})
	}

	return &pb.ListTagsResponse{Tags: protoTags}, nil
}

func (h *TaskHandler) AddTaskTag(ctx context.Context, req *pb.AddTaskTagRequest) (*pb.Empty, error) {
	err := h.tagUC.AddTaskTag(ctx, req.TaskId, req.TagId)
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (h *TaskHandler) RemoveTaskTag(ctx context.Context, req *pb.RemoveTaskTagRequest) (*pb.Empty, error) {
	err := h.tagUC.RemoveTaskTag(ctx, req.TaskId, req.TagId)
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

// --- Helpers ---

func mapTaskToProto(t *entity.Task) *pb.Task {
	var subtasks []*pb.Subtask
	for _, s := range t.Subtasks {
		subtasks = append(subtasks, mapSubtaskToProto(s))
	}

	var tags []*pb.Tag
	for _, tag := range t.Tags {
		tags = append(tags, &pb.Tag{Id: tag.ID, Name: tag.Name})
	}

	var dueDate *timestamppb.Timestamp
	if t.DueDate != nil {
		dueDate = timestamppb.New(*t.DueDate)
	}
	var assignedTo int64
	if t.AssignedTo != nil {
		assignedTo = *t.AssignedTo
	}


	return &pb.Task{
		Id:          t.ID,
		ProjectId:   t.ProjectID,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		Priority:    int32(t.Priority),
		AssignedTo:  assignedTo,
		DueDate:     dueDate,
		Subtasks:    subtasks,
		Tags:        tags,
		CreatedAt:   timestamppb.New(t.CreatedAt),
		UpdatedAt:   timestamppb.New(t.UpdatedAt),
	}
}

func mapSubtaskToProto(s *entity.Subtask) *pb.Subtask {
	var dueDate *timestamppb.Timestamp
	if s.DueDate != nil {
		dueDate = timestamppb.New(*s.DueDate)
	}

	return &pb.Subtask{
		Id:         s.ID,
		TaskId:     s.TaskID,
		Title:      s.Title,
		Status:     s.Status,
		AssignedTo: s.AssignedTo,
		DueDate:    dueDate,
		CreatedAt:  timestamppb.New(s.CreatedAt),
		UpdatedAt:  timestamppb.New(s.UpdatedAt),
	}
}
