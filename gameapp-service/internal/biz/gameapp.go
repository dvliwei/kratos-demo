/**
 * @Title
 * @Author: liwei
 * @Description:  TODO
 * @File:  gameapp
 * @Version: 1.0.0
 * @Date: 2026/05/29 23:06
 * @Update liwei 2026/5/29 23:06
 */

package biz

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// uint64 id = 1;
// int64 game_id = 2;
// string app_id = 3;
// string name= 4;
// string app_key = 5;
type GameApp struct {
	ID     uint64 `json:"id"`
	GameID int64  `json:"game_id"`
	AppID  string `json:"app_id"`
	Name   string `json:"name"`
	AppKey string `json:"app_key"`
}

type GameAppInfo struct {
	ID                        uint64         `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	GameID                    int16          `gorm:"column:game_id;not null;default:0;comment:游戏id" json:"game_id"`
	AppID                     string         `gorm:"column:app_id;size:16;not null;default:'';uniqueIndex:tab_game_app_app_id_unique;comment:app id" json:"app_id"`
	Name                      string         `gorm:"column:name;size:191;not null;default:'';comment:app name" json:"name"`
	AppKey                    string         `gorm:"column:app_key;size:32;not null;default:'';comment:app 唯一 key" json:"app_key"`
	Secret                    string         `gorm:"column:secret;size:16;not null;default:'';comment:sdk接口盐" json:"secret"`
	PayKey                    string         `gorm:"column:pay_key;size:32;not null;default:'';comment:支付通知key" json:"pay_key"`
	PayNotifyURL              string         `gorm:"column:pay_notify_url;size:191;not null;default:'';comment:支付通知接口" json:"pay_notify_url"`
	DefaultHeadImg            string         `gorm:"column:default_head_img;size:191;not null;default:'';comment:默认的头像地址" json:"default_head_img"`
	Comment                   *string        `gorm:"column:comment;type:text;comment:简介" json:"comment"`
	DeletedAt                 gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
	CreatedAt                 *time.Time     `gorm:"column:created_at" json:"created_at"`
	UpdatedAt                 *time.Time     `gorm:"column:updated_at" json:"updated_at"`
	AutiAddictionStatus       int8           `gorm:"column:auti_addiction_status;not null;default:0;index:tab_game_app_auti_addiction_status_index;comment:防沉迷开启状态" json:"auti_addiction_status"`
	GiftPackageExchangeStatus int8           `gorm:"column:gift_package_exchange_status;not null;default:0;index:tab_game_app_gift_package_exchange_status_index;comment:礼包cdky兑换开启状态 0未开启 1开启" json:"gift_package_exchange_status"`
	TypeOS                    int16          `gorm:"column:type_os;not null;default:0;index:tab_game_app_type_os_index;comment:平台类型" json:"type"`
	JudgeURL                  *string        `gorm:"column:judge_url;size:500;comment:隐私协议地址" json:"judge_url"`
	AdjustToken               *string        `gorm:"column:adjust_token;size:500;default:'';comment:adjust配置项" json:"adjust_token"`
	AdjustSecrets             *string        `gorm:"column:adjust_secrets;size:500;default:'';comment:adjust配置项2" json:"adjust_secrets"`
	FirebaseJSONURL           *string        `gorm:"column:firebase_json_url;size:500;comment:firebase配置文件地址" json:"firebase_json_url"`
	PackageName               *string        `gorm:"column:package_name;size:100;default:'';comment:firebase包名称" json:"package_name"`
	PackageID                 *string        `gorm:"column:package_id;size:100;default:'';comment:firebase包id" json:"package_id"`
	RegisterStatus            int16          `gorm:"column:register_status;not null;default:0;index:tab_game_app_register_status_index;comment:注册状态1为开启0为关闭" json:"register_status"`
	RefundCountLock           int32          `gorm:"column:refund_count_lock;not null;default:0;comment:退款次数锁定阀值" json:"refund_count_lock"`
	RefundCostLock            float64        `gorm:"column:refund_cost_lock;type:decimal(8,2);not null;default:0.00;comment:退款锁定金额阀值" json:"refund_cost_lock"`
	HioAppID                  *string        `gorm:"column:hio_app_id;size:191;comment:hio app id配置文件" json:"hio_app_id"`
	HioAPIToken               *string        `gorm:"column:hio_api_token;size:191;comment:hio api token配置文件" json:"hio_api_token"`
	HioSDKToken               *string        `gorm:"column:hio_sdk_token;size:191;comment:hio sdk token配置文件" json:"hio_sdk_token"`
	CommunityURL              *string        `gorm:"column:community_url;size:191;default:'';comment:社区地址" json:"community_url"`
	AppICOURL                 *string        `gorm:"column:app_ico_url;size:191;default:'';comment:appico地址" json:"app_ico_url"`
	VerCodeLoginStatus        int16          `gorm:"column:ver_code_login_status;not null;default:0;index:tab_game_app_ver_code_login_status_index;comment:验证码登录开启或者关闭" json:"ver_code_login_status"`
	PropsAmountVerifyURL      *string        `gorm:"column:props_amount_verify_url;size:191;default:'';comment:商品id与金额查询地址" json:"props_amount_verify_url"`
	RefundTotalCountLock      *int32         `gorm:"column:refund_total_count_lock;default:0;comment:累计退款次数锁定阀值" json:"refund_total_count_lock"`
	RefundTotalCostLock       *float64       `gorm:"column:refund_total_cost_lock;type:decimal(8,2);default:0.00;comment:累计退款锁定金额阀值" json:"refund_total_cost_lock"`
	RefundNotifyURL           *string        `gorm:"column:refund_notify_url;size:191;default:'';comment:退款通知回调地址" json:"refund_notify_url"`
	SendIsSandBox             int32          `gorm:"column:send_is_sand_box;not null;default:0;index:tab_game_app_send_is_sand_box_index;comment:发货通知同步发送订单类型" json:"send_is_sand_box"`
	AdjustPushGameURL         *string        `gorm:"column:adjust_push_game_url;size:500;default:'';comment:adjust数据同步游戏地址" json:"adjust_push_game_url"`
	NetworkPublicKey          *string        `gorm:"column:network_public_key;type:text;comment:通讯密钥对公钥" json:"network_public_key"`
	NetworkPrivateKey         *string        `gorm:"column:network_private_key;type:text;comment:通讯密钥对私钥" json:"network_private_key"`
	JSPayNotifyURL            *string        `gorm:"column:js_pay_notify_url;type:text;comment:商城支付通知地址" json:"js_pay_notify_url"`
	CCPPayNotifyURL           *string        `gorm:"column:ccp_pay_notify_url;type:text;comment:cpp支付通知地址" json:"ccp_pay_notify_url"`
	GameOpenPayURL            *string        `gorm:"column:game_open_pay_url;type:text;comment:游戏下单请求地址" json:"game_open_pay_url"`
	NetworkConfig             *string        `gorm:"column:network_config;type:longtext;comment:网络配置" json:"network_config"`
	PayStatus                 int16          `gorm:"column:pay_status;not null;default:1;index:tab_game_app_pay_status_index;comment:支付状态1为开启0为关闭" json:"pay_status"`
}

type GameAppRepo interface {
	//id查询
	FindByID(ctx context.Context, id uint64) (*GameApp, error)
	//分页查询
	ListGameAppsWithPage(ctx context.Context, pageNum, pageSize int, search *GameAppsSearch) ([]*GameAppInfo, int64, error)
	//统计游戏APP数量
	CountGameApps(ctx context.Context) (uint64, error)
}
type GameAppUseCase struct {
	repo GameAppRepo
}

func NewGameAppUseCase(repo GameAppRepo) *GameAppUseCase {
	return &GameAppUseCase{repo: repo}
}
func (u *GameAppUseCase) GetByIDApp(ctx context.Context, id uint64) (*GameApp, error) {
	return u.repo.FindByID(ctx, id)
}

func (u *GameAppUseCase) CountGameApps(ctx context.Context) (uint64, error) {
	return u.repo.CountGameApps(ctx)
}

type GameAppsSearch struct {
	Name   string `json:"name"`    // 应用名称
	TypeOs *int32 `json:"type_os"` // 系统类型，nil 表示未传，0 表示明确查询 0
}

func (u *GameAppUseCase) ListGameAppsWithPage(ctx context.Context, pageNum, pageSize int, search *GameAppsSearch) ([]*GameAppInfo, int64, error) {
	return u.repo.ListGameAppsWithPage(ctx, pageNum, pageSize, search)
}
