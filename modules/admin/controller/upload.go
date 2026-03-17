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

// UploadGroupController 文件分组控制器
type UploadGroupController struct {
	req  *request.Request
	serv *service.UploadGroupService
}

// NewUploadGroupController 创建文件分组控制器
func NewUploadGroupController(req *request.Request, serv *service.UploadGroupService) *UploadGroupController {
	return &UploadGroupController{req: req, serv: serv}
}

// List 获取分组列表(树形结构)
func (c *UploadGroupController) List(ctx *gin.Context) {
	result, err := c.serv.FindTreeList(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "请求成功", result)
}

// Create 创建分组
func (c *UploadGroupController) Create(ctx *gin.Context) {
	req := &dto.UploadGroupCreateRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	req.StoreId = request.GetStoreId(ctx)
	if err := c.serv.Create(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "创建成功", nil)
}

// Update 更新分组
func (c *UploadGroupController) Update(ctx *gin.Context) {
	req := &dto.UploadGroupUpdateRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	if err := c.serv.Update(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "更新成功", nil)
}

// Delete 删除分组
func (c *UploadGroupController) Delete(ctx *gin.Context) {
	req := &dto.UploadGroupDeleteRequest{}
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
