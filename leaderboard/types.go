package leaderboard

import "github.com/ShimmerGames-Co-Ltd/sdk-go/v3/core"

// JSON 键与 leaderboard proto 字段名一致（proto 本身为 camelCase，如 subId）。

type RegisterRequest struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	SubID           core.JSONInt64  `json:"subId,omitempty"`
	StartAt         core.JSONInt64  `json:"startAt"`
	ReportDuration  core.JSONInt64  `json:"reportDuration"`
	AliveDuration   core.JSONInt64  `json:"aliveDuration"`
	LoopMax         int32  `json:"loopMax,omitempty"`
	IsAsc           bool   `json:"isAsc,omitempty"`
	ScoreMin        core.JSONInt64  `json:"scoreMin,omitempty"`
	ScoreMax        core.JSONInt64  `json:"scoreMax,omitempty"`
	MembersMax      uint32 `json:"membersMax,omitempty"`
	PageSize        uint32 `json:"pageSize,omitempty"`
	ReportLimitADay uint32 `json:"reportLimitADay,omitempty"`
}

type AskReply struct {
	AppID           string `json:"appID"`
	ID              string `json:"id"`
	Name            string `json:"name"`
	SubID           core.JSONInt64  `json:"subId"`
	StartAt         core.JSONInt64  `json:"startAt"`
	StartCD         core.JSONInt64  `json:"startCD"`
	ReportDuration  core.JSONInt64  `json:"reportDuration"`
	ReportCD        core.JSONInt64  `json:"reportCD"`
	AliveDuration   core.JSONInt64  `json:"aliveDuration"`
	AliveCD         core.JSONInt64  `json:"aliveCD"`
	LoopMax         int32  `json:"loopMax"`
	LoopTimes       int32  `json:"loopTimes"`
	IsAsc           bool   `json:"isAsc"`
	ScoreMin        core.JSONInt64  `json:"scoreMin"`
	ScoreMax        core.JSONInt64  `json:"scoreMax"`
	MembersMax      uint32 `json:"membersMax"`
	PageSize        uint32 `json:"pageSize"`
	ReportLimitADay uint32 `json:"reportLimitADay"`
}

type RegisterReply struct {
	Detail AskReply `json:"detail"`
}

type DeleteRequest struct {
	ID string `json:"id"`
}

type DeleteAppRequest struct {
	DeleteBan bool `json:"deleteBan,omitempty"`
}

type AllOfAPPReply struct {
	List []AskReply `json:"list"`
}

type AskRequest struct {
	ID string `json:"id"`
}

type Member struct {
	UID        string `json:"uid"`
	Score      core.JSONInt64  `json:"score"`
	Rank       core.JSONInt64  `json:"rank"`
	Attachment string `json:"attachment"`
}

type ReportLimit struct {
	Count      uint32 `json:"count"`
	TomorrowCD core.JSONInt64  `json:"tomorrowCD"`
}

type SetScoreRequest struct {
	ID         string `json:"id"`
	SubID      core.JSONInt64  `json:"subId,omitempty"`
	UID        string `json:"uid"`
	Score      core.JSONInt64  `json:"score"`
	Attachment string `json:"attachment,omitempty"`
}

type SetScoreReply struct {
	Member      Member      `json:"member"`
	ReportLimit ReportLimit `json:"reportLimit"`
}

type IncrScoreRequest struct {
	ID         string `json:"id"`
	SubID      core.JSONInt64  `json:"subId,omitempty"`
	UID        string `json:"uid"`
	Incr       core.JSONInt64  `json:"incr"`
	Attachment string `json:"attachment,omitempty"`
}

type IncrScoreReply struct {
	Member      Member      `json:"member"`
	ReportLimit ReportLimit `json:"reportLimit"`
}

type GetScoreRequest struct {
	ID    string `json:"id"`
	SubID core.JSONInt64  `json:"subId,omitempty"`
	UID   string `json:"uid"`
}

type GetScoreReply struct {
	Member Member `json:"member"`
}

type ListRequest struct {
	ID        string `json:"id"`
	SubID     core.JSONInt64  `json:"subId,omitempty"`
	Offset    uint32 `json:"offset,omitempty"`
	Limit     uint32 `json:"limit,omitempty"`
	CallerUID string `json:"callerUid,omitempty"`
}

type ListReply struct {
	Members []Member `json:"members"`
	Caller  Member   `json:"caller"`
}

type BanAddRequest struct {
	UID    string `json:"uid"`
	StopAt core.JSONInt64  `json:"stopAt"`
}

type BanRemoveRequest struct {
	UID string `json:"uid"`
}

type BanGetRequest struct {
	UID string `json:"uid"`
}

type BanScoreSnapshot struct {
	LeaderboardID string `json:"leaderboardID"`
	SubID         core.JSONInt64  `json:"subId"`
	Score         core.JSONInt64  `json:"score"`
}

type BanDetail struct {
	UID            string             `json:"uid"`
	StopAt         core.JSONInt64              `json:"stopAt"`
	StopCD         core.JSONInt64              `json:"stopCD"`
	ScoreSnapshots []BanScoreSnapshot `json:"scoreSnapshots"`
}

type BanGetReply struct {
	Detail BanDetail `json:"detail"`
}

type BanListRequest struct {
	Offset uint32 `json:"offset,omitempty"`
	Limit  uint32 `json:"limit,omitempty"`
}

type BanListReply struct {
	List []BanDetail `json:"list"`
}

type MemberInfoRequest struct {
	ID  string `json:"id"`
	UID string `json:"uid"`
}

type MemberInfoReply struct {
	ReportLimit ReportLimit `json:"reportLimit"`
}

type HistoryMemberInfoRequest struct {
	ID    string `json:"id"`
	SubID core.JSONInt64  `json:"subId"`
	UID   string `json:"uid"`
}

type HistoryMemberInfoReply struct {
	Exists bool   `json:"exists"`
	Member Member `json:"member"`
}
