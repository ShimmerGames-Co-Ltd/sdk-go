package mail

import (
	"encoding/json"

	"github.com/ShimmerGames-Co-Ltd/sdk-go/v3/internal/jsoncompat"
)

// UnmarshalJSON 将传输格式兼容限制在解码边界，公开字段保持 []int64。
func (r *ReadReply) UnmarshalJSON(data []byte) error {
	type plain ReadReply
	value := plain(*r)
	wire := struct {
		*plain
		IDs jsoncompat.Int64List `json:"read_mail_ids"`
	}{plain: &value, IDs: jsoncompat.Int64List(r.ReadMailIDs)}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	value.ReadMailIDs = wire.IDs
	*r = ReadReply(value)
	return nil
}

// UnmarshalJSON 与 ReadReply 复用内部整数列表解析。
func (r *RemoveReply) UnmarshalJSON(data []byte) error {
	type plain RemoveReply
	value := plain(*r)
	wire := struct {
		*plain
		IDs jsoncompat.Int64List `json:"removed_mail_ids"`
	}{plain: &value, IDs: jsoncompat.Int64List(r.RemovedMailIDs)}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	value.RemovedMailIDs = wire.IDs
	*r = RemoveReply(value)
	return nil
}
