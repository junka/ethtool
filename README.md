# ethtool

A native go implementation of ethtool. 

Command line parameters are compatible with original version 5.4

Take care that no data[0] structure in go


### Compile & Run

 ```make``` 

```
./dist/ethtool -i eth0
./dist/ethtool -a eth0
./dist/ethtool -A eth0 tx on
...
```
