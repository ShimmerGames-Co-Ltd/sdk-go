package leaderboard

// RegisterRequest / RegisterReply 注册排行榜
type (
	// RegisterRequest 注册排行榜请求参数
	RegisterRequest struct {
		// 游戏ID
		AppID string `json:"appID"`
		// 排行榜主键
		Id string `json:"id"`
		// 排行榜名字
		Name string `json:"name"`
		// 排行榜子主键（例如赛季编号）
		SubId int64 `json:"subId"`
		// 排行榜生效时间
		StartAt int64 `json:"startAt"`
		// 数据上报有效期（单位：秒，超过有效期后无法上报数据）
		ReportDuration int64 `json:"reportDuration"`
		// 排行榜生存有效期（单位：秒，超过有效期后自动删除,该时间应该大于reportDuration）
		AliveDuration int64 `json:"aliveDuration"`
		// 最大循环次数（默认为1，如果需要多次循环则需要记录就填写几次）
		LoopMax int32 `json:"loopMax"`
		// 排序规则(true=升序，false=降序)
		IsAsc bool `json:"isAsc"`
		// 上报分数(分数的变化值)的最小值
		ScoreMin int64 `json:"scoreMin"`
		// 上报分数(分数变化值)的最大值
		ScoreMax int64 `json:"scoreMax"`
		// 最大排名人数(0=不限制人数)
		MembersMax uint32 `json:"membersMax"`
		// 分页拉取排行榜时，一次最多拉取条数（默认值50，最大不超过200）
		PageSize uint32 `json:"pageSize"`
		// 每日分数上报（修改）次数限制(0=没有限制)
		ReportLimitADay uint32 `json:"reportLimitADay"`
	}

	// RegisterReply 注册排行榜的返回值
	RegisterReply struct {
		Detail AskReply `json:"detail"`
	}
)

// DeleteRequest / DeleteReply 删除排行榜
type (
	// DeleteRequest 删除排行榜请求参数
	DeleteRequest struct {
		// 游戏ID
		AppID string `json:"appID"`
		// 排行榜主键
		Id string `json:"id"`
	}

	// DeleteReply 删除排行榜返回值
	DeleteReply struct{}
)

// DeleteAppRequest / DeleteAppReply 删除APPID下所有排行榜
type (
	// DeleteAppRequest 删除APPID下所有排行榜请求参数
	DeleteAppRequest struct {
		// 游戏ID
		AppID string `json:"appID"`
		// 是否删除所有对玩家的封禁记录
		DeleteBan bool `json:"deleteBan"`
	}

	// DeleteAppReply 删除APPID下所有排行榜返回值
	DeleteAppReply struct{}
)

// AllOfAPPRequest / AllOfAPPReply 获取 APP 下所有排行榜
type (
	// AllOfAPPRequest 获取当前APP下的所有排行榜
	AllOfAPPRequest struct {
		// 游戏ID
		AppID string `json:"appID"`
	}

	// AllOfAPPReply 获取当前APP下的所有排行榜返回值
	AllOfAPPReply struct {
		List []AskReply `json:"list"`
	}
)

// AskRequest / AskReply 获取排行榜状态
type (
	// AskRequest 获取当前排行榜状态信息请求参数
	AskRequest struct {
		// 游戏ID
		AppID string `json:"appID"`
		// 排行榜主键
		Id string `json:"id"`
	}

	// AskReply 获取当前排行榜状态信息返回值
	AskReply struct {
		// 游戏ID
		AppID string `json:"appID"`
		// 排行榜主键
		Id string `json:"id"`
		// 排行榜名字
		Name string `json:"name"`
		// 排行榜子主键（例如赛季编号）
		SubId int64 `json:"subId"`
		// 排行榜开始时间
		StartAt int64 `json:"startAt"`
		// 排行榜生效倒计时
		StartCD int64 `json:"startCD"`
		// 分数上报周期（秒）
		ReportDuration int64 `json:"reportDuration"`
		// 分数上报有效期倒计时（单位：秒，超过有效期后无法上报数据）
		ReportCD int64 `json:"reportCD"`
		// 排行榜生存有效周期（秒）
		AliveDuration int64 `json:"aliveDuration"`
		// 排行榜生存有效期倒计时（单位：秒，超过有效期后自动删除）
		AliveCD int64 `json:"aliveCD"`
		// 自动循环次数上限
		LoopMax int32 `json:"loopMax"`
		// 当前自动循环次数
		LoopTimes int32 `json:"loopTimes"`
		// 排序规则(true=升序，false=降序)
		IsAsc bool `json:"isAsc"`
		// 上报分数的最小值
		ScoreMin int64 `json:"scoreMin"`
		// 上报分数的最大值
		ScoreMax int64 `json:"scoreMax"`
		// 最大排名人数
		MembersMax uint32 `json:"membersMax"`
		// 分页拉取排行榜时，一次最多拉取条数
		PageSize uint32 `json:"pageSize"`
		// 每天上传分数的次数限制
		ReportLimitADay uint32 `json:"reportLimitADay"`
	}
)

// SetScoreRequest / SetScoreReply 上报玩家绝对分
type (
	// SetScoreRequest 上报玩家数据请求参数
	SetScoreRequest struct {
		// 游戏ID
		AppID string `json:"appID"`
		// 排行榜主键
		Id string `json:"id"`
		// 排行榜子主键
		SubId int64 `json:"subId"`
		// 玩家唯一标识
		Uid string `json:"uid"`
		// 玩家分数
		Score int64 `json:"score"`
		// 透传信息
		Attachment string `json:"attachment"`
	}

	// SetScoreReply 玩家上报数据返回值
	SetScoreReply struct {
		// 当前积分和排名
		Member Member `json:"member"`
		// 剩余每日上报次数
		ReportLimit ReportLimit `json:"reportLimit"`
	}
)

// IncrScoreRequest / IncrScoreReply 上报增量分
type (
	// IncrScoreRequest 上报玩家数据请求参数
	IncrScoreRequest struct {
		// 游戏ID
		AppID string `json:"appID"`
		// 排行榜主键
		Id string `json:"id"`
		// 排行榜子主键
		SubId int64 `json:"subId"`
		// 玩家唯一标识
		Uid string `json:"uid"`
		// 玩家分数的变化值
		Incr int64 `json:"incr"`
		// 透传信息
		Attachment string `json:"attachment"`
	}

	// IncrScoreReply 玩家上报数据返回值
	IncrScoreReply struct {
		// 当前积分和排名
		Member Member `json:"member"`
		// 剩余每日上报次数
		ReportLimit ReportLimit `json:"reportLimit"`
	}
)

// GetScoreRequest / GetScoreReply 查询玩家分
type (
	// GetScoreRequest 获取玩家当前得分请求参数
	GetScoreRequest struct {
		// 游戏ID
		AppID string `json:"appID"`
		// 排行榜主键
		Id string `json:"id"`
		// 排行榜子主键
		SubId int64 `json:"subId"`
		// 玩家唯一标识
		Uid string `json:"uid"`
	}

	// GetScoreReply 获取玩家当前得分返回值
	GetScoreReply struct {
		// 玩家当前得分
		Member Member `json:"member"`
	}
)

// ListRequest / ListReply 分页拉取排行榜
type (
	// ListRequest 分页拉取排行榜数据请求参数
	ListRequest struct {
		// 游戏ID
		AppID string `json:"appID"`
		// 排行榜主键
		Id string `json:"id"`
		// 排行榜子主键
		SubId int64 `json:"subId"`
		// 偏移量
		Offset uint32 `json:"offset"`
		// 拉取条数
		Limit uint32 `json:"limit"`
		// 如果需要返回查询者的排名，请传入该值
		CallerUid string `json:"callerUid"`
	}

	// ListReply 分页拉取排行榜数据返回值
	ListReply struct {
		// 拉取到的成员
		Members []Member `json:"members"`
		// 查询者的排名
		Caller Member `json:"caller"`
	}
)

// BanAddRequest / BanAddReply 玩家封禁
type (
	// BanAddRequest 玩家封禁
	BanAddRequest struct {
		// 游戏ID
		AppID string `json:"appID"`
		// 玩家ID
		Uid string `json:"uid"`
		// 截止时间(过了这个时间ban自动失效,-1表示永久)
		StopAt int64 `json:"stopAt"`
	}

	// BanAddReply 玩家封禁返回值
	BanAddReply struct{}
)

// BanRemoveRequest / BanRemoveReply 解除封禁
type (
	// BanRemoveRequest 玩家封禁解除
	BanRemoveRequest struct {
		// 游戏ID
		AppID string `json:"appID"`
		// 玩家ID
		Uid string `json:"uid"`
	}

	// BanRemoveReply 玩家封禁解除返回值
	BanRemoveReply struct{}
)

// BanCleanRequest / BanCleanReply APP 下全部封禁解除
type (
	// BanCleanRequest APP下所有玩家封禁解除
	BanCleanRequest struct {
		// 游戏ID
		AppID string `json:"appID"`
	}

	// BanCleanReply APP下所有玩家封禁解除返回值
	BanCleanReply struct{}
)

// BanGetRequest / BanGetReply 拉取封禁详情
type (
	// BanGetRequest 拉取玩家封禁详情
	BanGetRequest struct {
		// 游戏ID
		AppID string `json:"appID"`
		// 玩家ID
		Uid string `json:"uid"`
	}

	// BanGetReply 拉取玩家封禁详情返回值
	BanGetReply struct {
		// 封禁记录详情
		Detail BanDetail `json:"detail"`
	}
)

// BanListRequest / BanListReply 拉取封禁列表
type (
	// BanListRequest 根据APPID拉取封禁列表
	BanListRequest struct {
		// 游戏ID
		AppID string `json:"appID"`
		// 偏移量
		Offset uint32 `json:"offset"`
		// 拉取条数
		Limit uint32 `json:"limit"`
	}

	// BanListReply 根据APPID拉取封禁列表返回值
	BanListReply struct {
		// 封禁记录列表
		List []BanDetail `json:"list"`
	}
)

// MemberInfoRequest / MemberInfoReply 上报次数查询
type (
	// MemberInfoRequest 获取成员信息请求参数
	MemberInfoRequest struct {
		// 游戏ID
		AppID string `json:"appID"`
		// 排行榜ID
		Id string `json:"id"`
		// 玩家ID
		Uid string `json:"uid"`
	}

	// MemberInfoReply 获取成员信息返回参数
	MemberInfoReply struct {
		// 上报次数查询
		ReportLimit ReportLimit `json:"reportLimit"`
	}
)

// ---------- 其他复合类型 ---------------------------------------------------------

// ReportLimit 每日上报次数限制
type ReportLimit struct {
	// 剩余上报次数
	Count uint32 `json:"count"`
	// 跨天倒计时
	TomorrowCD int64 `json:"tomorrowCD"`
}

// Member 排行榜成员
type Member struct {
	// 玩家ID
	Uid string `json:"uid"`
	// 玩家得分
	Score int64 `json:"score"`
	// 玩家排名
	Rank int64 `json:"rank"`
	// 透传信息
	Attachment string `json:"attachment"`
}

// BanDetail 玩家封禁详情
type BanDetail struct {
	// 玩家ID
	Uid string `json:"uid"`
	// 结束封禁时间（-1表示永久）
	StopAt int64 `json:"stopAt"`
	// 结束封禁倒计时
	StopCD int64 `json:"stopCD"`
}
