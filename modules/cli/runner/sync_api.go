package runner

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"zero-backend/internal/apperror"
	"zero-backend/internal/logger"
	"zero-backend/internal/model"
	"zero-backend/internal/repository"
)

// OpenAPISpec OpenAPI 3.x 规范结构
type OpenAPISpec struct {
	OpenAPI string `json:"openapi"`
	Info    struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Version     string `json:"version"`
	} `json:"info"`
	Paths map[string]PathItem `json:"paths"`
}

// PathItem 路径定义
type PathItem map[string]Operation

// Operation 操作定义
type Operation struct {
	Tags        []string `json:"tags"`
	Summary     string   `json:"summary"`
	Description string   `json:"description"`
	OperationID string   `json:"operationId"`
}

// SyncApiRunner API同步执行器
type SyncApiRunner struct {
	logger *logger.Logger
	repo   *repository.RbacApiRepository
}

// NewSyncApiRunner 创建API同步执行器
func NewSyncApiRunner(l *logger.Logger, repo *repository.RbacApiRepository) *SyncApiRunner {
	return &SyncApiRunner{
		logger: l,
		repo:   repo,
	}
}

// SyncResult 同步结果
type SyncResult struct {
	TotalInDoc int      // 文档中的API数量
	Added      int      // 新增数量
	Updated    int      // 更新数量
	Deleted    int      // 删除数量
	Skipped    int      // 跳过数量
	Errors     []string // 错误列表
}

// Run 执行同步
func (r *SyncApiRunner) Run(ctx context.Context, filePath string, deleteExtra bool) (*SyncResult, error) {
	openapiFile := filePath
	if openapiFile == "" {
		openapiFile = filepath.Join("docs", "admin", "index.json")
	}

	// 读取并解析 OpenAPI 文档
	spec, err := r.parseOpenAPI(openapiFile)
	if err != nil {
		return nil, err
	}

	result := &SyncResult{
		Errors: []string{},
	}

	// 用于记录文档中存在的URL
	docURLs := make(map[string]bool)

	tagToParentID := make(map[string]uint32)
	uniqueTags := r.collectUniqueTags(spec)

	r.logger.Info("发现分组数量", "count", len(uniqueTags))

	for _, tag := range uniqueTags {
		parentID, err := r.ensureParentCategory(ctx, tag)
		if err != nil {
			result.Errors = append(result.Errors, "创建父级目录失败 "+tag+": "+err.Error())
			continue
		}
		tagToParentID[tag] = parentID
	}

	// 遍历所有路径和方法
	for path, methods := range spec.Paths {
		for method, op := range methods {
			// 当前项目中定义的接口统一使用POST请求，所以过滤掉非POST请求
			if strings.ToUpper(method) != "POST" {
				continue
			}

			// 构建URL
			url := path

			// 获取接口名称
			name := op.Summary
			if name == "" {
				name = op.OperationID
			}
			if name == "" {
				name = path
			}

			// 获取父级目录ID
			var parentID uint32
			if len(op.Tags) > 0 && op.Tags[0] != "" {
				parentID = tagToParentID[op.Tags[0]]
			}

			docURLs[url] = true

			// 检查是否已存在
			existing, err := r.repo.GetAPIByPath(ctx, url)
			if err != nil {
				result.Errors = append(result.Errors, "查询API失败 "+url+": "+err.Error())
				continue
			}

			if existing != nil && existing.ID > 0 {
				// 已存在 - 检查是否需要更新
				needsUpdate := existing.Name != name || existing.ParentId != parentID
				if needsUpdate {
					if err := r.updateAPI(ctx, existing.ID, name, parentID); err != nil {
						result.Errors = append(result.Errors, "更新API失败 "+url+": "+err.Error())
						continue
					}
					result.Updated++
					r.logger.Info("更新API", "url", url, "name", name, "parent_id", parentID)
				} else {
					result.Skipped++
				}
			} else {
				// 不存在 - 新增
				if err := r.createAPI(ctx, url, name, parentID); err != nil {
					result.Errors = append(result.Errors, "创建API失败 "+url+": "+err.Error())
					continue
				}
				result.Added++
				r.logger.Info("新增API", "url", url, "name", name, "parent_id", parentID)
			}
		}
	}

	// 统计文档中的API数量（只统计POST请求）
	for _, methods := range spec.Paths {
		for method := range methods {
			if strings.ToUpper(method) == "POST" {
				result.TotalInDoc++
			}
		}
	}

	// 处理删除
	if deleteExtra {
		deleted, err := r.deleteExtraAPIs(ctx, docURLs)
		if err != nil {
			result.Errors = append(result.Errors, "删除多余API失败: "+err.Error())
		} else {
			result.Deleted = deleted
		}
	}

	r.logger.Info("同步完成",
		"文档API数", result.TotalInDoc,
		"新增", result.Added,
		"更新", result.Updated,
		"删除", result.Deleted,
		"跳过", result.Skipped,
		"错误", len(result.Errors),
	)

	return result, nil
}

// collectUniqueTags 收集所有唯一的 tag
func (r *SyncApiRunner) collectUniqueTags(spec *OpenAPISpec) []string {
	tagSet := make(map[string]bool)
	for _, methods := range spec.Paths {
		for _, op := range methods {
			if len(op.Tags) > 0 && op.Tags[0] != "" {
				tagSet[op.Tags[0]] = true
			}
		}
	}

	tags := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		tags = append(tags, tag)
	}
	return tags
}

// ensureParentCategory 确保父级目录存在，返回父级目录的ID
func (r *SyncApiRunner) ensureParentCategory(ctx context.Context, tag string) (uint32, error) {
	existing, err := r.repo.GetAPIByName(ctx, tag)
	if err != nil {
		return 0, err
	}

	if existing != nil && existing.ID > 0 {
		return existing.ID, nil
	}

	api := &model.RbacApi{
		Name:     tag,
		Url:      "-",
		ParentId: 0,
		Sort:     100,
	}

	if err := r.repo.Create(ctx, api); err != nil {
		return 0, err
	}

	r.logger.Info("创建父级目录", "name", tag)
	return api.ID, nil
}

// parseOpenAPI 解析OpenAPI文档
func (r *SyncApiRunner) parseOpenAPI(filePath string) (*OpenAPISpec, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, apperror.NewUserError("找不到 OpenAPI 文档文件：" + filePath)
		}
		return nil, apperror.NewSystemError(err, "读取 OpenAPI 文档失败")
	}

	var spec OpenAPISpec
	if err := json.Unmarshal(data, &spec); err != nil {
		return nil, apperror.NewSystemError(err, "解析 OpenAPI 文档失败")
	}

	if spec.OpenAPI == "" {
		return nil, apperror.NewUserError("无效的 OpenAPI 文档：缺少 openapi 字段")
	}

	if len(spec.Paths) == 0 {
		return nil, apperror.NewUserError("OpenAPI 文档中未找到任何路径")
	}

	r.logger.Info("解析 OpenAPI 文档成功", "版本", spec.OpenAPI, "路径数", len(spec.Paths))

	return &spec, nil
}

// createAPI 创建新API
func (r *SyncApiRunner) createAPI(ctx context.Context, url, name string, parentID uint32) error {
	api := &model.RbacApi{
		Name:     name,
		Url:      url,
		ParentId: parentID,
		Sort:     100,
	}
	return r.repo.Create(ctx, api)
}

// updateAPI 更新API
func (r *SyncApiRunner) updateAPI(ctx context.Context, id uint32, name string, parentID uint32) error {
	api := &model.RbacApi{ID: id}
	return r.repo.Updates(ctx, api, map[string]any{
		"name":      name,
		"parent_id": parentID,
	})
}

// deleteExtraAPIs 删除文档中不存在的API
func (r *SyncApiRunner) deleteExtraAPIs(ctx context.Context, docURLs map[string]bool) (int, error) {
	allAPIs, err := r.repo.FindAll(ctx, nil, nil, nil)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, api := range allAPIs {
		if !docURLs[api.Url] {
			r.logger.Info("删除多余API", "url", api.Url, "name", api.Name)
			if err := r.repo.Delete(ctx, api.ID); err != nil {
				return count, err
			}
			count++
		}
	}

	return count, nil
}
