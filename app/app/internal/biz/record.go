package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type EthUserRecord struct {
	ID       int64
	UserId   int64
	Hash     string
	Status   string
	Type     string
	Amount   string
	CoinType string
}

type Location struct {
	ID           int64
	UserId       int64
	Status       string
	CurrentLevel int64
	Current      int64
	CurrentMax   int64
	Row          int64
	Col          int64
}

type RecordUseCase struct {
	ethUserRecordRepo EthUserRecordRepo
	locationRepo      LocationRepo
	userBalanceRepo   UserBalanceRepo
	tx                Transaction
	log               *log.Helper
}

type EthUserRecordRepo interface {
	GetEthUserRecordListByHash(ctx context.Context, hash ...string) (map[string]*EthUserRecord, error)
}

type LocationRepo interface {
	CreateLocation(ctx context.Context, rel *Location) (*Location, error)
	GetLocationLast(ctx context.Context) (*Location, error)
	GetRunningLocationByUserId(ctx context.Context, userId int64) (*Location, error)
	GetRewardLocationByRowOrCol(ctx context.Context, row int64, col int64) ([]*Location, error)
	UpdateLocation(ctx context.Context, location *Location) (*Location, error)
}

func NewRecordUseCase(ethUserRecordRepo EthUserRecordRepo, locationRepo LocationRepo, userBalanceRepo UserBalanceRepo, tx Transaction, logger log.Logger) *RecordUseCase {
	return &RecordUseCase{
		ethUserRecordRepo: ethUserRecordRepo,
		locationRepo:      locationRepo,
		userBalanceRepo:   userBalanceRepo,
		tx:                tx,
		log:               log.NewHelper(logger),
	}
}

func (ruc *RecordUseCase) GetEthUserRecordByTxHash(ctx context.Context, txHash ...string) (map[string]*EthUserRecord, error) {
	return ruc.ethUserRecordRepo.GetEthUserRecordListByHash(ctx, txHash...)
}

func (ruc *RecordUseCase) EthUserRecordHandle(ctx context.Context, ethUserRecord ...*EthUserRecord) (bool, error) {

	// todo 加入数据库，判断复投，直推，单个点位

	for _, v := range ethUserRecord {
		var (
			lastLocation         *Location
			currentValue         int64
			locationCurrentLevel int64
			locationCurrent      int64
			locationCurrentMax   int64
			locationRow          int64
			locationCol          int64
			currentLocation      *Location
			rewardLocations      []*Location
			err                  error
		)

		if "DHB" == v.CoinType {
			continue
		}
		// 获取当前用户的占位信息，已经有运行中的跳过
		_, err = ruc.locationRepo.GetRunningLocationByUserId(ctx, v.UserId)
		if nil == err {
			continue
		}

		// 获取最后一行数据
		lastLocation, err = ruc.locationRepo.GetLocationLast(ctx)
		if nil == err {
			locationRow = 1
			locationCol = 1
		} else {
			if 3 > lastLocation.Col {
				locationCol += 1
				locationRow = lastLocation.Row
			} else {
				locationCol = 1
				locationRow = lastLocation.Row + 1
			}
		}

		if "100000000000000000000" == v.Amount {
			locationCurrentLevel = 1
			locationCurrentMax = 5000000000000
			currentValue = 1000000000000
		} else if "200000000000000000000" == v.Amount {
			locationCurrentLevel = 2
			locationCurrentMax = 10000000000000
			currentValue = 2000000000000
		} else if "500000000000000000000" == v.Amount {
			locationCurrentLevel = 3
			locationCurrentMax = 25000000000000
			currentValue = 5000000000000
		} else {
			continue
		}

		// 占位分红人
		rewardLocations, err = ruc.locationRepo.GetRewardLocationByRowOrCol(ctx, locationRow, locationCol)

		if err = ruc.tx.ExecTx(ctx, func(ctx context.Context) error { // 事务
			currentLocation, err = ruc.locationRepo.CreateLocation(ctx, &Location{
				UserId:       v.UserId,
				Status:       "running",
				CurrentLevel: locationCurrentLevel,
				Current:      locationCurrent,
				CurrentMax:   locationCurrentMax,
				Row:          locationRow,
				Col:          locationCol,
			})
			if nil != err {
				return err
			}

			// 占位分红人分红
			if nil != rewardLocations {
				for _, vRewardLocations := range rewardLocations {
					if "running" != vRewardLocations.Status {
						continue
					}
					tmpAmount := currentValue / 10000000000
					tmpBalanceAmount := tmpAmount
					tmpCurrent := vRewardLocations.Current
					vRewardLocations.Current += tmpAmount

					if vRewardLocations.Current >= vRewardLocations.CurrentMax {
						tmpBalanceAmount = vRewardLocations.CurrentMax - tmpCurrent
						vRewardLocations.Current = vRewardLocations.CurrentMax
						vRewardLocations.Status = "stop"
					}

					_, err = ruc.locationRepo.UpdateLocation(ctx, vRewardLocations) // 分红占位数据修改
					if nil != err {
						return err
					}
					_, err = ruc.userBalanceRepo.LocationReward(ctx, vRewardLocations.UserId, tmpBalanceAmount, currentLocation.ID, vRewardLocations.ID) // 分红信息修改
					if nil != err {
						return err
					}

				}

			}

			return nil
		}); nil != err {
			continue
		}
	}

	// todo 占位

	// todo 直推人分红，推荐人等级调整，推荐人每月五人

	// todo 分红,分红人金额检测满额

	//

	return true, nil
}
