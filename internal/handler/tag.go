package handler

import (
	"context"

	"github.com/VaneZ444/forum-service/internal/handler/converter"
	"github.com/VaneZ444/forum-service/internal/usecase"
	ssov1 "github.com/VaneZ444/golang-forum-protos/gen/go/forum"
)

type TagHandler struct {
	ssov1.UnimplementedForumServer
	usecase usecase.TagUseCase
}

func NewTagHandler(tagUC usecase.TagUseCase) *TagHandler {
	return &TagHandler{
		usecase: tagUC,
	}
}

func (h *TagHandler) CreateTag(ctx context.Context, req *ssov1.CreateTagRequest) (*ssov1.CreateTagResponse, error) {
	id, err := h.usecase.CreateTag(ctx, req.GetName())
	if err != nil {
		return nil, err
	}
	return &ssov1.CreateTagResponse{TagId: id}, nil
}

func (h *TagHandler) GetTagByID(ctx context.Context, req *ssov1.GetTagByIDRequest) (*ssov1.GetTagByIDResponse, error) {
	tag, err := h.usecase.GetTagByID(ctx, req.GetTagId())
	if err != nil {
		return nil, err
	}
	return &ssov1.GetTagByIDResponse{
		Tag: converter.TagToProto(tag),
	}, nil
}

func (h *TagHandler) ListTags(ctx context.Context, _ *ssov1.ListTagsRequest) (*ssov1.ListTagsResponse, error) {
	tags, err := h.usecase.ListTags(ctx)
	if err != nil {
		return nil, err
	}

	var pbTags []*ssov1.Tag
	for _, t := range tags {
		pbTags = append(pbTags, converter.TagToProto(t))
	}

	return &ssov1.ListTagsResponse{Tags: pbTags}, nil
}

func (h *TagHandler) ListTagsByPostID(ctx context.Context, req *ssov1.ListTagsByPostIDRequest) (*ssov1.ListTagsByPostIDResponse, error) {
	tags, err := h.usecase.ListTagsByPostID(ctx, req.GetPostId())
	if err != nil {
		return nil, err
	}

	var pbTags []*ssov1.Tag
	for _, t := range tags {
		pbTags = append(pbTags, converter.TagToProto(t))
	}

	return &ssov1.ListTagsByPostIDResponse{Tags: pbTags}, nil
}
