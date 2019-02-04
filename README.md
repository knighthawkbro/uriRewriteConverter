# How to use uriRewriteConverter
## Source code to my small little project.
uriRewriteConverter takes in a flag and a file and converts it from ht.acl
rewrite format to Microsoft web.config XML or vice versa.

### Usage:
```
$ ./uriRewriteConverter [-v] {-a|-x} <FileName>
  -a    Converts a HT.ACL file to Microsoft Web.Config XML
  -v    Sends output to STDOUT instead of a file
  -x    Converts a Web.Config XML back to a HT.ACL
```
