package converter

import (
	"github.com/VaneZ444/forum-service/internal/entity"
	ssov1 "github.com/VaneZ444/golang-forum-protos/gen/go/forum"
)

func TagToProto(t *entity.Tag) *ssov1.Tag {
	return &ssov1.Tag{
		Id:   t.ID,
		Name: t.Name,
	}
}
