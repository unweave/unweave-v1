package lambdalabs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const ApiUrl = "https://cloud.lambdalabs.com/api/v1"

type InstanceType string

const (
	gpu1xA100 InstanceType = "gpu_1x_a100"
)

type Region string

const (
	RegionUSTX1 Region = "us-tx-1"
)

type Instance struct {
	ID              string       `json:"id"`
	Name            string       `json:"name"`
	IP              string       `json:"ip,omitempty"`
	Status          string       `json:"status"`
	SSHKeyNames     []string     `json:"ssh_key_names"`
	FileSystemNames []string     `json:"file_system_names"`
	Region          Region       `json:"region"`
	InstanceType    InstanceType `json:"instance_type"`
	Hostname        string       `json:"hostname,omitempty"`
	JupyterToken    string       `json:"jupyter_token,omitempty"`
	JupyterURL      string       `json:"jupyter_url,omitempty"`
}

type LaunchInstanceRequest struct {
	RegionName      string   `json:"region_name"`
	InstanceType    string   `json:"instance_type_name"`
	SSHKeyNames     []string `json:"ssh_key_names"`
	FileSystemNames []string `json:"file_system_names"`
	Quantity        int      `json:"quantity"`
	Name            string   `json:"name"`
}

type LaunchInstanceResponse struct {
	Data struct {
		InstanceIDs []string `json:"instance_ids"`
	} `json:"data"`
}

func LaunchInstance(req LaunchInstanceRequest) (*LaunchInstanceResponse, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	r, err := http.NewRequest("POST", ApiUrl+"/instance-operations/launch", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", "application/json")

	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var launchRes LaunchInstanceResponse
	if err := json.Unmarshal(resBody, &launchRes); err != nil {
		return nil, err
	}

	return &launchRes, nil
}

func GetInstance(id string) (*Instance, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf(ApiUrl+"/instances/%s", id), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var instance Instance
	if err = json.Unmarshal(resBody, &instance); err != nil {
		return nil, err
	}
	return &instance, nil
}

type AddSSHKeyRequest struct {
	SSHKey
}

type AddSSHKeyResponse struct {
	Data struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		PublicKey  string `json:"public_key"`
		PrivateKey string `json:"private_key"`
	} `json:"data"`
}

func AddSSHKey(req AddSSHKeyRequest) (*AddSSHKeyResponse, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	r, err := http.NewRequest("POST", ApiUrl+"/ssh-keys", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", "application/json")

	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var askRes AddSSHKeyResponse
	if err := json.Unmarshal(resBody, &askRes); err != nil {
		return nil, err
	}

	return &askRes, nil
}

type TerminateInstanceRequest struct {
	Instances []Instance `json:"instances"`
}

type TerminateInstanceResponse struct {
	Data struct {
		TerminatedInstances []Instance `json:"terminated_instances"`
	} `json:"data"`
}

func TerminateInstance(req TerminateInstanceRequest) (*TerminateInstanceResponse, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	r, err := http.NewRequest("POST", ApiUrl+"/instance-operations/terminate", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", "application/json")

	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var terminateRes TerminateInstanceResponse
	if err := json.Unmarshal(resBody, &terminateRes); err != nil {
		return nil, err
	}

	return &terminateRes, nil
}
