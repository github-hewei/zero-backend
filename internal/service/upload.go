package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"zero-backend/internal/apperror"
	"zero-backend/internal/ctxkeys"
	"zero-backend/internal/dto"
	"zero-backend/internal/model"
	"zero-backend/internal/repository"
	"zero-backend/internal/service/uploader"
)

// UploadGroupService 文件分组服务
type UploadGroupService struct {
	repo *repository.UploadGroupRepository
}

// NewUploadGroupService 创建文件分组服务
func NewUploadGroupService(repo *repository.UploadGroupRepository) *UploadGroupService {
	return &UploadGroupService{repo: repo}
}

// FindTreeList 获取分组树形列表
func (s *UploadGroupService) FindTreeList(ctx context.Context, storeId uint32) ([]*model.UploadGroup, error) {
	filter := &repository.UploadGroupFilterField{
		StoreId: storeId,
	}
	list, err := s.repo.FindAll(ctx, filter, nil, nil)

	if err != nil {
		return nil, apperror.NewSystemError(err, "查询分组列表失败")
	}

	if len(list) > 0 {
		uploadGroupList := model.UploadGroupList{}
		for _, group := range list {
			uploadGroupList = append(uploadGroupList, group)
		}
		list = uploadGroupList.Tree()
	}

	return list, nil
}

// Create 创建分组
func (s *UploadGroupService) Create(ctx context.Context, req *dto.UploadGroupCreateRequest) error {
	if err := s.checkName(ctx, req.Name, req.StoreId); err != nil {
		return err
	}

	group := &model.UploadGroup{
		StoreId:  req.StoreId,
		Name:     req.Name,
		ParentId: req.ParentId,
		Sort:     req.Sort,
	}

	if err := s.repo.Create(ctx, group); err != nil {
		return apperror.NewSystemError(err, "创建分组失败")
	}

	return nil
}

// checkName 检查分组名称是否已存在
func (s *UploadGroupService) checkName(ctx context.Context, name string, storeId uint32) error {
	filter := &repository.UploadGroupFilterField{Name: name, StoreId: storeId}
	group, err := s.repo.FindOne(ctx, filter)

	if err != nil {
		return apperror.NewSystemError(err, "检查分组名称失败")
	}

	if group.ID > 0 {
		return apperror.NewUserError("分组名称已存在")
	}

	return nil
}

// Update 更新分组
func (s *UploadGroupService) Update(ctx context.Context, req *dto.UploadGroupUpdateRequest) error {
	// 获取现有分组
	filter := &repository.UploadGroupFilterField{
		Id:      req.ID,
		StoreId: req.StoreId,
	}
	item, err := s.repo.FindOne(ctx, filter)
	if err != nil {
		return apperror.NewSystemError(err, "查询分组失败")
	}
	if item.ID == 0 {
		return apperror.NewUserError("分组不存在或无权限访问")
	}

	// 检查名称是否重复(排除自身)
	if item.Name != req.Name {
		if err := s.checkName(ctx, req.Name, req.StoreId); err != nil {
			return err
		}
	}

	// 更新分组信息
	updateData := map[string]any{
		"name":      req.Name,
		"parent_id": req.ParentId,
		"sort":      req.Sort,
		"store_id":  req.StoreId,
	}

	if err := s.repo.Updates(ctx, item, updateData); err != nil {
		return apperror.NewSystemError(err, "更新分组失败")
	}

	return nil
}

// Delete 删除分组
func (s *UploadGroupService) Delete(ctx context.Context, req *dto.UploadGroupDeleteRequest) error {
	// 检查分组是否存在
	filter := &repository.UploadGroupFilterField{
		Id:      req.ID,
		StoreId: req.StoreId,
	}
	item, err := s.repo.FindOne(ctx, filter)
	if err != nil {
		return apperror.NewSystemError(err, "查询分组失败")
	}
	if item.ID == 0 {
		return apperror.NewUserError("分组不存在或无权限访问")
	}

	// 执行删除
	if err := s.repo.Delete(ctx, req.ID); err != nil {
		return apperror.NewSystemError(err, "删除分组失败")
	}

	return nil
}

// UploadFileService 文件服务
type UploadFileService struct {
	repo        *repository.UploadFileRepository
	settingServ *SettingService
}

// NewUploadFileService 创建文件服务
func NewUploadFileService(
	repo *repository.UploadFileRepository,
	settingServ *SettingService,
) *UploadFileService {
	return &UploadFileService{
		repo:        repo,
		settingServ: settingServ,
	}
}

// FindList 获取文件列表
func (s *UploadFileService) FindList(ctx context.Context, req *dto.UploadFileListRequest) (*dto.ListResult, error) {
	result := &dto.ListResult{
		List:  []*model.UploadFile{},
		Total: 0,
	}

	filter := &repository.UploadFileFilterField{
		StoreId:  req.StoreId,
		GroupId:  req.GroupId,
		FileType: req.FileType,
		FileName: req.FileName,
	}

	pagination := &repository.Pagination{
		Page:  req.Page,
		Limit: req.Limit,
	}

	orders := repository.Orders{
		{Field: "id", Sort: "desc"},
	}

	total, err := s.repo.Count(ctx, filter)
	if err != nil {
		return nil, apperror.NewSystemError(err, "查询文件数量失败")
	}

	if total == 0 {
		return result, nil
	}

	result.Total = total

	list, err := s.repo.FindAll(ctx, filter, pagination, orders)
	if err != nil {
		return nil, apperror.NewSystemError(err, "查询文件列表失败")
	}

	result.List = list
	return result, nil
}

// getUploadConfig 获取上传配置
func (s *UploadFileService) getUploadConfig(ctx context.Context) (*dto.UploadConfig, error) {
	config := &dto.UploadConfig{}
	err := s.settingServ.GetSettingValue(ctx, "upload", config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// sanitizeFilename 安全处理文件名
func (s *UploadFileService) sanitizeFilename(filename string) string {
	// 移除路径信息
	filename = filepath.Base(filename)
	// 替换特殊字符
	filename = strings.ReplaceAll(filename, " ", "_")
	filename = strings.ReplaceAll(filename, "..", "")
	// 只保留字母数字和.-_
	reg := regexp.MustCompile(`[^\w\-\.]`)
	filename = reg.ReplaceAllString(filename, "")
	return filename
}

// validateFileContent 验证文件内容类型是否匹配
func (s *UploadFileService) validateFileContent(fileHeader []byte, fileExt string, allowedTypes []string) error {
	detectedType := http.DetectContentType(fileHeader)

	// 通用安全验证
	if detectedType == "application/octet-stream" {
		return apperror.NewUserError("无法识别的文件类型")
	}

	// 检查图片类型
	if slices.Contains(allowedTypes, "image") {
		if strings.HasPrefix(detectedType, "image/") {
			return nil
		}
	}

	// 检查文档类型
	if slices.Contains(allowedTypes, "document") {
		switch {
		case strings.HasPrefix(detectedType, "application/pdf"):
			return nil
		case strings.HasPrefix(detectedType, "application/msword"):
			return nil
		case strings.HasPrefix(detectedType, "application/vnd.ms-excel"):
			return nil
		}
	}

	// 检查视频类型
	if slices.Contains(allowedTypes, "video") {
		if strings.HasPrefix(detectedType, "video/") {
			return nil
		}
	}

	// 检查压缩包类型
	if slices.Contains(allowedTypes, "archive") {
		switch fileExt {
		case "zip", "rar":
			return nil
		}
	}

	return apperror.NewUserError("文件内容类型与扩展名不匹配")
}

// generateFilePath 生成文件存储路径
func (s *UploadFileService) generateFilePath(file *multipart.FileHeader) (string, error) {
	// 计算文件MD5
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, src); err != nil {
		return "", err
	}
	md5Str := hex.EncodeToString(hash.Sum(nil))

	// 获取文件扩展名
	ext := filepath.Ext(file.Filename)
	if ext == "" {
		ext = ".dat"
	}

	// 生成存储路径(固定到uploads目录)
	path := filepath.Join("uploads", md5Str[:2], md5Str[2:4], md5Str[4:]+ext)

	// 创建目录
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return "", err
	}

	return path, nil
}

// Upload 文件上传
func (s *UploadFileService) Upload(ctx context.Context, req *dto.UploadFileRequest) (*model.UploadFile, error) {
	// 获取上传配置
	config, err := s.getUploadConfig(ctx)
	if err != nil {
		return nil, apperror.NewSystemError(err, "获取上传配置失败")
	}

	// 验证文件大小
	maxSize, _ := strconv.Atoi(config.MaxSize)
	if req.File.Size > int64(maxSize)*1024*1024 {
		return nil, apperror.NewUserError(fmt.Sprintf("文件大小不能超过%dMB", maxSize))
	}

	// 验证文件类型和内容安全
	fileExt := strings.ToLower(strings.TrimPrefix(filepath.Ext(req.File.Filename), "."))

	if !s.checkFileExt(config.AllowedTypes, fileExt) {
		return nil, apperror.NewUserError("不支持的文件类型")
	}

	// 验证文件内容类型
	fileHeader := make([]byte, 512)
	src, err := req.File.Open()
	if err != nil {
		return nil, apperror.NewSystemError(err, "文件打开失败")
	}
	if _, err = src.Read(fileHeader); err != nil {
		src.Close()
		return nil, apperror.NewSystemError(err, "文件读取失败")
	}
	src.Close()

	// 验证文件内容类型是否匹配
	if err := s.validateFileContent(fileHeader, fileExt, config.AllowedTypes); err != nil {
		return nil, err
	}

	// 安全处理文件名
	safeFilename := s.sanitizeFilename(req.File.Filename)
	if safeFilename == "" {
		return nil, apperror.NewUserError("无效的文件名")
	}

	// 生成存储路径
	savePath, err := s.generateFilePath(req.File)
	if err != nil {
		return nil, apperror.NewSystemError(err, "生成文件路径失败")
	}

	// 创建文件记录
	fileType := int8(20) // 默认附件
	if strings.HasPrefix(http.DetectContentType(fileHeader), "image/") {
		fileType = 10 // 图片
	} else if strings.HasPrefix(http.DetectContentType(fileHeader), "video/") {
		fileType = 30 // 视频
	}

	filePath := "/" + strings.ReplaceAll(savePath, "\\", "/")
	uploadFile := &model.UploadFile{
		GroupId:    req.GroupId,
		Channel:    10, // 后台上传
		Storage:    config.StorageType,
		Domain:     "", // 本地存储留空
		FileType:   fileType,
		FilePath:   filePath,
		FileName:   req.File.Filename,
		FileSize:   uint32(req.File.Size),
		FileExt:    fileExt,
		UploaderId: req.UploaderId,
		StoreId:    req.StoreId,
	}

	// 如果通过七牛云上传，在上下文中添加七牛配置
	if config.StorageType == "qiniu" {
		qiniu := &dto.QiniuConfig{}
		if err := s.settingServ.GetSettingValue(ctx, "qiniu", qiniu); err != nil {
			return nil, err
		}
		if !qiniu.IsEnabled {
			return nil, apperror.NewUserError("七牛云存储未启用")
		}

		ctx = context.WithValue(ctx, ctxkeys.QiniuConfigKey{}, qiniu)
	}

	// 创建上传方法实现
	uploader, err := uploader.NewUploader(config.StorageType, ctx)
	if err != nil {
		return nil, err
	}

	domain, err := uploader.Upload(ctx, req.File, savePath)
	if err != nil {
		return nil, err
	}

	uploadFile.Domain = domain

	if err := s.repo.Create(ctx, uploadFile); err != nil {
		return nil, apperror.NewSystemError(err, "文件记录创建失败")
	}

	return uploadFile, nil
}

// checkFileExt 检测文件扩展名
func (s *UploadFileService) checkFileExt(allowedTypes []string, fileExt string) bool {
	for _, allowedType := range allowedTypes {
		switch allowedType {
		case "image":
			if slices.Contains([]string{"jpg", "png", "gif"}, fileExt) {
				return true
			}
		case "document":
			if slices.Contains([]string{"pdf", "doc", "xls", "xlsx"}, fileExt) {
				return true
			}
		case "video":
			if slices.Contains([]string{"mp4", "mov"}, fileExt) {
				return true
			}
		case "archive":
			if slices.Contains([]string{"zip", "rar"}, fileExt) {
				return true
			}
		}
	}

	return false
}

// Delete 删除文件
func (s *UploadFileService) Delete(ctx context.Context, req *dto.UploadFileDeleteRequest) error {
	// 检查文件是否存在
	filter := &repository.UploadFileFilterField{
		Id:      req.ID,
		StoreId: req.StoreId,
	}
	item, err := s.repo.FindOne(ctx, filter)
	if err != nil {
		return apperror.NewSystemError(err, "查询文件失败")
	}
	if item.ID == 0 {
		return apperror.NewUserError("文件不存在")
	}

	// 执行删除
	if err := s.repo.Delete(ctx, req.ID); err != nil {
		return apperror.NewSystemError(err, "删除文件失败")
	}

	return nil
}
