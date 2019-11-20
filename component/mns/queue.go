package mns

import (
	"github.com/aliyun/aliyun-mns-go-sdk"
	"github.com/ebar-go/ego/log"
	"github.com/ebar-go/ego/library"
	"encoding/json"
)

// Queue 队列结构体
type Queue struct {
	Name string // 队列名称
	instance ali_mns.AliMNSQueue // 队列实例
	Handler QueueHandler // 处理方式
	WaitSecond int
}

// QueueHandler 队列消息的处理器
type QueueHandler func(params Params) error

// SetHandler 设置队列消息处理器
func (queue *Queue) SetHandler(handler QueueHandler)  {
	queue.Handler = handler
}

// SendMessage 发送消息
func (queue *Queue) SendMessage(message string) (ali_mns.MessageSendResponse, error) {
	msg := ali_mns.MessageSendRequest{
		MessageBody:  message,
		DelaySeconds: 0,
		Priority:     8}

	resp, err := queue.instance.SendMessage(msg)
	return resp, err
}

// ReceiveMessage 接收消息并处理
func (queue *Queue) ReceiveMessage(waitSeconds int64) {

	if waitSeconds == 0 {
		waitSeconds = 30
	}
	endChan := make(chan int)
	respChan := make(chan ali_mns.MessageReceiveResponse)
	errChan := make(chan error)
	go func() {
		select {
		case resp := <-respChan:
			{
				var params Params

				// 解析消息
				if err := json.Unmarshal([]byte(library.DecodeBase64Str(resp.MessageBody)), &params); err != nil {
					library.Debug("消息结构异常:",queue.Name, err.Error(), resp.MessageBody )
				}else {

					log.Mq().Info("receiveMessage", log.Context{
						"receiveTime" : "receiveMessage",
						"queue_name" : queue.Name,
						"messageBody" : params.Content,
						"tag" : params.Tag,
						"trace_id" : params.TraceId,
					})

					if err := queue.Handler(params); err != nil {
						library.Debug("处理消息失败:",queue.Name, err.Error() )

						// TODO ChangeMessageVisibility
					}else {
						// 处理成功，删除消息
						err := queue.DeleteMessage(resp.ReceiptHandle)
						log.Mq().Info("deleteMessage", log.Context{
							"receiveTime" : "deleteMessage",
							"queue_name" : queue.Name,
							"messageBody" : params.Content,
							"tag" : params.Tag,
							"trace_id" : params.TraceId,
							"err" : err,
						})

						endChan <- 1
					}
				}

			}
		case err := <-errChan:
			{
				library.Debug(err)
				endChan <- 1
			}
		}
	}()

	// 通过chan去接收数据
	queue.instance.ReceiveMessage(respChan, errChan, waitSeconds)
	<-endChan
}

// DeleteMessage 删除消息
func (queue *Queue) DeleteMessage(receiptHandler string ) error{
	return queue.instance.DeleteMessage(receiptHandler)
}
