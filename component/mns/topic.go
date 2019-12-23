package mns

import (
	"encoding/base64"
	"encoding/json"
	"github.com/aliyun/aliyun-mns-go-sdk"
	"github.com/ebar-go/ego/app"
	"github.com/ebar-go/ego/component/log"
	"github.com/ebar-go/ego/helper"
)

// Topic
type Topic interface {
	// 发布消息
	PublishMessage(params Params, filterTag string) (*ali_mns.MessageSendResponse, error)
}

// topic
type topic struct {
	Name     string
	instance ali_mns.AliMNSTopic
}

// PublishMessage 发布消息
func (t *topic) PublishMessage(params Params, filterTag string) (*ali_mns.MessageSendResponse, error) {
	bytes, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	request := ali_mns.MessagePublishRequest{
		MessageBody: base64.StdEncoding.EncodeToString(bytes),
		MessageTag:  filterTag,
	}
	resp, err := t.instance.PublishMessage(request)
	if err != nil {
		return nil, err
	}

	app.LogManager().Mq().Info("publishMessage", log.Context{
		"action":          "publishMessage",
		"publish_time":    helper.GetTimeStr(),
		"msectime":        helper.GetTimeStampFloatStr(),
		"message_id":      resp.MessageId,
		"status_code":     resp.Code,
		"topic_name":      t.Name,
		"message_tag":     params.Tag,
		"global_trace_id": helper.NewTraceId(),
		"trace_id":        params.TraceId,
		"filter_tag":      filterTag,
		"sign":            params.Sign,
	})

	return &resp, nil
}
