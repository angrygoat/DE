### (Loose) Instructions for Modifying Discovery Environment to Point to an External Subject Source

###On LDAP machine: 
point /etc/openldap/ldap.conf at external LDAP

###On Grouper machine:
point /etc/grouper/client.properties at external LDAP

point /etc/grouper/grouper.client.properties at external LDAP

point /etc/grouper/sources.xml SubjectResolve at external LDAP

###Re-initialize Grouper:
$ bin/gsh -registry -init

###Each user must exist in iRODS
[irods@visr-de-irods1 ~]$ ./add_irods_user.sh
