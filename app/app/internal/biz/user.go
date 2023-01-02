package biz

import (
	"context"
	v1 "dhb/app/app/api"
	"encoding/base64"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"strconv"
	"strings"
	"time"
)

type User struct {
	ID        int64
	Address   string
	CreatedAt time.Time
}

type UserInfo struct {
	ID               int64
	UserId           int64
	Vip              int64
	HistoryRecommend int64
}

type UserRecommend struct {
	ID            int64
	UserId        int64
	RecommendCode string
}

type UserCurrentMonthRecommend struct {
	ID              int64
	UserId          int64
	RecommendUserId int64
	Date            time.Time
}

type Config struct {
	ID      int64
	KeyName string
	Name    string
	Value   string
}

type UserBalance struct {
	ID          int64
	UserId      int64
	BalanceUsdt int64
	BalanceDhb  int64
}

type Withdraw struct {
	ID              int64
	UserId          int64
	Amount          int64
	RelAmount       int64
	BalanceRecordId int64
	Status          string
	Type            string
	CreatedAt       time.Time
}

type UserUseCase struct {
	repo                          UserRepo
	urRepo                        UserRecommendRepo
	configRepo                    ConfigRepo
	uiRepo                        UserInfoRepo
	ubRepo                        UserBalanceRepo
	locationRepo                  LocationRepo
	userCurrentMonthRecommendRepo UserCurrentMonthRecommendRepo
	tx                            Transaction
	log                           *log.Helper
}

type Reward struct {
	ID               int64
	UserId           int64
	Amount           int64
	BalanceRecordId  int64
	Type             string
	TypeRecordId     int64
	Reason           string
	ReasonLocationId int64
	LocationType     string
	CreatedAt        time.Time
}

type ConfigRepo interface {
	GetConfigByKeys(ctx context.Context, keys ...string) ([]*Config, error)
}

type UserBalanceRepo interface {
	CreateUserBalance(ctx context.Context, u *User) (*UserBalance, error)
	LocationReward(ctx context.Context, userId int64, amount int64, locationId int64, myLocationId int64, locationType string) (int64, error)
	RecommendReward(ctx context.Context, userId int64, amount int64, locationId int64) (int64, error)
	FirstRecommendReward(ctx context.Context, userId int64, amount int64, locationId int64) (int64, error)
	Deposit(ctx context.Context, userId int64, amount int64) (int64, error)
	GetUserBalance(ctx context.Context, userId int64) (*UserBalance, error)
	GetUserRewardByUserId(ctx context.Context, userId int64) ([]*Reward, error)
	GetUserRewards(ctx context.Context) ([]*Reward, error)
	GetUserBalanceByUserIds(ctx context.Context, userIds ...int64) (map[int64]*UserBalance, error)
	Withdraw(ctx context.Context, userId int64, amount int64, coinType string) (*Withdraw, error)
	GetWithdrawByUserId(ctx context.Context, userId int64) ([]*Withdraw, error)
	GetWithdraws(ctx context.Context) ([]*Withdraw, error)
}

type UserRecommendRepo interface {
	GetUserRecommendByUserId(ctx context.Context, userId int64) (*UserRecommend, error)
	CreateUserRecommend(ctx context.Context, u *User, recommendUser *UserRecommend) (*UserRecommend, error)
	GetUserRecommendByCode(ctx context.Context, code string) ([]*UserRecommend, error)
}

type UserCurrentMonthRecommendRepo interface {
	GetUserCurrentMonthRecommendByUserId(ctx context.Context, userId int64) ([]*UserCurrentMonthRecommend, error)
	CreateUserCurrentMonthRecommend(ctx context.Context, u *UserCurrentMonthRecommend) (*UserCurrentMonthRecommend, error)
	GetUserCurrentMonthRecommendCountByUserIds(ctx context.Context, userIds ...int64) (map[int64]int64, error)
}

type UserInfoRepo interface {
	CreateUserInfo(ctx context.Context, u *User) (*UserInfo, error)
	GetUserInfoByUserId(ctx context.Context, userId int64) (*UserInfo, error)
	UpdateUserInfo(ctx context.Context, u *UserInfo) (*UserInfo, error)
	GetUserInfoByUserIds(ctx context.Context, userIds ...int64) (map[int64]*UserInfo, error)
}

type UserRepo interface {
	GetUserById(ctx context.Context, Id int64) (*User, error)
	GetUserByAddresses(ctx context.Context, Addresses ...string) (map[string]*User, error)
	GetUserByAddress(ctx context.Context, address string) (*User, error)
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByUserIds(ctx context.Context, userIds ...int64) (map[int64]*User, error)
	GetUsers(ctx context.Context) ([]*User, error)
}

func NewUserUseCase(repo UserRepo, tx Transaction, configRepo ConfigRepo, uiRepo UserInfoRepo, urRepo UserRecommendRepo, locationRepo LocationRepo, userCurrentMonthRecommendRepo UserCurrentMonthRecommendRepo, ubRepo UserBalanceRepo, logger log.Logger) *UserUseCase {
	return &UserUseCase{
		repo:                          repo,
		tx:                            tx,
		configRepo:                    configRepo,
		locationRepo:                  locationRepo,
		userCurrentMonthRecommendRepo: userCurrentMonthRecommendRepo,
		uiRepo:                        uiRepo,
		urRepo:                        urRepo,
		ubRepo:                        ubRepo,
		log:                           log.NewHelper(logger),
	}
}

func (uuc *UserUseCase) GetUserByAddress(ctx context.Context, Addresses ...string) (map[string]*User, error) {
	return uuc.repo.GetUserByAddresses(ctx, Addresses...)
}

func (uuc *UserUseCase) GetExistUserByAddressOrCreate(ctx context.Context, u *User, req *v1.EthAuthorizeRequest) (*User, error) {
	var (
		user          *User
		recommendUser *UserRecommend
		userRecommend *UserRecommend
		userInfo      *UserInfo
		userBalance   *UserBalance
		err           error
		userId        int64
		decodeBytes   []byte
	)

	user, err = uuc.repo.GetUserByAddress(ctx, u.Address) // 查询用户
	if nil == user || nil != err {
		code := req.SendBody.Code // 查询推荐码 abf00dd52c08a9213f225827bc3fb100 md5 dhbmachinefirst
		if "abf00dd52c08a9213f225827bc3fb100" != code {
			decodeBytes, err = base64.StdEncoding.DecodeString(code)
			code = string(decodeBytes)
			if 1 >= len(code) {
				return nil, errors.New(500, "USER_ERROR", "无效的推荐码")
			}
			if userId, err = strconv.ParseInt(code[1:], 10, 64); 0 >= userId || nil != err {
				return nil, errors.New(500, "USER_ERROR", "无效的推荐码")
			}

			// 查询推荐人的相关信息
			recommendUser, err = uuc.urRepo.GetUserRecommendByUserId(ctx, userId)
			if err != nil {
				return nil, errors.New(500, "USER_ERROR", "无效的推荐码")
			}
		}

		if err = uuc.tx.ExecTx(ctx, func(ctx context.Context) error { // 事务
			user, err = uuc.repo.CreateUser(ctx, u) // 用户创建
			if err != nil {
				return err
			}

			userInfo, err = uuc.uiRepo.CreateUserInfo(ctx, user) // 创建用户信息
			if err != nil {
				return err
			}

			userRecommend, err = uuc.urRepo.CreateUserRecommend(ctx, user, recommendUser) // 创建用户信息
			if err != nil {
				return err
			}

			userBalance, err = uuc.ubRepo.CreateUserBalance(ctx, user) // 创建余额信息
			if err != nil {
				return err
			}

			return nil
		}); err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (uuc *UserUseCase) UserInfo(ctx context.Context, user *User) (*v1.UserInfoReply, error) {
	var (
		myUser                     *User
		userInfo                   *UserInfo
		locations                  []*Location
		userBalance                *UserBalance
		userRecommend              *UserRecommend
		userRecommends             []*UserRecommend
		userRewards                []*Reward
		rewardLocations            []*Location
		userCurrentMonthRecommends []*UserCurrentMonthRecommend
		userRewardTotal            int64
		encodeString               string
		myUserRecommendUserId      int64
		myRecommendUser            *User
		myRow                      int64
		rowNum                     int64
		colNum                     int64
		myCol                      int64
		recommendTeamNum           int64
		recommendTotal             int64
		locationTotal              int64
		myCode                     string
		inviteUserAddress          string
		amount                     string
		status                     string
		currentMonthRecommendNum   int64
		configs                    []*Config
		level1Dhb                  string
		level2Dhb                  string
		level3Dhb                  string
		err                        error
	)

	myUser, err = uuc.repo.GetUserById(ctx, user.ID)
	if nil != err {
		return nil, err
	}

	userInfo, err = uuc.uiRepo.GetUserInfoByUserId(ctx, myUser.ID)
	if nil != err {
		return nil, err
	}

	locations, err = uuc.locationRepo.GetLocationsByUserId(ctx, myUser.ID)
	if nil != err {
		return nil, err
	}
	for _, v := range locations {
		if "running" == v.Status {
			status = "running"
			amount = fmt.Sprintf("%.2f", float64(v.CurrentMax-v.Current)/float64(10000000000))
			myCol = v.Col
			myRow = v.Row
			break
		}
	}

	userBalance, err = uuc.ubRepo.GetUserBalance(ctx, myUser.ID)
	if nil != err {
		return nil, err
	}

	userRecommend, err = uuc.urRepo.GetUserRecommendByUserId(ctx, myUser.ID)
	if nil == userRecommend {
		return nil, err
	}

	myCode = "D" + strconv.FormatInt(myUser.ID, 10)
	codeByte := []byte(myCode)
	encodeString = base64.StdEncoding.EncodeToString(codeByte)

	if "" != userRecommend.RecommendCode {
		tmpRecommendUserIds := strings.Split(userRecommend.RecommendCode, "D")
		if 2 <= len(tmpRecommendUserIds) {
			myUserRecommendUserId, _ = strconv.ParseInt(tmpRecommendUserIds[len(tmpRecommendUserIds)-1], 10, 64) // 最后一位是直推人
		}
		myRecommendUser, err = uuc.repo.GetUserById(ctx, myUserRecommendUserId)
		if nil != err {
			return nil, err
		}
		inviteUserAddress = myRecommendUser.Address
		myCode = userRecommend.RecommendCode + myCode
	}

	// 团队
	userRecommends, err = uuc.urRepo.GetUserRecommendByCode(ctx, myCode)
	if nil != userRecommends {
		recommendTeamNum = int64(len(userRecommends))
	}

	// 累计奖励
	userRewards, err = uuc.ubRepo.GetUserRewardByUserId(ctx, myUser.ID)
	if nil != userRewards {
		for _, vUserReward := range userRewards {
			userRewardTotal += vUserReward.Amount
			if "recommend" == vUserReward.Reason || "recommend_vip" == vUserReward.Reason {
				recommendTotal += vUserReward.Amount
			} else if "location" == vUserReward.Reason {
				locationTotal += vUserReward.Amount
			}
		}
	}

	// 位置
	if 0 < myRow && 0 < myCol {
		rewardLocations, err = uuc.locationRepo.GetRewardLocationByRowOrCol(ctx, myRow, myCol)
		if nil != rewardLocations {
			for _, vRewardLocation := range rewardLocations {
				if myRow == vRewardLocation.Row {
					colNum++
				}
				if myCol == vRewardLocation.Col {
					rowNum++
				}
			}
		}
	}

	// 当月推荐人数
	userCurrentMonthRecommends, err = uuc.userCurrentMonthRecommendRepo.GetUserCurrentMonthRecommendByUserId(ctx, myUser.ID)
	if nil != userCurrentMonthRecommends {
		now := time.Now().UTC().Add(8 * time.Hour)
		for _, vUserCurrentMonthRecommend := range userCurrentMonthRecommends {
			if vUserCurrentMonthRecommend.Date.After(time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)) {
				currentMonthRecommendNum++
			}
		}
	}

	// 配置
	configs, err = uuc.configRepo.GetConfigByKeys(ctx, "level1Dhb", "level2Dhb", "level3Dhb")
	if nil != configs {
		for _, vConfig := range configs {
			if "level1Dhb" == vConfig.KeyName {
				level1Dhb = vConfig.Value
			} else if "level2Dhb" == vConfig.KeyName {
				level2Dhb = vConfig.Value
			} else if "level3Dhb" == vConfig.KeyName {
				level3Dhb = vConfig.Value
			}
		}
	}

	return &v1.UserInfoReply{
		Address:                  myUser.Address,
		Level:                    userInfo.Vip,
		Status:                   status,
		Amount:                   amount,
		BalanceUsdt:              fmt.Sprintf("%.2f", float64(userBalance.BalanceUsdt)/float64(10000000000)),
		BalanceDhb:               fmt.Sprintf("%.2f", float64(userBalance.BalanceDhb)/float64(10000000000)),
		InviteUrl:                encodeString,
		InviteUserAddress:        inviteUserAddress,
		RecommendNum:             userInfo.HistoryRecommend,
		RecommendTeamNum:         recommendTeamNum,
		Total:                    fmt.Sprintf("%.2f", float64(userRewardTotal)/float64(10000000000)),
		Row:                      rowNum,
		Col:                      colNum,
		CurrentMonthRecommendNum: currentMonthRecommendNum,
		RecommendTotal:           fmt.Sprintf("%.2f", float64(recommendTotal)/float64(10000000000)),
		LocationTotal:            fmt.Sprintf("%.2f", float64(locationTotal)/float64(10000000000)),
		Level1Dhb:                level1Dhb,
		Level2Dhb:                level2Dhb,
		Level3Dhb:                level3Dhb,
		Usdt:                     "0x337610d27c682E347C9cD60BD4b3b107C9d34dDd",
		Dhb:                      "0x96BD81715c69eE013405B4005Ba97eA1f420fd87",
		Account:                  "0xe865f2e5ff04b8b7952d1c0d9163a91f313b158f",
	}, nil
}

func (uuc *UserUseCase) RewardList(ctx context.Context, req *v1.RewardListRequest, user *User) (*v1.RewardListReply, error) {
	var (
		userRewards    []*Reward
		locationIdsMap map[int64]int64
		locations      map[int64]*Location
		err            error
	)
	res := &v1.RewardListReply{
		Rewards: make([]*v1.RewardListReply_List, 0),
	}

	userRewards, err = uuc.ubRepo.GetUserRewardByUserId(ctx, user.ID)
	if nil != err {
		return res, nil
	}

	locationIdsMap = make(map[int64]int64, 0)
	if nil != userRewards {
		for _, vUserReward := range userRewards {
			if "location" == vUserReward.Reason && req.Type == vUserReward.LocationType && 1 <= vUserReward.ReasonLocationId {
				locationIdsMap[vUserReward.ReasonLocationId] = vUserReward.ReasonLocationId
			}
		}

		var tmpLocationIds []int64
		for _, v := range locationIdsMap {
			tmpLocationIds = append(tmpLocationIds, v)
		}
		if 0 >= len(tmpLocationIds) {
			return res, nil
		}

		locations, err = uuc.locationRepo.GetRewardLocationByIds(ctx, tmpLocationIds...)

		for _, vUserReward := range userRewards {
			if "location" == vUserReward.Reason && req.Type == vUserReward.LocationType {
				if _, ok := locations[vUserReward.ReasonLocationId]; !ok {
					continue
				}

				res.Rewards = append(res.Rewards, &v1.RewardListReply_List{
					CreatedAt:      vUserReward.CreatedAt.Format("2006-01-02 15:04:05"),
					Amount:         fmt.Sprintf("%.2f", float64(vUserReward.Amount)/float64(10000000000)),
					LocationStatus: locations[vUserReward.ReasonLocationId].Status,
					Type:           vUserReward.Type,
				})
			}
		}
	}

	return res, nil
}

func (uuc *UserUseCase) RecommendRewardList(ctx context.Context, user *User) (*v1.RecommendRewardListReply, error) {
	var (
		userRewards []*Reward
		err         error
	)
	res := &v1.RecommendRewardListReply{
		Rewards: make([]*v1.RecommendRewardListReply_List, 0),
	}

	userRewards, err = uuc.ubRepo.GetUserRewardByUserId(ctx, user.ID)
	if nil != err {
		return res, nil
	}

	for _, vUserReward := range userRewards {
		if "recommend" == vUserReward.Reason || "recommend_vip" == vUserReward.Reason {
			res.Rewards = append(res.Rewards, &v1.RecommendRewardListReply_List{
				CreatedAt: vUserReward.CreatedAt.Format("2006-01-02 15:04:05"),
				Amount:    fmt.Sprintf("%.2f", float64(vUserReward.Amount)/float64(10000000000)),
				Type:      vUserReward.Type,
				Reason:    vUserReward.Reason,
			})
		}
	}

	return res, nil
}

func (uuc *UserUseCase) FeeRewardList(ctx context.Context, user *User) (*v1.FeeRewardListReply, error) {
	var (
		userRewards []*Reward
		err         error
	)
	res := &v1.FeeRewardListReply{
		Rewards: make([]*v1.FeeRewardListReply_List, 0),
	}

	userRewards, err = uuc.ubRepo.GetUserRewardByUserId(ctx, user.ID)
	if nil != err {
		return res, nil
	}

	for _, vUserReward := range userRewards {
		if "fee" == vUserReward.Reason {
			res.Rewards = append(res.Rewards, &v1.FeeRewardListReply_List{
				CreatedAt: vUserReward.CreatedAt.Format("2006-01-02 15:04:05"),
				Amount:    fmt.Sprintf("%.2f", float64(vUserReward.Amount)/float64(10000000000)),
			})
		}
	}

	return res, nil
}

func (uuc *UserUseCase) WithdrawList(ctx context.Context, user *User) (*v1.WithdrawListReply, error) {

	var (
		withdraws []*Withdraw
		err       error
	)

	res := &v1.WithdrawListReply{
		Withdraw: make([]*v1.WithdrawListReply_List, 0),
	}

	withdraws, err = uuc.ubRepo.GetWithdrawByUserId(ctx, user.ID)
	if nil != err {
		return res, err
	}

	for _, v := range withdraws {
		res.Withdraw = append(res.Withdraw, &v1.WithdrawListReply_List{
			CreatedAt: v.CreatedAt.Format("2006-01-02 15:04:05"),
			Amount:    fmt.Sprintf("%.2f", float64(v.Amount)/float64(10000000000)),
			Status:    v.Status,
			Type:      v.Type,
		})
	}

	return res, nil
}

func (uuc *UserUseCase) Withdraw(ctx context.Context, req *v1.WithdrawRequest, user *User) (*v1.WithdrawReply, error) {

	if "dhb" != req.SendBody.Type && "usdt" != req.SendBody.Type {
		return &v1.WithdrawReply{
			Status: "fail",
		}, nil
	}

	amountFloat, _ := strconv.ParseFloat(req.SendBody.Amount, 10)
	amountFloat *= 10000000000
	amount, _ := strconv.ParseInt(strconv.FormatFloat(amountFloat, 'f', -1, 64), 10, 64)
	if 0 >= amount {
		return &v1.WithdrawReply{
			Status: "fail",
		}, nil
	}

	_, err := uuc.ubRepo.Withdraw(ctx, user.ID, amount, req.SendBody.Type)
	if nil != err {
		return nil, err
	}

	return &v1.WithdrawReply{
		Status: "ok",
	}, nil
}

func (uuc *UserUseCase) AdminRewardList(ctx context.Context, req *v1.AdminRewardListRequest) (*v1.AdminRewardListReply, error) {
	var (
		userRewards []*Reward
		users       map[int64]*User
		userIdsMap  map[int64]int64
		userIds     []int64
		err         error
	)
	res := &v1.AdminRewardListReply{
		Rewards: make([]*v1.AdminRewardListReply_List, 0),
	}

	userRewards, err = uuc.ubRepo.GetUserRewards(ctx)
	if nil != err {
		return res, nil
	}

	userIdsMap = make(map[int64]int64, 0)
	for _, vUserReward := range userRewards {
		userIdsMap[vUserReward.UserId] = vUserReward.UserId
	}
	for _, v := range userIdsMap {
		userIds = append(userIds, v)
	}

	users, err = uuc.repo.GetUserByUserIds(ctx, userIds...)
	if nil != err {
		return res, nil
	}

	for _, vUserReward := range userRewards {
		if _, ok := users[vUserReward.UserId]; !ok {
			continue
		}

		res.Rewards = append(res.Rewards, &v1.AdminRewardListReply_List{
			CreatedAt: vUserReward.CreatedAt.Format("2006-01-02 15:04:05"),
			Amount:    fmt.Sprintf("%.2f", float64(vUserReward.Amount)/float64(10000000000)),
			Type:      vUserReward.Type,
			Address:   users[vUserReward.UserId].Address,
			Reason:    vUserReward.Reason,
		})
	}

	return res, nil
}

func (uuc *UserUseCase) AdminUserList(ctx context.Context, req *v1.AdminUserListRequest) (*v1.AdminUserListReply, error) {
	var (
		users                          []*User
		userIds                        []int64
		userBalances                   map[int64]*UserBalance
		userInfos                      map[int64]*UserInfo
		userCurrentMonthRecommendCount map[int64]int64
		err                            error
	)
	res := &v1.AdminUserListReply{
		Users: make([]*v1.AdminUserListReply_UserList, 0),
	}

	users, err = uuc.repo.GetUsers(ctx)
	if nil != err {
		return res, nil
	}
	for _, vUsers := range users {
		userIds = append(userIds, vUsers.ID)
	}

	userBalances, err = uuc.ubRepo.GetUserBalanceByUserIds(ctx, userIds...)
	if nil != err {
		return res, nil
	}

	userInfos, err = uuc.uiRepo.GetUserInfoByUserIds(ctx, userIds...)
	if nil != err {
		return res, nil
	}

	userCurrentMonthRecommendCount, err = uuc.userCurrentMonthRecommendRepo.GetUserCurrentMonthRecommendCountByUserIds(ctx, userIds...)

	for _, v := range users {
		if _, ok := userBalances[v.ID]; !ok {
			continue
		}
		if _, ok := userInfos[v.ID]; !ok {
			continue
		}

		var tmpCount int64
		if nil != userCurrentMonthRecommendCount {
			if _, ok := userCurrentMonthRecommendCount[v.ID]; ok {
				tmpCount = userCurrentMonthRecommendCount[v.ID]
			}
		}

		res.Users = append(res.Users, &v1.AdminUserListReply_UserList{
			CreatedAt:        v.CreatedAt.Format("2006-01-02 15:04:05"),
			Address:          v.Address,
			BalanceUsdt:      fmt.Sprintf("%.2f", float64(userBalances[v.ID].BalanceUsdt)/float64(10000000000)),
			BalanceDhb:       fmt.Sprintf("%.2f", float64(userBalances[v.ID].BalanceDhb)/float64(10000000000)),
			Vip:              userInfos[v.ID].Vip,
			MonthRecommend:   tmpCount,
			HistoryRecommend: userInfos[v.ID].HistoryRecommend,
		})
	}

	return res, nil
}

func (uuc *UserUseCase) AdminLocationList(ctx context.Context, req *v1.AdminLocationListRequest) (*v1.AdminLocationListReply, error) {
	var (
		locations  []*Location
		userIds    []int64
		userIdsMap map[int64]int64
		users      map[int64]*User
		err        error
	)

	res := &v1.AdminLocationListReply{
		Locations: make([]*v1.AdminLocationListReply_LocationList, 0),
	}

	locations, err = uuc.locationRepo.GetLocations(ctx)
	if nil != err {
		return res, nil
	}

	userIdsMap = make(map[int64]int64, 0)
	for _, vLocations := range locations {
		userIdsMap[vLocations.UserId] = vLocations.UserId
	}
	for _, v := range userIdsMap {
		userIds = append(userIds, v)
	}

	users, err = uuc.repo.GetUserByUserIds(ctx, userIds...)
	if nil != err {
		return res, nil
	}

	for _, v := range locations {
		if _, ok := users[v.UserId]; !ok {
			continue
		}

		res.Locations = append(res.Locations, &v1.AdminLocationListReply_LocationList{
			CreatedAt:    v.CreatedAt.Format("2006-01-02 15:04:05"),
			Address:      users[v.UserId].Address,
			Row:          v.Row,
			Col:          v.Col,
			Status:       v.Status,
			CurrentLevel: v.CurrentLevel,
			Current:      fmt.Sprintf("%.2f", float64(v.Current)/float64(10000000000)),
			CurrentMax:   fmt.Sprintf("%.2f", float64(v.CurrentMax)/float64(10000000000)),
		})
	}

	return res, nil

}

func (uuc *UserUseCase) AdminWithdrawList(ctx context.Context, req *v1.AdminWithdrawListRequest) (*v1.AdminWithdrawListReply, error) {
	var (
		withdraws  []*Withdraw
		userIds    []int64
		userIdsMap map[int64]int64
		users      map[int64]*User
		err        error
	)

	res := &v1.AdminWithdrawListReply{
		Withdraw: make([]*v1.AdminWithdrawListReply_List, 0),
	}

	withdraws, err = uuc.ubRepo.GetWithdraws(ctx)
	if nil != err {
		return res, err
	}

	userIdsMap = make(map[int64]int64, 0)
	for _, vWithdraws := range withdraws {
		userIdsMap[vWithdraws.UserId] = vWithdraws.UserId
	}
	for _, v := range userIdsMap {
		userIds = append(userIds, v)
	}

	users, err = uuc.repo.GetUserByUserIds(ctx, userIds...)
	if nil != err {
		return res, nil
	}

	for _, v := range withdraws {
		if _, ok := users[v.UserId]; !ok {
			continue
		}
		res.Withdraw = append(res.Withdraw, &v1.AdminWithdrawListReply_List{
			CreatedAt: v.CreatedAt.Format("2006-01-02 15:04:05"),
			Amount:    fmt.Sprintf("%.2f", float64(v.Amount)/float64(10000000000)),
			Status:    v.Status,
			Type:      v.Type,
			Address:   users[v.UserId].Address,
			RelAmount: fmt.Sprintf("%.2f", float64(v.RelAmount)/float64(10000000000)),
		})
	}

	return res, nil

}
