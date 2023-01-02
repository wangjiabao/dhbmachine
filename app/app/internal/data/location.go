package data

import (
	"context"
	"dhb/app/app/internal/biz"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"time"
)

type Location struct {
	ID           int64     `gorm:"primarykey;type:int"`
	UserId       int64     `gorm:"type:int;not null"`
	Row          int64     `gorm:"type:int;not null"`
	Col          int64     `gorm:"type:int;not null"`
	Status       string    `gorm:"type:varchar(45);not null"`
	CurrentLevel int64     `gorm:"type:int;not null"`
	Current      int64     `gorm:"type:bigint;not null"`
	CurrentMax   int64     `gorm:"type:bigint;not null"`
	CreatedAt    time.Time `gorm:"type:datetime;not null"`
	UpdatedAt    time.Time `gorm:"type:datetime;not null"`
}

type LocationRepo struct {
	data *Data
	log  *log.Helper
}

func NewLocationRepo(data *Data, logger log.Logger) biz.LocationRepo {
	return &LocationRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// CreateLocation .
func (lr *LocationRepo) CreateLocation(ctx context.Context, rel *biz.Location) (*biz.Location, error) {
	var location Location
	location.Col = rel.Col
	location.Row = rel.Row
	location.Status = rel.Status
	location.Current = rel.Current
	location.CurrentMax = rel.CurrentMax
	location.CurrentLevel = rel.CurrentLevel
	location.UserId = rel.UserId
	res := lr.data.DB(ctx).Table("location").Create(&location)
	if res.Error != nil {
		return nil, errors.New(500, "CREATE_LOCATION_ERROR", "占位信息创建失败")
	}

	return &biz.Location{
		ID:           location.ID,
		UserId:       location.UserId,
		Status:       location.Status,
		CurrentLevel: location.CurrentLevel,
		Current:      location.Current,
		CurrentMax:   location.CurrentMax,
		Row:          location.Row,
		Col:          location.Col,
	}, nil
}

// GetLocationLast .
func (lr *LocationRepo) GetLocationLast(ctx context.Context) (*biz.Location, error) {
	var location Location
	if err := lr.data.db.Table("location").Order("id desc").First(&location).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("LOCATION_NOT_FOUND", "location not found")
		}

		return nil, errors.New(500, "LOCATION ERROR", err.Error())
	}

	return &biz.Location{
		ID:           location.ID,
		UserId:       location.UserId,
		Status:       location.Status,
		CurrentLevel: location.CurrentLevel,
		Current:      location.Current,
		CurrentMax:   location.CurrentMax,
		Row:          location.Row,
		Col:          location.Col,
	}, nil
}

// GetLocationsByUserId .
func (lr *LocationRepo) GetLocationsByUserId(ctx context.Context, userId int64) ([]*biz.Location, error) {
	var locations []*Location
	res := make([]*biz.Location, 0)
	if err := lr.data.db.Table("location").
		Where("user_id=?", userId).Find(&locations).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res, errors.NotFound("LOCATION_NOT_FOUND", "location not found")
		}

		return nil, errors.New(500, "LOCATION ERROR", err.Error())
	}

	for _, location := range locations {
		res = append(res, &biz.Location{
			ID:           location.ID,
			UserId:       location.UserId,
			Status:       location.Status,
			CurrentLevel: location.CurrentLevel,
			Current:      location.Current,
			CurrentMax:   location.CurrentMax,
			Row:          location.Row,
			Col:          location.Col,
		})
	}

	return res, nil
}

// UpdateLocation .
func (lr *LocationRepo) UpdateLocation(ctx context.Context, l *biz.Location) (*biz.Location, error) {
	var location Location
	location.Status = l.Status
	location.Current = l.Current
	if err := lr.data.db.Table("location").Where("id=?", l.ID).
		Updates(&location).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("LOCATION_NOT_FOUND", "location not found")
		}

		return nil, errors.New(500, "LOCATION ERROR", err.Error())
	}

	return &biz.Location{
		ID:           location.ID,
		UserId:       location.UserId,
		Status:       location.Status,
		CurrentLevel: location.CurrentLevel,
		Current:      location.Current,
		CurrentMax:   location.CurrentMax,
		Row:          location.Row,
		Col:          location.Col,
	}, nil
}

// GetRewardLocationByRowOrCol .
func (lr *LocationRepo) GetRewardLocationByRowOrCol(ctx context.Context, row int64, col int64) ([]*biz.Location, error) {
	var (
		rowMin    int64 = 1
		rowMax    int64
		locations []*Location
	)
	if row > 25 {
		rowMin = row - 25
	}
	rowMax = row + 25

	if err := lr.data.db.Table("location").
		Where("row=? or (col=? and row>=? and row<=?)", row, col, rowMin, rowMax).
		Find(&locations).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("LOCATION_NOT_FOUND", "location not found")
		}

		return nil, errors.New(500, "LOCATION ERROR", err.Error())
	}

	res := make([]*biz.Location, 0)
	for _, location := range locations {
		res = append(res, &biz.Location{
			ID:           location.ID,
			UserId:       location.UserId,
			Status:       location.Status,
			CurrentLevel: location.CurrentLevel,
			Current:      location.Current,
			CurrentMax:   location.CurrentMax,
			Row:          location.Row,
			Col:          location.Col,
		})
	}

	return res, nil
}

// GetRewardLocationByIds .
func (lr *LocationRepo) GetRewardLocationByIds(ctx context.Context, ids ...int64) (map[int64]*biz.Location, error) {
	var locations []*Location
	if err := lr.data.db.Table("location").
		Where("id IN (?)", ids).
		Find(&locations).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("LOCATION_NOT_FOUND", "location not found")
		}

		return nil, errors.New(500, "LOCATION ERROR", err.Error())
	}

	res := make(map[int64]*biz.Location, 0)
	for _, location := range locations {
		res[location.ID] = &biz.Location{
			ID:           location.ID,
			UserId:       location.UserId,
			Status:       location.Status,
			CurrentLevel: location.CurrentLevel,
			Current:      location.Current,
			CurrentMax:   location.CurrentMax,
			Row:          location.Row,
			Col:          location.Col,
		}
	}

	return res, nil
}
