package mns

import (
	"github.com/aliyun/aliyun-mns-go-sdk"
	"os"
	"github.com/ebar-go/ego/http/constant"
	"github.com/ebar-go/ego/library"
	"encoding/base64"
	"github.com/ebar-go/ego/log"
	"encoding/json"
)

type Topic struct {
	Name string
	instance ali_mns.AliMNSTopic
}

// PublishMessage 发布消息
func (topic *Topic) PublishMessage(params Params, filterTag string) (*ali_mns.MessageSendResponse, error) {
	bytes , err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	if params.ReferServiceName == "" {
		params.ReferServiceName = os.Getenv(constant.EnvSystemName)
	}

	if params.TraceId == "" {
		params.TraceId = library.UniqueId()
	}

	if params.Sign == "" {
		params.Sign = params.GenerateSign(client.conf.AccessKeySecret)
	}

	request := ali_mns.MessagePublishRequest{
		MessageBody: base64.StdEncoding.EncodeToString(bytes),
		MessageTag: filterTag,
	}
	resp, err := topic.instance.PublishMessage(request)
	if err != nil {
		return nil, err
	}

	log.Mq().Info("publishMessage", log.Context{
		"action" : "publishMessage",
		"publish_time" : library.GetTimeStr(),
		"msectime" : library.GetTimeStampFloatStr(),
		"message_id" : resp.MessageId,
		"status_code" : resp.Code,
		"topic_name" : topic.Name,
		"message_tag" : params.Tag,
		"global_trace_id" : library.GetTraceId(),
		"trace_id": params.TraceId,
		"filter_tag" : filterTag,
		"sign" : params.Sign,
	})

	return &resp, nil
}