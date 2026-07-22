package user

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&User{}, &UserPointsLog{}))
	return db
}

func seedUsers(t *testing.T, db *gorm.DB, count int) {
	t.Helper()
	for i := 0; i < count; i++ {
		username := "user" + string(rune('a'+i))
		db.Create(&User{Username: username, Password: "hashed", Mobile: "13800000001", Status: 1})
	}
}

func TestService_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)
	svc := NewService(db, repo, nil)

	ctx := context.Background()
	err := svc.Create(ctx, &CreateRequest{
		Username: "testuser",
		Password: "123456",
		Mobile:   "13800138000",
		NickName: "测试",
		Status:   1,
	})
	require.NoError(t, err)

	var count int64
	db.Model(&User{}).Where("username = ?", "testuser").Count(&count)
	assert.Equal(t, int64(1), count)

	err = svc.Create(ctx, &CreateRequest{
		Username: "testuser",
		Password: "123456",
		Mobile:   "13800138001",
		NickName: "重复",
		Status:   1,
	})
	assert.Error(t, err)
}

func TestService_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)
	svc := NewService(db, repo, nil)

	ctx := context.Background()
	seedUsers(t, db, 5)

	result, err := svc.List(ctx, &ListRequest{Page: 1, Limit: 10})
	require.NoError(t, err)
	assert.Equal(t, int64(5), result.Total)
	assert.Len(t, result.List, 5)

	result, err = svc.List(ctx, &ListRequest{Username: "usera", Page: 1, Limit: 10})
	require.NoError(t, err)
	assert.Equal(t, int64(1), result.Total)
}

func TestService_List_Empty(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)
	svc := NewService(db, repo, nil)

	ctx := context.Background()
	result, err := svc.List(ctx, &ListRequest{Page: 1, Limit: 10})
	require.NoError(t, err)
	assert.Equal(t, int64(0), result.Total)
}

func TestService_Detail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)
	svc := NewService(db, repo, nil)

	ctx := context.Background()
	db.Create(&User{Username: "detailuser", Password: "hashed", Mobile: "13800000001", Status: 1})

	user, err := svc.Detail(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, "detailuser", user.Username)
}

func TestService_Detail_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)
	svc := NewService(db, repo, nil)

	ctx := context.Background()
	_, err := svc.Detail(ctx, 999)
	assert.Error(t, err)
}

func TestService_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)
	svc := NewService(db, repo, nil)

	ctx := context.Background()
	db.Create(&User{Username: "updateuser", Password: "hashed", Mobile: "13800000001", NickName: "旧名", Status: 1})

	err := svc.Update(ctx, &UpdateRequest{
		Id:       1,
		Username: "updateuser",
		Mobile:   "13800000002",
		NickName: "新名",
		Status:   2,
	})
	require.NoError(t, err)

	var user User
	db.First(&user, 1)
	assert.Equal(t, "新名", user.NickName)
	assert.Equal(t, int8(2), user.Status)
}

func TestService_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)
	svc := NewService(db, repo, nil)

	ctx := context.Background()
	db.Create(&User{Username: "deluser", Password: "hashed", Mobile: "13800000001", Status: 1})

	err := svc.Delete(ctx, &DeleteRequest{Id: 1})
	require.NoError(t, err)

	var user User
	err = db.Unscoped().First(&user, 1).Error
	require.NoError(t, err)
	assert.NotZero(t, user.DeletedAt)
}

func TestService_GetPointsLogs(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)
	pointsLogRepo := NewPointsLogRepository(db)
	svc := NewService(db, repo, pointsLogRepo)

	ctx := context.Background()
	db.Create(&User{Username: "pointsuser", Password: "hashed", Mobile: "13800000001", Status: 1})
	for i := 0; i < 3; i++ {
		db.Create(&UserPointsLog{UserId: 1, Points: 10, ChangeType: 1, SourceType: 10})
	}

	result, err := svc.GetPointsLogs(ctx, &PointsLogListRequest{UserId: 1, Page: 1, Limit: 10})
	require.NoError(t, err)
	assert.Equal(t, int64(3), result.Total)
}

func TestService_ChangePoints(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)
	pointsLogRepo := NewPointsLogRepository(db)
	svc := NewService(db, repo, pointsLogRepo)

	ctx := context.Background()
	db.Create(&User{Username: "pointsuser", Password: "hashed", Mobile: "13800000001", Status: 1, Points: 0})

	err := svc.ChangePoints(ctx, &PointsChangeRequest{
		UserId:     1,
		Points:     50,
		ChangeType: 1,
		SourceType: 20,
		SourceId:   "order_001",
		Remark:     "充值",
	})
	require.NoError(t, err)

	var user User
	db.First(&user, 1)
	assert.Equal(t, uint32(50), user.Points)

	var log UserPointsLog
	db.First(&log, 1)
	assert.Equal(t, "order_001", log.SourceId)
}

func TestService_ChangePoints_Reduce(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)
	pointsLogRepo := NewPointsLogRepository(db)
	svc := NewService(db, repo, pointsLogRepo)

	ctx := context.Background()
	db.Create(&User{Username: "pointsuser", Password: "hashed", Mobile: "13800000001", Status: 1, Points: 100})

	err := svc.ChangePoints(ctx, &PointsChangeRequest{
		UserId:     1,
		Points:     30,
		ChangeType: PointsChangeTypeReduce,
		SourceType: 10,
		SourceId:   "order_002",
		Remark:     "消费",
	})
	require.NoError(t, err)

	var user User
	db.First(&user, 1)
	assert.Equal(t, uint32(70), user.Points)
}

func TestService_ChangePoints_InvalidInput(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)
	pointsLogRepo := NewPointsLogRepository(db)
	svc := NewService(db, repo, pointsLogRepo)

	ctx := context.Background()

	err := svc.ChangePoints(ctx, &PointsChangeRequest{
		UserId:     1,
		Points:     0,
		ChangeType: 1,
		SourceType: 20,
	})
	assert.Error(t, err)

	err = svc.ChangePoints(ctx, &PointsChangeRequest{
		UserId:     1,
		Points:     10,
		ChangeType: 1,
		SourceType: 5,
	})
	assert.Error(t, err)
}
