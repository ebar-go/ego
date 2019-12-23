package mns

import (
	"encoding/json"
	"github.com/aliyun/aliyun-mns-go-sdk"
	"github.com/ebar-go/ego/helper"
	"github.com/ebar-go/ego/log"
)

// Queue 队列接口
type Queue interface {
	// 发送消息
	SendMessage(message string) (ali_mns.MessageSendResponse, error)

	// 接收消息
	ReceiveMessage()

	// 删除消息
	DeleteMessage(receiptHandler string) error

	// 是否设置处理方法
	HasHandler() bool
}

// Queue 队列
type queue struct {
	Name       string              // 队列名称
	instance   ali_mns.AliMNSQueue // 队列实例
	handler    QueueHandler        // 处理方式
	WaitSecond int
}

// QueueHandler 队列消息的处理器
type QueueHandler func(params Params) error

// HasHandler 是否有处理方法
func (q *queue) HasHandler() bool {
	return q.handler != nil
}

// SendMessage 发送消息
func (q *queue) SendMessage(message string) (ali_mns.MessageSendResponse, error) {
	msg := ali_mns.MessageSendRequest{
		MessageBody:  message,
		DelaySeconds: 0,
		Priority:     8}

	resp, err := q.instance.SendMessage(msg)
	return resp, err
}

// ReceiveMessage 接收消息并处理
func (q *queue) ReceiveMessage() {

	if q.WaitSecond == 0 {
		q.WaitSecond = 30
	}
	endChan := make(chan int)
	respChan := make(chan ali_mns.MessageReceiveResponse)
	errChan := make(chan error)
	go func() {
		select {
		case resp := <-respChan:
			{
				var params Params
				log.System().Info("mqParams", log.Context{
					"receiveTime": helper.GetTimeStr(),
					"queue_name":  q.Name,
					"resp":        resp,
					"message_id":  resp.MessageId,
					"body":        resp.MessageBody,
				})

				// 解析消息
				if err := json.Unmarshal([]byte(helper.DecodeBase64Str(resp.MessageBody)), &params); err != nil {
					log.System().Error("invalidMessageBody", log.Context{
						"err":   err.Error(),
						"trace": helper.Trace(),
					})
				} else {

					log.Mq().Info("receiveMessage", log.Context{
						"receiveTime": helper.GetTimeStr(),
						"queue_name":  q.Name,
						"messageBody": params.Content,
						"tag":         params.Tag,
						"trace_id":    params.TraceId,
					})

					if err := q.handler(params); err != nil {
						log.System().Warn("processMessageFailed", log.Context{
							"err":   err.Error(),
							"trace": helper.Trace(),
						})

					} else {
						// 处理成功，删除消息
						err := q.DeleteMessage(resp.ReceiptHandle)
						log.Mq().Info("deleteMessage", log.Context{
							"receiveTime": helper.GetTimeStr(),
							"queue_name":  q.Name,
							"messageBody": params.Content,
							"tag":         params.Tag,
							"trace_id":    params.TraceId,
							"err":         err,
						})

						endChan <- 1
					}
				}

			}
		case err := <-errChan:
			{
				log.System().Info("receiveMessageFailed", log.Context{
					"err":   err.Error(),
					"trace": helper.Trace(),
				})
				endChan <- 1
			}
		}
	}()

	// 通过chan去接收数据
	q.instance.ReceiveMessage(respChan, errChan, int64(q.WaitSecond))
	<-endChan
}

// DeleteMessage 删除消息
func (q *queue) DeleteMessage(receiptHandler string) error {
	return q.instance.DeleteMessage(receiptHandler)
}
