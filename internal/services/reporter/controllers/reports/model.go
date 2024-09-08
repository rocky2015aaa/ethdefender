package reports

import "time"

type Resp struct {
	Data       interface{} `json:"data"`
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
}

type ContractReportResp struct {
	ContractName string      `json:"contract_name"`
	Report       interface{} `json:"report"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}
