# Install procedures

For the purposes of this fork of DE, we standardized on CentOS 7 VMs, though the number required will be up to the organization. Some services (jex, condor log monitor, condor) need to run on the same machine but most may be safely run on their own node or on shared hardware. Note that if you intend to run all internal services on one node, you'll want to throw a fair amount of resources at that box. Elasticsearch in particular may require a larger RAM allocation, and Condor will need the CPU, RAM, and disk space (/var) to accommodate the needs of your users.

### Ansible Setup
* Create a privileged Ansible user on all boxes, select a head node, generate ssh key, distribute public keys, grant sudo to Ansible user.
* on ansible head node, create a de directory under the ansible home dir (this is just suggested).  Under this, 
create an ansible-vars dir.  In here you will place your inventory file, and your Ansible group vars.  You can use the Annotated
group vars and annotated inventory (TODO) as a start
* under the de directory, check out this git repo, which will result in a DE directory
* under the de directory, create a directory for holding key files that will be distributed to other nodes, such as key
pairs.  It is suggested to call this localdata.

Doing this results in

/home/ansibleuser/de
    - DE (the git repo)
    - ansible-vars
    - localdata

* install ansible on each machine
``` yum install ansible ```

### Prerequisite Playbooks
* install CentOS library prereqs: **$ ansible-playbook -i inventory -e @group_vars -s -K playbooks/prereqs.yaml**
* configure iptables: **$ ansible-playbook -i inventory -e @group_vars -s -K iptables.yaml**

### Docker and Local Docker Repo

Configure (need to make this optional) a secure private docker repo, install docker, configure if need be for proxies

* create the local copy of the ssl keypair on the ansible head node

```

mkdir /home/ansible/de/ansible-vars/localdata/dockercerts

openssl genrsa -out server.key

openssl req -new -x509 -key server.key -out server.crt -days 365

```
* update the docker.registry group vars, per the AnnotatedGroupVar example

* run the combined docker playbook: **$ ansible-playbook -i inventory -e @group_vars -s -K docker.yaml**


* install the data container by building and pushing to the private docker repo

```
docker images
  936  docker tag 12239e411d9f xxx.renci.org:443/de-data
  937  docker push xxx.renci.org:443/de-data

```

### Misc Prereqs

* install openJDK7 on services VM: **$ ansible-playbook -i inventory -e @group_vars -s -K playbooks/java7.yaml**
* install timezone packages: **$ ansible-playbook -i inventory -e @group_vars -s -K playbooks/timezone.yaml**

### CAS/LDAP

Note group vars for cas: and ldap:.  The ldap playbook is for using the canned openLDAP dedicated to DE, this may vary if you go against your own existing LDAP.  The cas role can be used to install a cas server, configured to look at LDAP.  This may be sufficiant for simple cases, or serve as a template
for your particular install environment.  We orient cas here for our 'reference implementation'.  Consult the AnnotatedGroupVars for details of the cas and ldap settings!

* Generate or provision the SSL public/private key pair in a secure location on the ansible head node.  This location is configured by the cas.ssl_cert_file and cas.ssl_key_file group vars.  This source location on the ansible head node is used to copy to a target directory on the inventoried cas box.

For our purposes, this can be under the de/localdata directory described above, for this one, we can call it cascerts.  Be sure
the directory has proper visibility so that others cannot discover the private key!

example statements for self-signing.  First generate a private SSL key

```
openssl genrsa -out server.key

```
Now generate a public key

```
openssl req -new -x509 -key server.key -out server.crt -days 365

```

The private and public key are entered in the group vars for the cas, and will be copied and configured into tomcat for CAS.  Once the keys are in place, run:

**$ ansible-playbook -i inventory -e @group_vars -s -K playbooks/cas-ldap.yaml** to install the reference implementation
LDAP.  Beware it sets up testde1, 2, and 3 users and some test groups by default.  You can also run the cas.yaml playbook to just install cas 
pointed to the configured LDAP.  Note that you may need to update the cas overlay, depending on your particular settings,
in which case you can alter the git url from which the cas overlay will be downloaded.

### Grouper

Grouper is used to maintain authorization for the Apps framework.  It relies on a Grouper database, and needs a Grouper
user entry in LDAP for the Grouper DE user.  As users are added in DE they are automatically added to a default Grouper group.
DE expects a prefabricated tree for Apps, and this is preconfigured using the Sharkbait tool, in a fashion similar to how
the Facepalm tool is used to initialize the Postgres databases.

* generate SSL certs for nginx on the Grouper box and configure in the grouper: stanza of the group vars

```
cd /home/ansibleuser/de/ansible-vars/localdata

mkdir groupercerts

cd groupercerts

```
Now gen the keys in grouper ssl

```
openssl genrsa -out server.key

```
Now generate a public key

```
openssl req -new -x509 -key server.key -out server.crt -days 365

```

* Run Grouper install playbooks: **$ ansible-playbook -i inventory -e @group_vars -s -K grouper-all.yaml**

* Run the Sharkbait tool to initialize Grouper structures for DE

This installs and configures Grouper.  Once Grouper is running, it needs to be initialized with some basic DE group structures,
using the DE tool 'sharkbait', which can be found in tools/sharkbait.  The standalong jar file can be built via clojore, or
downloaded from the DE releases page, where available.  

This utility expects the Grouper configuration files to be present in `/etc/grouper` on the local host before it can be
executed. It is generally best to run this utility on the Grouper host itself. This utility also expects an account by
the name of `de_grouper` to exist in the LDAP directory that is used as Grouper's subject source. (NB the reference
implementation cas-ldap playbook will add this de_grouper user to the LDAP directory).

wget the sharkbait jar, or build it from the DE source, and then execute the command with the following options. Note
the database coordinates map to the de database setting in the AnnotatedGroupVars


```

java -jar /path/to/sharkbait-standalone.jar  OPTIONS

  -h, --host HOST                localhost  The database hostname.
  -p, --port PORT                5432       The database port number.
  -d, --database DATABASE        de         The database name.
  -U, --user USER                de         The database username.
  -v, --version                             Show the sharkbait version.
  -e, --environment ENVIRONMENT  dev        The name of the DE environment.


```


A sample run looks as follows:


```

root@dfc-de-grouper ~]# java -jar sharkbait-standalone.jar -h host.example.org -d de -U de -e de
WARNING: update already refers to: #'clojure.core/update in namespace: sharkbait.db, being replaced by: #'honeysql.helpers/update
WARNING: read already refers to: #'clojure.core/read in namespace: sharkbait.permissions, being replaced by: #'sharkbait.permissions/read
de password: 
Grouper starting up: version: 2.2.1, build date: null, env: <no label configured>
grouper.properties read from: /root/file:/root/sharkbait-standalone.jar!/grouper.properties
Grouper current directory is: /root
log4j.properties read from:   /root/file:/root/sharkbait-standalone.jar!/log4j.properties
Grouper is logging to file:   console, at min level WARN for package: edu.internet2.middleware.grouper, based on log4j.properties
grouper.hibernate.properties: /root/file:/root/sharkbait-standalone.jar!/grouper.hibernate.properties
grouper.hibernate.properties: null@null
sources.xml read from:         [cant find sources.xml]
sources.xml ldap source id:   dfc: cn=Manager,dc=xxx,dc=example,dc=org@ldap://ldap.example.org:389
sources.xml groupersource id: g:gsa
sources.xml groupersource id: grouperEntities
2016-08-08 14:32:18,326: [main] WARN  GrouperStartup.printConfigOnce(167) -  - Grouper starting up: version: 2.2.1, build date: null, env: <no label configured>
Grouper warning: Not checking configuration integrity due to grouper.properties: configuration.detect.errors
2016-08-08 14:32:20,338: [main] WARN  GrouperCheckConfig.configCheckDisabled(296) -  - Not checking configuration integrity due to grouper.properties: configuration.detect.errors
Loading DE subjects...
"Elapsed time: 46.252702 msecs"
Registering DE users...
"Elapsed time: 374.047179 msecs"
Creating DE permission definitions...
Creating permission def:app-permission-def...
Creating permission def:analysis-permission-def...
"Elapsed time: 388.3188 msecs"
Registering DE apps...
"Elapsed time: 204.134627 msecs"



```



### Deploy Discovery Environment
* pull the trigger: **$ ansible-playbook -i inventory -e @group_vars -s -K deploy-all.yaml**

### Other Dependencies
* install Condor: **$ ansible-playbook -i inventory -e @group_vars -s -K playbooks/condor.yaml**

### Setup databases (if doing this manually, otherwise, skip ahead to the data container)

* install postgres: **$ ansible-playbook -i inventory -e @group_vars -s -K playbooks/postgres.yaml**
* create DBs: **$ ansible-playbook -i inventory -e @group_vars -s -K playbooks/db-creator.yaml**

Note that Discovery Environment currently uses PostgresQL-9.4, which the above role should install. Try your luck with facepalm, but connectivity errors forced us to build it and run it at the command line. Locations for current database tarballs:

```
db_targz: https://everdene.iplantcollaborative.org/jenkins/job/databases-dev/lastSuccessfulBuild/artifact/databases/de-database-schema/database.tar.gz
jexdb_targz: https://everdene.iplantcollaborative.org/jenkins/job/databases-dev/lastSuccessfulBuild/artifact/databases/jex-db/jex-db.tar.gz
metadata_db_targz: https://everdene.iplantcollaborative.org/jenkins/job/databases-dev/lastSuccessfulBuild/artifact/databases/metadata/metadata-db.tar.gz
notification_db_targz: https://everdene.iplantcollaborative.org/jenkins/job/databases-dev/lastSuccessfulBuild/artifact/databases/notification-db/notification-db.tar.gz
```
See also https://everdene.iplantcollaborative.org/jenkins/job/databases-dev/

###Facepalm commands to run:

Template for reference:
```
java -jar /tmp/facepalm-standalone.jar -m {{fmode}} {{db_version_flag|default('')}} -h {{curr_db_host}} -p {{curr_db_port}} -d {{curr_db_name}} -U {{curr_db_user}} -A {{curr_db_admin}} {{tgz_flag}}
```

de db
```
java -jar target/facepalm-standalone.jar -m  init -h dfc-test-vmlab3.edc.renci.org -p 5432 -d de -U de -A postgres -f database.tar.gz
```

jex db
```
java -jar target/facepalm-standalone.jar -m  init -h dfc-test-vmlab3.edc.renci.org -p 5432 -d jex -U jex_user -A postgres -f jex-db.tar.gz
```

metadata db
```
java -jar target/facepalm-standalone.jar -m  init -h dfc-test-vmlab3.edc.renci.org -p 5432 -d metadata -U metadata_db -A postgres -f metadata-db.tar.gz
```

notification db
```
java -jar target/facepalm-standalone.jar -m  init -h dfc-test-vmlab3.edc.renci.org -p 5432 -d notifications -U notification_user -A postgres -f notification-db.tar.gz
```

Note that the facepalm-standalone jar can be found as a file in the de releases, and is in a zip with the schema files that go with each release point.

For use on machines where proxies are involved, there is now a --proxy-host and --proxy-port flag, as facepalm will try and load some jars from clojars.  For example:

```
java -jar facepalm-standalone.jar -m  init -h localhost -p 5432 -d de -U de -A postgres -f database.tar.gz --proxy-host gateway.example.org --proxy-port 8080

```


## Data Container

Data container work is in a private repo, which we need to formalize (mc), in the meantime, we build our own data container docker image and serve it out of a private docker registry, per this [nice setup instructions](https://www.digitalocean.com/community/tutorials/how-to-set-up-a-private-docker-registry-on-ubuntu-14-04)


## Build configs

Edit the generate_configs.sh to point to the correct location for variables and inventory files.  This will generate a set of .properties files for each service, typically in the /etc/iplant/de directory like this:

```

[ansible@dfc-test-vmlab1 de]$ ls
anon-files.properties          data-info.properties    info-typer.properties     jex.properties       notificationagent.properties  user-preferences.properties
apps.properties                dewey.properties        iplant-email.properties   kifshare.properties  saved-searches.properties     user-sessions.properties
clockwork.properties           exim-sender.properties  iplant-groups.properties  metadata.properties  terrain.properties
condor-log-monitor.properties  infosquito.properties   jexevents.properties      monkey.properties    tree-urls.properties


```
## Generate Configs  

Run the ansible playbook generate-configs.yaml to create the consolidated set of properties and configuration that is to be used in the next step to create the config docker images for each service.  

```
ansible-playbook -i /home/ansible/ansible-vars/inventory -e @/home/ansible/ansible-vars/group_vars.yaml -s -K -vvvv playbooks/generate-configs.yaml

```

Then run local services config

```

 ansible-playbook -i /home/ansible/ansible-vars/inventory -e @/home/ansible/ansible-vars/group_vars.yaml -s  -vvvv playbooks/local-services-cfg.yml --extra-vars="group=ansible owner=ansible"


```

This generates additional configuration information under /etc (e.g. docker-gc/)

## Generate config containers and push to a docker repo

Once the configs are in place, run the build_config_images.sh to copy the configs into a docker image and push it to a private docker repo.

To set up a private repo consult this link: https://docs.docker.com/registry/deploying/

The script may need to be configured for your particular Docker repo
