# goalarm

Goalarm run alarm server. It receive a command, and a response its status.

This was made to work with other tools.
By using this tool, you can use alarm on other tool without managing alarm logic.

## Dependencies

https://github.com/hajimehoshi/oto#prerequisite


## install

```shell
$ go get github.com/komem3/goalarm/cmd/goalarm
```

## The plugin using this

https://github.com/komem3/goalarm.el

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

### Examples

#### 5 min timer
```shell
$ goalarm -file ./bell.mp3 -min 5
```

#### 15 o'clock alarm
```shell
$ goalarm -file ./bell.mp3 -time 15:00:00
```

#### looping 5 min timer
```shell
$ goalarm -file ./bell.mp3 -min 5 -loop
```

#### get time from alarm server by `get command`
```shell
$ goalarm -file ./bell.mp3 -min 5
get
{"status":"running","left":"4m58s","error":"","task":{"index":0,"range":"5m0s","name":"alarm"}}
```

#### describe commands and statuses.

```shell
$ goalarm -describe command
$ goalarm -describe status
```

#### start routine
```shell
$ goalarm -file ./bell.mp3 -routine '[{"range":20,"name":"working"},{"range":5,"name":"break"}]'
```

In the above case, after a 20-minute timer named `woriking` runs, a 5-minute timer named `break` runs.

## Author

komem3

## License

[MIT](./LICENSE)
