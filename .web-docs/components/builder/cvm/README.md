Type: `tencentcloud-cvm`
Artifact BuilderId: `tencent.cloud`

The `tencentcloud-cvm` Packer builder plugin provide the capability to build
customized images based on an existing base images.

## Configuration Reference

The following configuration options are available for building Tencentcloud images.
In addition to the options listed here,
a [communicator](/packer/docs/templates/legacy_json_templates/communicator) can be configured for this builder.

### Required:

- `secret_id` (string) - Tencentcloud secret id. You should set it directly,
  or set the `TENCENTCLOUD_SECRET_ID` environment variable.

- `secret_key` (string) - Tencentcloud secret key. You should set it directly,
  or set the `TENCENTCLOUD_SECRET_KEY` environment variable.

- `region` (string) - The region where your cvm will be launch. You should
  reference [Region and Zone](https://intl.cloud.tencent.com/document/product/213/6091)
  for parameter taking.

- `zone` (string) - The zone where your cvm will be launch. You should
  reference [Region and Zone](https://intl.cloud.tencent.com/document/product/213/6091)
  for parameter taking.

- `instance_type` (string) - The instance type your cvm will be launched by.
  You should reference [Instance Type](https://intl.cloud.tencent.com/document/product/213/11518)
  for parameter taking.

- `source_image_id` (string) - The base image id of Image you want to create
  your customized image from.

- `image_name` (string) - The name you want to create your customize image,
  it should be composed of no more than 60 characters, of letters, numbers
  or minus sign.

### Optional:

- `force_poweroff` (boolean) - Indicates whether to perform a forced shutdown to
  create an image when soft shutdown fails. Default value is `false`.

- `image_description` (string) - Image description. It should no more than 60 characters.

- `reboot` (boolean, **deprecated**) - Whether shutdown cvm to create Image.
  Please refer to parameter `force_poweroff`.

- `sysprep` (boolean) - Whether enable Sysprep during creating windows image.

- `image_copy_regions` (array of strings) - Regions that will be copied to after
  your image created.

- `image_share_accounts` (array of strings) - Accounts that will be shared to
  after your image created.

- `skip_region_validation` (boolean) - Do not check region and zone when validate.

- `associate_public_ip_address` (boolean) - Whether allocate public ip to your cvm.
  Default value is `false`.

  If not set, you could access your cvm from the same vpc.

- `internet_max_bandwidth_out` (number) - Max bandwidth out your cvm will be launched by(in MB).
  values can be set between 1 ~ 100.

- `instance_name` (string) - Instance name.

- `disk_type` (string) - Root disk type your cvm will be launched by, default is `CLOUD_PREMIUM`. you could
  reference [Disk Type](https://intl.cloud.tencent.com/document/product/213/15753#SystemDisk)
  for parameter taking.

- `disk_size` (number) - Root disk size your cvm will be launched by. values range(in GB):

  - LOCAL_BASIC: 50
  - Other: 50 ~ 1000 (need whitelist if > 50)

- `data_disks` (array of data disks) - Add one or more data disks to the instance before creating the
  image. Note that if the source image has data disk snapshots, this argument will be ignored, and
  the running instance will use source image data disk settings, in such case, `disk_type`
  argument will be used as disk type for all data disks, and each data disk size will use the
  origin value in source image.
  The data disks allow for the following argument:

  - `disk_type` - Type of the data disk. Valid choices: `CLOUD_BASIC`, `CLOUD_PREMIUM` and `CLOUD_SSD`.
  - `disk_size` - Size of the data disk.
  - `disk_snapshot_id` - Id of the snapshot for a data disk.

- `vpc_id` (string) - Specify vpc your cvm will be launched by.

- `vpc_name` (string) - Specify vpc name you will create. if `vpc_id` is not set, Packer will
  create a vpc for you named this parameter.

- `cidr_block` (boolean) - Specify cider block of the vpc you will create if `vpc_id` is not set.

- `subnet_id` (string) - Specify subnet your cvm will be launched by.

- `subnet_name` (string) - Specify subnet name you will create. if `subnet_id` is not set, Packer will
  create a subnet for you named this parameter.

- `subnect_cidr_block` (boolean) - Specify cider block of the subnet you will create if
  `subnet_id` is not set.

- `security_group_id` (string) - Specify security group your cvm will be launched by.

- `security_group_name` (string) - Specify security name you will create if `security_group_id` is not set.

- `user_data` (string) - userdata.

- `user_data_file` (string) - userdata file.

- `host_name` (string) - host name.

- `run_tags` (map of strings) - Tags to apply to the instance that is _launched_ to create the image.
  These tags are _not_ applied to the resulting image.

- `cvm_endpoint` (string) - The endpoint you want to reach the cloud endpoint,
  if tce cloud you should set a tce cvm endpoint.

- `vpc_endpoint` (string) - The endpoint you want to reach the cloud endpoint,
  if tce cloud you should set a tce vpc endpoint.

### Communicator Configuration

In addition to the above options, a communicator can be configured
for this builder.

#### Optional:

<!-- Code generated from the comments of the Config struct in communicator/config.go; DO NOT EDIT MANUALLY -->

- `communicator` (string) - Packer currently supports three kinds of communicators:
  
  -   `none` - No communicator will be used. If this is set, most
      provisioners also can't be used.
  
  -   `ssh` - An SSH connection will be established to the machine. This
      is usually the default.
  
  -   `winrm` - A WinRM connection will be established.
  
  In addition to the above, some builders have custom communicators they
  can use. For example, the Docker builder has a "docker" communicator
  that uses `docker exec` and `docker cp` to execute scripts and copy
  files.

- `pause_before_connecting` (duration string | ex: "1h5m2s") - We recommend that you enable SSH or WinRM as the very last step in your
  guest's bootstrap script, but sometimes you may have a race condition
  where you need Packer to wait before attempting to connect to your
  guest.
  
  If you end up in this situation, you can use the template option
  `pause_before_connecting`. By default, there is no pause. For example if
  you set `pause_before_connecting` to `10m` Packer will check whether it
  can connect, as normal. But once a connection attempt is successful, it
  will disconnect and then wait 10 minutes before connecting to the guest
  and beginning provisioning.

<!-- End of code generated from the comments of the Config struct in communicator/config.go; -->


<!-- Code generated from the comments of the SSH struct in communicator/config.go; DO NOT EDIT MANUALLY -->

- `ssh_host` (string) - The address to SSH to. This usually is automatically configured by the
  builder.

- `ssh_port` (int) - The port to connect to SSH. This defaults to `22`.

- `ssh_username` (string) - The username to connect to SSH with. Required if using SSH.

- `ssh_password` (string) - A plaintext password to use to authenticate with SSH.

- `ssh_ciphers` ([]string) - This overrides the value of ciphers supported by default by Golang.
  The default value is [
    "aes128-gcm@openssh.com",
    "chacha20-poly1305@openssh.com",
    "aes128-ctr", "aes192-ctr", "aes256-ctr",
  ]
  
  Valid options for ciphers include:
  "aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com",
  "chacha20-poly1305@openssh.com",
  "arcfour256", "arcfour128", "arcfour", "aes128-cbc", "3des-cbc",

- `ssh_clear_authorized_keys` (bool) - If true, Packer will attempt to remove its temporary key from
  `~/.ssh/authorized_keys` and `/root/.ssh/authorized_keys`. This is a
  mostly cosmetic option, since Packer will delete the temporary private
  key from the host system regardless of whether this is set to true
  (unless the user has set the `-debug` flag). Defaults to "false";
  currently only works on guests with `sed` installed.

- `ssh_key_exchange_algorithms` ([]string) - If set, Packer will override the value of key exchange (kex) algorithms
  supported by default by Golang. Acceptable values include:
  "curve25519-sha256@libssh.org", "ecdh-sha2-nistp256",
  "ecdh-sha2-nistp384", "ecdh-sha2-nistp521",
  "diffie-hellman-group14-sha1", and "diffie-hellman-group1-sha1".

- `ssh_certificate_file` (string) - Path to user certificate used to authenticate with SSH.
  The `~` can be used in path and will be expanded to the
  home directory of current user.

- `ssh_pty` (bool) - If `true`, a PTY will be requested for the SSH connection. This defaults
  to `false`.

- `ssh_timeout` (duration string | ex: "1h5m2s") - The time to wait for SSH to become available. Packer uses this to
  determine when the machine has booted so this is usually quite long.
  Example value: `10m`.
  This defaults to `5m`, unless `ssh_handshake_attempts` is set.

- `ssh_disable_agent_forwarding` (bool) - If true, SSH agent forwarding will be disabled. Defaults to `false`.

- `ssh_handshake_attempts` (int) - The number of handshakes to attempt with SSH once it can connect.
  This defaults to `10`, unless a `ssh_timeout` is set.

- `ssh_bastion_host` (string) - A bastion host to use for the actual SSH connection.

- `ssh_bastion_port` (int) - The port of the bastion host. Defaults to `22`.

- `ssh_bastion_agent_auth` (bool) - If `true`, the local SSH agent will be used to authenticate with the
  bastion host. Defaults to `false`.

- `ssh_bastion_username` (string) - The username to connect to the bastion host.

- `ssh_bastion_password` (string) - The password to use to authenticate with the bastion host.

- `ssh_bastion_interactive` (bool) - If `true`, the keyboard-interactive used to authenticate with bastion host.

- `ssh_bastion_private_key_file` (string) - Path to a PEM encoded private key file to use to authenticate with the
  bastion host. The `~` can be used in path and will be expanded to the
  home directory of current user.

- `ssh_bastion_certificate_file` (string) - Path to user certificate used to authenticate with bastion host.
  The `~` can be used in path and will be expanded to the
  home directory of current user.

- `ssh_file_transfer_method` (string) - `scp` or `sftp` - How to transfer files, Secure copy (default) or SSH
  File Transfer Protocol.
  
  **NOTE**: Guests using Windows with Win32-OpenSSH v9.1.0.0p1-Beta, scp
  (the default protocol for copying data) returns a a non-zero error code since the MOTW
  cannot be set, which cause any file transfer to fail. As a workaround you can override the transfer protocol
  with SFTP instead `ssh_file_transfer_protocol = "sftp"`.

- `ssh_proxy_host` (string) - A SOCKS proxy host to use for SSH connection

- `ssh_proxy_port` (int) - A port of the SOCKS proxy. Defaults to `1080`.

- `ssh_proxy_username` (string) - The optional username to authenticate with the proxy server.

- `ssh_proxy_password` (string) - The optional password to use to authenticate with the proxy server.

- `ssh_keep_alive_interval` (duration string | ex: "1h5m2s") - How often to send "keep alive" messages to the server. Set to a negative
  value (`-1s`) to disable. Example value: `10s`. Defaults to `5s`.

- `ssh_read_write_timeout` (duration string | ex: "1h5m2s") - The amount of time to wait for a remote command to end. This might be
  useful if, for example, packer hangs on a connection after a reboot.
  Example: `5m`. Disabled by default.

- `ssh_remote_tunnels` ([]string) - 

- `ssh_local_tunnels` ([]string) - 

<!-- End of code generated from the comments of the SSH struct in communicator/config.go; -->


<!-- Code generated from the comments of the SSHTemporaryKeyPair struct in communicator/config.go; DO NOT EDIT MANUALLY -->

- `temporary_key_pair_type` (string) - `dsa` | `ecdsa` | `ed25519` | `rsa` ( the default )
  
  Specifies the type of key to create. The possible values are 'dsa',
  'ecdsa', 'ed25519', or 'rsa'.
  
  NOTE: DSA is deprecated and no longer recognized as secure, please
  consider other alternatives like RSA or ED25519.

- `temporary_key_pair_bits` (int) - Specifies the number of bits in the key to create. For RSA keys, the
  minimum size is 1024 bits and the default is 4096 bits. Generally, 3072
  bits is considered sufficient. DSA keys must be exactly 1024 bits as
  specified by FIPS 186-2. For ECDSA keys, bits determines the key length
  by selecting from one of three elliptic curve sizes: 256, 384 or 521
  bits. Attempting to use bit lengths other than these three values for
  ECDSA keys will fail. Ed25519 keys have a fixed length and bits will be
  ignored.
  
  NOTE: DSA is deprecated and no longer recognized as secure as specified
  by FIPS 186-5, please consider other alternatives like RSA or ED25519.

<!-- End of code generated from the comments of the SSHTemporaryKeyPair struct in communicator/config.go; -->


- `ssh_keypair_name` (string) - If specified, this is the key that will be used for SSH with the
  machine. The key must match a key pair name loaded up into the remote.
  By default, this is blank, and Packer will generate a temporary keypair
  unless [`ssh_password`](#ssh_password) is used.
  [`ssh_private_key_file`](#ssh_private_key_file) or
  [`ssh_agent_auth`](#ssh_agent_auth) must be specified when
  [`ssh_keypair_name`](#ssh_keypair_name) is utilized.


- `ssh_private_key_file` (string) - Path to a PEM encoded private key file to use to authenticate with SSH.
  The `~` can be used in path and will be expanded to the home directory
  of current user.


- `ssh_agent_auth` (bool) - If true, the local SSH agent will be used to authenticate connections to
  the source instance. No temporary keypair will be created, and the
  values of [`ssh_password`](#ssh_password) and
  [`ssh_private_key_file`](#ssh_private_key_file) will be ignored. The
  environment variable `SSH_AUTH_SOCK` must be set for this option to work
  properly.


## Basic Example

Here is a basic example for Tencentcloud.

```json
{
  "variables": {
    "secret_id": "{{env `TENCENTCLOUD_SECRET_ID`}}",
    "secret_key": "{{env `TENCENTCLOUD_SECRET_KEY`}}"
  },
  "builders": [
    {
      "type": "tencentcloud-cvm",
      "secret_id": "{{user `secret_id`}}",
      "secret_key": "{{user `secret_key`}}",
      "region": "ap-guangzhou",
      "zone": "ap-guangzhou-4",
      "instance_type": "S4.SMALL1",
      "source_image_id": "img-oikl1tzv",
      "ssh_username": "root",
      "image_name": "PackerTest",
      "disk_type": "CLOUD_PREMIUM",
      "packer_debug": true,
      "associate_public_ip_address": true,
      "run_tags": {
        "good": "luck"
      }
    }
  ],
  "provisioners": [
    {
      "type": "shell",
      "inline": ["sleep 30", "yum install redis.x86_64 -y"]
    }
  ]
}
```

See the
[examples/tencentcloud](https://github.com/hashicorp/packer-plugin-tencentcloud/tree/master/builder/tencentcloud/examples)
folder in the Packer project for more examples.
