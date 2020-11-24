# goalarm

Goalarm run alarm server. It receive a command, and a response its status.

## Dependencies

https://github.com/hajimehoshi/oto#prerequisite


## install

```shell
$ go get github.com/komem3/goalarm/cmd/goalarm
```

## Usage

```shell
$ goalarm -h
Usage of goalarm:
  -describe
    	Describe command or status.
  -file string
    	Path of sound file.
  -hour int
    	Wait hour.
  -loop
    	Loop Alarm.
  -min int
    	Wait minute.
  -routine string
    	Alarm routine. Format is json array. [{"range":20,"name":"working"},{"range":5,"name":"break"}]
  -sec int
    	Wait second.
  -time string
    	Call time.(15:00:01)
  -v	Ouput verbose.
```

## Author

komem3

## License

[MIT](./LICENSE)
