package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/portfolio/project-service/internal/domain/entity"
	"github.com/portfolio/project-service/internal/usecase"
	pb "github.com/portfolio/proto/project"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ProjectHandler handles gRPC requests for project service
type ProjectHandler struct {
	pb.UnimplementedProjectServiceServer
	projectUC      *usecase.ProjectUseCase
	skillUC        *usecase.SkillUseCase
	projectSkillUC *usecase.ProjectSkillUseCase
	techUC         *usecase.TechUseCase
	imageUC        *usecase.ImageUseCase
	linkUC         *usecase.LinkUseCase
}

// NewProjectHandler creates a new ProjectHandler
func NewProjectHandler(
	projectUC *usecase.ProjectUseCase,
	skillUC *usecase.SkillUseCase,
	projectSkillUC *usecase.ProjectSkillUseCase,
	techUC *usecase.TechUseCase,
	imageUC *usecase.ImageUseCase,
	linkUC *usecase.LinkUseCase,
) *ProjectHandler {
	return &ProjectHandler{
		projectUC:      projectUC,
		skillUC:        skillUC,
		projectSkillUC: projectSkillUC,
		techUC:         techUC,
		imageUC:        imageUC,
		linkUC:         linkUC,
	}
}

// --- Project CRUD ---

func (h *ProjectHandler) CreateProject(ctx context.Context, req *pb.CreateProjectRequest) (*pb.ProjectResponse, error) {
	startDate := req.StartDate.AsTime()
	endDate := req.EndDate.AsTime()

	project, err := h.projectUC.CreateProject(ctx, req.Name, req.Description, req.Status, &startDate, &endDate)
	if err != nil {
		return nil, err
	}

	return &pb.ProjectResponse{Project: mapProjectToProto(project)}, nil
}

func (h *ProjectHandler) GetProject(ctx context.Context, req *pb.GetProjectRequest) (*pb.ProjectResponse, error) {
	fmt.Println("GetProject")
	fmt.Println(req.Id)
	project, err := h.projectUC.GetProject(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.ProjectResponse{Project: mapProjectToProto(project)}, nil
}

func (h *ProjectHandler) UpdateProject(ctx context.Context, req *pb.UpdateProjectRequest) (*pb.ProjectResponse, error) {
	var startDate, endDate *time.Time
	if req.StartDate != nil {
		t := req.StartDate.AsTime()
		startDate = &t
	}
	if req.EndDate != nil {
		t := req.EndDate.AsTime()
		endDate = &t
	}

	project, err := h.projectUC.UpdateProject(ctx, req.Id, req.Name, req.Description, req.Status, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return &pb.ProjectResponse{Project: mapProjectToProto(project)}, nil
}

func (h *ProjectHandler) DeleteProject(ctx context.Context, req *pb.DeleteProjectRequest) (*pb.Empty, error) {
	err := h.projectUC.DeleteProject(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (h *ProjectHandler) ListProjects(ctx context.Context, req *pb.ListProjectsRequest) (*pb.ListProjectsResponse, error) {
	projects, total, err := h.projectUC.ListProjects(ctx, int(req.Page), int(req.Limit), req.Status)
	if err != nil {
		return nil, err
	}

	var protoProjects []*pb.Project
	for _, p := range projects {
		protoProjects = append(protoProjects, mapProjectToProto(p))
	}

	return &pb.ListProjectsResponse{
		Projects: protoProjects,
		Total:    int32(total),
	}, nil
}

// --- Skills ---

func (h *ProjectHandler) CreateSkill(ctx context.Context, req *pb.CreateSkillRequest) (*pb.SkillResponse, error) {
	skill, err := h.skillUC.CreateSkill(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	return &pb.SkillResponse{Skill: &pb.Skill{Id: skill.ID, Name: skill.Name}}, nil
}

func (h *ProjectHandler) ListSkills(ctx context.Context, req *pb.Empty) (*pb.ListSkillsResponse, error) {
	skills, err := h.skillUC.ListSkills(ctx)
	if err != nil {
		return nil, err
	}

	var protoSkills []*pb.Skill
	for _, s := range skills {
		protoSkills = append(protoSkills, &pb.Skill{Id: s.ID, Name: s.Name})
	}

	return &pb.ListSkillsResponse{Skills: protoSkills}, nil
}

func (h *ProjectHandler) AddProjectSkill(ctx context.Context, req *pb.AddProjectSkillRequest) (*pb.Empty, error) {
	err := h.projectSkillUC.AddSkill(ctx, req.ProjectId, req.SkillId)
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (h *ProjectHandler) RemoveProjectSkill(ctx context.Context, req *pb.RemoveProjectSkillRequest) (*pb.Empty, error) {
	err := h.projectSkillUC.RemoveSkill(ctx, req.ProjectId, req.SkillId)
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

// --- Tech Stack ---

func (h *ProjectHandler) AddProjectTech(ctx context.Context, req *pb.AddProjectTechRequest) (*pb.Empty, error) {
	err := h.techUC.AddTech(ctx, req.ProjectId, req.TechName)
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (h *ProjectHandler) RemoveProjectTech(ctx context.Context, req *pb.RemoveProjectTechRequest) (*pb.Empty, error) {
	err := h.techUC.RemoveTech(ctx, req.ProjectId, req.TechName)
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

// --- Images ---

func (h *ProjectHandler) AddProjectImage(ctx context.Context, req *pb.AddProjectImageRequest) (*pb.ProjectImageResponse, error) {
	image, err := h.imageUC.AddImage(ctx, req.ProjectId, req.ImageUrl, req.Description)
	if err != nil {
		return nil, err
	}

	return &pb.ProjectImageResponse{
		Image: &pb.ProjectImage{
			Id:          image.ID,
			ProjectId:   image.ProjectID,
			ImageUrl:    image.ImageURL,
			Description: image.Description,
			UploadedAt:  timestamppb.New(image.UploadedAt),
		},
	}, nil
}

func (h *ProjectHandler) RemoveProjectImage(ctx context.Context, req *pb.RemoveProjectImageRequest) (*pb.Empty, error) {
	err := h.imageUC.RemoveImage(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (h *ProjectHandler) ListProjectImages(ctx context.Context, req *pb.ListProjectImagesRequest) (*pb.ListProjectImagesResponse, error) {
	images, err := h.imageUC.GetImages(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}

	var protoImages []*pb.ProjectImage
	for _, img := range images {
		protoImages = append(protoImages, &pb.ProjectImage{
			Id:          img.ID,
			ProjectId:   img.ProjectID,
			ImageUrl:    img.ImageURL,
			Description: img.Description,
			UploadedAt:  timestamppb.New(img.UploadedAt),
		})
	}

	return &pb.ListProjectImagesResponse{Images: protoImages}, nil
}

// --- Links ---

func (h *ProjectHandler) AddProjectLink(ctx context.Context, req *pb.AddProjectLinkRequest) (*pb.ProjectLinkResponse, error) {
	link, err := h.linkUC.AddLink(ctx, req.ProjectId, req.LinkUrl, req.LinkType)
	if err != nil {
		return nil, err
	}

	return &pb.ProjectLinkResponse{
		Link: &pb.ProjectLink{
			Id:        link.ID,
			ProjectId: link.ProjectID,
			LinkUrl:   link.LinkURL,
			LinkType:  link.LinkType,
		},
	}, nil
}

func (h *ProjectHandler) RemoveProjectLink(ctx context.Context, req *pb.RemoveProjectLinkRequest) (*pb.Empty, error) {
	err := h.linkUC.RemoveLink(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (h *ProjectHandler) ListProjectLinks(ctx context.Context, req *pb.ListProjectLinksRequest) (*pb.ListProjectLinksResponse, error) {
	links, err := h.linkUC.GetLinks(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}

	var protoLinks []*pb.ProjectLink
	for _, l := range links {
		protoLinks = append(protoLinks, &pb.ProjectLink{
			Id:        l.ID,
			ProjectId: l.ProjectID,
			LinkUrl:   l.LinkURL,
			LinkType:  l.LinkType,
		})
	}

	return &pb.ListProjectLinksResponse{Links: protoLinks}, nil
}

// --- Helpers ---

func mapProjectToProto(p *entity.Project) *pb.Project {
	var skills []*pb.Skill
	for _, s := range p.Skills {
		skills = append(skills, &pb.Skill{Id: s.ID, Name: s.Name})
	}

	var techStack []string
	for _, t := range p.TechStack {
		techStack = append(techStack, t)
	}

	var images []*pb.ProjectImage
	for _, i := range p.Images {
		images = append(images, &pb.ProjectImage{
			Id:          i.ID,
			ProjectId:   i.ProjectID,
			ImageUrl:    i.ImageURL,
			Description: i.Description,
			UploadedAt:  timestamppb.New(i.UploadedAt),
		})
	}

	var links []*pb.ProjectLink
	for _, l := range p.Links {
		links = append(links, &pb.ProjectLink{
			Id:        l.ID,
			ProjectId: l.ProjectID,
			LinkUrl:   l.LinkURL,
			LinkType:  l.LinkType,
		})
	}

	var startDate *timestamppb.Timestamp
	if p.StartDate != nil {
		startDate = timestamppb.New(*p.StartDate)
	}

	var endDate *timestamppb.Timestamp
	if p.EndDate != nil {
		endDate = timestamppb.New(*p.EndDate)
	}

	return &pb.Project{
		Id:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		StartDate:   startDate,
		EndDate:     endDate,
		Status:      p.Status,
		Skills:      skills,
		TechStack:   techStack,
		Images:      images,
		Links:       links,
		CreatedAt:   timestamppb.New(p.CreatedAt),
		UpdatedAt:   timestamppb.New(p.UpdatedAt),
	}
}
