package cron

import (
	"fmt"
	"github.com/open-falcon/alarm/g"
	"github.com/open-falcon/common/model"
	"github.com/open-falcon/common/utils"
)

func BuildCommonSMSContent(event *model.Event) string {
	return fmt.Sprintf(
		"[P%d][%s][%s][][%s %s %s %s %s%s%s][O%d %s]",
		event.Priority(),
		event.Status,
		event.Endpoint,
		event.Note(),
		event.Func(),
		event.Metric(),
		utils.SortedTags(event.PushedTags),
		utils.ReadableFloat(event.LeftValue),
		event.Operator(),
		utils.ReadableFloat(event.RightValue()),
		event.CurrentStep,
		event.FormattedTime(),
	)
}

func BuildCommonMailContent(event *model.Event) string {
	link := g.Link(event)
	tdtl := `style="border: 1px solid #ccc; background: #FFF4F4;"`
	tdtr := `style="border: 1px solid #ccc; border-left: none;"`
	tdl := `style="border: 1px solid #ccc; border-top: none; background: #FFF4F4;"`
	tdr := `style="border: 1px solid #ccc; border-top: none; border-left: none;"`
	return fmt.Sprintf(
		`<html><head><meta charset="utf-8"></head>
		<body>
			<table border="0" cellpadding="5" cellspacing="0">
                                <tr>
                                        <td %s >%s</td>
                                        <td %s >%d</td></tr>
                                <tr>
                                        <td %s>Endpoint:</td>
                                        <td %s>%s</td>
                                </tr>
                                <tr>
                                        <td %s>Metric:</td>
                                        <td %s>%s</td>
                                </tr>
                                <tr>
                                        <td %s>Tags:</td>
                                        <td %s>%s</td>
                                </tr>
                                <tr>
                                        <td %s>%s</td>
                                        <td %s>%s%s%s</td>
                                </tr>
                                <tr>
                                        <td %s>Note:</td>
                                        <td %s>%s</td>
                                </tr>
                                <tr>
                                        <td %s>Max:</td>
                                        <td %s>%d</td>
                                </tr>
                                <tr>
                                        <td %s>Current:</td>
                                        <td %s>%d</td>
                                </tr>
                                <tr>
                                        <td %s>Timesramp:</td>
                                        <td %s>%s</td>
                                </tr>
                        </table>
			<br>
			<a href="%s">%s</a>
		</body></html>`,

		tdtl, event.Status, tdtr, event.Priority(),
		tdl, tdr, event.Endpoint,
		tdl, tdr, event.Metric(),
		tdl, tdr, utils.SortedTags(event.PushedTags),
		tdl, event.Func(), tdr, utils.ReadableFloat(event.LeftValue), event.Operator(),	utils.ReadableFloat(event.RightValue()),
		tdl, tdr, event.Note(),
		tdl, tdr, event.MaxStep(),
		tdl, tdr, event.CurrentStep,
		tdl, tdr, event.FormattedTime(),
		link,
		link,
	)
}

func GenerateSmsContent(event *model.Event) string {
	return BuildCommonSMSContent(event)
}

func GenerateMailContent(event *model.Event) string {
	return BuildCommonMailContent(event)
}
