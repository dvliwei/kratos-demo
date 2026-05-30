/**
 * @Title
 * @Author: liwei
 * @Description:  TODO
 * @File:  gameapp
 * @Version: 1.0.0
 * @Date: 2026/05/29 23:04
 * @Update liwei 2026/5/29 23:04
 */

package data

import (
	"context"
	stderrors "errors"
	"gameapp-service/internal/biz"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type gameAppRepo struct {
	data *Data
	log  *log.Helper
}
type GameAppModel struct {
	ID                        uint64         `gorm:"column:id;primaryKey;autoIncrement"`
	GameID                    int16          `gorm:"column:game_id;not null;default:0;comment:游戏id"`
	AppID                     string         `gorm:"column:app_id;size:16;not null;default:'';uniqueIndex:tab_game_app_app_id_unique;comment:app id"`
	Name                      string         `gorm:"column:name;size:191;not null;default:'';comment:app name"`
	AppKey                    string         `gorm:"column:app_key;size:32;not null;default:'';comment:app 唯一 key"`
	Secret                    string         `gorm:"column:secret;size:16;not null;default:'';comment:sdk接口盐"`
	PayKey                    string         `gorm:"column:pay_key;size:32;not null;default:'';comment:支付通知key"`
	PayNotifyURL              string         `gorm:"column:pay_notify_url;size:191;not null;default:'';comment:支付通知接口"`
	DefaultHeadImg            string         `gorm:"column:default_head_img;size:191;not null;default:'';comment:默认的头像地址"`
	Comment                   *string        `gorm:"column:comment;type:text;comment:简介"`
	DeletedAt                 gorm.DeletedAt `gorm:"column:deleted_at;index"`
	CreatedAt                 *time.Time     `gorm:"column:created_at"`
	UpdatedAt                 *time.Time     `gorm:"column:updated_at"`
	AutiAddictionStatus       int8           `gorm:"column:auti_addiction_status;not null;default:0;index:tab_game_app_auti_addiction_status_index;comment:防沉迷开启状态"`
	GiftPackageExchangeStatus int8           `gorm:"column:gift_package_exchange_status;not null;default:0;index:tab_game_app_gift_package_exchange_status_index;comment:礼包cdky兑换开启状态 0未开启 1开启"`
	TypeOS                    int16          `gorm:"column:type_os;not null;default:0;index:tab_game_app_type_os_index;comment:平台类型"`
	JudgeURL                  *string        `gorm:"column:judge_url;size:500;comment:隐私协议地址"`
	AdjustToken               *string        `gorm:"column:adjust_token;size:500;default:'';comment:adjust配置项"`
	AdjustSecrets             *string        `gorm:"column:adjust_secrets;size:500;default:'';comment:adjust配置项2"`
	FirebaseJSONURL           *string        `gorm:"column:firebase_json_url;size:500;comment:firebase配置文件地址"`
	PackageName               *string        `gorm:"column:package_name;size:100;default:'';comment:firebase包名称"`
	PackageID                 *string        `gorm:"column:package_id;size:100;default:'';comment:firebase包id"`
	RegisterStatus            int16          `gorm:"column:register_status;not null;default:0;index:tab_game_app_register_status_index;comment:注册状态1为开启0为关闭"`
	RefundCountLock           int32          `gorm:"column:refund_count_lock;not null;default:0;comment:退款次数锁定阀值"`
	RefundCostLock            float64        `gorm:"column:refund_cost_lock;type:decimal(8,2);not null;default:0.00;comment:退款锁定金额阀值"`
	HioAppID                  *string        `gorm:"column:hio_app_id;size:191;comment:hio app id配置文件"`
	HioAPIToken               *string        `gorm:"column:hio_api_token;size:191;comment:hio api token配置文件"`
	HioSDKToken               *string        `gorm:"column:hio_sdk_token;size:191;comment:hio sdk token配置文件"`
	CommunityURL              *string        `gorm:"column:community_url;size:191;default:'';comment:社区地址"`
	AppICOURL                 *string        `gorm:"column:app_ico_url;size:191;default:'';comment:appico地址"`
	VerCodeLoginStatus        int16          `gorm:"column:ver_code_login_status;not null;default:0;index:tab_game_app_ver_code_login_status_index;comment:验证码登录开启或者关闭"`
	PropsAmountVerifyURL      *string        `gorm:"column:props_amount_verify_url;size:191;default:'';comment:商品id与金额查询地址"`
	RefundTotalCountLock      *int32         `gorm:"column:refund_total_count_lock;default:0;comment:累计退款次数锁定阀值"`
	RefundTotalCostLock       *float64       `gorm:"column:refund_total_cost_lock;type:decimal(8,2);default:0.00;comment:累计退款锁定金额阀值"`
	RefundNotifyURL           *string        `gorm:"column:refund_notify_url;size:191;default:'';comment:退款通知回调地址"`
	SendIsSandBox             int32          `gorm:"column:send_is_sand_box;not null;default:0;index:tab_game_app_send_is_sand_box_index;comment:发货通知同步发送订单类型"`
	AdjustPushGameURL         *string        `gorm:"column:adjust_push_game_url;size:500;default:'';comment:adjust数据同步游戏地址"`
	NetworkPublicKey          *string        `gorm:"column:network_public_key;type:text;comment:通讯密钥对公钥"`
	NetworkPrivateKey         *string        `gorm:"column:network_private_key;type:text;comment:通讯密钥对私钥"`
	JSPayNotifyURL            *string        `gorm:"column:js_pay_notify_url;type:text;comment:商城支付通知地址"`
	CCPPayNotifyURL           *string        `gorm:"column:ccp_pay_notify_url;type:text;comment:cpp支付通知地址"`
	GameOpenPayURL            *string        `gorm:"column:game_open_pay_url;type:text;comment:游戏下单请求地址"`
	NetworkConfig             *string        `gorm:"column:network_config;type:longtext;comment:网络配置"`
	PayStatus                 int16          `gorm:"column:pay_status;not null;default:1;index:tab_game_app_pay_status_index;comment:支付状态1为开启0为关闭"`
}

func (GameAppModel) TableName() string {
	return "tab_game_app"
}

/**
 * @Description:  创建游戏应用仓库
 * @param data 数据
 * @param logger 日志记录器
 * @return biz.GameAppRepo
 * */
func NewGameAppRepo(data *Data, logger log.Logger) biz.GameAppRepo {
	return &gameAppRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

/**
 * @Description:  根据ID查询游戏应用
 */
func (r *gameAppRepo) FindByID(ctx context.Context, id uint64) (*biz.GameApp, error) {
	var gameApp GameAppModel
	err := r.data.db.WithContext(ctx).
		Select("id", "game_id", "app_id", "name", "app_key").
		First(&gameApp, "id = ?", id).
		Error
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("GAME_APP_NOT_FOUND", "game app not found")
		}
		r.log.Errorf("query game app failed: %v", err)
		return nil, err
	}
	return &biz.GameApp{
		ID:     gameApp.ID,
		GameID: int64(gameApp.GameID),
		AppID:  gameApp.AppID,
		Name:   gameApp.Name,
		AppKey: gameApp.AppKey,
	}, nil
}

func (r *gameAppRepo) CountGameApps(ctx context.Context) (uint64, error) {
	var count int64
	err := r.data.db.WithContext(ctx).Model(&GameAppModel{}).Count(&count).Error
	if err != nil {
		r.log.Errorf("count game apps failed: %v", err)
		return 0, err
	}
	return uint64(count), nil
}
func (r *gameAppRepo) ListGameAppsWithPage(ctx context.Context, pageNum, pageSize int, search *biz.GameAppsSearch) ([]*biz.GameAppInfo, int64, error) {
	var gameApps []*GameAppModel
	var total int64
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	query := r.data.db.WithContext(ctx).Model(&GameAppModel{})
	if search != nil {
		if search.Name != "" {
			query = query.Where("name like ?", "%"+search.Name+"%")
		}
		if search.TypeOs != nil {
			query = query.Where("type_os = ?", *search.TypeOs)
		}
	}
	err := query.Count(&total).Error
	if err != nil {
		r.log.Errorf("count game apps failed: %v", err)
		return nil, 0, err
	}
	err = query.
		Order("id DESC").
		Offset((pageNum - 1) * pageSize).
		Limit(pageSize).
		Find(&gameApps).
		Error
	if err != nil {
		r.log.Errorf("list game apps failed: %v", err)
		return nil, 0, err
	}
	gameAppInfos := make([]*biz.GameAppInfo, 0, len(gameApps))
	for _, gameApp := range gameApps {
		gameAppInfos = append(gameAppInfos, &biz.GameAppInfo{
			ID:           gameApp.ID,
			GameID:       gameApp.GameID,
			AppID:        gameApp.AppID,
			Name:         gameApp.Name,
			AppKey:       gameApp.AppKey,
			TypeOS:       gameApp.TypeOS,
			CreatedAt:    gameApp.CreatedAt,
			UpdatedAt:    gameApp.UpdatedAt,
			PayStatus:    gameApp.PayStatus,
			AppICOURL:    gameApp.AppICOURL,
			PayKey:       gameApp.PayKey,
			PayNotifyURL: gameApp.PayNotifyURL,
		})
	}
	return gameAppInfos, total, nil
}
