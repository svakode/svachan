# Installation

This assumes you already have a working Go environment, if not please see
[this page](https://golang.org/doc/install) first.

After you clone your project, you need to run
```
make copy-config
``` 

Open the generated `application.yml` and fill in the required configuration. The configuration is
not required for external parties (e.g Twitter & Google). If empty then we will not register the command to Svachan. 