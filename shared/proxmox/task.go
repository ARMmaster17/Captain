package proxmox

import (
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"log"
	"strings"
	"time"
)

const TaskStatusRetry = 1 * time.Second

type Task struct {
	Node string
	UPID string
}

func NewTask(node string, body string) (Task, error) {
	upid := gjson.Get(body, "data").String()
	if upid == "" || strings.HasPrefix(upid, "UPID:") {
		return Task{}, errors.New("the Proxmox API response did not return a task UPID")
	}
	return Task{
		Node: node,
		UPID: upid,
	}, nil
}

func (t *Task) IsRunning(p *Proxmox) (bool, error) {
	err := p.Authenticate()
	if err != nil {
		log.Println(err)
		return false, errors.New("unable to authenticate with Proxmox API")
	}
	body, err := p.Client.Get(fmt.Sprintf("nodes/%s/tasks/%s/status", t.Node, t.UPID))
	if err != nil {
		log.Println(err)
		return false, errors.New("unable to check on status of task " + t.UPID)
	}
	isRunning := gjson.Get(body, "data.status").String()
	if isRunning == "" {
		return false, errors.New("the Proxmox API did not return a task status")
	}
	return isRunning == "running", nil
}

func (t *Task) WaitForTaskCompletion(p *Proxmox) error {
	for {
		status, err := t.IsRunning(p)
		if err != nil {
			log.Println(err)
			return errors.New("unable to wait for task to complete")
		}
		if !status {
			return nil
		}
		time.Sleep(TaskStatusRetry)
	}
}
