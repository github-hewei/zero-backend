package controller

import (
	"strconv"
	"zero-backend/internal/apperror"
	"zero-backend/internal/dto"
	"zero-backend/internal/request"
	"zero-backend/internal/response"
	"zero-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// UploadFileController 文件控制器
type UploadFileController struct {
	req  *request.Request
	serv *service.UploadFileService
}

// NewUploadFileController 创建文件控制器
func NewUploadFileController(req *request.Request, serv *service.UploadFileService) *UploadFileController {
	return &UploadFileController{req: req, serv: serv}
}

// List 获取文件列表
func (c *UploadFileController) List(ctx *gin.Context) {
	req := &dto.UploadFileListRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	result, err := c.serv.FindList(ctx.Request.Context(), req)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "请求成功", result)
}

// Upload 文件上传
func (c *UploadFileController) Upload(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		response.Error(ctx, apperror.NewUserError("请选择上传文件"))
		return
	}

	groupId, _ := strconv.Atoi(ctx.PostForm("group_id"))
	req := &dto.UploadFileRequest{
		File:    file,
		GroupId: uint32(groupId),
	}

	result, err := c.serv.Upload(ctx.Request.Context(), req)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "上传成功", result)
}

// Delete 删除文件
func (c *UploadFileController) Delete(ctx *gin.Context) {
	req := &dto.UploadFileDeleteRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	if err := c.serv.Delete(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "删除成功", nil)
}
