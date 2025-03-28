## first, add binariers to req. binaries
you can either do:
ex: ```/bin/ls ```
or directories like ```/bin```
ex file:
```
../../iso/toRun
/bin
/usr/bin/
/bin/tree
```
all of these are valid, note: ```/bin``` and ```/usr/bin``` do the same thing

## Now you can run your command:
```  Usage: go run ./cmd <command> ```

ex:
```
sudo go run . ./pwd 
sudo go run . ./tree /exDir
```
or you can compile:
```sudo go build ```
and then
```./cmd ./cat hi.txt ```
