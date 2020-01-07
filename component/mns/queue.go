package mns

import ali_mns "github.com/aliyun/aliyun-mns-go-sdk"


// Queue 队列
type Queue struct {
	Name       string              // 队列名称
	instance   ali_mns.AliMNSQueue // 队列实例
	handler    QueueHandler        // 处理方式
	WaitSecond int
}

// QueueHandler 队列消息的处理器
type QueueHandler func(params Params) error

// HasHandler 是否有处理方法
func (q *Queue) HasHandler() bool {
	return q.handler != nil
}

// SendMessage 发送消息
func (q *Queue) SendMessage(message string) (ali_mns.MessageSendResponse, error) {
	msg := ali_mns.MessageSendRequest{
		MessageBody:  message,
		DelaySeconds: 0,
		Priority:     8}

	resp, err := q.instance.SendMessage(msg)
	return resp, err
}

// DeleteMessage 删除消息
func (q *Queue) DeleteMessage(receiptHandler string) error {
	return q.instance.DeleteMessage(receiptHandler)
}
