I wrote a Unix pipeline analog, something like:
```
grep 127.0.0.1 | awk '{print $2}' | sort | uniq -c | sort -nr
```

When the STDOUT of one program is passed as STDIN to another program.

But in our case, these roles are played by pipes, which we pass from one function to another.

The hash calculation is implemented using the following chain:
* SingleHash calculates the value crc32(data)+"~"+crc32(md5(data)), where data is the input (essentially the numbers from the first function)
* MultiHash calculates the value crc32(th+data)), where data is the input (and the output of SingleHash)
* CombineResults receives all results, sorts them, and concatenates the sorted results into a single string using the _ (underscore) character.
* crc32 is calculated using the DataSignerCrc32 function
* md5 is calculated using DataSignerMd5

What's the problems:
* DataSignerMd5 can only be called once at a time, and takes 10 ms. If several are launched simultaneously, there will be a 1-second overheat.
* DataSignerCrc32, calculated as 1 second.
* Calculations for all tests take 3.433 seconds.

## How to run
1. Clone the repo: `git clone https://github.com/ugadanjemrtvima/async-pipeline`
2. Run tests: `go test -v -race .`