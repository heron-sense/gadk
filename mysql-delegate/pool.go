package mysql_delegate

import (
	"database/sql"
	fsc "github.com/heron-sense/gadk/flow-state-code"
)

type SqlConnPool struct {
	Schema string
}

func (p SqlConnPool) NewSession() (*HSession, fsc.FlowStateCode) {
	antiFraudSession, err := sql.Open("mysql", p.Schema)
	if err != nil {
		return nil, 3233435
	}
	return &HSession{db: antiFraudSession}, 0
}
