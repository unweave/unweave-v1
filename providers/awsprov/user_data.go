package awsprov

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"text/template"

	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/types"
)

type volume struct {
	DeviceName string
	VolumeID   string
	MountPath  string
}

type userDataInput struct {
	DeviceName string
	Region     string
	PubKeys    []string
	Volumes    []volume
}

const userDataTemplate = `#!/bin/bash
set -x #echo on
##
## Mount volumes, and create file systems
##
TOKEN=$(curl -X PUT "http://169.254.169.254/latest/api/token" -H "X-aws-ec2-metadata-token-ttl-seconds: 21600")
OUTPUT=$(curl -H "X-aws-ec2-metadata-token: $TOKEN" http://169.254.169.254/latest/meta-data/instance-id)
{{range .Volumes}}
aws ec2 attach-volume --volume-id {{.VolumeID}} --device {{.DeviceName}} --instance-id $OUTPUT --region {{$.Region}}
{{end}}
{{range .Volumes}}
while [[ ! -b $(readlink -f {{.DeviceName}}) ]]; do
    echo "waiting for the disk {{.DeviceName}} to appear..">&2;
    sleep 5;
done
blkid $(readlink -f {{.DeviceName}}) || mkfs -t ext4 $(readlink -f {{.DeviceName}})
mkdir -p {{.MountPath}} || true
mount {{.DeviceName}} {{.MountPath}}
echo "Mounted {{.DeviceName}} to {{.MountPath}}">&2;
{{end}}
##
## Add unweave user, make no password sudoer
##
sudo adduser unweave --home-dir /home/unweave --create-home --shell /bin/bash
echo "unweave ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers
sudo su - unweave
sudo mkdir /logs || true
sudo chown -R unweave:unweave /logs
mkdir /home/unweave/.ssh || true
chmod 700 /home/unweave/.ssh
touch /home/unweave/.ssh/authorized_keys
chmod 600 /home/unweave/.ssh/authorized_keys
sudo chown -R unweave:unweave /home/unweave
##
## Add public keys for unweave and ec2-user
##
{{ range .PubKeys }}
echo "{{.}}" >> /home/ec2-user/.ssh/authorized_keys
echo "{{.}}" >> /home/unweave/.ssh/authorized_keys
{{end}}

`

var (
	tmpl     = template.Must(template.New("user-data").Parse(userDataTemplate))
	alphabet = []rune("fghijklmnop")
)

func UserData(region string, pubKeys []string, volumes []types.ExecVolume) (string, error) {
	userData := &bytes.Buffer{}
	base64Enc := base64.NewEncoder(base64.StdEncoding, userData)

	userDataVolumes := make([]volume, len(volumes))
	for i, vol := range volumes {
		userDataVolumes[i] = volume{
			VolumeID:   vol.VolumeID,
			MountPath:  vol.MountPath,
			DeviceName: fmt.Sprintf("/dev/sd%c", alphabet[i]),
		}
	}

	input := userDataInput{
		Region:  region,
		PubKeys: pubKeys,
		Volumes: userDataVolumes,
	}

	if err := tmpl.Execute(base64Enc, input); err != nil {
		return "", fmt.Errorf("template userdata: %w", err)
	}

	log.Debug().Str("user_data", userData.String()).Msg("base64 user data script")

	return userData.String(), nil
}
