/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

dispatcher.go

Base dispatcher functions
*/
package dispatchers

import (
	"errors"
	"sandboxer/pkg/config"
	"sandboxer/pkg/sandbox"
	"sandboxer/pkg/task"

	"github.com/mpkondrashin/vone"
)

type Dispatcher interface {
	InboundChannel() task.Channel
	ProcessTask(tsk *task.Task) error
}

type BaseDispatcher struct {
	conf     *config.Configuration
	channels *task.Channels
	list     *task.TaskList
}

func NewBaseDispatcher(conf *config.Configuration, channels *task.Channels, list *task.TaskList) BaseDispatcher {
	return BaseDispatcher{conf, channels, list}
}

func (d *BaseDispatcher) Channel(ch task.Channel) task.IDChannel {
	return d.channels.TaskChannel[ch]
}

func (d *BaseDispatcher) Sandbox() (sandbox.Sandbox, error) {
	token := d.conf.VisionOne.Token
	if token == "" {
		return nil, errors.New("token is not set")
	}
	domain := d.conf.VisionOne.Domain
	if domain == "" {
		return nil, errors.New("domain is not set")
	}
	return sandbox.NewVOneSandbox(vone.NewVOne(domain, token)), nil
}
