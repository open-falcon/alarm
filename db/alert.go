package db

import (
	"fmt"
	cmodel "github.com/open-falcon/common/model"
	"github.com/open-falcon/common/utils"
	"github.com/open-falcon/alarm/api"
	"log"
)

func AddAlert(event *cmodel.Event, action *api.Action) {
	sql := fmt.Sprintf("insert into alerts(event_id, endpoint, counter, max_step, current_step, priority, expression_id, strategy_id, content, note, status, team, event_time) values ('%s', '%s', '%s', %d, %d, %d, %d, %d, '%s', '%s', '%s', '%s', '%s')",
		event.Id,
		event.Endpoint,
		event.Counter(),
		event.MaxStep(),
		event.CurrentStep,
		event.Priority(),
		event.ExpressionId(),
		event.StrategyId(),
		getEventContent(event),
		event.Note(),
		event.Status,
		action.Uic,
		utils.UnixTsFormat(event.EventTime))

	_, err := DB.Exec(sql)
	if err != nil {
		log.Println("exec", sql, "failed", err)
	}
}

func UpdateAlert(event *cmodel.Event, action *api.Action) {
	sql := ""
	if event.Status == "OK" {
		sql = fmt.Sprintf("update alerts set status = 'OK', recovery_time = NOW() where event_id='%s' and status = 'PROBLEM'", event.Id)
		_, err := DB.Exec(sql)
		if err != nil {
			log.Println("exec", sql, "failed", err)
		}
	} else {
		sql := fmt.Sprintf("insert into alerts(event_id, endpoint, counter, max_step, current_step, priority, expression_id, strategy_id, content, note, status, team, event_time) values ('%s', '%s', '%s', %d, %d, %d, %d, %d, '%s', '%s', '%s', '%s', '%s')",
			event.Id,
			event.Endpoint,
			event.Counter(),
			event.MaxStep(),
			event.CurrentStep,
			event.Priority(),
			event.ExpressionId(),
			event.StrategyId(),
			getEventContent(event),
			event.Note(),
			event.Status,
			action.Uic,
			utils.UnixTsFormat(event.EventTime))
		_, err := DB.Exec(sql)
		if err != nil {
			log.Println("exec", sql, "failed", err)
		}
	}
}

func MarkSolvedAlert(event_id string) {
	sql := fmt.Sprintf("update alerts set status = 'OK', recovery_time = NOW() where event_id='%s' and status = 'PROBLEM'", event_id)
	_, err := DB.Exec(sql)
	if err != nil {
		log.Println("exec", sql, "failed", err)
	}
}
