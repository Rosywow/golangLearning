package cmn

import (
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx/types"
)

type ReplyProto struct {
	//Status, 0: success, others: fault
	Status int `json:"status"`

	//Msg, Action result describe by literal
	Msg string `json:"msg,omitempty"`

	//Data, operand
	Data types.JSONText `json:"data,omitempty"`

	// RowCount, just row count
	RowCount int64 `json:"rowCount,omitempty"`

	//API, call target
	API string `json:"API,omitempty"`

	//Method, using http method
	Method string `json:"method,omitempty"`

	//SN, call order
	SN int `json:"SN,omitempty"`
}

type ReqProto struct {
	Action string `json:"action,omitempty"`

	Sets    []string            `json:"sets,omitempty"`
	OrderBy []map[string]string `json:"orderBy,omitempty"`

	//***页码从第零页开始***
	Page     int64 `json:"page,omitempty"`
	PageSize int64 `json:"pageSize,omitempty"`

	Data   json.RawMessage `json:"data,omitempty"`
	Filter interface{}     `json:"filter,omitempty"`

	AuthFilter interface{} `json:"authFilter,omitempty"`
}

//QNearTime for customize json.Unmarshal
type QNearTime struct {
	time.Time
}

//QNearTimeLayout customize time layout
const QNearTimeLayout = `2006-01-02T15:04:05.000`

//Unmarshal customize json Unmarshal
func (v *QNearTime) Unmarshal(d []byte) (err error) {
	if d == nil || len(d) == 0 {
		return
	}

	if d[0] == '"' && d[len(d)-1] == '"' {
		d = d[1 : len(d)-1]
	}

	v.Time, err = time.Parse(QNearTimeLayout, string(d))
	return
}

func (v *ReplyProto) MarshalJSON() ([]byte, error) {
	type Alias ReplyProto
	return json.Marshal(&struct {
		Data string `json:"data,omitempty"`
		*Alias
	}{
		Data:  "",
		Alias: (*Alias)(v),
	})
}
