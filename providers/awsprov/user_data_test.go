package awsprov_test

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/providers/awsprov"
)

const expected = `#!/bin/bash
set -x #echo on
##
## Mount volumes, and create file systems
##
TOKEN=$(curl -X PUT "http://169.254.169.254/latest/api/token" -H "X-aws-ec2-metadata-token-ttl-seconds: 21600")
OUTPUT=$(curl -H "X-aws-ec2-metadata-token: $TOKEN" http://169.254.169.254/latest/meta-data/instance-id)

aws ec2 attach-volume --volume-id abc123 --device /dev/sdf --instance-id $OUTPUT --region us-west-1

aws ec2 attach-volume --volume-id xyz123 --device /dev/sdg --instance-id $OUTPUT --region us-west-1


while [[ ! -b $(readlink -f /dev/sdf) ]]; do
    echo "waiting for the disk /dev/sdf to appear..">&2;
    sleep 5;
done
blkid $(readlink -f /dev/sdf) || mkfs -t ext4 $(readlink -f /dev/sdf)
mkdir -p /data/foo || true
mount /dev/sdf /data/foo
echo "Mounted /dev/sdf to /data/foo">&2;

while [[ ! -b $(readlink -f /dev/sdg) ]]; do
    echo "waiting for the disk /dev/sdg to appear..">&2;
    sleep 5;
done
blkid $(readlink -f /dev/sdg) || mkfs -t ext4 $(readlink -f /dev/sdg)
mkdir -p /data/bar || true
mount /dev/sdg /data/bar
echo "Mounted /dev/sdg to /data/bar">&2;

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

echo "ssh-key abc==" >> /home/ec2-user/.ssh/authorized_keys
echo "ssh-key abc==" >> /home/unweave/.ssh/authorized_keys

echo "ssh-key def==" >> /home/ec2-user/.ssh/authorized_keys
echo "ssh-key def==" >> /home/unweave/.ssh/authorized_keys


`

func TestUserData(t *testing.T) {
	t.Parallel()

	data, _ := awsprov.UserData("us-west-1", []string{"ssh-key abc==", "ssh-key def=="}, []types.ExecVolume{
		{
			VolumeID:  "abc123",
			MountPath: "/data/foo",
		},
		{
			VolumeID:  "xyz123",
			MountPath: "/data/bar",
		},
	})

	u, _ := base64.StdEncoding.DecodeString(data)

	t.Log(string(u))

	assert.Equal(t, expected, string(u))
}
