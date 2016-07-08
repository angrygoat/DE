## Proxies can be a pain 

If you are running on VMs you may run into proxy errors.  Here we document some of the proxy errors that may bite 
you at various points in time and how to deal with them.

### Git Proxy

If you run into issues accessing GIT from a machine, add the appropriate proxy like this:

```
git config --global http.proxy http://proxyuser:proxypwd@proxy.server.com:8080

```

### yum mirrors

At RENCI, if you run into yum errors contacting mirrors, check in the /etc/yum.repos.d entries to make sure internal
mirrors are not indicated for boxes that are outside of the vm firewall.  That is, a box with a public ip needs to not
refer to internal yum repos.




