### (Loose) Instructions for Modifying Discovery Environment to Point to an External Subject Source

###On LDAP machine: 
point /etc/openldap/ldap.conf at external LDAP

###On Grouper machine:
point /etc/grouper/client.properties at external LDAP

point /etc/grouper/grouper.client.properties at external LDAP

point /etc/grouper/sources.xml SubjectResolve at external LDAP

###Then re-initialize Grouper:
$ bin/gsh -registry -init
