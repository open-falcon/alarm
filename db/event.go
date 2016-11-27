package db

import (
	"fmt"
	cmodel "github.com/open-falcon/common/model"
	"github.com/open-falcon/common/utils"
	"github.com/open-falcon/alarm/api"
	"github.com/open-falcon/alarm/g"
	"log"
	"time"
)

func AddEvent(event *cmodel.Event, action *api.Action) {
	sql := fmt.Sprintf("insert into events(event_id, endpoint, counter, max_step, current_step, priority, expression_id, strategy_id, content, note, status, team, event_time) values ('%s', '%s', '%s', %d, %d, %d, %d, %d, '%s', '%s', '%s', '%s', '%s')",
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

func MarkSolvedEvent(event *g.EventDto) {
	sql := fmt.Sprintf("insert into events(event_id, endpoint, counter, max_step, current_step, priority, expression_id, strategy_id, content, note, status, team, event_time) values ('%s', '%s', '%s', %d, %d, %d, %d, %d, '%s', '%s', '%s', '%s', '%s')",
		event.Id,
		event.Endpoint,
		event.Counter,
		event.MaxStep,
		1,
		event.Priority,
		event.ExpressionId,
		event.StrategyId,
		"mark it solved",
		event.Note,
		"OK",
		"",
		utils.UnixTsFormat(time.Now().Unix()))

	_, err := DB.Exec(sql)
	if err != nil {
		log.Println("exec", sql, "failed", err)
	}
}

func getEventContent(event *cmodel.Event) string {
	priority := 0
	maxStep := 0
	if event.Strategy != nil {
		priority = event.Strategy.Priority
		maxStep = event.Strategy.MaxStep
	} else {
		priority = event.Expression.Priority
		maxStep = event.Expression.MaxStep
	}
	content := fmt.Sprintf("[P%d #%d/%d] %s %s %s%s%s",
		priority,
		event.CurrentStep,
		maxStep,
		event.Counter(),
		event.Func(),
		utils.ReadableFloat(event.LeftValue),
		event.Operator(),
		utils.ReadableFloat(event.RightValue()),
	)
	return content
}
