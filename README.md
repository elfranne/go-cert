This utility is designed to check the validity of certificates. It will check: 
- expiration dates (both expired and not yet valid, also for CA)
- certificate and key are matching
- CA is matching

if any error it will exit 1.

Background:
Managing certifcates is painful and when dealing with semi automtic process it can lead to human errors. This tool can be added to your automation software as a pre check before changing certificates. For example I am using Puppet and this allows me the use the [validate_cmd](https://www.puppet.com/docs/puppet/latest/types/file.html#file-attribute-validate_cmd)

```puppet

file { '/etc/ssl/example.pem':
  content      => 'insert cert here',
  validate_cmd => '/usr/local/sbin/go-ssl --cert %',
}
```

While I am not not familiar with Ansible you should be able to use the [validate function](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/template_module.html#parameter-validate)