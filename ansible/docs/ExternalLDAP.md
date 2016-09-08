### (Loose) Instructions for Modifying Discovery Environment to Point to an External Subject Source

###On LDAP machine: 
point /etc/openldap/ldap.conf at external LDAP

###On Grouper machine:
point /etc/grouper/client.properties at external LDAP

point /etc/grouper/grouper.client.properties at external LDAP

point /etc/grouper/sources.xml "Entity Subject Resolver" section at external LDAP
point /etc/grouper/sources.xml base

###On CAS machine:
point /etc/cas/cas.properties ldap.url to external LDAP

###Re-initialize Grouper:
$ bin/gsh -registry -init

###Re-run sharkbait on DE database:
$ java -jar sharkbait-standalone.jar -h <dbhost> -d <dbname> -U <dbuser> -e de

###Each user must exist in iRODS
[irods@visr-de-irods1 ~]$ ./add_irods_user.sh
