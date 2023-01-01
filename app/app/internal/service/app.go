package service

import (
	"context"
	"dhb/app/app/internal/pkg/middleware/auth"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	jwt2 "github.com/golang-jwt/jwt/v4"
	"io"
	"net/url"
	"strconv"

	v1 "dhb/app/app/api"
	"dhb/app/app/internal/biz"
	"dhb/app/app/internal/conf"
	"io/ioutil"
	"net/http"
	"time"
)

// AppService service.
type AppService struct {
	v1.UnimplementedAppServer

	uuc *biz.UserUseCase
	ruc *biz.RecordUseCase
	log *log.Helper
	ca  *conf.Auth
}

// NewAppService new a service.
func NewAppService(uuc *biz.UserUseCase, ruc *biz.RecordUseCase, logger log.Logger, ca *conf.Auth) *AppService {
	return &AppService{uuc: uuc, ruc: ruc, log: log.NewHelper(logger), ca: ca}
}

// EthAuthorize ethAuthorize.
func (a *AppService) EthAuthorize(ctx context.Context, req *v1.EthAuthorizeRequest) (*v1.EthAuthorizeReply, error) {
	// TODO 有效的参数验证
	userAddress := req.SendBody.Address // 以太坊账户
	if "" == userAddress || 20 > len(userAddress) {
		return nil, errors.New(500, "AUTHORIZE_ERROR", "账户地址参数错误")
	}

	// TODO 验证签名

	// 根据地址查询用户，不存在时则创建
	user, err := a.uuc.GetExistUserByAddressOrCreate(ctx, &biz.User{
		Address: userAddress,
	}, req)
	if err != nil {
		return nil, err
	}

	claims := auth.CustomClaims{
		UserId:   user.ID,
		UserType: "user",
		StandardClaims: jwt2.StandardClaims{
			NotBefore: time.Now().Unix(),              // 签名的生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*7, // 7天过期
			Issuer:    "DHB",
		},
	}
	token, err := auth.CreateToken(claims, a.ca.JwtKey)
	if err != nil {
		return nil, errors.New(500, "AUTHORIZE_ERROR", "生成token失败")
	}

	userInfoRsp := v1.EthAuthorizeReply{
		Token: token,
	}
	return &userInfoRsp, nil
}

// Deposit deposit.
func (a *AppService) Deposit(ctx context.Context, req *v1.DepositRequest) (*v1.DepositReply, error) {

	var (
		depositUsdtResult     map[string]*eth
		depositDhbResult      map[string]*eth
		tmpDepositDhbResult   map[string]*eth
		userDepositDhbResult  map[string]map[string]*eth
		notExistDepositResult []*biz.EthUserRecord
		existEthUserRecords   map[string]*biz.EthUserRecord
		depositUsers          map[string]*biz.User
		fromAccount           []string
		hashKeys              []string
		err                   error
	)

	// 每次一共最多查2000条，所以注意好外层调用的定时查询的时间设置，当然都可以重新定义，
	// 在功能上调用者查询两种币的交易记录，每次都要把数据覆盖查询，是一个较大范围的查找防止遗漏数据，范围最起码要大于实际这段时间的入单量，不能边界查询容易掉单，这样的实现是因为简单
	for i := 1; i < 10; i++ {
		depositUsdtResult, err = requestEthDepositResult(200, int64(i), "0x337610d27c682E347C9cD60BD4b3b107C9d34dDd")
		// 辅助查询
		depositDhbResult, err = requestEthDepositResult(200, int64(i), "0x337610d27c682E347C9cD60BD4b3b107C9d34dDd")
		tmpDepositDhbResult, err = requestEthDepositResult(100, int64(i+1), "0x337610d27c682E347C9cD60BD4b3b107C9d34dDd")
		for kTmpDepositDhbResult, v := range tmpDepositDhbResult {
			if _, ok := tmpDepositDhbResult[kTmpDepositDhbResult]; !ok {
				depositDhbResult[kTmpDepositDhbResult] = v
			}
		}

		if 0 >= len(depositUsdtResult) {
			break
		}
		fmt.Println(depositUsdtResult, err)

		for hashKey, vDepositResult := range depositUsdtResult { // 主查询
			hashKeys = append(hashKeys, hashKey)
			fromAccount = append(fromAccount, vDepositResult.From)
		}
		userDepositDhbResult = make(map[string]map[string]*eth, 0) // 辅助数据
		for k, v := range depositDhbResult {
			hashKeys = append(hashKeys, k)
			fromAccount = append(fromAccount, v.From)
			userDepositDhbResult[v.From][k] = v
		}

		depositUsers, err = a.uuc.GetUserByAddress(ctx, fromAccount...)
		if nil != err || nil == depositUsers {
			continue
		}
		existEthUserRecords, err = a.ruc.GetEthUserRecordByTxHash(ctx, hashKeys...)

		// 统计开始
		notExistDepositResult = make([]*biz.EthUserRecord, 0)
		for _, vDepositUsdtResult := range depositUsdtResult { // 主查usdt
			if _, ok := existEthUserRecords[vDepositUsdtResult.Hash]; ok { // 记录已存在
				continue
			}
			if _, ok := depositUsers[vDepositUsdtResult.From]; !ok { // 用户不存在
				continue
			}
			if _, ok := userDepositDhbResult[vDepositUsdtResult.From]; !ok { // 没有dhb的充值记录
				continue
			}
			var (
				tmpDhbHash, tmpDhbHashValue string
			)
			// todo DHB config
			for _, vUserDepositDhbResult := range userDepositDhbResult[vDepositUsdtResult.From] { // 充值数额类型匹配
				if "100000000000000000000" == vDepositUsdtResult.Value && "100000000000000000000" == vUserDepositDhbResult.Value {

				} else if "200000000000000000000" == vDepositUsdtResult.Value && "2000000000000000000" == vUserDepositDhbResult.Value {

				} else {
					continue
				}

				tmpDhbHash = vUserDepositDhbResult.Hash
				tmpDhbHashValue = vUserDepositDhbResult.Value
			}

			notExistDepositResult = append(notExistDepositResult, &biz.EthUserRecord{ // 两种币的记录
				UserId:   depositUsers[vDepositUsdtResult.From].ID,
				Hash:     vDepositUsdtResult.Hash,
				Status:   "success",
				Type:     "deposit",
				Amount:   vDepositUsdtResult.Value,
				CoinType: "USDT",
			}, &biz.EthUserRecord{
				UserId:   depositUsers[vDepositUsdtResult.From].ID,
				Hash:     tmpDhbHash,
				Status:   "success",
				Type:     "deposit",
				Amount:   tmpDhbHashValue,
				CoinType: "DHB",
			})
		}

		_, err = a.ruc.EthUserRecordHandle(ctx, notExistDepositResult...)
		if nil != err {
			fmt.Println(err)
		}

		//time.Sleep(2 * time.Second)
	}

	return &v1.DepositReply{}, nil
}

type eth struct {
	Value       string
	Hash        string
	TokenSymbol string
	From        string
	To          string
}

func requestEthDepositResult(offset int64, page int64, contractAddress string) (map[string]*eth, error) {
	apiUrl := "https://api-testnet.bscscan.com/api"
	// URL param
	data := url.Values{}
	data.Set("module", "account")
	data.Set("action", "tokentx")
	data.Set("contractaddress", contractAddress)
	data.Set("address", "0xe865f2e5ff04b8b7952d1c0d9163a91f313b158f")
	data.Set("sort", "desc")
	data.Set("offset", strconv.FormatInt(offset, 10))
	data.Set("page", strconv.FormatInt(page, 10))

	u, err := url.ParseRequestURI(apiUrl)
	if err != nil {
		return nil, err
	}
	u.RawQuery = data.Encode() // URL encode
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	fmt.Println(u.String())

	resp, err := client.Get(u.String())
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var i struct {
		Message string `json:"message"`
		Result  []*eth `json:"Result"`
	}
	err = json.Unmarshal(b, &i)
	if err != nil {
		return nil, err
	}

	res := make(map[string]*eth, 0)
	for _, v := range i.Result {
		if "0xe865f2e5ff04b8b7952d1c0d9163a91f313b158f" == v.To { // 接收者
			res[v.Hash] = v
		}
	}

	return res, err
}
