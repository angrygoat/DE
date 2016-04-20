# publish

A pluggable REST API for publishing data in the CyVerse Discovery Environment.

## Usage

To compile a jar file for use, run `lein do clean, uberjar`. This will produce `target/iplant-groups-standalone.jar`. To
simply run a server for direct use, `lein run`.

iplant-groups listens on port 31311 rather than any configured port when run with lein-ring or a jar produced by
lein-ring, rather than the methods listed above.

## License

Copyright © 2016 iPlant Collaborative, Arizona Board of Regents

Distributed under the terms of the iPlant License:
http://iplantcollaborative.org/sites/default/files/iPLANT-LICENSE.txt
